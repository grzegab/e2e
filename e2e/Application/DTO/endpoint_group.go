package DTO

import (
	"e2e/Domain"
)

type EndpointGroups struct {
	groups map[string]EndpointGroup
}

type EndpointGroup struct {
	id       string
	url      string
	verbs    map[string]Domain.Endpoint
	subgroup *EndpointGroup
}

func (eg *EndpointGroup) AddSubGroup(e EndpointGroup) {
	eg.subgroup = &e
}

//func AddEndpoint(groups map[string]EndpointGroup, endpoint Domain.Endpoint) map[string]EndpointGroup {
//	var mainGroupUrl string
//	var currentGroupUrl string
//
//	//check if url contains { and }
//	if strings.Contains(endpoint.GetUrl(), "{") && strings.Contains(endpoint.GetUrl(), "}") {
//		regexPattern := "\\{[^\\}]+\\}"
//		regex, err := regexp.Compile(regexPattern)
//		if err != nil {
//			mainGroupUrl = "/"
//		}
//
//		mainGroupUrl = removeLastCharIfNecessary(regex.ReplaceAllString(endpoint.GetUrl(), ""))
//		currentGroupUrl = removeLastCharIfNecessary(endpoint.GetUrl())
//
//		_, mainGroupExists := groups[mainGroupUrl]
//		if mainGroupExists == true { // There is main group already and need to add subgroup verb
//			group := groups[mainGroupUrl]
//			subgroup := group.group
//
//			if subgroup != nil { // subgroup exists
//				subgroup.verbs[endpoint.GetVerb()] = endpoint
//			} else { //no subgroup
//				verb := make(map[string]Domain.Endpoint)
//				verb[endpoint.GetVerb()] = endpoint
//				newGroup := EndpointGroup{
//					url:   currentGroupUrl,
//					verbs: verb,
//					group: nil,
//				}
//
//				group.AddSubGroup(newGroup)
//			}
//		} else { // Need to add main group, sub goup and verb
//			verb := make(map[string]Domain.Endpoint)
//			newGroup := EndpointGroup{
//				url:   endpoint.GetUrl(),
//				verbs: verb,
//				group: nil,
//			}
//
//			groups[mainGroupUrl] = newGroup
//		}
//	} else { // only main group
//		mainGroupUrl = removeLastCharIfNecessary(endpoint.GetUrl())
//
//		_, mainGroupExists := groups[mainGroupUrl]
//		if mainGroupExists == true { // there is main group add only verb
//			groups[mainGroupUrl].verbs[endpoint.GetVerb()] = endpoint
//		} else { // there is no main group
//			verb := make(map[string]Domain.Endpoint)
//			verb[endpoint.GetVerb()] = endpoint
//			newGroup := EndpointGroup{
//				url:   endpoint.GetUrl(),
//				verbs: verb,
//				group: nil,
//			}
//
//			newGroups := make(map[string]EndpointGroup)
//			newGroups[mainGroupUrl] = newGroup
//
//			groups = newGroups
//		}
//	}
//
//	return groups
//}

func removeLastCharIfNecessary(checkString string) string {
	// Remove the last character if '/'
	if len(checkString) > 0 && checkString[len(checkString)-1] == '/' {
		checkString = checkString[:len(checkString)-1]
	}

	return checkString
}

func (eg *EndpointGroup) SetId(newId string) {
	eg.id = newId
}

func (eg *EndpointGroup) GetVerbs() map[string]Domain.Endpoint {
	return eg.verbs
}

func (eg *EndpointGroup) GetInfo() (string, map[string]Domain.Endpoint) {
	return eg.url, eg.verbs
}

func (eg *EndpointGroup) GetUrl() string {
	return eg.url
}
