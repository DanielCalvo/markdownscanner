package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

//Should I put the file config inside another struct as it was done in Prometheus?
//Don't forget the rename TODO as described on config.yaml
//This is awful to read!
type Config struct {
	Filesystem struct {
		ProjectFolder      string `yaml:"projectFolder"` //I don't like this -- get it from the $CURRENT_DIR if you can
		TmpFolder          string `yaml:"tmpFolder"`
		ScanMetadataFolder string `yaml:"scanMetadataFolder"`
	} `yaml:"filesystem"`

	S3 struct {
		Region       string `yaml:"region"`
		BucketName   string `yaml:"bucketName"`
		BuckerFolder string `yaml:"buckerFolder"`
	} `yaml:"s3"`
	GithubProjects []string         `yaml:"GithubProjects"`
	Repositories   []string         `yaml:"Repositories"`
	S3session      *session.Session //This does not need to be here

}

//This was put together by copying and pasting from other places. Some of these are by reference and others are not. Adopt a standard!
//Also can you merge "LoadFile" and "Initialize?" Or at organize them better, they're doing very similar things
func New(configFile string) (Config, error) {
	cfg, err := LoadFile(configFile)

	if err != nil {
		return cfg, err
	}

	log.Println("Initializing config file")
	err = Initialize(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Load parses the YAML input s into a Config.
func Load(s string) (Config, error) {
	cfg := Config{}
	err := yaml.UnmarshalStrict([]byte(s), &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (Config, error) {

	var cfg = Config{}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	cfg, err = Load(string(content))
	if err != nil {
		return cfg, errors.Wrapf(err, "parsing YAML file %s", filename)
	}
	return cfg, nil
}

//Hmmm
func ValidateGitRepos(s string) error {
	return nil
}

//check that the s3 bucket exists and you can upload files to it
func ValidateS3Bucket(c *Config) error {
	//import something from the S3 package for this, do later
	fmt.Println(c.S3.BucketName)
	fmt.Println(c.S3.Region)
	return nil
}

//somehow check that you can write to the tmp folder?
func ValidateDir(s string) error {
	//Remove the contents of the temporary directory on startup (but not the directory itself)
	dir, err := ioutil.ReadDir(s)
	for _, d := range dir {
		_ = os.RemoveAll(path.Join([]string{"tmp", d.Name()}...))
	}

	f, err := os.Stat(s)

	if err != nil {
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(s, 0755)
			if errDir != nil {
				return errDir
			}
		} else {
			return fmt.Errorf("unexpected filesystem error initiating configuration: %s", err)
		}
	} else {
		if !f.IsDir() {
			return fmt.Errorf(s + " is a file, needs to be a directory")
		}
	}
	return nil
}

//This function may be doing to many things -- maybe do the S3 check elsewhere
func Initialize(c *Config) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	c.Filesystem.ProjectFolder = dir
	c.Filesystem.TmpFolder = dir + string(os.PathSeparator) + "tmp"
	c.Filesystem.ScanMetadataFolder = dir + string(os.PathSeparator) + "metadata"

	//tmp folder
	err = ValidateDir(c.Filesystem.ProjectFolder)
	if err != nil {
		return err
	}

	//metadata folder
	err = ValidateDir(c.Filesystem.ScanMetadataFolder)
	if err != nil {
		return err
	}

	c.S3session, err = session.NewSession(&aws.Config{Region: aws.String(c.S3.Region)})
	if err != nil {
		return err
	}

	//err = ValidateS3Bucket(c)
	//if err != nil {
	//	return err
	//}

	return nil
}
