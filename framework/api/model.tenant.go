package api

import (
	"github.com/aoticombr/golang/dbconndataset"
)

type Tenant struct {
	Name       string
	Connection *dbconndataset.ConnDataSet
}

type Tenants []*Tenant

func (ts *Tenants) Len() int {
	return len(*ts)
}

func (ts *Tenants) FindByName(value string) *Tenant {
	for _, tenant := range *ts {
		if tenant.Name == value {
			return tenant
		}
	}
	return nil
}

func (ts *Tenants) Add(name string, conn *dbconndataset.ConnDataSet) {
	t := &Tenant{
		Name:       name,
		Connection: conn,
	}
	*ts = append(*ts, t)

}
