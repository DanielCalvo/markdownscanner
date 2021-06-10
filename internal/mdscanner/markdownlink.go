package mdscanner

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type MarkdownLink struct {
	FileName      string
	LocalFilePath string
	HTTPFilePath  string
	Name          string
	Destination   string
	Type          string
	Status        string
}

func (m *MarkdownLink) IsHTTP() bool {
	_, err := url.ParseRequestURI(m.Destination)
	if err == nil && strings.HasPrefix(m.Destination, "http") {
		return true
	}
	return false
}

func (m *MarkdownLink) IsFile() bool {
	if strings.HasPrefix(m.Destination, ".") || strings.HasPrefix(m.Destination, "/") {
		return true
	}
	matched, _ := regexp.Match(`[A-Za-z0-9_]+\..+`, []byte(m.Destination))
	if matched {
		return true
	}
	return false
}

//It's starting to look like you could use a map here
func (m *MarkdownLink) IsIgnored() bool {

	//If changelog file
	if strings.Contains(strings.ToLower(m.LocalFilePath), "changelog") {
		return true
	}
	//If maintainers file
	if strings.Contains(strings.ToLower(m.LocalFilePath), "maintainers") {
		return true
	}
	//If minutes file
	if strings.Contains(strings.ToLower(m.LocalFilePath), "minutes") {
		return true
	}
	//If Link to an email
	if strings.HasPrefix(m.Destination, "mailto") {
		return true
	}
	//If inside a vendors folder
	if strings.Contains(strings.ToLower(m.LocalFilePath), "/vendor/") {
		return true
	}

	if strings.Contains(strings.ToLower(m.LocalFilePath), "/releases/") {
		return true
	}

	//if link to a Github pull request
	if strings.Contains(m.Destination, "github.com") && strings.Contains(m.Destination, "/pull/") {
		return true
	}
	//if link to a Github issue
	if strings.Contains(m.Destination, "github.com") && strings.Contains(m.Destination, "/issues/") {
		return true
	}

	return false
}

func (m *MarkdownLink) CheckHTTP() {
	m.Type = "HTTP"
	resp, err := http.Head(m.Destination)
	if err != nil {
		m.Status = "ERR"
		return
	}
	m.Status = strconv.Itoa(resp.StatusCode)
}

func (m *MarkdownLink) CheckFile() {
	m.Type = "FILE"

	dir := filepath.Dir(m.LocalFilePath)
	err := os.Chdir(dir)
	if err != nil {
		m.Status = "ERR"
		return
	}

	//Still can't check things like: /app_management/secrets_and_configmaps.md#secrets-from-files (yet!)
	if strings.HasPrefix(m.Destination, "#") || strings.Contains(m.Destination, "#") {
		m.Status = "N/A"
		return
	}

	if _, err := os.Stat(m.Destination); os.IsNotExist(err) {
		m.Status = "404"
	} else {
		m.Status = "200"
	}
}

func (m *MarkdownLink) SetIgnored() {
	m.Type = "IGNORED"
	m.Status = "IGNORED"
}

func (m *MarkdownLink) CheckLink() {
	switch {
	case m.IsIgnored():
		m.SetIgnored()
	case m.IsHTTP():
		m.CheckHTTP()
	case m.IsFile():
		m.CheckFile()
	default:
		m.Type = "UNKNOWN"
		m.Status = "N/A"
	}
}
