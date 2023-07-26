package main

import (
	json "github.com/aoticombr/golang/jsonconfig"
)

func main() {

	json.NewJsonConfigGlobal().Name = "config.json"
	json.NewJsonConfigGlobal().Load()

}
