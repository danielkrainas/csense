package main

import (
	"fmt"
	"net/http"

	"github.com/danielkrainas/csense/api/client"
	"github.com/danielkrainas/csense/api/v1"
)

func main() {
	const ENDPOINT = "http://localhost:9181"

	// Create a new client
	c := client.New(ENDPOINT, http.DefaultClient)
	fmt.Printf("created new client to %q\n", ENDPOINT)

	// Check V1 endpoint is good and healthy
	//=====================================
	err := c.Ping()
	if err != nil {
		panic("error sending ping")
	}

	fmt.Println("sent ping")

	// Create a hook
	//=====================================
	err = c.Hooks().CreateHook(&v1.NewHookRequest{
		Name:   "Foo Hook",
		TTL:    0,
		Url:    "http://localhost:9181/v1",
		Format: v1.FormatJSON,
		Events: []v1.EventType{v1.EventCreate},
		Criteria: &v1.Criteria{
			ImageName: &v1.Condition{
				Op:    v1.OperandEqual,
				Value: "registry",
			},
		},
	})

	if err != nil {
		panic("error creating hook")
	}
}
