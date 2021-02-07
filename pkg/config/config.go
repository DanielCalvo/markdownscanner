package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

//Should I put the file config inside another struct as it was done in Prometheus?
//Don't forget the rename TODO as described on config.yaml
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
	GithubProjects []string `yaml:"GithubProjects"`
	Repositories   []string `yaml:"Repositories"`
	S3session      *session.Session
}

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}
	// If the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have to set the DefaultConfig at the entry
	// point as well.

	//Note from Dani: Is setting a default config useful here?
	//*cfg = DefaultConfig

	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	if filename == "" {
		return nil, fmt.Errorf("config.file flag required but not passed")
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
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
	err := ValidateDir(c.Filesystem.TmpFolder)
	if err != nil {
		return err
	}

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
