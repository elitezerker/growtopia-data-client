package httpclient

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	StatusCodeOk = 200
)

type HttpClient struct {
	Client *http.Client
}

func (httpClient HttpClient) ExecuteRequestAndGetResponse(request *http.Request) (*http.Response, error) {
	response, err := httpClient.Client.Do(request)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return nil, err
	}
	return response, nil
}

func (httpClient HttpClient) HandleResponse(response *http.Response) (*http.Response, error) {
	if response.StatusCode != StatusCodeOk {
		log.Printf("Received status code %d", response.StatusCode)
		return nil, fmt.Errorf("status code is not OK: %d", response.StatusCode)
	}
	return response, nil
}

func (httpClient HttpClient) FetchResponseBody(request *http.Request) (*http.Response, error) {
	response, err := httpClient.ExecuteRequestAndGetResponse(request)
	if err != nil {
		return nil, err
	}
	response, err = httpClient.HandleResponse(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (httpClient HttpClient) ReadResponseBody(response *http.Response) string {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return ""
	}
	return string(responseBody)
}
