package main

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v2"
)

type SAMTemplate struct {
	Globals struct {
		Function FunctionSetting `yaml:"Function,omitempty"`
	} `yaml:"Globals,omitempty"`
	Resources map[string]Function `yaml:"Resources,omitempty"`
}

func (t SAMTemplate) Functions() map[string]Function {
	functions := make(map[string]Function)
	globalEnvVars := t.Globals.Function.Environment.ParseMap()
	for name, res := range t.Resources {
		if res.Type == "AWS::Serverless::Function" {
			res.Properties.Environment.Variables = res.Properties.Environment.ParseMap()
			if res.Properties.CodeUri == "" {
				res.Properties.CodeUri = t.Globals.Function.CodeUri
			}
			if res.Properties.Runtime == "" {
				res.Properties.Runtime = t.Globals.Function.Runtime
			}
			if res.Properties.MemorySize == 0 {
				res.Properties.MemorySize = t.Globals.Function.MemorySize
			}
			if res.Properties.Timeout == 0 {
				res.Properties.Timeout = t.Globals.Function.Timeout
			}
			if res.Properties.Handler == "" {
				res.Properties.Handler = t.Globals.Function.Handler
			}
			if res.Properties.Tracing == "" {
				res.Properties.Tracing = t.Globals.Function.Tracing
			}
			for k, v := range globalEnvVars {
				if _, ok := res.Properties.Environment.Variables[k]; !ok {
					res.Properties.Environment.Variables[k] = v
				}
			}
			res.Properties.Environment.InnerVars = nil
			functions[name] = res
		}
	}
	return functions
}

type FunctionSetting struct {
	CodeUri     string `yaml:"CodeUri,omitempty"`
	Runtime     string `yaml:"Runtime,omitempty"`
	MemorySize  int    `yaml:"MemorySize,omitempty"`
	Timeout     int    `yaml:"Timeout,omitempty"`
	Handler     string `yaml:"Handler,omitempty"`
	Tracing     string `yaml:"Tracing,omitempty"`
	Environment Env    `yaml:"Environment,omitempty"`
}

type Env struct {
	InnerVars yaml.MapSlice `yaml:"Variables,omitempty"`
	Variables map[string]string
}

func (e Env) ParseMap() map[string]string {
	m := make(map[string]string)
	for _, item := range e.InnerVars {
		key := ""
		switch v := item.Key.(type) {
		case string:
			key = v
		default:
			panic(fmt.Sprintf("Type of key: %s", reflect.TypeOf(item.Key)))
		}
		switch v := item.Value.(type) {
		case string:
			m[key] = v
		case bool:
			if v {
				m[key] = "true"
			} else {
				m[key] = "false"
			}
		default:
			fmt.Printf("Non string type of key %s: %v (%T)\n", key, v, v)
			m[key] = "dummy"
		}
	}
	return m
}

type Function struct {
	Type       string          `yaml:"Type,omitempty"`
	Properties FunctionSetting `yaml:"Properties,omitempty"`
}
