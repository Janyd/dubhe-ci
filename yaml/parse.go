package yaml

import (
	"github.com/buildkite/yaml"
	"io/ioutil"
)

func ParseString(s string) (Resource, error) {
	return ParseByte([]byte(s))
}

func ParseFile(p string) (Resource, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return ParseByte(data)
}

func ParseByte(b []byte) (Resource, error) {
	resource := &RawResource{}
	err := yaml.Unmarshal(b, resource)
	if err != nil {
		return nil, err
	}
	resource.Data = b

	res, err := parseRaw(resource)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func parseRaw(r *RawResource) (Resource, error) {
	var obj Resource
	switch r.Kind {
	default:
		obj = new(Pipeline)
	}
	err := yaml.Unmarshal(r.Data, obj)
	return obj, err
}
