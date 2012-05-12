package sts

import (
	"github.com/justinbarry/goamz-aws"
	"github.com/stathat/jconfig"
	"testing"
)

func TestGetSessionToken(t *testing.T) {
	conf := jconfig.LoadConfig("../config.json")
	auth := aws.Auth{conf.GetString("AccessKeyId"), conf.GetString("SecretAccessKey")}
	region := aws.USEast
	sts := New(auth, region)
	_, err := sts.GetSessionToken()
	if err != nil {
		t.Error("There was an error getting the session token.")
	}
}
