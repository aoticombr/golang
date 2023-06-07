package pg

var cfgglobal *ConfigPG

type ConfigPG struct {
	host     string
	user     string
	pass     string
	port     string
	database string
}

func GetConfigPG() *ConfigPG {
	if cfgglobal == nil {
		cfgglobal = &ConfigPG{}
	}
	return cfgglobal
}
