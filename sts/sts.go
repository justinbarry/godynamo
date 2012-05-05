package sts

import (
	"encoding/xml"
	"github.com/justinbarry/goamz-aws"
	"net/http"
	"net/url"
	"time"
)

type STS struct {
	aws.Auth
	aws.Region
}

type CredentialResp struct {
	SessionToken    string `xml:"GetSessionTokenResult>Credentials>SessionToken"`
	SecretAccessKey string `xml:"GetSessionTokenResult>Credentials>SecretAccessKey"`
	AccessKeyId     string `xml:"GetSessionTokenResult>Credentials>AccessKeyId"`
}

func New(auth aws.Auth, region aws.Region) *STS {
	return &STS{auth, region}
}

func (sts *STS) GetSessionToken() (resp *CredentialResp, err error) {
	resp = &CredentialResp{}
	err = sts.GetSessionCredentials(resp)
	return
}

func (sts *STS) GetSessionCredentials(resp interface{}) error {
	host := sts.Region.STSEndpoint
	auth := sts.Auth
	path := "/"
	method := "GET"
	tm := time.Now()

	params := map[string]string{
		"Action":    "GetSessionToken",
		"Version":   "2011-06-15",
		"Timestamp": tm.Format(time.RFC3339),
	}

	endPoint, _ := url.Parse(host)

	sign(auth, method, path, params, endPoint.Host)

	endPoint.RawQuery = multimap(params).Encode()
	r, err := http.Get(endPoint.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	err = xml.NewDecoder(r.Body).Decode(resp)
	return err
}

func multimap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}
