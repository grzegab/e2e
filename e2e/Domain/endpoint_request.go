package Domain

import "net/http"

type EndpointRequest struct {
	HttpRequest     *http.Request
	HttpError       error
	ExpectedResults []EndpointResult
	RequestHeader   []Params
	Response        EndpointResult
	Passed          bool
	TestInformation []string
}

func (er *EndpointRequest) AddResult(statusCode int, body []byte, header http.Header) {
	var result EndpointResult

	result.StatusCode = statusCode
	result.Body = string(body)
	result.Header = header

	er.Response = result
}

func (er *EndpointRequest) GetHttpRequest() *http.Request {
	return er.HttpRequest
}

func (er *EndpointRequest) GetHttpError() error {
	return er.HttpError
}

func (er *EndpointRequest) AddHeader(header Params) {
	er.RequestHeader = append(er.RequestHeader, header)
}

func (er *EndpointRequest) AddExpectedResult(expected EndpointResult) {
	er.ExpectedResults = append(er.ExpectedResults, expected)
}

func (er *EndpointRequest) DeleteHeader(key string) {
	for i, param := range er.RequestHeader {
		if param.Value == key {
			er.RequestHeader = append(er.RequestHeader[:i], er.RequestHeader[i+1:]...)
		}
	}
}

func CreateEndpointResult(statusCode int, body []byte, header http.Header) EndpointResult {
	var result EndpointResult

	result.StatusCode = statusCode
	result.Body = string(body)
	result.Header = header

	return result
}

func BuildEndpointRequest(httpRequest *http.Request, httpError error) EndpointRequest {
	var request EndpointRequest

	request.HttpRequest = httpRequest
	request.HttpError = httpError

	return request
}
