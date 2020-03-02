package markdownlink

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type MarkdownLink struct {
	File        string
	Name        string
	Destination string
	Type        string
	Status      string
}

func (m *MarkdownLink) IsHTTP() bool {
	_, err := url.ParseRequestURI(m.Destination)
	if err == nil {
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

func (m *MarkdownLink) IsEmail() bool {
	if strings.HasPrefix(m.Destination, "mailto") {
		return true
	}
	return false
}

//Don't forge to handle links in the docs that you may have to login into first (see if you can handle a forbidden or something)
//Double check all HTTP codes one by one to make sure you're handling them correctly
//Perhaps you need to handle timeouts!
func (m *MarkdownLink) CheckHTTP() {
	m.Type = "HTTP"
	resp, err := http.Head(m.Destination)
	if err != nil {
		m.Status = "BROKEN"
		return
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 208 || resp.StatusCode >= 300 && resp.StatusCode <= 308 || resp.StatusCode == 405 {
		m.Status = "OK"
	} else {
		m.Status = "BROKEN"
	}
}

func (m *MarkdownLink) CheckFile() {
	m.Type = "FILE"

	mDestination := filepath.Dir(m.File) + string(os.PathSeparator) + m.Destination

	//Still can't check things like: /app_management/secrets_and_configmaps.md#secrets-from-files (yet!)
	if strings.HasPrefix(m.Destination, "#") || strings.Contains(m.Destination, "#") {
		m.Status = "NOT IMPLEMENTED"
		return
	}

	_, err := os.Stat(mDestination)
	if os.IsNotExist(err) {
		m.Status = "BROKEN"
	} else {
		m.Status = "OK"
	}

	if _, err := os.Stat(mDestination); os.IsNotExist(err) {
		m.Status = "BROKEN"
	} else {
		m.Status = "OK"
	}
}

func (m *MarkdownLink) CheckEmail() {
	m.Type = "EMAIL"
	m.Status = "NOT IMPLEMENTED"
}

func (m *MarkdownLink) CheckLink() {
	switch {
	case m.IsHTTP():
		m.CheckHTTP()
	case m.IsFile():
		m.CheckFile()
	case m.IsEmail():
		m.CheckEmail()
	default:
		m.Type = "UNKNOWN"
		m.Status = "NOT IMPLEMENTED"
	}
}
