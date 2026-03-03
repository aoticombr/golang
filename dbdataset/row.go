package dbdataset

import "github.com/aoticombr/golang/variant"

type Row struct {
	List map[string]*variant.Variant
}

func NewRow() Row {
	row := Row{
		List: make(map[string]*variant.Variant),
	}
	return row
}
