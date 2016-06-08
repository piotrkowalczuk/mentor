package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type variable struct {
	Default     interface{}
	Description string
	Fields      []string `hcl:",decodedFields"`
}

// Alphasfile ...
type Alphasfile struct {
	Variables map[string]*variable `hcl:"variable,"`
	Service   []*Service           `hcl:"service,expand"`
}

type Service struct {
	Name      string                 `hcl:"name,key"`
	Import    string                 `hcl:"import"`
	Branch    string                 `hlc:"branch"`
	Install   string                 `hlc:"install"`
	Arguments map[string]interface{} `hcl:"arguments"`
	Log       string                 `hcl:"log"`
	Color     string                 `hcl:"color"`
	Fields    []string               `hcl:",decodedFields"`
}

// JoinArgs ...
func JoinArgs(args map[string]interface{}) []string {
	r := make([]string, 0, len(args))
	for flag, value := range args {
		r = append(r, fmt.Sprintf("-%s=%v", flag, value))
	}

	return r
}

func openAlphasfile(path string) (*Alphasfile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("alpha: Alphasfile read error: %s", err.Error())
	}

	var af Alphasfile
	if err = hcl.Unmarshal(b, &af); err != nil {
		return nil, fmt.Errorf("alpha: Alphasfile parsing error: %s", err.Error())
	}

	for i, r := range af.Service {
		if r.Color == "" {
			r.Color = colors[i%len(colors)]
		}
	}

	return &af, nil
}
