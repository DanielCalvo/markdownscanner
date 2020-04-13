package markdownlink

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

func (m *MarkdownLink) IsIgnored() bool {
	//If Link to an email
	if strings.HasPrefix(m.Destination, "mailto") {
		return true
	}
	//if link to a Github pull request
	if strings.Contains(m.Destination, "github.com") && strings.Contains(m.Destination, "/pull/") {
		return true
	}
	//If changelog file
	if strings.Contains(strings.ToLower(m.LocalFilePath), "changelog.md") {
		return true
	}
	//If maintainers file
	if strings.Contains(strings.ToLower(m.LocalFilePath), "maintainers.md") {
		return true
	}
	//If minutes file
	if strings.Contains(strings.ToLower(m.LocalFilePath), "minutes") {
		return true
	}
	//If inside a vendors folter
	if strings.Contains(strings.ToLower(m.LocalFilePath), "/vendor/") {
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

	mDestination := filepath.Dir(m.LocalFilePath) + string(os.PathSeparator) + m.Destination

	//Still can't check things like: /app_management/secrets_and_configmaps.md#secrets-from-files (yet!)
	if strings.HasPrefix(m.Destination, "#") || strings.Contains(m.Destination, "#") {
		m.Status = "N/A"
		return
	}

	_, err := os.Stat(mDestination)
	if os.IsNotExist(err) {
		m.Status = "404"
	} else {
		m.Status = "200"
	}

	if _, err := os.Stat(mDestination); os.IsNotExist(err) {
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

//func (m *MarkdownLink) CheckLink() {
//	if m.IsIgnored() {
//		m.SetIgnored()
//		return
//	}
//	if m.IsHTTP() {
//		m.CheckHTTP()
//		return
//	}
//	if m.IsFile() {
//		m.CheckFile()
//		return
//	}
//	m.Type = "UNKNOWN"
//	m.Status = "N/A"
//	return
//}
