package mdscanner

import "testing"

var TestData = []struct {
	name        string
	destination string
	desiredType string
}{
	{
		name:        "HTTP URL",
		destination: "http://example.com",
		desiredType: "http",
	},
	{
		name:        "HTTPS URL",
		destination: "https://example.com",
		desiredType: "http",
	},
	{
		name:        "HTTPS URL with path",
		destination: "https://example.com/something",
		desiredType: "http",
	},
	{
		name:        "Relative file path",
		destination: "./file.md",
		desiredType: "file",
	},
	{
		name:        "Absolute file path",
		destination: "/file.md",
		desiredType: "file",
	},
	{
		name:        "Filename with extension",
		destination: "file.md",
		desiredType: "file",
	},
	{
		name:        "Changelog file 1",
		destination: "./CHANGELOG.md",
		desiredType: "ignored",
	},
	{
		name:        "Changelog file 2",
		destination: "./changelog.md",
		desiredType: "ignored",
	},
	{
		name:        "Maintainers file",
		destination: "./MAINTAINERS.md",
		desiredType: "ignored",
	},
	{
		name:        "Minutes file",
		destination: "./2022-01-01-minutes.md",
		desiredType: "ignored",
	},
	{
		name:        "Email link",
		destination: "mailto:user@example.com",
		desiredType: "ignored",
	},
	{
		name:        "Vendors folder",
		destination: "./vendor/somefile.md",
		desiredType: "ignored",
	},
	{
		name:        "Releases folder",
		destination: "./releases/somefile.md",
		desiredType: "ignored",
	},
	{
		name:        "GitHub pull request",
		destination: "https://github.com/user/repo/pull/123",
		desiredType: "ignored",
	},
	{
		name:        "GitHub issue",
		destination: "https://github.com/user/repo/issues/456",
		desiredType: "ignored",
	},
	{
		name:        "Anchor link",
		destination: "#section-title",
		desiredType: "unknown",
	},
	{
		name:        "File with anchor",
		destination: "file.md#section-title",
		desiredType: "file",
	},
	{
		name:        "FTP URL",
		destination: "ftp://example.com",
		desiredType: "unknown",
	},
}

func TestMarkdownlinkTestLinkType(t *testing.T) {
	for _, tt := range TestData {
		t.Run(tt.name, func(t *testing.T) {
			mdLink := &MarkdownLink{
				Destination: tt.destination,
			}
			if !mdLink.IsHTTP() && tt.desiredType == "http" {
				t.Errorf("Mismatching isHTTP(), got isHTTP() == true on: %v", tt.destination)
			}
			if !mdLink.IsFile() && tt.desiredType == "file" {
				t.Errorf("Mismatching isFile(), got isFile() == true on: %v", tt.destination)
			}
			if !mdLink.IsIgnored() && tt.desiredType == "ignored" {
				t.Errorf("Mismatching isIgnored(), got isIgnored() == true on: %v", tt.destination)
			}
		})
	}
}

var markdown string = `Links to a file:
- [file.md on the testdata subfolder](./testdata/file.md)
- [some title](./testdata/file.md##Some Title!)
`

//How do I test checkFile and checkHTTP?

func TestCheckHTTP(t *testing.T) {
	m := MarkdownLink{
		FileName:      "",
		LocalFilePath: "",
		HTTPFilePath:  "",
		Name:          "",
		Destination:   "testdata/file.md",
		Type:          "",
		Status:        "",
	}

}
