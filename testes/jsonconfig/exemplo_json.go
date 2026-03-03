package main

import (
	json "github.com/aoticombr/golang/jsonconfig"
)

func main() {

	js := json.NewJsonConfigGlobal()
	js.Name = "config.json"
	js.Load()

	for key, boot := range js.GetConfig().Boots {
		println("Boot ", key, ": ", boot.Name)
		for key, schema := range boot.GetSchemas() {
			println("Schema ", key, ": ", schema.Schema)
		}
	}

}
