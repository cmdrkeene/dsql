// Communicate with the service
package dsql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/aws4"
	"io"
	"net/http"
	"net/url"
	"time"
)

var Clients = map[string]Client{} // for testing

type Client interface {
	Post(Operation, Request) (io.ReadCloser, error)
}

type DynamoClient struct {
	keys    *aws4.Keys
	service *aws4.Service
	urlStr  string
}

func GetClient(name string) Client {
	c := Clients[name]
	if c == nil {
		c = NewClient(name)
		Clients[name] = c
	}
	return c
}

func NewClient(name string) Client {
	u, err := url.Parse(name)
	if err != nil {
		panic(err)
	}

	accessKey := u.User.Username()
	secretKey, _ := u.User.Password()
	region := u.Host

	return &DynamoClient{
		keys:    &aws4.Keys{AccessKey: accessKey, SecretKey: secretKey},
		service: &aws4.Service{Name: "dynamodb", Region: region},
		urlStr:  fmt.Sprintf("https://dynamodb.%s.amazonaws.com", region),
	}
}

func (d *DynamoClient) Post(op Operation, req Request) (io.ReadCloser, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", d.urlStr, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-amz-json-1.0")
	request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	request.Header.Set("X-Amz-Target", string(op))

	err = aws4.Sign(d.keys, request)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

type MockClient struct {
	OnPost func(Operation, Request) (io.ReadCloser, error)
}

func (m MockClient) Post(o Operation, r Request) (io.ReadCloser, error) {
	return m.OnPost(o, r)
}
