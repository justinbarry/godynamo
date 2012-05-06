package dynamo

/*
BatchGetItem
CreateTable
DeleteItem
DeleteTable
GetItem
PutItem
Query
Scan
UpdateItem
UpdateTable
*/

import (
	"bytes"
	"encoding/json"
	"github.com/justinbarry/goamz-aws"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type DYNAMO struct {
	aws.Auth
	aws.Region
	SessionToken string
}

func New(auth aws.Auth, region aws.Region, sessionToken string) *DYNAMO {
	return &DYNAMO{auth, region, sessionToken}
}

func (dynamo *DYNAMO) ListTables() map[string]interface{} {
	return dynamo.Request("ListTables", "{}")
}

func (dynamo *DYNAMO) DescribeTable(table string) map[string]interface{} {
	return dynamo.Request("DescribeTable", "{\"TableName\":\""+table+"\"}")
}

func (dynamo *DYNAMO) PutItem() {

}

func (dynamo *DYNAMO) Request(op string, body string) (m map[string]interface{}) {
	host := dynamo.Region.DynamoDBEndPoint
	endPoint, _ := url.Parse(host)
	headers := map[string]string{
		"host":                 endPoint.Host,
		"x-amz-date":           time.Now().Format(time.RFC1123Z),
		"x-amz-target":         "DynamoDB_20111205." + op,
		"x-amz-security-token": dynamo.SessionToken,
		"content-type":         "application/x-amz-json-1.0",
	}

	sign(dynamo.Auth, "POST", "/", headers, endPoint.Host, body)

	thebody := ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
	req, _ := http.NewRequest("POST", host, thebody)
	req.Header = multimap(headers)

	r, _ := http.DefaultClient.Do(req)
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&m)
	return
}

// This is clever little function from the canonical EC2 packages
func multimap(p map[string]string) http.Header {
	q := make(http.Header, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}
