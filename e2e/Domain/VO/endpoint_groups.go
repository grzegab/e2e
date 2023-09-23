package VO

import (
	"e2e/Domain"
	"regexp"
	"strings"
)

// EXAMPLE:
// groups => http://localhost/v1/tempaltes => ( k => v )
//
//	    ID: nil
//		url: http://localhost/v1/tempaltes/
//		verbs =>
//			POST: ...
//			PATCH: ...
//			PUT: ...
//			DELETE: ...
//			GET: ...
//		subgroups => http://localhost/v1/tempaltes/{templateID} => ( k => v)
//			ID: 123
//			url: http://localhost/v1/tempaltes/{templateID}/
//			verbs =>
//				POST:...
//				DELETE:...
//			group =>
//				ID: ...
//				POST: ..

type EndpointGroups struct {
	groups map[string]*Domain.EndpointGroup
}

func (egs *EndpointGroups) Init() {
	egs.groups = make(map[string]*Domain.EndpointGroup)
}

func (egs *EndpointGroups) GetGroups() map[string]*Domain.EndpointGroup {
	return egs.groups
}

func (egs *EndpointGroups) AddEndpoint(e Domain.Endpoint) {
	var groupUrl string

	//check if url contains { and }
	if strings.Contains(e.GetUrl(), "{") && strings.Contains(e.GetUrl(), "}") {
		groupUrlOriginal := retrieveOriginalUrl(e.GetUrl())
		groupUrl = removeLastCharIfNecessary(groupUrlOriginal)

		groupExists, group := checkForUrlInGroups(egs, groupUrl)
		if groupExists { // main group exists
			if group.HasSubgroup() == true { // subgroup exists
				subgroup := group.GetSubgroup()
				subgroup.UpdateVerb(e.GetVerb(), e)
			} else { // main group there, create subgroup
				verbs := make(map[string]Domain.Endpoint)
				verbs[e.GetVerb()] = e

				newSubgroup := Domain.CreateEndpointGroup()
				newSubgroup.SetUrl(e.GetUrl())
				newSubgroup.SetVerbs(verbs)

				group.SetSubgroup(&newSubgroup)
			}
		} else { // create group + subgoup
			verbs := make(map[string]Domain.Endpoint)
			verbs[e.GetVerb()] = e

			subgroup := Domain.CreateEndpointGroup()
			subgroup.SetUrl(e.GetUrl())
			subgroup.SetVerbs(verbs)

			verbsMain := make(map[string]Domain.Endpoint)
			newGroup := Domain.CreateEndpointGroup()
			newGroup.SetUrl(groupUrlOriginal)
			newGroup.SetVerbs(verbsMain)
			newGroup.SetSubgroup(&subgroup)

			egs.AddNewGroup(groupUrl, &newGroup)
		}
	} else { // only main group
		groupUrl = removeLastCharIfNecessary(e.GetUrl())
		groupExists, group := checkForUrlInGroups(egs, groupUrl)

		if groupExists == true { // there is main group add only verb
			group.UpdateVerb(e.GetVerb(), e)
		} else {
			verbs := make(map[string]Domain.Endpoint)
			verbs[e.GetVerb()] = e
			newGroup := Domain.CreateEndpointGroup()
			newGroup.SetUrl(e.GetUrl())
			newGroup.SetVerbs(verbs)

			egs.AddNewGroup(groupUrl, &newGroup)
		}
	}
}

func (egs *EndpointGroups) AddNewGroup(groupName string, eg *Domain.EndpointGroup) {
	if len(egs.groups) == 0 {
		egs.Init()
	}

	egs.groups[groupName] = eg
}

func checkForUrlInGroups(egs *EndpointGroups, mainGroupUrl string) (bool, *Domain.EndpointGroup) {
	var group *Domain.EndpointGroup

	if g, ok := egs.groups[mainGroupUrl]; ok {
		return true, g
	}

	return false, group
}

func removeLastCharIfNecessary(checkString string) string {
	// Remove the last character if '/'
	if len(checkString) > 0 && checkString[len(checkString)-1] == '/' {
		checkString = checkString[:len(checkString)-1]
	}

	return checkString
}

func retrieveOriginalUrl(urlWithId string) string {
	var checkString string

	regexPattern := "\\{[^\\}]+\\}"
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return "/"
	}

	checkString = regex.ReplaceAllString(urlWithId, "")

	if len(checkString) > 0 && checkString[len(checkString)-2] == '/' && checkString[len(checkString)-1] == '/' {
		checkString = checkString[:len(checkString)-1]
	}

	return checkString
}
