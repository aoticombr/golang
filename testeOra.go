package main

import (
	"fmt"

	ora "github.com/aoticombr/godataset/connection/ora"
	ds "github.com/aoticombr/godataset/dataset"
)

func main() {

	conn := ora.GetConnOra()
	q := ds.GetDataSetNew(conn.GetDB())
	q.AddSql("SELECT * FROM CONTASPAGAR where rownum <= 10")
	q.Open()
	q.First()
	for !q.Eof {
		fmt.Println(
			//	q.FieldByName("cdcontaspagar").AsInt64(),
			q.FieldByName("valor").AsFloat(),
			//q.FieldByName("juros").AsString()
		)
		q.Next()
	}

}
