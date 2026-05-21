package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Resume struct {
	Name           string              `json:"Name"`
	Location       string              `json:"Location"`
	Phone          string              `json:"Phone"`
	Email          string              `json:"Email"`
	Website        string              `json:"Website"`
	Github         string              `json:"Github"`
	Summary        string              `json:"Summary"`
	Skills         map[string][]string `json:"Skills"`
	Experience     []Job               `json:"Experience"`
	Education      []Education         `json:"Education"`
	Certifications []string            `json:"Certifications"`
}

type Job struct {
	Title    string   `json:"Title"`
	Company  string   `json:"Company"`
	Location string   `json:"Location"`
	Date     string   `json:"Date"`
	Details  []string `json:"Details"`
}

type Education struct {
	Degree   string `json:"Degree"`
	School   string `json:"School"`
	Location string `json:"Location"`
	Date     string `json:"Date"`
}

type TemplateOption struct {
	File        string
	Name        string
	Description string
}

type SiteChrome struct {
	SiteName    string
	SiteURL     string
	GitHubRepo  string
	ActivePage  string
	FooterClass string
	FooterExtra string
}

type ContributeData struct {
	Chrome    SiteChrome
	Templates []TemplateOption
}

type FormData struct {
	Chrome      SiteChrome
	Templates   []TemplateOption
	SiteTagline string
}

type NotFoundData struct {
	Chrome      SiteChrome
	RequestPath string
}

type DocsData struct {
	Chrome                SiteChrome
	GuideMarkdown         string
	JSONSchemaMarkdown    string
	PromptMarkdown        string
	VerboseResumeMarkdown string
}

type ResumePageData struct {
	Chrome               SiteChrome
	Resume               Resume
	Content              template.HTML
	TemplateName         string
	SelectedTemplateName string
	Templates            []TemplateOption
	ResumeJSON           string
}

type UploadData struct {
	ResumeJSON         []byte
	TemplateName       string
	CustomTemplate     []byte
	CustomTemplateName string
}

type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type mcpToolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

var templateOptions = []TemplateOption{
	{
		File:        "classic.html",
		Name:        "Classic",
		Description: "A more traditional serif layout with a lighter visual style.",
	},
	{
		File:        "clean.html",
		Name:        "Clean",
		Description: "The default layout with skill pills and careful page-break rules.",
	},
	{
		File:        "compact.html",
		Name:        "Compact",
		Description: "A tighter layout for fitting more content on fewer pages.",
	},
}

const maxUploadBytes = 5 << 20

const siteName = "Verbose Resume"
const siteURL = "https://verboseresume.com"
const githubRepoURL = "https://github.com/Jmainguy/verboseresume"
const siteTagline = "Keep one detailed career record. Tailor it for every role."

const resumeGeneratorGuide = `# Verbose Resume guide

Verbose Resume (https://verboseresume.com) is a free community tool for turning structured resume JSON into a polished HTML resume that can be reviewed in the browser and printed to PDF.

The site is designed around two resume layers:

1. Keep a verbose resume in Markdown (for example verboseResume.md) as your private source of truth.
2. Use an LLM to tailor that verbose resume into a concise upload JSON for one job or audience.

Update the verbose resume continually, about once a month. People forget project details, production issues, metrics, and decisions quickly. A good agent can help by reviewing recent prompts, git history, Jira tickets, pull requests, docs, chat threads, and other work records, then suggesting factual additions under the correct role and Notable work headings.

Do not upload the verbose Markdown file to the site. Upload only the tailored JSON.

The app does not store uploaded resume JSON or custom templates. Uploads are processed in memory for the current request.`

const uploadJSONGuide = `# Upload JSON format

The upload file should be valid JSON with these top-level fields:

- Name: string
- Location: string
- Phone: string
- Email: string
- Website: string
- Github: string
- Summary: string. Keep it as a single concise tailored summary.
- Skills: object where each key is a display group and each value is an array of short skill labels.
- Experience: array of jobs with Title, Company, Location, Date, and Details.
- Education: array of schools with Degree, School, Location, and Date.
- Certifications: array of strings.

Keep the uploaded JSON concise. It is the final resume content, not the long source document.`

const llmPromptGuide = `# LLM workflow (verbose Markdown → upload JSON)

Keep verboseResume.md as the long source of truth. For each job, iterate on upload JSON with your assistant—do not maintain a separate tailored Markdown resume.

Loop:
1. Provide verboseResume.md, the job description, and the upload JSON spec.
2. Ask for only valid upload JSON. Save and upload to Verbose Resume to preview.
3. Paste the current JSON back with feedback; ask for revised JSON. Re-upload until satisfied.
4. Print or save PDF from the browser preview.

First pass prompt:

You are helping me tailor my resume for the job description below. Use only facts from my verbose resume (verboseResume.md). Do not invent employers, dates, tools, credentials, outcomes, metrics, citizenship status, clearance, or experience. Select the most relevant experience, compress it, and output only valid JSON matching the Verbose Resume upload format.

Prioritize relevance, truthful impact, short specific bullets, skills supported by experience, and a resume that fits cleanly into one to two pages.

Inputs: verbose resume, job description, upload JSON spec, any constraints I add.

Output only the final JSON.

Revision pass prompt:

Here is my current tailored resume JSON and feedback on the preview. Revise the JSON only. Keep facts truthful and supported by my verbose resume. Output only the complete updated JSON in the upload format—no commentary outside the JSON.`

func loadVerboseResumeFormatDoc() string {
	if strings.TrimSpace(verboseResumeFormatMarkdown) == "" {
		return "# Verbose Resume (Markdown)\n\nMaintain verboseResume.md as your source of truth.\n"
	}
	return verboseResumeFormatMarkdown
}

func loadVerboseResumeQuestionsDoc() string {
	if strings.TrimSpace(verboseResumeQuestionsMarkdown) == "" {
		return "# Sample questions\n\nSee docs/VERBOSE-RESUME-QUESTIONS.md\n"
	}
	return verboseResumeQuestionsMarkdown
}

func loadExampleVerboseResumeDoc() string {
	if strings.TrimSpace(exampleVerboseResumeMarkdown) == "" {
		return "# Example verbose resume\n\nSee docs/example-verbose-resume.md\n"
	}
	return exampleVerboseResumeMarkdown
}

func serveEmbeddedMarkdown(w http.ResponseWriter, content, filename string) {
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if strings.TrimSpace(content) == "" {
		http.Error(w, "embedded document missing; run: go build -o resumeGen .", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(content))
}

func renderTemplate(w http.ResponseWriter, templateName string, data any, failMsg string) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, templateName, data); err != nil {
		log.Printf("render %s: %v", templateName, err)
		http.Error(w, failMsg, http.StatusInternalServerError)
		return
	}
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("write %s response: %v", templateName, err)
	}
}

func summaryParagraphs(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return []string{strings.Join(strings.Fields(s), " ")}
}

const summaryRoleLead = "Senior platform engineer."

func hasSummaryRoleLead(p string) bool {
	return strings.HasPrefix(strings.TrimSpace(p), summaryRoleLead)
}

func summaryAfterRoleLead(p string) string {
	p = strings.TrimSpace(p)
	if !strings.HasPrefix(p, summaryRoleLead) {
		return p
	}
	return strings.TrimSpace(strings.TrimPrefix(p, summaryRoleLead))
}

func safeTemplateName(name string) string {
	for _, option := range templateOptions {
		if name == option.File {
			return name
		}
	}
	return templateOptions[0].File
}

var templates = loadTemplates()

func newSiteChrome(activePage, footerClass, footerExtra string) SiteChrome {
	return SiteChrome{
		SiteName:    siteName,
		SiteURL:     siteURL,
		GitHubRepo:  githubRepoURL,
		ActivePage:  activePage,
		FooterClass: footerClass,
		FooterExtra: footerExtra,
	}
}

func listenAddr() string {
	if addr := strings.TrimSpace(os.Getenv("LISTEN_ADDR")); addr != "" {
		return addr
	}
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

func main() {
	staticHandler, err := staticFileServer()
	if err != nil {
		log.Fatalf("static assets: %v", err)
	}

	http.HandleFunc("/", uploadFormHandler)
	http.HandleFunc("/docs", docsHandler)
	http.HandleFunc("/contribute", contributeHandler)
	http.HandleFunc("/docs/embed/verbose-resume-questions.md", verboseResumeQuestionsHandler)
	http.HandleFunc("/docs/embed/example-verbose-resume.md", exampleVerboseResumeHandler)
	http.HandleFunc("/mcp", mcpHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/static/", staticHandler)

	addr := listenAddr()
	if strings.TrimSpace(exampleVerboseResumeMarkdown) == "" || strings.TrimSpace(verboseResumeQuestionsMarkdown) == "" {
		log.Printf("warning: embedded docs markdown is empty; run: go build -o resumeGen . && ./resumeGen")
	}
	log.Printf("Listening on %s (docs: /docs) ...", displayListenURL(addr))
	log.Fatal(http.ListenAndServe(addr, nil))
}

func displayListenURL(addr string) string {
	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr
	}
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return addr
	}
	return "http://" + addr
}

func uploadFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	renderTemplate(w, "form.html", FormData{
		Chrome:      newSiteChrome("home", "", ""),
		Templates:   templateOptions,
		SiteTagline: siteTagline,
	}, "Failed to render upload form")
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	renderTemplate(w, "docs.html", DocsData{
		Chrome:                newSiteChrome("docs", "", ""),
		GuideMarkdown:         resumeGeneratorGuide,
		JSONSchemaMarkdown:    uploadJSONGuide,
		PromptMarkdown:        llmPromptGuide,
		VerboseResumeMarkdown: loadVerboseResumeFormatDoc(),
	}, "Failed to render docs")
}

func contributeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	renderTemplate(w, "contribute.html", ContributeData{
		Chrome:    newSiteChrome("contribute", "", ""),
		Templates: templateOptions,
	}, "Failed to render contribute page")
}

func verboseResumeQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/docs/embed/verbose-resume-questions.md" {
		notFoundHandler(w, r)
		return
	}
	serveEmbeddedMarkdown(w, loadVerboseResumeQuestionsDoc(), "verbose-resume-questions.md")
}

func exampleVerboseResumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/docs/embed/example-verbose-resume.md" {
		notFoundHandler(w, r)
		return
	}
	serveEmbeddedMarkdown(w, loadExampleVerboseResumeDoc(), "example-verbose-resume.md")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	data := NotFoundData{
		Chrome:      newSiteChrome("404", "", ""),
		RequestPath: r.URL.Path,
	}
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, "404.html", data); err != nil {
		log.Printf("render 404.html: %v", err)
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	_, _ = buf.WriteTo(w)
}

func mcpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "mcp-matrix.html", nil, "Failed to render MCP easter egg")
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "MCP endpoint expects POST JSON-RPC requests", http.StatusMethodNotAllowed)
		return
	}

	var req mcpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMCPError(w, nil, -32700, "Parse error")
		return
	}

	if req.Method == "notifications/initialized" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch req.Method {
	case "initialize":
		writeMCPResult(w, req.ID, map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "verboseresume",
				"version": "0.1.0",
			},
		})
	case "tools/list":
		writeMCPResult(w, req.ID, map[string]any{
			"tools": []map[string]any{
				mcpTool("get_resume_generator_guide", "Explains Verbose Resume and the Markdown verbose-resume workflow."),
				mcpTool("get_verbose_resume_format", "Returns the Markdown structure for verboseResume.md (source of truth)."),
				mcpTool("get_verbose_resume_questions", "Returns sample interview questions and an LLM interviewer prompt for filling verboseResume.md."),
				mcpTool("get_example_verbose_resume", "Returns a fictional, very verbose example verboseResume.md for structure reference."),
				mcpTool("get_upload_json_format", "Returns the upload JSON format expected by Verbose Resume."),
				mcpTool("get_llm_prompt_guide", "Returns prompt guidance for turning a verbose Markdown resume into final upload JSON."),
				mcpCreateArtifactTool(),
			},
		})
	case "tools/call":
		handleMCPToolCall(w, req)
	default:
		writeMCPError(w, req.ID, -32601, "Method not found")
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize())

	upload, err := parseUpload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if upload.TemplateName == "custom" && len(upload.CustomTemplate) == 0 {
		http.Error(w, "Choose an HTML file for the Custom template.", http.StatusBadRequest)
		return
	}

	var resume Resume
	if err := json.Unmarshal(upload.ResumeJSON, &resume); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	templateName := upload.TemplateName
	if templateName != "custom" {
		templateName = safeTemplateName(templateName)
	} else {
		templateName = templateOptions[0].File
	}

	content, displayTemplateName, err := renderResumeContent(resume, templateName, upload.CustomTemplateName, upload.CustomTemplate)
	if err != nil {
		http.Error(w, "Failed to render resume template", http.StatusBadRequest)
		return
	}

	pageData := ResumePageData{
		Chrome:               newSiteChrome("resume", "app-shell app-footer no-print", ""),
		Resume:               resume,
		Content:              content,
		TemplateName:         displayTemplateName,
		SelectedTemplateName: templateName,
		Templates:            templateOptions,
		ResumeJSON:           string(upload.ResumeJSON),
	}
	renderTemplate(w, "resume-page.html", pageData, "Failed to render resume page")
}

func parseUpload(r *http.Request) (UploadData, error) {
	reader, err := r.MultipartReader()
	if err != nil {
		return UploadData{}, fmt.Errorf("expected multipart form upload")
	}

	var upload UploadData
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return UploadData{}, fmt.Errorf("failed to read upload")
		}

		switch part.FormName() {
		case "resume":
			data, err := readPart(part)
			if err != nil {
				return UploadData{}, fmt.Errorf("failed to read resume JSON")
			}
			upload.ResumeJSON = data
		case "resume_json":
			data, err := readPart(part)
			if err != nil {
				return UploadData{}, fmt.Errorf("failed to read resume JSON")
			}
			if strings.TrimSpace(string(data)) != "" {
				upload.ResumeJSON = data
			}
		case "template":
			data, err := readPart(part)
			if err != nil {
				return UploadData{}, fmt.Errorf("failed to read template selection")
			}
			upload.TemplateName = strings.TrimSpace(string(data))
		case "custom_template":
			data, err := readPart(part)
			if err != nil {
				return UploadData{}, fmt.Errorf("failed to read custom template")
			}
			if strings.TrimSpace(string(data)) == "" {
				continue
			}
			upload.CustomTemplate = data
			upload.CustomTemplateName = safeUploadFileName(part.FileName())
		default:
			_, _ = io.Copy(io.Discard, part)
		}
	}

	if len(upload.ResumeJSON) == 0 {
		return UploadData{}, fmt.Errorf("resume JSON is required")
	}

	return upload, nil
}

func readPart(part *multipart.Part) ([]byte, error) {
	defer func() { _ = part.Close() }()
	return io.ReadAll(part)
}

func safeUploadFileName(filename string) string {
	name := filepath.Base(filename)
	if name == "." || name == string(filepath.Separator) {
		name = "custom template"
	}
	return name
}

func renderResumeContent(resume Resume, templateName, customTemplateName string, customTemplateData []byte) (template.HTML, string, error) {
	var buf bytes.Buffer
	if len(customTemplateData) > 0 {
		tmpl, err := template.New("custom-resume").Funcs(templateFuncs()).Parse(string(customTemplateData))
		if err != nil {
			return "", "", err
		}
		if err := executeCustomTemplate(tmpl, customTemplateName, &buf, resume); err != nil {
			return "", "", err
		}
		return template.HTML(buf.String()), customTemplateName, nil
	}

	if err := templates.ExecuteTemplate(&buf, templateName, resume); err != nil {
		return "", "", err
	}
	return template.HTML(buf.String()), templateName, nil
}

func executeCustomTemplate(tmpl *template.Template, customTemplateName string, buf *bytes.Buffer, resume Resume) error {
	if named := tmpl.Lookup(customTemplateName); named != nil {
		return named.Execute(buf, resume)
	}
	return tmpl.Execute(buf, resume)
}

func mcpTool(name, description string) map[string]any {
	return map[string]any{
		"name":        name,
		"description": description,
		"inputSchema": map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
	}
}

func mcpCreateArtifactTool() map[string]any {
	return map[string]any{
		"name":        "create_resume_artifact",
		"description": "Create editable resume artifacts from upload-format JSON.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"resume_json": map[string]any{
					"type":        "string",
					"description": "Resume JSON matching the Verbose Resume upload format.",
				},
				"template": map[string]any{
					"type":        "string",
					"description": "Built-in template file: classic.html, clean.html, or compact.html. Defaults to classic.html.",
				},
				"paper_tone": map[string]any{
					"type":        "string",
					"description": "Optional paper background color hex for standalone HTML.",
				},
				"filename_base": map[string]any{
					"type":        "string",
					"description": "Optional base name for suggested local files. Defaults to resume.",
				},
			},
			"required": []string{"resume_json"},
		},
	}
}

func handleMCPToolCall(w http.ResponseWriter, req mcpRequest) {
	var params mcpToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		writeMCPError(w, req.ID, -32602, "Invalid params")
		return
	}

	var text string
	switch params.Name {
	case "get_resume_generator_guide":
		text = resumeGeneratorGuide
	case "get_verbose_resume_format":
		text = loadVerboseResumeFormatDoc()
	case "get_verbose_resume_questions":
		text = loadVerboseResumeQuestionsDoc()
	case "get_example_verbose_resume":
		text = loadExampleVerboseResumeDoc()
	case "get_upload_json_format":
		text = uploadJSONGuide
	case "get_llm_prompt_guide":
		text = llmPromptGuide
	case "create_resume_artifact":
		handleMCPCreateArtifact(w, req.ID, params.Arguments)
		return
	default:
		writeMCPError(w, req.ID, -32602, "Unknown tool")
		return
	}

	writeMCPResult(w, req.ID, map[string]any{
		"content": []map[string]string{
			{"type": "text", "text": text},
		},
	})
}

func handleMCPCreateArtifact(w http.ResponseWriter, id any, args map[string]any) {
	rawResume, ok := args["resume_json"].(string)
	if !ok || strings.TrimSpace(rawResume) == "" {
		writeMCPError(w, id, -32602, "create_resume_artifact requires a resume_json string argument")
		return
	}

	var resume Resume
	if err := json.Unmarshal([]byte(rawResume), &resume); err != nil {
		writeMCPError(w, id, -32602, "resume_json is not valid upload-format JSON")
		return
	}

	templateName := templateOptions[0].File
	if requested, ok := args["template"].(string); ok && strings.TrimSpace(requested) != "" {
		templateName = safeTemplateName(requested)
	}

	paperTone := "#ffffff"
	if requested, ok := args["paper_tone"].(string); ok && strings.TrimSpace(requested) != "" {
		paperTone = safePaperTone(requested)
	}

	baseName := "resume"
	if requested, ok := args["filename_base"].(string); ok && strings.TrimSpace(requested) != "" {
		baseName = safeArtifactBaseName(requested)
	}

	html, err := renderStandaloneResumeHTML(resume, templateName, paperTone)
	if err != nil {
		writeMCPError(w, id, -32602, "Failed to render resume HTML artifact")
		return
	}

	normalizedJSON, err := json.MarshalIndent(resume, "", "  ")
	if err != nil {
		writeMCPError(w, id, -32602, "Failed to normalize resume JSON")
		return
	}

	markdown := renderResumeMarkdown(resume)
	suggestedFiles := map[string]string{
		"markdown": baseName + ".md",
		"html":     baseName + ".html",
		"json":     baseName + ".json",
	}
	pdfInstructions := artifactPDFInstructions(baseName, suggestedFiles)

	writeMCPResult(w, id, map[string]any{
		"content": []map[string]string{
			{
				"type": "text",
				"text": "Created editable resume artifacts. Save structuredContent.markdown and structuredContent.html locally, revise them with prompts if needed, then generate a PDF using structuredContent.pdfInstructions.",
			},
		},
		"structuredContent": map[string]any{
			"markdown":        markdown,
			"html":            html,
			"resume_json":     string(normalizedJSON),
			"pdfInstructions": pdfInstructions,
			"suggestedFiles":  suggestedFiles,
			"template":        templateName,
			"paper_tone":      paperTone,
		},
	})
}

func renderResumeMarkdown(resume Resume) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", resume.Name)
	fmt.Fprintf(&b, "%s | %s | %s\n\n", resume.Location, resume.Phone, resume.Email)
	if resume.Website != "" || resume.Github != "" {
		fmt.Fprintf(&b, "%s | %s\n\n", resume.Website, resume.Github)
	}
	if s := strings.TrimSpace(resume.Summary); s != "" {
		b.WriteString("## Summary\n\n")
		b.WriteString(s)
		b.WriteString("\n\n")
	}
	if len(resume.Skills) > 0 {
		b.WriteString("## Skills\n\n")
		groups := make([]string, 0, len(resume.Skills))
		for group := range resume.Skills {
			groups = append(groups, group)
		}
		sort.Strings(groups)
		for _, group := range groups {
			fmt.Fprintf(&b, "### %s\n\n", group)
			for _, skill := range resume.Skills[group] {
				fmt.Fprintf(&b, "- %s\n", skill)
			}
			b.WriteString("\n")
		}
	}
	if len(resume.Experience) > 0 {
		b.WriteString("## Experience\n\n")
		for _, job := range resume.Experience {
			fmt.Fprintf(&b, "### %s — %s\n\n", job.Company, job.Title)
			fmt.Fprintf(&b, "%s | %s\n\n", job.Date, job.Location)
			for _, detail := range job.Details {
				fmt.Fprintf(&b, "- %s\n", detail)
			}
			b.WriteString("\n")
		}
	}
	if len(resume.Education) > 0 {
		b.WriteString("## Education\n\n")
		for _, edu := range resume.Education {
			fmt.Fprintf(&b, "- %s | %s | %s | %s\n", edu.Degree, edu.School, edu.Location, edu.Date)
		}
		b.WriteString("\n")
	}
	if len(resume.Certifications) > 0 {
		b.WriteString("## Certifications\n\n")
		for _, cert := range resume.Certifications {
			fmt.Fprintf(&b, "- %s\n", cert)
		}
	}
	return b.String()
}

func renderStandaloneResumeHTML(resume Resume, templateName, paperTone string) (string, error) {
	content, _, err := renderResumeContent(resume, templateName, "", nil)
	if err != nil {
		return "", err
	}

	title := template.HTMLEscapeString(resume.Name)
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s - Resume</title>
  <style>
    :root {
      --resume-paper: %s;
    }
    body {
      margin: 0;
      background: var(--resume-paper);
      print-color-adjust: exact;
      -webkit-print-color-adjust: exact;
    }
    main {
      background: var(--resume-paper);
      print-color-adjust: exact;
      -webkit-print-color-adjust: exact;
    }
    @page {
      background: var(--resume-paper);
    }
    @media print {
      html, body, main {
        background: var(--resume-paper) !important;
        print-color-adjust: exact;
        -webkit-print-color-adjust: exact;
      }
    }
  </style>
</head>
<body>
<main>
%s
</main>
</body>
</html>`, title, paperTone, string(content)), nil
}

func artifactPDFInstructions(baseName string, files map[string]string) string {
	return fmt.Sprintf(`Save the artifacts locally, edit them with prompts if needed, then create a PDF.

Suggested files:
- %s: prompt-editable resume source
- %s: printable layout using the selected template
- %s: upload-format JSON for Verbose Resume

Workflow:
1. Write structuredContent.markdown to %s for quick text edits.
2. Write structuredContent.html to %s for layout and print styling.
3. Revise either file locally. Markdown is best for wording; HTML is best for spacing and print layout.
4. Generate a PDF from the HTML file.

PDF options:
- Open %s in Chrome and use Print -> Save as PDF. Disable browser headers and footers.
- Or run headless Chromium locally:
  chromium --headless --disable-gpu --print-to-pdf=%s.pdf %s

To iterate again, edit %s or %s, regenerate the PDF, and only re-upload JSON to Verbose Resume if you want the hosted web preview.`,
		files["markdown"],
		files["html"],
		files["json"],
		files["markdown"],
		files["html"],
		files["html"],
		baseName,
		files["html"],
		files["markdown"],
		files["html"],
	)
}

func safePaperTone(value string) string {
	allowed := map[string]struct{}{
		"#ffffff": {},
		"#fbfaf6": {},
		"#f7f8f2": {},
		"#f8fbff": {},
		"#f4f8fb": {},
		"#f3f7f7": {},
		"#f6f8fb": {},
	}
	value = strings.TrimSpace(value)
	if _, ok := allowed[value]; ok {
		return value
	}
	return "#ffffff"
}

func safeArtifactBaseName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "resume"
	}
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		}
	}
	if b.Len() == 0 {
		return "resume"
	}
	return b.String()
}

func writeMCPResult(w http.ResponseWriter, id any, result any) {
	writeJSON(w, map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	})
}

func writeMCPError(w http.ResponseWriter, id any, code int, message string) {
	writeJSON(w, map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"summaryParagraphs":    summaryParagraphs,
		"hasSummaryRoleLead":   hasSummaryRoleLead,
		"summaryAfterRoleLead": summaryAfterRoleLead,
	}
}
