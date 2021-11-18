// +build linux darwin
package main

// This package is a plugin. Build it with `go build -buildmode=plugin -o shout.so`

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jvmatl/go-plugindemo/processors"
)

// ShoutProcessor will capitalize any byte slices passed in
type ShoutProcessor struct {
	Configured    bool `yaml:"Configured"`
	LogEverything bool `yaml:"LogEverything"`
}

// NewProcessor is more strongly typed, and a better way to go if you expect to have many plugins
func NewProcessor() processors.Processor {
	return &ShoutProcessor{}
}

// GenericNew is the quick and dirty way to do this, without needing the separate processors package
func GenericNew() interface{} {
	return &ShoutProcessor{}
}

// Init accepts configuration information for your processor object
func (p *ShoutProcessor) Init(config map[string]interface{}) error {
	var ok bool
	if p.LogEverything, ok = config["Log_everything"].(bool); !ok {
		return errors.New("invalid config")
	}

	p.Configured = true
	data, _ := json.Marshal(p)
	fmt.Println(string(data))
	return nil
}

// Process will take in a []byte and do something cool with it. :)
func (p *ShoutProcessor) Process(buf []byte) []byte {
	if p == nil || !p.Configured {
		panic(fmt.Sprintf("Unconfigured %T!", p))
	}

	if p.LogEverything {
		fmt.Printf("  Shouter got data: %v\n", buf)
	}

	return []byte(strings.ToUpper(string(buf)))
}
