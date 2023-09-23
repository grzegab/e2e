package Domain

type Params struct {
	Key   string
	Value string
}

func CreateParams(key string, value string) Params {
	var params Params

	params.Key = key
	params.Value = value

	return params
}

func (p *Params) GetValue() string {
	return p.Value
}

func (p *Params) GetKey() string {
	return p.Key
}
