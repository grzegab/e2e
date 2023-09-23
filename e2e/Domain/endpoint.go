package Domain

import (
	"e2e/Interface"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"strconv"
	"strings"
)

type Endpoint struct {
	name            string
	url             string
	verb            string
	body            string
	header          []Params
	token           string
	authRequired    bool
	tested          bool
	passed          bool
	possibleResults []EndpointResult
	Requests        []EndpointRequest
}

type requestBody struct {
	values []Params
}

func CreateEndpointFromFile(file Interface.ConfigFile, token string) Endpoint {
	endpoint := createDefaultEndpoint()

	endpoint.name = file.Name
	endpoint.url = file.Path
	endpoint.verb = file.Verb
	endpoint.body = file.RequestBody
	endpoint.token = token
	endpoint.authRequired = true // @TODO: set auth same way as in while creating from url

	if file.RequestHeader != nil {
		for _, header := range file.RequestHeader {
			endpoint.header = append(endpoint.header, Params{
				Key:   header.Parameter,
				Value: header.Value,
			})
		}
	}

	for _, availableResult := range file.AvailableResults {
		var possibleResult EndpointResult

		possibleResult.StatusCode = availableResult.Code
		endpoint.possibleResults = append(endpoint.possibleResults, possibleResult)
	}

	return endpoint
}

func CreateEndpointFromUrl(path string, openApiOperation *openapi3.Operation, verb string, token string) Endpoint {
	endpoint := createDefaultEndpoint()

	endpoint.url = path
	endpoint.verb = verb
	endpoint.token = token
	endpoint.name = openApiOperation.OperationID
	endpoint.authRequired = true // @TODO: find a way to read if there is auth to the endpoint, error in openApi3 parser

	//header
	if openApiOperation.Parameters != nil {
		for _, parameter := range openApiOperation.Parameters {
			if parameter.Value.In == "header" {
				if str, ok := parameter.Value.Schema.Value.Default.(string); ok {
					endpoint.header = append(endpoint.header, Params{
						Key:   parameter.Value.Name,
						Value: str,
					})
				}
			}
		}
	}

	//body
	if openApiOperation.RequestBody != nil {
		var rawBody requestBody

		for _, contentType := range openApiOperation.RequestBody.Value.Content {
			// @TODO: get all parameters to make more tests
			if len(contentType.Schema.Value.Required) > 0 {
				for _, requiredField := range contentType.Schema.Value.Required {
					for name, property := range contentType.Schema.Value.Properties {
						if name == requiredField {
							param := Params{Key: name, Value: property.Value.Example.(string)}
							rawBody.values = append(rawBody.values, param)
						}
					}
				}
			}
		}

		var jsonStrings []string
		for _, v := range rawBody.values {
			json := "\"" + v.Key + "\":\"" + v.Value + "\""
			jsonStrings = append(jsonStrings, json)
		}

		joinedString := strings.Join(jsonStrings, ", ")
		jsonContentString := "{" + joinedString + "}"

		endpoint.body = jsonContentString
	}

	//responses
	for code := range openApiOperation.Responses {
		var possibleResult EndpointResult

		statusCode, err := strconv.Atoi(code)
		if err != nil {
			continue
		}

		possibleResult.StatusCode = statusCode
		endpoint.possibleResults = append(endpoint.possibleResults, possibleResult)
	}

	return endpoint
}

func (e *Endpoint) Validate() {
	// Iterate through all request made to compare with possible results
	for i, r := range e.Requests {
		e.tested = true
		testPassed := false
		var possibleStatusCodes []int

		gotRequiredStatusCode := false
		for _, expectedResult := range r.ExpectedResults {
			possibleStatusCodes = append(possibleStatusCodes, expectedResult.StatusCode)
			if expectedResult.StatusCode == r.Response.StatusCode {
				gotRequiredStatusCode = true
				break
			}
		}

		if r.HttpError != nil {
			e.Requests[i].TestInformation = append(r.TestInformation, "There was an error while request: "+r.HttpError.Error())
		}

		if !gotRequiredStatusCode {
			possibleStatuses := fmt.Sprintf("possible statuses: %v", possibleStatusCodes)
			e.Requests[i].TestInformation = append(r.TestInformation, "Wrong status code! "+possibleStatuses+", got: "+strconv.Itoa(r.Response.StatusCode))
		}

		if gotRequiredStatusCode && r.HttpError == nil {
			testPassed = true
		}

		e.Requests[i].Passed = testPassed
	}

	for i, r := range e.Requests {
		e.tested = true
		e.Requests[i].Passed = false

		for _, expectedResult := range r.ExpectedResults {
			if expectedResult.StatusCode == r.Response.StatusCode {
				e.Requests[i].Passed = true
				break
			}
		}
	}
}

func (e *Endpoint) TestResult() bool {
	success := true

	for _, request := range e.Requests {
		if !request.Passed {
			success = false
		}
	}

	return success
}

func (e *Endpoint) GetName() string {
	return e.name
}

func (e *Endpoint) GetUrl() string {
	return e.url
}

func (e *Endpoint) UpdateUrl(newUrl string) {
	e.url = newUrl
}

func (e *Endpoint) GetVerb() string {
	return strings.ToUpper(e.verb)
}

func (e *Endpoint) GetBody() string {
	return e.body
}

func (e *Endpoint) GetHeader() []Params {
	return e.header
}

func (e *Endpoint) GetToken() string {
	return e.token
}

func (e *Endpoint) GetPossibleResults() []EndpointResult {
	return e.possibleResults
}

func (e *Endpoint) GetAllRequests() []EndpointRequest {
	return e.Requests
}

func (e *Endpoint) GetAuthRequired() bool {
	return e.authRequired
}

func (e *Endpoint) SetAllRequests(requests []EndpointRequest) {
	e.Requests = requests
}

func (e *Endpoint) AddRequest(r EndpointRequest) {
	e.Requests = append(e.Requests, r)
}

func (e *Endpoint) AtLeastOneRequestOk() bool {
	have200Status := false

	for _, r := range e.Requests {
		if r.Response.StatusCode == 200 || r.Response.StatusCode == 201 || r.Response.StatusCode == 204 {
			have200Status = true
			break
		}
	}

	return have200Status
}

func createDefaultEndpoint() Endpoint {
	endpoint := Endpoint{
		tested: false,
		passed: false,
	}

	return endpoint
}
