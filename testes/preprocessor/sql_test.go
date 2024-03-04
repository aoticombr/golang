package main

import (
	"fmt"
	"testing"

	"github.com/aoticombr/golang/preprocesssql"
)

func Test_Sql(t *testing.T) {
	sql := `select
   'S' AS GERAR_ENTIDADE_INTEGRADA
  ,NULL AS ID_ENTIDADE_INTEGRADA
  ,os.cod_empresa
  ,nvl(os.numero_os_fabrica, os.numero_os) as numero_os_fabrica
  ,trunc(os.data_encerrada) as data_encerrada
  ,os.hora_encerrada as hora_encerrada
  from os
  where 1 = 1
    &macroxxx
	and os.cod_empresa = :COD_EMPRESA
	and os.numero_os > 0
	and os.status_os = 1
	and nvl(os.orcamento, 'N') = 'N'
	and rownum <= 10 
	and trunc(os.data_encerrada) >= trunc(nvl('hh:mm:ss',sysdate))
	and trunc(os.data_encerrada) >= trunc(nvl(:DT_INICIO,sysdate))
	and trunc(os.data_encerrada) < trunc(sysdate + 1)
	-- and trunc(os.data_encerrada) = :data_encerrada
	and not exists(SELECT 1 
				  FROM FAB_EI_OP5_MB_MGT a1
				  where a1.cod_empresa = os.cod_empresa
					and a1.numero_os_fabrica = nvl(os.numero_os_fabrica, os.numero_os)
					and trunc(a1.data_encerrada) = trunc(os.data_encerrada )
					and a1.hora_encerrada = os.hora_encerrada)`

	//sql = "select * from dual where id_xxx = :id_sss and xxx = :bbb_ff and yyy =:ccc_uu &hhhh');"
	Params, MacrosUpd, MacrosRead, err := preprocesssql.PreprocessSQL(sql, true, true, true, true, true)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, p := range Params.Items {
			fmt.Println("param.Name:", p.Name)
		}
		for _, m := range MacrosUpd.Items {
			fmt.Println("macro.Name:", m.Name)
		}
		for _, m := range MacrosRead.Items {
			fmt.Println("macro.Name:", m.Name)
		}
	}

}
