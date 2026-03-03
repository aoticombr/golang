package config

import "strings"

type Cert struct {
	Name         string
	Ativo        bool
	Path_crt     string
	Path_private string
	Pass         string
}

func (c *Certs) GetCert(value string) (bool, *Cert) {
	for _, v := range *c {
		if v.Ativo {
			if strings.ToUpper(v.Name) == strings.ToUpper(value) {
				return true, v
			}
		}
	}
	return false, nil
}

type Certs []*Cert
