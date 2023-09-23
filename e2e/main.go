package main

import (
	"e2e/Application"
	"flag"
	"os"
)

func main() {
	var listen bool

	flag.BoolVar(&listen, "l", false, "Web server fo testing")
	flag.Parse()

	if listen { // start webserver and show endpoint status, with option to test things, edit config
		Application.StartTestWebserver()
	} else { // get all endpoints and check them (e.g. in pipeline)
		Application.RunTests()
	}

	os.Exit(0)
}
