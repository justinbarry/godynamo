package dynamo

import (
	"../sts"
	//"fmt"
	"github.com/justinbarry/goamz-aws"
	"github.com/stathat/jconfig"
	"testing"
)

var dynamo *DYNAMO

func setupNewDynamo() {
	if dynamo == nil {
		conf := jconfig.LoadConfig("../config.json")
		auth := aws.Auth{conf.GetString("AccessKeyId"), conf.GetString("SecretAccessKey")}
		region := aws.USEast
		sts := sts.New(auth, region)
		resp, _ := sts.GetSessionToken()
		auth = aws.Auth{resp.AccessKeyId, resp.SecretAccessKey}
		dynamo = New(auth, region, resp.SessionToken)
	}
}

func TestListTables(t *testing.T) {
	setupNewDynamo()
	r := dynamo.ListTables()
	if _, ok := r["TableNames"]; !ok {
		t.Error("Invalid Response")
	}
}

func TestDescribeTable(t *testing.T) {
	setupNewDynamo()
	r := dynamo.DescribeTable("school")
	if _, ok := r["Table"]; !ok {
		t.Error("Invalid Response")
	}
}
