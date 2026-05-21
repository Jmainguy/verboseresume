package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testResumeJSON = `{
  "Name": "Test Person",
  "Location": "Remote",
  "Phone": "555-0100",
  "Email": "test@example.com",
  "Website": "https://example.com",
  "Github": "https://github.com/example",
  "Summary": "Senior platform engineer. Builds reliable platforms.",
  "Skills": {
    "Cloud": ["Kubernetes"]
  },
  "Experience": [
    {
      "Title": "Engineer",
      "Company": "Example Co",
      "Location": "Remote",
      "Date": "2024 - Present",
      "Details": ["Built platform tooling."]
    }
  ],
  "Education": [],
  "Certifications": []
}`

func TestUploadHandlerWrapsBuiltInTemplate(t *testing.T) {
	body := renderUpload(t, "", nil)

	for _, want := range []string{
		`class="app-shell app-header no-print"`,
		`WYSIWYG edit mode`,
		`Headers and footers`,
		`Test Person`,
		`id="template-switch"`,
		`id="paper-tone"`,
		`name="resume_json"`,
		`class="resume-template-compact"`,
		`Template: compact.html`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("rendered response did not contain %q", want)
		}
	}
}

func TestUploadHandlerSwitchesTemplateFromHiddenJSON(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("resume_json", testResumeJSON); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("template", "classic.html"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	uploadHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	for _, want := range []string{
		`class="resume-template-classic"`,
		`Template: classic.html`,
		`<option value="classic.html" selected>Classic</option>`,
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("rendered response did not contain %q", want)
		}
	}
}

func TestUploadHandlerWrapsCustomTemplate(t *testing.T) {
	customTemplate := []byte(`<style>.custom-resume { font-family: sans-serif; }</style><article class="custom-resume"><h1>{{ .Name }}</h1></article>`)
	body := renderUpload(t, "custom.html", customTemplate)

	for _, want := range []string{
		`class="app-shell app-header no-print"`,
		`WYSIWYG edit mode`,
		`Template: custom.html`,
		`class="custom-resume"`,
		`Test Person`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("rendered response did not contain %q", want)
		}
	}
}

func TestEmbeddedVerboseDocs(t *testing.T) {
	if len(strings.TrimSpace(exampleVerboseResumeMarkdown)) < 500 {
		t.Fatal("embedded example verbose resume is empty or too short; rebuild with current docs/")
	}
	if len(strings.TrimSpace(verboseResumeQuestionsMarkdown)) < 500 {
		t.Fatal("embedded verbose resume questions doc is empty or too short")
	}
	if !strings.Contains(exampleVerboseResumeMarkdown, "Alex Rivera") {
		t.Fatal("embedded example missing expected fictional name")
	}
}

func TestDocsEmbedMarkdownHandlers(t *testing.T) {
	for _, path := range []string{
		"/docs/embed/example-verbose-resume.md",
		"/docs/embed/verbose-resume-questions.md",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		switch path {
		case "/docs/embed/example-verbose-resume.md":
			exampleVerboseResumeHandler(rec, req)
		case "/docs/embed/verbose-resume-questions.md":
			verboseResumeQuestionsHandler(rec, req)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d: %s", path, rec.Code, rec.Body.String())
		}
		if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "text/plain") {
			t.Fatalf("%s: expected text/plain, got %q", path, ct)
		}
	}
	body := httptest.NewRecorder()
	exampleVerboseResumeHandler(body, httptest.NewRequest(http.MethodGet, "/docs/embed/example-verbose-resume.md", nil))
	if !strings.Contains(body.Body.String(), "Alex Rivera") {
		t.Fatal("example embed missing Alex Rivera")
	}
	q := httptest.NewRecorder()
	verboseResumeQuestionsHandler(q, httptest.NewRequest(http.MethodGet, "/docs/embed/verbose-resume-questions.md", nil))
	if !strings.Contains(q.Body.String(), "What company do you work at") {
		t.Fatal("questions embed missing starter prompt")
	}
}

func TestDocsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	docsHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	for _, want := range []string{
		`Verbose Resume`,
		`verboseResume.md`,
		`The verbose resume idea`,
		`verbose-resume-questions`,
		`example-verbose-resume`,
		`/docs/embed/example-verbose-resume.md`,
		`/docs/embed/verbose-resume-questions.md`,
		`Alex Rivera`,
		`get_verbose_resume_questions`,
		`MCP endpoint`,
		`Jonathan Mainguy`,
		`staff augmentation`,
		`Harborvane`,
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("rendered docs did not contain %q", want)
		}
	}
}

func TestContributeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/contribute", nil)
	rec := httptest.NewRecorder()

	contributeHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	body := rec.Body.String()
	for _, want := range []string{
		`Contribute`,
		`pull request`,
		`fork-only`,
		`do not add collaborators`,
		`New resume templates`,
		`templateOptions`,
		`GPLv2`,
		`https://github.com/Jmainguy/verboseresume`,
		`classic.html`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("rendered contribute page did not contain %q", want)
		}
	}
}

func TestFaviconHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	rec := httptest.NewRecorder()

	faviconHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %q", ct)
	}
	if len(rec.Body.Bytes()) == 0 {
		t.Fatal("expected non-empty favicon body")
	}
}

func TestMCPHandlerGETShowsMatrix(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	rec := httptest.NewRecorder()

	mcpHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	body := rec.Body.String()
	for _, want := range []string{
		`<canvas id="matrix"`,
		`POST /mcp`,
		`JSON-RPC`,
		`verboseResume.json`,
		`href="/static/brand/favicon.svg"`,
		`href="/favicon.ico"`,
		`get_resume_generator_guide`,
		`Give your agent`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("MCP GET response did not contain %q", want)
		}
	}
}

func TestNotFoundHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/definitely-not-a-route", nil)
	rec := httptest.NewRecorder()

	notFoundHandler(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		`404`,
		`Bullet not found`,
		`/definitely-not-a-route`,
		`Verbose Resume`,
		`ascii-drift`,
		`Enter the Matrix`,
		`GET /mcp`,
		`id="matrix-bg"`,
		`href="/static/brand/favicon.svg"`,
		`href="/favicon.ico"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("404 response did not contain %q", want)
		}
	}
}

func TestUploadFormHandlerUnknownPath404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nope", nil)
	rec := httptest.NewRecorder()

	uploadFormHandler(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestMCPHandlerListsAndCallsTools(t *testing.T) {
	listBody := postMCP(t, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	for _, want := range []string{
		`get_resume_generator_guide`,
		`get_verbose_resume_format`,
		`get_verbose_resume_questions`,
		`get_example_verbose_resume`,
		`get_upload_json_format`,
		`get_llm_prompt_guide`,
		`create_resume_artifact`,
	} {
		if !strings.Contains(listBody, want) {
			t.Fatalf("MCP tools/list response did not contain %q", want)
		}
	}

	callBody := postMCP(t, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_llm_prompt_guide"}}`)
	if !strings.Contains(callBody, `Use only facts from my verbose resume`) {
		t.Fatalf("MCP tools/call response did not include prompt guide: %s", callBody)
	}
}

func TestMCPHandlerCreatesResumeArtifact(t *testing.T) {
	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]any{
			"name": "create_resume_artifact",
			"arguments": map[string]any{
				"resume_json":   testResumeJSON,
				"template":      "compact.html",
				"filename_base": "test-person",
				"paper_tone":    "#f8fbff",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	body := postMCP(t, string(payload))
	var response struct {
		Result struct {
			StructuredContent struct {
				Markdown        string            `json:"markdown"`
				HTML            string            `json:"html"`
				ResumeJSON      string            `json:"resume_json"`
				PDFInstructions string            `json:"pdfInstructions"`
				SuggestedFiles  map[string]string `json:"suggestedFiles"`
			} `json:"structuredContent"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatal(err)
	}

	sc := response.Result.StructuredContent
	if !strings.Contains(sc.Markdown, "# Test Person") {
		t.Fatalf("markdown missing name: %q", sc.Markdown)
	}
	if !strings.Contains(sc.HTML, "<!DOCTYPE html>") || !strings.Contains(sc.HTML, "resume-template-compact") {
		t.Fatalf("html missing expected template wrapper: %q", sc.HTML[:120])
	}
	if !strings.Contains(sc.ResumeJSON, `"Name": "Test Person"`) {
		t.Fatalf("resume_json missing name: %q", sc.ResumeJSON)
	}
	if sc.SuggestedFiles["html"] != "test-person.html" {
		t.Fatalf("unexpected suggested html file %q", sc.SuggestedFiles["html"])
	}
	if !strings.Contains(sc.PDFInstructions, "test-person.html") || !strings.Contains(sc.PDFInstructions, "Save as PDF") {
		t.Fatalf("pdf instructions missing expected guidance: %q", sc.PDFInstructions)
	}
}

func postMCP(t *testing.T, body string) string {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mcpHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var decoded map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("MCP response was not valid JSON: %v", err)
	}

	return rec.Body.String()
}

func renderUpload(t *testing.T, customTemplateName string, customTemplate []byte) string {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	resumePart, err := writer.CreateFormFile("resume", "resume.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := resumePart.Write([]byte(testResumeJSON)); err != nil {
		t.Fatal(err)
	}

	if err := writer.WriteField("template", "compact.html"); err != nil {
		t.Fatal(err)
	}

	if len(customTemplate) > 0 {
		templatePart, err := writer.CreateFormFile("custom_template", customTemplateName)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := templatePart.Write(customTemplate); err != nil {
			t.Fatal(err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	uploadHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	return rec.Body.String()
}
