package api

import "github.com/aoticombr/golang/dbconndataset"

type Processo struct {
	Tenants *Tenants
}

func (pr *Processo) Carregar(conn *dbconndataset.ConnDataSet) {
	pr.Tenants = &Tenants{}
}
