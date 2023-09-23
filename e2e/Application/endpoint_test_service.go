package Application

import (
	"e2e/Domain"
	"e2e/Domain/VO"
	"e2e/Interface"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func RunTests() {
	applicationConfig := getApplicationConfig()             // Get config for tests
	endpoints := getAllPossibleEndpoints(applicationConfig) // Get endpoints from file & openApi
	groupedEndpoints := groupAndSortEndpoints(endpoints)    // Group request to make logical chain of requests
	requestEndpointsData(groupedEndpoints)                  // Execute requests to get test results
	validateResult(groupedEndpoints)                        // Verify results
	showResults(groupedEndpoints)                           // Test sum up - show to user results
}

func getApplicationConfig() Config {
	applicationConfig := GetConfig() // get config from .evn file
	// @TODO: other config from other sources e.g. command line

	return applicationConfig
}

func getAllPossibleEndpoints(applicationConfig Config) []Domain.Endpoint {
	var allEndpoints []Domain.Endpoint

	Interface.PrintSimpleText("Looking for endpoints...")

	// endpoints from files
	manualEndpoints := getManualEndpoints(applicationConfig)
	allEndpoints = append(allEndpoints, manualEndpoints...)

	// openapi endpoints
	openApiEndpoints := getOpenApiEndpoints(applicationConfig)
	allEndpoints = append(allEndpoints, openApiEndpoints...)

	if len(allEndpoints) > 0 {
		Interface.PrintSuccessText("Found " + strconv.Itoa(len(allEndpoints)))
		Interface.EndMsgLine()
	} else {
		Interface.WarningMsg("No endpoints to execute! Check documentation and files", []any{})
		os.Exit(0)
	}

	return allEndpoints
}

func requestData(eg *Domain.EndpointGroup, id string) {
	endpoints := eg.GetVerbs()
	id = getEndpointsResponse(&endpoints, id)

	if eg.HasSubgroup() == true {
		requestData(eg.GetSubgroup(), id)
	}
}

func getEndpointsResponse(endpoints *map[string]Domain.Endpoint, id string) string {

	// @TODO: get all verbs and order them from existing pool, replace the order below
	//var verbList []string
	//for existingVerb, _ := range *endpoints {
	//	verbList = append(verbList, strings.ToUpper(existingVerb))
	//}

	// sort request by endpoints:
	order := make(map[int]string)
	order[0] = "POST"
	order[1] = "GET"
	order[2] = "PUT"
	order[3] = "PATCH"
	order[4] = "DELETE"

	for i := 0; i <= len(order); i++ {
		verb := order[i]
		e := *endpoints

		_, exists := e[verb]
		if exists {
			updatedEndpoint, newId := makeEndpointRequest(e[verb], id)
			e[verb] = updatedEndpoint

			if newId != "" {
				id = newId
			}
		}
	}

	return id
}

func makeEndpointRequest(e Domain.Endpoint, id string) (Domain.Endpoint, string) {
	//if contains {} replace with id!
	url := e.GetUrl()
	if strings.Contains(url, "{") && strings.Contains(url, "}") {
		regexPattern := `\{.*?\}`
		regex, err := regexp.Compile(regexPattern)

		if err != nil {
			panic(err)
		}
		url = regex.ReplaceAllString(url, id)
		e.UpdateUrl(url)
	}

	//build http request
	httpRequest, httpError := Interface.CreateRequest(e.GetUrl(), e.GetVerb(), e.GetBody())
	request := Domain.BuildEndpointRequest(httpRequest, httpError)

	if e.GetHeader() != nil {
		for _, header := range e.GetHeader() {
			request.HttpRequest.Header.Add(header.Key, header.Value)
		}
	}

	if e.GetToken() != "" { // There is a token used in requests // @TODO: make this read from auth not JWT token!
		noTokenHttpRequest, noTokenHttpError := Interface.CreateRequest(e.GetUrl(), e.GetVerb(), e.GetBody())
		noTokenRequest := Domain.BuildEndpointRequest(noTokenHttpRequest, noTokenHttpError)

		// Build request with random token (if there is a token)
		randomTokenHttpRequest, randomTokenHttpError := Interface.CreateRequest(e.GetUrl(), e.GetVerb(), e.GetBody())
		randomTokenRequest := Domain.BuildEndpointRequest(randomTokenHttpRequest, randomTokenHttpError)

		tokenString := fmt.Sprintf("Bearer %s", e.GetToken())
		request.HttpRequest.Header.Add("Authorization", tokenString)

		randomTokenString := fmt.Sprintf("Bearer %s", GenerateRandomJWT())
		randomTokenRequest.HttpRequest.Header.Add("Authorization", randomTokenString)

		for _, possibleResult := range e.GetPossibleResults() {
			statusCodeString := strconv.Itoa(possibleResult.StatusCode)
			if len(statusCodeString) > 0 && statusCodeString[0] == '2' {
				request.AddExpectedResult(possibleResult)
				noTokenRequest.AddExpectedResult(possibleResult)
			}

			if len(statusCodeString) > 0 && statusCodeString[0] == '4' {
				request.AddExpectedResult(possibleResult)
				noTokenRequest.AddExpectedResult(possibleResult)
				randomTokenRequest.AddExpectedResult(possibleResult)
			}
		}

		e.AddRequest(request)
		e.AddRequest(noTokenRequest)
		e.AddRequest(randomTokenRequest)
	} else {
		for _, possibleResult := range e.GetPossibleResults() {
			statusCodeString := strconv.Itoa(possibleResult.StatusCode)
			if len(statusCodeString) > 0 && statusCodeString[0] == '2' {
				request.AddExpectedResult(possibleResult)
			}
		}

		e.AddRequest(request)
	}

	//make request
	Interface.StartMsgLine()
	Interface.PrintSimpleText("[" + e.GetVerb() + "]" + e.GetUrl() + "...")
	for i, er := range e.Requests {

		endpointTestStatusCode, endpointBody, endpointHeader := Interface.FullRequest(er.HttpRequest)

		idTmp := readIdFromJson(endpointBody)
		if idTmp != "" && e.GetVerb() == "POST" {
			id = idTmp
		}

		e.Requests[i].Response = Domain.CreateEndpointResult(endpointTestStatusCode, endpointBody, endpointHeader)

		Interface.PrintInfoText(strconv.Itoa(endpointTestStatusCode))
	}

	return e, id
}

func requestEndpointsData(endpointGroups VO.EndpointGroups) {
	Interface.StartMsgLine()
	Interface.PrintSimpleText("Executing requests...")

	// @TODO: For each group do concurency requests, POST first to get ID, if no POST, try to GET id from GET?
	for _, eg := range endpointGroups.GetGroups() {
		var id string

		requestData(eg, id)
	}

	Interface.EndMsgLine()
}

func readIdFromJson(j []byte) string {
	c := make(map[string]json.RawMessage)
	e := json.Unmarshal(j, &c)

	if e != nil {
		return ""
	}

	for s, v := range c {
		if s == "id" {
			return strings.ReplaceAll(string(v), "\"", "")
		}
	}

	return ""
}

func groupAndSortEndpoints(endpoints []Domain.Endpoint) VO.EndpointGroups {
	voGroups := VO.EndpointGroups{}

	for _, endpoint := range endpoints {
		voGroups.AddEndpoint(endpoint)
	}
	return voGroups
}

func validateResult(eg VO.EndpointGroups) {
	Interface.StartMsgLine()
	Interface.PrintSimpleText("Validating results...")

	// @TODO: go routine
	for _, g := range eg.GetGroups() {
		for _, e := range g.GetEndpoints() {
			e.Validate()
		}

		if g.HasSubgroup() {
			for _, sg := range g.GetSubgroup().GetVerbs() {
				sg.Validate()
			}
		}
	}

	Interface.PrintSuccessText("OK")
	Interface.EndMsgLine()
}

func showResults(eg VO.EndpointGroups) {
	var warningTestNames []string // No 200 responses
	var passedTestCount int
	var failedTestCount int

	Interface.StartMsgLine()
	Interface.PrintSimpleText("Sum up results...")
	Interface.EndMsgLine()

	for _, g := range eg.GetGroups() {
		for _, endpoint := range g.GetEndpoints() {
			passed := endpoint.TestResult()

			if !passed {
				Interface.PrintSimpleText(endpoint.GetUrl())
				Interface.PrintFailText("FAILED")
				Interface.EndMsgLine()
				Interface.PrintSimpleText("----------- " + endpoint.GetVerb() + " -----------")
				Interface.EndMsgLine()
				for _, r := range endpoint.Requests {
					if r.TestInformation != nil {
						for _, t := range r.TestInformation {
							Interface.PrintSimpleText(t)
							Interface.EndMsgLine()
						}
					}
				}
				Interface.PrintSimpleText("----------- ---- -----------")
				Interface.EndMsgLine()

				//failedTestNames = append(failedTestNames, endpoint.GetUrl())
				failedTestCount++
			} else {
				passedTestCount++
			}

			if !endpoint.AtLeastOneRequestOk() {
				warningTestNames = append(warningTestNames, "["+endpoint.GetVerb()+"]"+endpoint.GetUrl())
			}
		}

		if g.HasSubgroup() {
			for _, subgroupEndpoint := range g.GetSubgroup().GetVerbs() {
				passed := subgroupEndpoint.TestResult()

				if !passed {
					Interface.PrintSimpleText(subgroupEndpoint.GetUrl())
					Interface.PrintFailText("FAILED")
					Interface.EndMsgLine()
					Interface.PrintSimpleText("----------- " + subgroupEndpoint.GetVerb() + " -----------")
					Interface.EndMsgLine()
					for _, r := range subgroupEndpoint.Requests {
						if r.TestInformation != nil {
							for _, t := range r.TestInformation {
								Interface.PrintSimpleText(t)
								Interface.EndMsgLine()
							}
						}
					}
					Interface.PrintSimpleText("----------- ---- -----------")
					Interface.EndMsgLine()

					//failedTestNames = append(failedTestNames, subgroupEndpoint.GetUrl())
					failedTestCount++
				} else {
					passedTestCount++
				}

				if !subgroupEndpoint.AtLeastOneRequestOk() {
					warningTestNames = append(warningTestNames, "["+subgroupEndpoint.GetVerb()+"]"+subgroupEndpoint.GetUrl())
				}
			}
		}
	}

	// Warnings
	if len(warningTestNames) > 0 {
		Interface.StartMsgLine()
		Interface.PrintSimpleText("Warning no 200 status present in these requests...")
		Interface.EndMsgLine()

		for _, w := range warningTestNames {
			Interface.PrintFailText(w)
			Interface.EndMsgLine()
		}
	}

	Interface.StartMsgLine()
	Interface.PrintSimpleText("Passed...")
	Interface.PrintSuccessText(strconv.Itoa(passedTestCount))
	Interface.EndMsgLine()
	Interface.PrintSimpleText("Failed...")
	Interface.PrintFailText(strconv.Itoa(failedTestCount))
	Interface.EndMsgLine()

	if failedTestCount > 0 {
		os.Exit(1)
	}

	os.Exit(0)
}

func getManualEndpoints(config Config) []Domain.Endpoint {
	var endpoints []Domain.Endpoint

	for _, file := range config.manualEndpoints {
		configFilesEndpointList, err := Interface.ReadTestFile(file)
		if err != nil {
			Interface.ErrorMsg("There was an error while testing status codes: %s", []any{err.Error()})
			os.Exit(2)
		}

		for _, configFileEndpoint := range configFilesEndpointList {
			endpoint := Domain.CreateEndpointFromFile(configFileEndpoint, config.JWT.Token)

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

func getOpenApiEndpoints(config Config) []Domain.Endpoint {
	var endpoints []Domain.Endpoint

	if config.DocUrl != "" {
		urls, readingError := Interface.ReadOpenapiDocs(config.DocUrl)

		if readingError != nil {
			Interface.ErrorMsg("There was an error while reading OpenApi documentation: %s", []any{readingError.Error()})
			os.Exit(2)
		}

		for urlPath, openApiData := range urls {
			fullPath := fmt.Sprintf("%s%s", config.TestingUrl, urlPath)

			//Get endpoint
			if openApiData.Get != nil {
				endpoint := Domain.CreateEndpointFromUrl(fullPath, openApiData.Get, "Get", config.JWT.Token)
				endpoints = append(endpoints, endpoint)
			}

			//Post endpoint
			if openApiData.Post != nil {
				endpoint := Domain.CreateEndpointFromUrl(fullPath, openApiData.Post, "Post", config.JWT.Token)
				endpoints = append(endpoints, endpoint)
			}

			//Put endpoint
			if openApiData.Put != nil {
				endpoint := Domain.CreateEndpointFromUrl(fullPath, openApiData.Put, "Put", config.JWT.Token)
				endpoints = append(endpoints, endpoint)
			}

			//Patch endpoint
			if openApiData.Patch != nil {
				endpoint := Domain.CreateEndpointFromUrl(fullPath, openApiData.Patch, "Patch", config.JWT.Token)
				endpoints = append(endpoints, endpoint)
			}

			//Delete endpoint
			if openApiData.Delete != nil {
				endpoint := Domain.CreateEndpointFromUrl(fullPath, openApiData.Delete, "Delete", config.JWT.Token)
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	return endpoints
}
