package Domain

type EndpointResult struct {
	StatusCode int
	Body       string
	Header     map[string][]string
}
