package main

import (
	"fmt"

	ora "github.com/aoticombr/godataset/connection/ora"
	ds "github.com/aoticombr/godataset/dataset"
)

func main() {
	cf := ora.GetConfigOra().Load()
	conn := ora.GetConnOra()
	q := ds.GetDataSetNew(conn.GetDB())
	q.AddSql("SELECT CDCLIENTE FROM CLIENTE	where rownum <= 10")
	q.Open()
	q.First()
	for !q.Eof {
		fmt.Println(q.FieldByName("CDCLIENTE").AsInt64())
		q.Next()
	}

}
