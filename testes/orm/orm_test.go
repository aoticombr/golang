package main

import (
	"fmt"
	"testing"

	orm "github.com/aoticombr/golang/orm"
)

type tb0001 struct {
	Id    string  `json:"id" column:"id,insert,primarykey" table:"tb0001" `
	Email string  `json:"email" column:"email,insert,update" `
	Nome  string  `json:"nome" column:"nome,insert,update"`
	Cpf   *string `json:"cpf" column:"cpf,insert,update,omitempty"`
	CRUD  string  `json:"CRUD,crud" `
}

type tb0002 struct {
	Id    string  `json:"id" column:"id,insert,primarykey" table:"tb0001"`
	Email string  `json:"email" column:"email,insert,update" `
	Nome  string  `json:"nome" column:"nome,insert,update"`
	Cpf   *string `json:"cpf" column:"cpf,insert,update,omitempty"`
	CRUD  string  `json:"CRUD,crud" `
}
type tb0003 struct {
	Id    string  `json:"id" column:"id,insert,primarykey" table:"tb0001"`
	Email string  `json:"email" column:"email,insert,update" `
	Nome  string  `json:"nome" column:"nome,insert,update"`
	Cpf   *string `json:"cpf" column:"cpf,insert,update,omitempty"`
	CRUD  string  `json:"CRUD,crud" `
}

func Test_tipo1(t *testing.T) {
	// dados := &tb0001{
	// 	Id:    "1",
	// 	Email: "paulo@example.com",
	// 	Nome:  "Paulo",
	// }

	tb := orm.NewTable(&tb0001{})

	fmt.Println(tb)
	fmt.Println(tb.SqlInsert())
	fmt.Println(tb.SqlUpdate())
	fmt.Println(tb.SqlDelete())
	fmt.Println(tb.SqlStatus())

}

func Test_tipo2(t *testing.T) {
	var dados []*tb0001
	for i := 0; i < 4; i++ {
		dado := &tb0001{
			Id:    fmt.Sprintf("%d", i),
			Email: fmt.Sprintf("%d", i) + "@example.com",
			Nome:  "Cliente " + fmt.Sprintf("%d", i),
		}
		if i%2 == 0 {
			cpf := "cpf " + fmt.Sprintf("%d", i)
			dado.Cpf = &cpf
			dado.CRUD = "new"
		} else if i%3 == 0 {
			cpf := "cpf " + fmt.Sprintf("%d", i)
			dado.Cpf = &cpf
			dado.CRUD = "del"
		} else {
			dado.CRUD = "old"
		}
		dados = append(dados, dado)
	}
	for _, dado := range dados {

		tb := orm.NewTable(dado)
		tb.Options.Delete = orm.D_Disable
		//fmt.Println(dado)
		fmt.Println(tb)
		//fmt.Println(tb.SqlInsert())
		//fmt.Println(tb.SqlUpdate())
		//fmt.Println(tb.SqlDelete())
		fmt.Println(tb.SqlStatus())
		fmt.Println("----------------------------")
	}

}

func Test_tipoValid(t *testing.T) {
	fmt.Println(
		`################################
Validando Primary Key
################################`)
	var dados1 []*tb0002

	for i := 0; i < 1; i++ {
		dado := &tb0002{
			Id:    fmt.Sprintf("%d", i),
			Email: fmt.Sprintf("%d", i) + "@example.com",
			Nome:  "Cliente " + fmt.Sprintf("%d", i),
		}
		if i%2 == 0 {
			cpf := "cpf " + fmt.Sprintf("%d", i)
			dado.Cpf = &cpf
			dado.CRUD = "new"
		} else {
			dado.CRUD = "old"
		}
		dados1 = append(dados1, dado)
	}
	for _, dado := range dados1 {

		tb := orm.NewTable(dado)

		//fmt.Println(dado)
		fmt.Println(tb)
		fmt.Println(tb.SqlInsert())
		fmt.Println(tb.SqlUpdate())
		fmt.Println(tb.SqlDelete())
		fmt.Println(tb.SqlStatus())
		fmt.Println("----------------------------")
	}
	fmt.Println(
		`################################
Validando Table
################################`)
	var dados2 []*tb0003
	for i := 0; i < 1; i++ {
		dado := &tb0003{
			Id:    fmt.Sprintf("%d", i),
			Email: fmt.Sprintf("%d", i) + "@example.com",
			Nome:  "Cliente " + fmt.Sprintf("%d", i),
		}
		if i%2 == 0 {
			cpf := "cpf " + fmt.Sprintf("%d", i)
			dado.Cpf = &cpf
			dado.CRUD = "new"
		} else {
			dado.CRUD = "old"
		}
		dados2 = append(dados2, dado)
	}
	for _, dado := range dados2 {

		tb := orm.NewTable(dado)
		//fmt.Println(dado)
		fmt.Println(tb)
		fmt.Println(tb.SqlInsert())
		fmt.Println(tb.SqlUpdate())
		fmt.Println(tb.SqlDelete())
		fmt.Println(tb.SqlStatus())
		fmt.Println("----------------------------")
	}

}
