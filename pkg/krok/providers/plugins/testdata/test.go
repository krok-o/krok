package main

import "fmt"

// Execute is the entrypoint for this plugin.
func Execute(payload string, opts ...interface{}) (string, bool, error) {
	fmt.Println("running plugin Test...")
	fmt.Println("got raw: ", payload)
	o := opts[0].(string)
	return payload + ":" + o, true, nil
}
