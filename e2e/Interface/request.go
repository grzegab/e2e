package Interface

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
)

func CreateRequest(url string, method string, body string) (*http.Request, error) {
	loginBodyBytes := []byte(body)
	bodyReader := bytes.NewReader(loginBodyBytes)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return req, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return req, nil
}

func FullRequest(r *http.Request) (int, []byte, http.Header) {
	httpClient := createHttpClient()

	resp, err := httpClient.Do(r)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return resp.StatusCode, body, resp.Header
}

func PingRequest(r *http.Request) bool {
	httpClient := createHttpClient()
	_, err := httpClient.Do(r)
	if err != nil {
		return false
	}

	return true
}

func createHttpClient() *http.Client {
	httpClient := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 100

	return httpClient
}
