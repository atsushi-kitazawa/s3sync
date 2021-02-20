package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Setting struct {
    BucketName string `yaml:"bucket_name"`
    LocalDir string `yaml:"local_dir"`
    S3Dir string `yaml:"s3_dir"`
    Credential Credential `yaml:"credential"`
}

type Credential struct {
    Region string `yaml:"region"`
    Apikey string `yaml:"apikey"`
    Secretkey string `yaml:"secretkey"`
}

func Load(settingPath string) *Setting {
    var s *Setting
    buf, err := ioutil.ReadFile(settingPath)
    if err != nil {
	panic(err)
    }

    err = yaml.UnmarshalStrict(buf, &s)
    if err != nil {
	panic(err)
    }

    return s
}
