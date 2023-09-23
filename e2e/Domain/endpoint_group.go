package Domain

type EndpointGroup struct {
	id       string
	url      string
	verbs    map[string]Endpoint
	subgroup *EndpointGroup
}

func CreateEndpointGroup() EndpointGroup {
	e := EndpointGroup{
		id:       "",
		url:      "",
		verbs:    nil,
		subgroup: nil,
	}

	return e
}

func (eg *EndpointGroup) GetVerbs() map[string]Endpoint {
	return eg.verbs
}

func (eg *EndpointGroup) GetEndpoints() []*Endpoint {
	var verbs []*Endpoint

	for _, e := range eg.verbs {
		verbs = append(verbs, &e)
	}

	return verbs
}

func (eg *EndpointGroup) HasSubgroup() bool {
	return eg.subgroup != nil
}

func (eg *EndpointGroup) GetSubgroup() *EndpointGroup {
	return eg.subgroup
}

func (eg *EndpointGroup) UpdateVerb(v string, e Endpoint) {
	eg.verbs[v] = e // @TODO: what if no verb exists -> check for error
}

func (eg *EndpointGroup) SetUrl(url string) {
	eg.url = url
}

func (eg *EndpointGroup) SetVerbs(v map[string]Endpoint) {
	eg.verbs = v
}

func (eg *EndpointGroup) SetSubgroup(s *EndpointGroup) {
	eg.subgroup = s
}
