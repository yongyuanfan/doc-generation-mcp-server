package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr                string
	APIPrefix               string
	MCPPath                 string
	RequestTimeout          time.Duration
	MCPServerName           string
	MCPServerVersion        string
	DocxTempDir             string
	DocxTemplateDir         string
	DocxDefaultAuthor       string
	DocxDefaultFont         string
	DocxDefaultFontSize     int
	DocxMaxRequestBodyBytes int64
	DocxMaxFileAge          time.Duration
	DocumentTypeTemplateMap map[string]string
}

func Load() (Config, error) {
	if err := loadDotEnvIfPresent(); err != nil {
		return Config{}, err
	}

	timeoutSeconds, err := intFromEnv("REQUEST_TIMEOUT_SECONDS", 120)
	if err != nil {
		return Config{}, err
	}

	fontSize, err := intFromEnv("DOCX_DEFAULT_FONT_SIZE", 22)
	if err != nil {
		return Config{}, err
	}

	maxRequestBytes, err := int64FromEnv("DOCX_MAX_REQUEST_BODY_BYTES", 2<<20)
	if err != nil {
		return Config{}, err
	}

	maxFileAgeMinutes, err := intFromEnv("DOCX_MAX_FILE_AGE_MINUTES", 60)
	if err != nil {
		return Config{}, err
	}

	tempDir := strings.TrimSpace(os.Getenv("DOCX_TEMP_DIR"))
	if tempDir == "" {
		tempDir = filepath.Join(os.TempDir(), "doc-generation-mcp-server")
	}

	templateDir := strings.TrimSpace(os.Getenv("DOCX_TEMPLATE_DIR"))
	if templateDir == "" {
		templateDir = filepath.Join(".", "templates")
	}

	cfg := Config{
		HTTPAddr:                envOrDefault("HTTP_ADDR", ":9103"),
		APIPrefix:               normalizePath(envOrDefault("API_PREFIX", "/api/v1")),
		MCPPath:                 normalizePath(envOrDefault("MCP_PATH", "/mcp")),
		RequestTimeout:          time.Duration(timeoutSeconds) * time.Second,
		MCPServerName:           envOrDefault("MCP_SERVER_NAME", "docx-generation-server"),
		MCPServerVersion:        envOrDefault("MCP_SERVER_VERSION", "0.1.0"),
		DocxTempDir:             tempDir,
		DocxTemplateDir:         templateDir,
		DocxDefaultAuthor:       envOrDefault("DOCX_DEFAULT_AUTHOR", "doc-generation-mcp-server"),
		DocxDefaultFont:         envOrDefault("DOCX_DEFAULT_FONT", "Calibri"),
		DocxDefaultFontSize:     fontSize,
		DocxMaxRequestBodyBytes: maxRequestBytes,
		DocxMaxFileAge:          time.Duration(maxFileAgeMinutes) * time.Minute,
		DocumentTypeTemplateMap: documentTypeTemplateMapFromEnv(),
	}

	return cfg, nil
}

func documentTypeTemplateMapFromEnv() map[string]string {
	mapping := map[string]string{
		"business_letter": "business-letter.docx",
	}
	raw := strings.TrimSpace(os.Getenv("DOCX_DOCUMENT_TYPE_TEMPLATES"))
	if raw == "" {
		return mapping
	}
	for _, item := range strings.Split(raw, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(item), "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		mapping[key] = value
	}
	return mapping
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func intFromEnv(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	if parsed <= 0 {
		return 0, fmt.Errorf("invalid %s: must be greater than 0", key)
	}
	return parsed, nil
}

func int64FromEnv(key string, fallback int64) (int64, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	if parsed <= 0 {
		return 0, fmt.Errorf("invalid %s: must be greater than 0", key)
	}
	return parsed, nil
}

func normalizePath(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "/"
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	return strings.TrimRight(trimmed, "/")
}

func loadDotEnvIfPresent() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.Open(filepath.Join(wd, ".env"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
