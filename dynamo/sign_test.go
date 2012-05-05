package dynamo

import (
	"github.com/justinbarry/goamz-aws"
	"strings"
	"testing"
)

func TestSign(t *testing.T) {
	auth := aws.Auth{"key", "secret"}
	headers := map[string]string{
		"x-amz-date":           "Mon, 09 Apr 2012 05:43:51 -0700",
		"x-amz-target":         "DynamoDB_20111205.Method",
		"x-amz-security-token": "security_token",
	}

	expected := "AWS3 AWSAccessKeyId=key,Algorithm=HmacSHA256," +
		"SignedHeaders=host;x-amz-date;x-amz-security-token;" +
		"x-amz-target,Signature=ZPibXX2FSARKOwRP6lyE8C1bPym+4TwlcpSwC4a2Mhw="
	sign(auth, "POST", "/", headers, "dynamodb.us-east-1.amazonaws.com", "{}")

	if headers["x-amzn-authorization"] != expected {
		t.Error("Signature Does not match.\nExpected: " + expected + "\nprovided: " + headers["x-amzn-authorization"])
	}
}

func TestHeadersToSign(t *testing.T) {
	var headers = map[string]string{
		"x-amz-date":           "",
		"x-amz-target":         "",
		"x-amz-security-token": "",
		"should-not-be-added":  "",
	}

	r := headersToSign(headers)
	if strings.Join(r, ";") != "x-amz-date;x-amz-security-token;x-amz-target" {
		t.Error(r)
	}
}

func TestCanonicalHeaders(t *testing.T) {
	// The is lowercase here because this is how Amazon seems to  prefer it.
	var headers = map[string]string{
		"host":                 "dynamodb.us-east-1.amazonaws.com",
		"x-amz-date":           "Mon, 09 Apr 2012 05:50:51 -0700",
		"x-amz-target":         "DynamoDB_20111205.Method",
		"x-amz-security-token": "security_token",
		"should-not-be-added":  "",
	}

	expected := "host:dynamodb.us-east-1.amazonaws.com" +
		"\nx-amz-date:Mon, 09 Apr 2012 05:50:51 -0700" +
		"\nx-amz-security-token:security_token" +
		"\nx-amz-target:DynamoDB_20111205.Method" +
		"\n"

	r := canonicalHeaders(headers)
	if r != expected {
		t.Error("\nProvided:\n" + r + "\nExpected:\n" + expected)
	}
}

func TestStringToSign(t *testing.T) {
	h := "Host:dynamodb.us-east-1.amazonaws.com" +
		"\nx-amz-date:Mon, 09 Apr 2012 05:50:51 -0700" +
		"\nx-amz-security-token:security_token" +
		"\nx-amz-target:DynamoDB_20111205.Method\n"
	s := stringToSign("POST", "/", h, "{}")
	if s != "POST\n/\n\n"+h+"\n{}" {
		t.Error("\nString to sign does not return proper format\n" + s)
	}
}

func TestSignString(t *testing.T) {
	s := "POST" +
		"\n/" +
		"\n" +
		"\nHost:dynamodb.us-east-1.amazonaws.com" +
		"\nx-amz-date:Mon, 09 Apr 2012 05:43:51 -0700" +
		"\nx-amz-security-token:security_token" +
		"\nx-amz-target:DynamoDB_20111205.Method" +
		"\n" +
		"\n{}"

	auth := aws.Auth{"key", "secret"}
	expected := "hx9JetaCUTZKmW4Sq86C3upmc9TFMIEn4tI+YhODcLs="

	r := signString(s, auth)
	if string(r) != expected {
		t.Error("\nSignature did not match\nProvided:" + string(r) + "\nExpected: " + expected)
	}
}
