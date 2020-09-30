package mdscanner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DanielCalvo/markdownscanner/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type GithubProjectApiResponse []struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	HTMLURL          string    `json:"html_url"`
	Description      string    `json:"description"`
	Fork             bool      `json:"fork"`
	URL              string    `json:"url"`
	ForksURL         string    `json:"forks_url"`
	KeysURL          string    `json:"keys_url"`
	CollaboratorsURL string    `json:"collaborators_url"`
	TeamsURL         string    `json:"teams_url"`
	HooksURL         string    `json:"hooks_url"`
	IssueEventsURL   string    `json:"issue_events_url"`
	EventsURL        string    `json:"events_url"`
	AssigneesURL     string    `json:"assignees_url"`
	BranchesURL      string    `json:"branches_url"`
	TagsURL          string    `json:"tags_url"`
	BlobsURL         string    `json:"blobs_url"`
	GitTagsURL       string    `json:"git_tags_url"`
	GitRefsURL       string    `json:"git_refs_url"`
	TreesURL         string    `json:"trees_url"`
	StatusesURL      string    `json:"statuses_url"`
	LanguagesURL     string    `json:"languages_url"`
	StargazersURL    string    `json:"stargazers_url"`
	ContributorsURL  string    `json:"contributors_url"`
	SubscribersURL   string    `json:"subscribers_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	CommitsURL       string    `json:"commits_url"`
	GitCommitsURL    string    `json:"git_commits_url"`
	CommentsURL      string    `json:"comments_url"`
	IssueCommentURL  string    `json:"issue_comment_url"`
	ContentsURL      string    `json:"contents_url"`
	CompareURL       string    `json:"compare_url"`
	MergesURL        string    `json:"merges_url"`
	ArchiveURL       string    `json:"archive_url"`
	DownloadsURL     string    `json:"downloads_url"`
	IssuesURL        string    `json:"issues_url"`
	PullsURL         string    `json:"pulls_url"`
	MilestonesURL    string    `json:"milestones_url"`
	NotificationsURL string    `json:"notifications_url"`
	LabelsURL        string    `json:"labels_url"`
	ReleasesURL      string    `json:"releases_url"`
	DeploymentsURL   string    `json:"deployments_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PushedAt         time.Time `json:"pushed_at"`
	GitURL           string    `json:"git_url"`
	SSHURL           string    `json:"ssh_url"`
	CloneURL         string    `json:"clone_url"`
	SvnURL           string    `json:"svn_url"`
	Homepage         string    `json:"homepage"`
	Size             int       `json:"size"`
	StargazersCount  int       `json:"stargazers_count"`
	WatchersCount    int       `json:"watchers_count"`
	Language         string    `json:"language"`
	HasIssues        bool      `json:"has_issues"`
	HasProjects      bool      `json:"has_projects"`
	HasDownloads     bool      `json:"has_downloads"`
	HasWiki          bool      `json:"has_wiki"`
	HasPages         bool      `json:"has_pages"`
	ForksCount       int       `json:"forks_count"`
	MirrorURL        string    `json:"mirror_url"`
	Archived         bool      `json:"archived"`
	Disabled         bool      `json:"disabled"`
	OpenIssuesCount  int       `json:"open_issues_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
		NodeID string `json:"node_id"`
	} `json:"license"`
	Forks         int    `json:"forks"`
	OpenIssues    int    `json:"open_issues"`
	Watchers      int    `json:"watchers"`
	DefaultBranch string `json:"default_branch"`
	Permissions   struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

type Repository struct {
	URL string
	//Leave a comment here explaining the bucket path!
	Name               string
	FilesystemPath     string
	LastScanned        string
	LinksScanned       int
	Links404           int
	JSONReportPath     string
	HTMLReportPath     string
	MetadataReportPath string
	MarkdownFiles      []MarkdownFile
	MarkdownLinks      []MarkdownLink
}

type MarkdownFile struct {
	FileName string
	FilePath string
	HTTPAddr string
}

type RepositoriesYaml struct {
	Repositories []string `yaml:"repositories"`
	Projects     []string `yaml:"projects"`
}

func NewRepository(c *config.Config, repoURL string) (Repository, error) {
	var r Repository

	url, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return r, err
	}

	//You can populate other fields in here as you progress with the project, such as Last Scanned
	r = Repository{
		URL:                repoURL,
		Name:               url.Path,
		FilesystemPath:     c.Filesystem.TmpFolder + url.Path,
		LastScanned:        "",
		LinksScanned:       0,
		Links404:           0,
		JSONReportPath:     url.Path + ".json",
		HTMLReportPath:     url.Path + ".html",
		MetadataReportPath: c.Filesystem.ScanMetadataFolder + "/" + GetRepoNameWithUnderscores(url.Path) + ".json",
		MarkdownFiles:      []MarkdownFile{},
		MarkdownLinks:      []MarkdownLink{},
	}

	return r, nil
}

func ValidateGitRepository(repository string) error {
	_, err := url.ParseRequestURI(repository)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(repository, "http") {
		return errors.New("Repository URL does not start with http(s)")
	}
	return nil
}

//leave comment about how scan metadata is just a subset of the Repository struct, except files and links
func SaveScanMetadata(r Repository) error {
	r.MarkdownFiles = nil
	r.MarkdownLinks = nil

	err := SaveStructToJson(r, r.MetadataReportPath)
	if err != nil {
		return err
	}
	return nil
}

func GetMarkdownFiles(r *Repository) {

	//validate if file is a broken symlink somehow (&& mdscanner.DoesExist(file.Name())
	err := filepath.Walk(r.FilesystemPath, func(path string, file os.FileInfo, err error) error {
		if strings.HasSuffix(file.Name(), ".md") {
			s := strings.Split(path, r.Name)
			mdFile := MarkdownFile{
				FilePath: path,
				FileName: file.Name(),
				HTTPAddr: r.URL + "/tree/master" + s[1],
			}
			r.MarkdownFiles = append(r.MarkdownFiles, mdFile)
		}
		return err
	})
	if err != nil {
		log.Printf("Error walking the path %q: %v\n", err)
	}
}

func GetMarkdownLinksFromFiles(r *Repository) {

	for _, mdFile := range r.MarkdownFiles {
		fileContents, err := ioutil.ReadFile(mdFile.FilePath)
		if err != nil {
			continue
		}

		//regex for footnote style MarkdownLinks
		re := regexp.MustCompile(`(\[.+\])\s*:\s*(.+)`)
		for _, matchedMarkdownLink := range re.FindAllStringSubmatch(string(fileContents), -1) {
			mdLink := MarkdownLink{
				FileName:      mdFile.FileName,
				LocalFilePath: mdFile.FilePath,
				HTTPFilePath:  mdFile.HTTPAddr,
				Name:          matchedMarkdownLink[1],
				Destination:   matchedMarkdownLink[2],
			}
			r.MarkdownLinks = append(r.MarkdownLinks, mdLink)
		}

		//regex for inline style links
		re = regexp.MustCompile(`(\[.+?\])((\()(.+?)(\)))`)
		for _, matchedMarkdownLink := range re.FindAllStringSubmatch(string(fileContents), -1) {
			mdLink := MarkdownLink{
				FileName:      mdFile.FileName,
				LocalFilePath: mdFile.FilePath,
				HTTPFilePath:  mdFile.HTTPAddr,
				Name:          matchedMarkdownLink[1],
				Destination:   matchedMarkdownLink[4],
			}
			r.MarkdownLinks = append(r.MarkdownLinks, mdLink)
		}
	}
}

//You couldn't figure out why an in-place operation in mdLink.CheckLink did not get the value out of this function
func CheckMarkdownLinksWithSleep(r *Repository, sleepTime time.Duration) []MarkdownLink {
	var scannedLinks []MarkdownLink

	for _, mdLink := range r.MarkdownLinks {
		mdLink.CheckLink()
		if mdLink.Type == "HTTP" {
			time.Sleep(sleepTime)
		}
		fmt.Println("Checked:", mdLink)
		scannedLinks = append(scannedLinks, mdLink)
		log.Println("Checked:", mdLink)
	}
	return scannedLinks
}

//func GetRepositoryUrlsFromYaml(yamlPath string) ([]string, error) {
//	var repositoriesYaml RepositoriesYaml
//	var repositoryUrls []string
//
//	repositoriesYamlFile, err := ioutil.ReadFile(yamlPath)
//	if err != nil {
//		return nil, err
//	}
//
//	err = yaml.Unmarshal(repositoriesYamlFile, &repositoriesYaml)
//	if err != nil {
//		return nil, err
//	}
//
//	fmt.Println(repositoriesYaml.Projects)
//	for _, githubProject := range repositoriesYaml.Projects {
//		repos, err := GetReposFromProject(githubProject)
//		if err == nil {
//			repositoryUrls = append(repositoryUrls, repos...)
//		} else {
//			log.Println("Error getting repositories for:"+githubProject+":", err)
//		}
//	}
//
//	for _, githubRepo := range repositoriesYaml.Repositories {
//		err = ValidateGitRepository(githubRepo)
//		if err == nil {
//			repositoryUrls = append(repositoryUrls, githubRepo)
//		} else {
//			log.Println("Error validating repository URL:"+githubRepo+":", err)
//		}
//	}
//	return repositoryUrls, nil
//}
//

//go-git can't clone large repositories without using very large amounts of memory: https://github.com/src-d/go-git/issues/761
func CloneRepository(r Repository) error {
	if !DoesExist(r.FilesystemPath) {
		cmd := exec.Command("git", "clone", r.URL, r.FilesystemPath)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteRepository(r Repository) error {
	err := os.RemoveAll(r.FilesystemPath)
	if err != nil {
		return err
	}
	return nil
}

func DoesExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func SaveStructToJson(s interface{}, jsonFilesystemPath string) error {
	resultJson, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jsonFilesystemPath, resultJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetRepoNameWithUnderscores(repoName string) string {
	return strings.ReplaceAll(repoName, "/", "_")
}

func GetRepositoriesScanMetadata(c config.Config) ([]Repository, error) {
	var repositories []Repository

	err := filepath.Walk(c.Filesystem.ScanMetadataFolder, func(path string, file os.FileInfo, err error) error {
		var repository Repository

		fmt.Println("Path:", path)
		//What you actually want is "starts with". Improve this!
		if strings.HasSuffix(file.Name(), ".json") {
			jsonFile, err := os.Open(path)
			if err != nil {
				log.Println(err)
			}
			defer jsonFile.Close()
			//Ignoring error? Questionable :o
			GitRepositoryJson, _ := ioutil.ReadAll(jsonFile)
			err = json.Unmarshal(GitRepositoryJson, &repository)
			if err != nil {
				return err
			}
			repositories = append(repositories, repository)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func GenerateAndUploadIndexHtml(c config.Config) error {

	repositories, err := GetRepositoriesScanMetadata(c)
	if err != nil {
		return err
	}

	indexTpl := template.Must(template.ParseFiles(c.Filesystem.ProjectFolder + "/assets/index.gohtml"))

	var buf bytes.Buffer

	err = indexTpl.Execute(&buf, repositories)
	if err != nil {
		return err
	}

	_, err = s3.New(c.S3session).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(c.S3.BucketName),
		Key:         aws.String("index.html"),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("text/html"),
		Body:        bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return err
	}
	return nil
}

//Divide smaller parts of this into subfunctions
func UploadResultsToS3(c config.Config, r Repository) error {
	err := UploadJSONToS3(c, r)
	if err != nil {
		return err
	}

	err = UploadHTMLToS3(c, r)
	if err != nil {
		return err
	}
	return nil
}

func UploadJSONToS3(c config.Config, r Repository) error {
	mdLinksJson, err := json.MarshalIndent(r.MarkdownLinks, "", "  ")
	if err != nil {
		return err
	}

	_, err = s3.New(c.S3session).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(c.S3.BucketName),
		Key:         aws.String(r.JSONReportPath),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("text/html"),
		Body:        bytes.NewReader(mdLinksJson),
	})

	if err != nil {
		return err
	}
	return nil
}

func UploadHTMLToS3(c config.Config, r Repository) error {
	//Dear lord I should store this string somewhere
	buf, err := TemplateHTMLReportToBuffer(c.Filesystem.ProjectFolder+string(os.PathSeparator)+"assets"+string(os.PathSeparator)+"results_table.gohtml", r)
	if err != nil {
		return err
	}

	_, err = s3.New(c.S3session).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(c.S3.BucketName),
		Key:         aws.String(r.HTMLReportPath),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("text/html"),
		Body:        bytes.NewReader(buf.Bytes()),
	})

	if err != nil {
		return err
	}

	return nil
}

func TemplateHTMLReportToBuffer(s string, r Repository) (bytes.Buffer, error) {
	tpl := template.Must(template.ParseFiles(s))

	var buf bytes.Buffer

	err := tpl.Execute(&buf, r.MarkdownLinks)
	if err != nil {
		return buf, err
	}
	return buf, err
}

func Count404MarkdownLinks(mdLinks []MarkdownLink) int {
	var c int
	for _, link := range mdLinks {
		if link.Status == "404" {
			c++
		}
	}
	return c
}

func GetRepoUrlsFromProject(project string) ([]string, error) {
	var repositoryUrls []string

	req, err := http.NewRequest("GET", "https://api.github.com/orgs/"+project+"/repos?per_page=2000", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.inertia-preview+json")
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err)
	}

	ghResponse := GithubProjectApiResponse{}

	err = json.Unmarshal(body, &ghResponse)
	if err != nil {
		return nil, err
	}

	for _, repo := range ghResponse {
		repositoryUrls = append(repositoryUrls, repo.HTMLURL)
	}
	return repositoryUrls, nil
}

//func SortRepositoriesByUnscannedFirst(repoUrls []string, staticFolder string) []string {
//	var notScanned []string
//	var alreadyScanned []string
//
//	for _, repo := range repoUrls {
//
//		HtmlFileName, _ := GetRepoFilenameWithExtension(repo, "html")
//		_, err := os.Stat(staticFolder + HtmlFileName)
//
//		if os.IsNotExist(err) {
//			notScanned = append(notScanned, repo)
//		} else {
//			alreadyScanned = append(alreadyScanned, repo)
//		}
//
//	}
//	return append(notScanned, alreadyScanned...)
//}
//

//If ran on a Unix system:
//If this function receives /tmp it will return /tmp/
//If it receives /tmp/ it will return /tmp/
//func CheckAndAddPathSeparatorSuffix(fsPath string) string {
//	if !strings.HasSuffix(fsPath, string(os.PathSeparator)) {
//		fsPath = fsPath + string(os.PathSeparator)
//		return fsPath
//	} else {
//		return fsPath
//	}
//}
//

//
//func SortLinksByStatus(mdLinks []MarkdownLink, status string) []MarkdownLink {
//	var tmpLinks []MarkdownLink
//
//	//returns a slice bounds out of error if 404 link is on the last element. Redo do this!
//	for i, v := range mdLinks {
//		if strings.HasPrefix(v.Status, "4") {
//			tmpLinks = append(tmpLinks, v)
//			mdLinks = append(mdLinks[:i], mdLinks[i+1:]...)
//		}
//	}
//	mdLinks = append(tmpLinks, mdLinks...)
//	return mdLinks
//}
//
