package dynamo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/justinbarry/goamz-aws"
	"regexp"
	"sort"
	"strings"
)

var b64 = base64.StdEncoding

func sign(auth aws.Auth, method string, path string, params map[string]string, host string, body string) {
	params["host"] = host
	defer delete(params, "host")

	authorization := map[string]string{}
	authorization["AWS3 AWSAccessKeyId"] = auth.AccessKey
	authorization["Algorithm"] = "HmacSHA256"
	authorization["SignedHeaders"] = strings.Join(headersToSign(params), ";")

	var autharray []string
	for k, v := range authorization {
		autharray = append(autharray, strings.Trim(k, " ")+"="+strings.Trim(v, " "))
	}

	sort.StringSlice(autharray).Sort()
	authjoined := strings.Join(autharray, ",")

	var sarray []string
	for k, v := range params {
		sarray = append(sarray, strings.Trim(k, " ")+":"+strings.Trim(v, " "))
	}

	headers := canonicalHeaders(params)
	payload := stringToSign(method, path, headers, body)
	signature := signString(payload, auth)

	params["x-amzn-authorization"] = authjoined + ",Signature=" + string(signature)
}

func canonicalHeaders(headers map[string]string) (s string) {
	headersToUse := headersToSign(headers)
	for _, v := range headersToUse {
		if _, ok := headers[v]; ok {
			s += v + ":" + headers[v] + "\n"
		}
	}
	return
}

func headersToSign(headers map[string]string) (sarray []string) {
	r, _ := regexp.Compile("^x-amz|content-encoding|host")
	for k, _ := range headers {
		if r.Match([]byte(k)) {
			sarray = append(sarray, strings.Trim(k, " "))
		}
	}
	sort.StringSlice(sarray).Sort()
	return
}

func signString(s string, auth aws.Auth) (signature []byte) {
	hash := sha256.New()
	hash.Write([]byte(s))
	request_hash := hash.Sum(nil)
	hash = hmac.New(sha256.New, []byte(auth.SecretKey))
	hash.Write([]byte(request_hash))
	signature = make([]byte, b64.EncodedLen(hash.Size()))
	b64.Encode(signature, hash.Sum(nil))
	return
}

/*
	Line 1: The HTTP method (POST), followed by a newline.
	Line 2: The request URI (/), followed by a newline.
	Line 3: An empty string. Typically, a query string goes here, but Amazon DynamoDB doesn't use a query string. Follow with a newline.
	Line 4-n: The string representing the canonicalized request headers you computed in step 1, followed by a newline.
	The request body. Do not follow the request body with a newline.
*/
func stringToSign(method string, path string, headers string, body string) (s string) {
	t := []string{method, path, "", headers, body}
	return strings.Join(t, "\n")
}
