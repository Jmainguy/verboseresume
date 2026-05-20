# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/verboseresume .

FROM cgr.dev/chainguard/static:latest
ENV PORT=8080
EXPOSE 8080
COPY --from=build /out/verboseresume /verboseresume
USER 65532:65532
ENTRYPOINT ["/verboseresume"]
