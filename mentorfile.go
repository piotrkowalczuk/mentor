package main

import "fmt"

type variable struct {
	Default     interface{}
	Description string
	Fields      []string `hcl:",decodedFields"`
}

// Mentorfile ...
type Mentorfile struct {
	Variables map[string]*variable `hcl:"variable,"`
	Services  []*Service           `hcl:"service,expand"`
}

type Service struct {
	Name      string                 `hcl:"name,key"`
	Path      string                 `hcl:"path"`
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
