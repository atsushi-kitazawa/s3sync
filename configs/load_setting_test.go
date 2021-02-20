package configs

import (
    "testing"
)

func TestLoad(t *testing.T) {
    var s *Setting
    s = Load("setting.yaml")

    t.Log(s.BucketName)
    t.Log(s.LocalDir)
    t.Log(s.S3Dir)
    t.Log(s.Credential.Region)
    t.Log(s.Credential.Apikey)
    //t.Log(s.Credential.SecretKey)
}
