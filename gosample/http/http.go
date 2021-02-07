// http is sample package
package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type httpClient interface {
	Request()
}

type httpClientImpl struct {
}

func (h *httpClientImpl) Request(
	endpoint, method, apipath string,
	header, query map[string]string,
	body []byte) (
	string,
	error,
	int,
) {
	// configure base URL
	baseURL, _ := url.Parse(endpoint)

	// join api path
	if apipath != "" {
		baseURL.Path = path.Join(baseURL.Path, apipath)
	}

	// set query
	q := baseURL.Query()
	for key, value := range query {
		q.Set(key, value)
	}
	baseURL.RawQuery = q.Encode()

	// create request URL
	reqURL := baseURL.String()

	// configure request body
	reqBody := bytes.NewReader(body)

	// create new http request
	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return "", err, 0
	}

	// configure header
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// create and execute http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err, 0
	}
	defer resp.Body.Close()

	// configure return value
	statusCode := resp.StatusCode
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err, statusCode
	}
	return string(respBodyBytes), nil, statusCode
}

func NewHttpClientImpl() *httpClientImpl {
	return &httpClientImpl{}
}
