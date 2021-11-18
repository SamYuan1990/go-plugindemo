package main

import (
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"reflect"

	"github.com/jvmatl/go-plugindemo/processors"
	"gopkg.in/yaml.v2"
)

// simple demo app to show a couple of ways structure an app that can load
// multiple plugins that meet a common interface.
// See: https://stackoverflow.com/q/56693941/3117035 for context

func main() {
	fmt.Println("Using generic interface{} constructor")
	easyWay()

	fmt.Println("--------")

	fmt.Println("Using constructor that returns a defined interface")
	betterWay()
}

func easyWay() {
	// for demo, we'll just look for the plugin binary in the directory where it was built.
	pluginName := "shout"
	pluginFile, _ := filepath.Abs(fmt.Sprintf("./processors/%s/%s.so", pluginName, pluginName))

	p, err := plugin.Open(pluginFile)
	if err != nil {
		log.Fatalf("Error opening plugin %s: %v", pluginName, err)
	}

	newIntf, err := p.Lookup("GenericNew")
	if err != nil {
		log.Fatalf("Error looking up GenericNew() func in plugin %s: %v", pluginName, err)
	}

	newProc, ok := newIntf.(func() interface{})
	if !ok {
		log.Fatalf("Error casting newIntf for plugin %s: %T", pluginName, newIntf)
	}

	shoutProc := newProc().(processors.Processor) // call the constructor, get a new ShoutProcessor

	// Initialize my new Processor
	shoutProc.Init(map[string]interface{}{"Log_everything": true})
	data, err := yaml.Marshal(&shoutProc)
	fmt.Println(string(data))

	// Process some bytes!
	message := "whisper"
	fmt.Printf("  Before processing: %s\n", message)
	output := shoutProc.Process([]byte(message))
	fmt.Printf("  After processing: %s\n", output)
}

func betterWay() {
	// for demo, we'll just look for the plugin binary in the directory where it was built.
	pluginName := "shout"
	pluginFile, _ := filepath.Abs(fmt.Sprintf("./processors/%s/%s.so", pluginName, pluginName))

	p, err := plugin.Open(pluginFile)
	if err != nil {
		log.Fatalf("Error opening plugin %s: %v", pluginName, err)
	}

	newProcIntf, err := p.Lookup("NewProcessor")
	if err != nil {
		log.Fatalf("Error looking up New() func in plugin %s: %v", pluginName, err)
	}

	newProc, ok := newProcIntf.(func() processors.Processor)
	if !ok {
		log.Fatalf("Error casting procNewIntf for plugin %s: %T", pluginName, newProcIntf)
	}

	shoutProc := newProc() // call the constructor, get a new ShoutProcessor

	// Initialize my new Processor
	shoutProc.Init(map[string]interface{}{"Log_everything": true})

	data, err := yaml.Marshal(&shoutProc)
	fmt.Println(string(data))

	shoutProc2 := newProc() // call the constructor, get a new ShoutProcessor
	// proof for struct init from raw
	err = yaml.Unmarshal(data, shoutProc2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("cool")
	// proof for get type from plugin
	fmt.Println(reflect.TypeOf(shoutProc2))
	// proof for use reflect with plugin
	test_map := make(map[reflect.Type]string)
	test_map[reflect.TypeOf(shoutProc2)] = "abc"
	str, ok := test_map[reflect.TypeOf(shoutProc)]
	if ok {
		fmt.Println(str)
	}
	// proof for method execute
	message := "whisper"
	fmt.Printf("  Before processing: %s\n", message)
	output := shoutProc2.Process([]byte(message))
	fmt.Printf("  After processing: %s\n", output)
}
