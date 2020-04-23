package markdownscanner

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RepositoriesYaml struct {
	Repositories []string `yaml:"repositories"`
	Projects     []string `yaml:"projects"`
}

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

func ValidateGitRepository(repository string) error {
	_, err := url.ParseRequestURI(repository)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(repository, "http") {
		return errors.New("Repository URL does not start with http")
	}
	return nil
}

func GetRepositoryUrlsFromYaml(yamlPath string) ([]string, error) {
	var repositoriesYaml RepositoriesYaml
	var repositoryUrls []string

	repositoriesYamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(repositoriesYamlFile, &repositoriesYaml)
	if err != nil {
		return nil, err
	}

	fmt.Println(repositoriesYaml.Projects)
	for _, githubProject := range repositoriesYaml.Projects {
		repos, err := GetReposFromProject(githubProject)
		if err == nil {
			repositoryUrls = append(repositoryUrls, repos...)
		} else {
			log.Println("Error getting repositories for:"+githubProject+":", err)
		}
	}

	for _, githubRepo := range repositoriesYaml.Repositories {
		err = ValidateGitRepository(githubRepo)
		if err == nil {
			repositoryUrls = append(repositoryUrls, githubRepo)
		} else {
			log.Println("Error validating repository URL:"+githubRepo+":", err)
		}
	}
	return repositoryUrls, nil
}

func GetReposFromProject(project string) ([]string, error) {
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

//go-git can't clone large repositories without using very large amounts of memory: https://github.com/src-d/go-git/issues/761
func CloneGitRepository(repositoryUrl, repositoriesFolder string) (string, error) {
	url, err := url.ParseRequestURI(repositoryUrl)
	if err != nil {
		return "", err
	}

	repositoryFolder := repositoriesFolder + url.Path

	if !DoesExist(repositoryFolder) {
		cmd := exec.Command("git", "clone", repositoryUrl, repositoryFolder)
		err = cmd.Run()
		if err != nil {
			return "", err
		}
	}

	return repositoryFolder, nil
}

//This function is used to check if a file is a broken symlink.
func DoesExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func SortRepositoriesByUnscannedFirst(repoUrls []string, staticFolder string) []string {
	var notScanned []string
	var alreadyScanned []string

	for _, repo := range repoUrls {

		HtmlFileName, _ := GetRepoFilenameWithExtension(repo, "html")
		_, err := os.Stat(staticFolder + HtmlFileName)

		if os.IsNotExist(err) {
			notScanned = append(notScanned, repo)
		} else {
			alreadyScanned = append(alreadyScanned, repo)
		}

	}
	return append(notScanned, alreadyScanned...)
}

func GetRepoFilenameWithExtension(repoUrl, extension string) (string, error) {
	url, err := url.ParseRequestURI(repoUrl)
	if err != nil {
		return "", err
	}
	urlPathUnderscores := strings.ReplaceAll(url.Path, "/", "_")
	return url.Host + urlPathUnderscores + "." + extension, nil
}
