package main

import (
	"fmt"
	"testing"

	gorm "github.com/aoticombr/golang/orm"
)

type tb0001 struct {
	Id    string  `json:"id" column:"id" table:"tb0001" primarykey:"true"`
	Email string  `json:"email" column:"email" primarykey:"true"`
	Nome  string  `json:"nome" column:"nome"`
	Cpf   *string `json:"cpf" column:"cpf,omitempty"`
}

func Test_tipo1(t *testing.T) {
	// dados := &tb0001{
	// 	Id:    "1",
	// 	Email: "paulo@example.com",
	// 	Nome:  "Paulo",
	// }

	tb := gorm.NewTable(&tb0001{})

	fmt.Println(tb)
	fmt.Println(tb.SqlInsert())
	fmt.Println(tb.SqlUpdate())
	fmt.Println(tb.SqlDelete())

}

func Test_tipo2(t *testing.T) {
	var dados []*tb0001
	for i := 0; i < 10; i++ {
		dado := &tb0001{
			Id:    fmt.Sprintf("%d", i),
			Email: fmt.Sprintf("%d", i) + "@example.com",
			Nome:  "Cliente " + fmt.Sprintf("%d", i),
		}
		if i%2 == 0 {
			cpf := "cpf " + fmt.Sprintf("%d", i)
			dado.Cpf = &cpf
		}
		dados = append(dados, dado)
	}
	for _, dado := range dados {

		tb := gorm.NewTable(dado)
		//fmt.Println(dado)
		fmt.Println(tb)
		fmt.Println(tb.SqlInsert())
		fmt.Println(tb.SqlInsertOmiteNil())
		fmt.Println(tb.SqlUpdate())
		fmt.Println(tb.SqlDelete())
		fmt.Println("----------------------------")
	}

}
