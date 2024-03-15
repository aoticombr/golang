package main

import (
	"fmt"
	"testing"

	"github.com/aoticombr/golang/preprocesssql"
)

func Test_Sql1(t *testing.T) {

	sql := "select * from dual where id_xxx = :id_sss and xxx = :bbb_ff and yyy =:ccc_uu &hhhh');"
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

func Test_Sql2(t *testing.T) {
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
func Test_Sql3(t *testing.T) {
	sql := `select
	 c.cod_empresa
	,c.documento
	,c.tipo
	,c.chassi
	, pk_util.Retirar_Nao_Alfa_Num(case when c.tipo = 3 then oa.placa else osd.placa end) licensePlateNumber
	, 'DMS' as initiator
	, 'KM' as distanceUnit
	, c.chassi as finOrVin
	, 'BRL' as currency
	, 'ACCEPTANCE_DATE_PENDING' as state
	from fab_mov_ei_cria_mb_xa mei
	left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
	left join xentry_controle_envio c on c.id = ei.id_controle
	left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
	left join os_dados_veiculos osd on c.cod_empresa = osd.cod_empresa and c.documento = osd.numero_os and c.tipo in (1,2)
	left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
	where mei.id_mov_mb_xa = :id_mov`

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
func Test_Sql4(t *testing.T) {
	sql := `select
			case c.tipo
			when 1 then 'OSID'
			when 2 then 'ORC'
			when 3 then 'AGE'
			end ||'-'|| abs( c.documento ) as  orderId
			, case
			when c.tipo in (1,2) then pkg_xentry.getDataISO8601(trunc(o.data_emissao),o.hora_emissao,c.cod_empresa,'S')
			when c.tipo = 3 then pkg_xentry.getDataISO8601(trunc(oa.data_agendada),to_char(oa.data_agendada,'hh24:mi'),c.cod_empresa,'S')
			end as acceptanceDate
			--, case when c.tipo in (1,2) and length(trim(o.hora_prometida)) = 5 and o.data_prometida >= o.data_emissao then
			--case when (o.data_prometida > o.data_emissao) or o.hora_prometida > o.hora_emissao
			--then pkg_xentry.getDataISO8601(trunc(o.data_prometida),o.hora_prometida,c.cod_empresa,'S') end
			--when c.tipo = 3 and length(trim(to_char(oa.hora_prometida,'hh24:mi'))) = 5 and oa.data_prometida >= trunc(oa.data_agendada) then
			--case when (oa.data_prometida > trunc(oa.data_agendada)) or to_char(oa.hora_prometida,'hh24:mi') > to_char(oa.data_agendada,'hh24:mi')
			--then pkg_xentry.getDataISO8601(trunc(oa.data_prometida),to_char(oa.hora_prometida,'hh24:mi'),c.cod_empresa,'S') end
			--end as deliveryDate
			, 'Empresa: ' || c.cod_empresa || ' - ' || case c.tipo when 1 then 'Numero da Ordem de Servico Fabrica: ' || nvl(o.numero_os_fabrica,o.numero_os)
			when 2 then 'Numero do Orcamento: ' || abs(o.numero_os)
			when 3 then 'Codigo do Agendamento: ' || oa.cod_os_agenda  end as dmsOrderInformation
			from fab_mov_ei_cria_mb_xa mei
			left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
			left join xentry_controle_envio c on c.id = ei.id_controle
			left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
			left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
			where mei.id_mov_mb_xa = :ID_MOVacTy`

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
func Test_Sql5(t *testing.T) {
	sql := `select
    substr(trim(eu.nome_completo),1, instr(trim(eu.nome_completo),' ') - 1) as firstName
  , substr(trim(eu.nome_completo), instr(trim(eu.nome_completo),' ') + 1 , length(trim(eu.nome_completo))) as lastName
  , eu.nome as identification
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
  left join empresas_usuarios eu on case when c.tipo in (1,2) then o.nome
                                         when c.tipo = 3 then nvl(oa.consultor,oa.nome) end = eu.nome
 where mei.id_mov_mb_xa = :ID_MOV`

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
func Test_Sql6(t *testing.T) {
	sql := ` select
    case when cl.cod_classe = 'J' then lpad(cd.cod_cliente,14,0) else lpad(cd.cod_cliente,11,0) end as cod_cliente
  , substr(trim(cd.nome),1, instr(trim(cd.nome),' ') - 1) as firstName
  , substr(trim(cd.nome), instr(trim(cd.nome),' ') + 1 , length(trim(cd.nome))) as lastName
  , coalesce(cl.endereco_eletronico, cl.email_nfe) as email
  , case when cl.cod_classe = 'J'
           then coalesce( cl.prefixo_com || cl.telefone_com
                        , cl.prefixo_cel || cl.telefone_cel
                        , cl.prefixo_res || cl.telefone_res )
           else coalesce( cl.prefixo_cel || cl.telefone_cel
                        , cl.prefixo_res || cl.telefone_res
                        , cl.prefixo_com || cl.telefone_com ) end as phone
  , remove_caracteres_json( case when cl.cod_classe = 'J' then cd.nome else cl.nome_empresa_trab end ) as companyName
  , remove_caracteres_json( case when cid_com.cod_cidades > 0 and length(trim(cl.rua_com)) > 3 then
        'End.:'||cl.rua_com || case when trim(cl.fachada_com) is not null then ' N.:'|| cl.fachada_com end ||' Bairro:'|| cl.bairro_com
        || case when trim(cl.complemento_com) is not null then ' Comp.:'|| cl.complemento_com end
        ||' Cep:'|| cl.cep_com ||' Cid.:'|| trim(cid_com.descricao) ||' UF:'|| cid_com.uf
       end ) as businessCustomerAddress
  , remove_caracteres_json( case when cid_res.cod_cidades > 0 and length(trim(cl.rua_res)) > 3 then
        'End.:'||cl.rua_res || case when trim(cl.fachada_res) is not null then ' N.:'|| cl.fachada_res end ||' Bairro:'|| cl.bairro_res
        || case when trim(cl.complemento_res) is not null then ' Comp.:'|| cl.complemento_res end
        ||' Cep:'|| cl.cep_res ||' Cid.:'|| trim(cid_res.descricao) ||' UF:'|| cid_res.uf
       end ) as privateCustomerAddress
  , mot.nome_do_motorista as vehiclePickupPerson
  , decode(o.cliente_aguardou,'S','true','false') as waiting
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
  left join cliente_diverso cd on case when c.tipo in (1,2) then o.cod_cliente
                                       when c.tipo = 3 then oa.cod_cliente end = cd.cod_cliente
  left join clientes cl on cd.cod_cliente = cl.cod_cliente
  left join motoristas mot on mot.codigo_motorista = o.codigo_motorista
  left join cidades cid_com on cid_com.cod_cidades = cl.cod_cid_com and cid_com.uf = cl.uf_com
  left join cidades cid_res on cid_res.cod_cidades = cl.cod_cid_res and cid_res.uf = cl.uf_res
 where mei.id_mov_mb_xa = 1 --:ID_MOV`

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
func Test_Sql7(t *testing.T) {
	sql := ` select
	'DMS' as odometerSource
   , case when osd.horimetro > 0 then osd.horimetro else null end as operatingTimeValueInHours
   , case when osd.combustivel >= 8 then 1
		  when nvl(osd.combustivel,0) <= 0 then 0
		   else osd.combustivel / 8 end as fuelIndicatorValue
   , pkg_xentry.getUltimaPassagemOficina(c.chassi,'N',c.cod_empresa,'S') as previousServiceDate
   from fab_mov_ei_cria_mb_xa mei
   left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
   left join xentry_controle_envio c on c.id = ei.id_controle
   left join os_dados_veiculos osd on c.cod_empresa = osd.cod_empresa and c.documento = osd.numero_os and c.tipo in (1,2)
   left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
  where mei.id_mov_mb_xa = 1 --:ID_MOV`

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
func Test_Sql8(t *testing.T) {
	sql := ` select
    case when c.tipo in (1,2) then osd.km
      else nvl(oa.km,0) end valor
  , 'KM' as unit
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os_dados_veiculos osd on c.cod_empresa = osd.cod_empresa and c.documento = osd.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
 where mei.id_mov_mb_xa = 1 --:ID_MOV`

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
func Test_Sql9(t *testing.T) {
	sql := ` select
    trim(remove_caracteres_json(ma.descricao_marca)) as brand
  , trim(remove_caracteres_json(p.descricao_produto)) ||' / '||trim(remove_caracteres_json(pm.descricao_modelo)) as familyText
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os_dados_veiculos osd on c.cod_empresa = osd.cod_empresa and c.documento = osd.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
  left join produtos_modelos pm on pm.cod_produto = case when c.tipo in (1,2) then osd.cod_produto else oa.cod_produto end
                               and pm.cod_modelo = case when c.tipo in (1,2) then osd.cod_modelo else oa.cod_modelo end
  left join produtos p on pm.cod_produto = p.cod_produto
  left join marcas ma on ma.cod_marca = p.cod_marca
 where mei.id_mov_mb_xa = 1 --:ID_MOV`

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
func Test_Sql10(t *testing.T) {
	sql := ` select
    case c.tipo when 3 then oar.item else ori.item end as item
  , 'DMS' as origin
  , remove_caracteres_json( case c.tipo when 3 then oar.descricao else ori.descricao end  ) as title
  , case when c.tipo <> 3 and upper(ori.cod_externo) = 'MAINTENANCE' then ori.cod_externo
         when c.tipo <> 3 and upper(ori.cod_externo) = 'COMPLAINT' then ori.cod_externo
         when c.tipo <> 3 and upper(ori.cod_externo) = 'CUSTOMER_REQUEST' then ori.cod_externo
      else 'UNASSIGNED' end as tipo
  , substr(remove_caracteres_json( case c.tipo when 3 then oar.comentario else ori.observacao end ),1,300) as customerStatementNote
  , case c.tipo when 3 then oar.item else ori.item end as dmsExternalId
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os_original ori on c.cod_empresa = ori.cod_empresa and c.documento = ori.numero_os and c.tipo in (1,2)
  left join os_agenda_reclamacao oar on c.cod_empresa = oar.cod_empresa and c.documento = oar.cod_os_agenda and c.tipo = 3
 where mei.id_mov_mb_xa = 1 --:ID_MOV
 order by 1`

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
func Test_Sql11(t *testing.T) {
	sql := ` select
    case c.tipo when 1 then oss.item
                when 2 then oso.item
                when 3 then oas.item
                  end item
 ,  'DMS' as origin
 ,  s.cod_servico as operationNumber
 ,  case when c.tipo = 1 and oss.status = 2 then 'true' else 'false' end completed
 ,  nvl(case c.tipo when 1 then oss.tempo_padrao
               when 2 then oso.tempo
               when 3 then oas.tempo_padrao
                 end,s.tempo_padrao) as operationValue
 ,  'H' as operationUnit
 ,  case c.tipo when 1 then oss.item
                when 2 then oso.item
                when 3 then oas.item
                  end ||'#'|| s.cod_servico as dmsExternalId
 ,  case c.tipo when 1 then o.tipo
                when 2 then o.tipo
                when 3 then coalesce(oa.tipo,oa.tipo_os)
                  end as invoiceCode
 ,  remove_caracteres_json(s.descricao_servico) as descricao
 ,  case c.tipo when 1 then oss.preco_venda
                when 2 then oso.preco_venda
                when 3 then oas.preco_venda
                  end as priceGrossValue
 ,  case c.tipo when 1 then oss.preco_venda - nvl(oss.desconto_por_serv,0)
                when 2 then oso.preco_venda - nvl(oso.desconto_por_serv,0)
                when 3 then oas.preco_venda - nvl(oas.valor_desconto_serv,0)
                  end as priceNetValue
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
  left join os_servicos oss on c.cod_empresa = oss.cod_empresa and c.documento = oss.numero_os and c.tipo = 1
  left join os_serv_orc oso on c.cod_empresa = oso.cod_empresa and c.documento = oso.numero_os and c.tipo = 2
  left join os_agenda_servicos oas on c.cod_empresa = oas.cod_empresa and c.documento = oas.cod_os_agenda and c.tipo = 3
  left join servicos s on s.cod_servico = case c.tipo when 1 then oss.cod_servico
                                                      when 2 then oso.cod_servico
                                                        else oas.cod_servico end
  left join os_original ori on ori.cod_empresa = c.cod_empresa and ori.numero_os = c.documento and ori.item = case c.tipo when 1 then oss.item
                                                                                                                          when 2 then oso.item
                                                                                                                            else null end
  left join os_agenda_reclamacao oar on oar.cod_empresa = c.cod_empresa and oar.cod_os_agenda = c.documento and oar.item = case c.tipo when 3 then oas.item
                                                                                                                            else null end
  left join reclamacao_servicos recs on recs.cod_reclamacao = case when c.tipo in (1,2) then ori.cod_reclamacao else oar.cod_reclamacao end
                                    and recs.cod_servico = case c.tipo when 1 then oss.cod_servico
                                                                       when 2 then oso.cod_servico
                                                                         else oas.cod_servico end
  left join reclamacoes_padroes rec on rec.cod_reclamacao = recs.cod_reclamacao
 where mei.id_mov_mb_xa = 1 --:ID_MOV
   and s.cod_servico is not null
   and nvl(rec.eh_kit,'N') != 'S'`

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
func Test_Sql12(t *testing.T) {
	sql := `   select
    case c.tipo when 1 then ore.item
                when 2 then nvl(ooi.item,1)
                when 3 then oap.item
                  end as item
 ,  i.cod_item
 ,  case c.tipo when 1 then ore.cod_fornecedor
                when 2 then ooi.cod_fornecedor
                when 3 then oap.cod_fornecedor
                  end as cod_fornecedor
 ,  case c.tipo when 1 then ore.requisicao
                when 2 then 2
                when 3 then 3
                  end as requisicao
 ,  'DMS' as origin
 ,  remove_caracteres_json(i.descricao) as name
 ,  case c.tipo when 1 then ore.item ||'#'||ore.cod_fornecedor
                when 2 then NVL(ooi.item,1) ||'#'||ooi.cod_fornecedor
                when 3 then oap.item ||'#'||oap.cod_fornecedor
                  end ||'#'|| i.cod_item as dmsExternalId
 ,  case c.tipo when 1 then o.tipo
                when 2 then o.tipo
                when 3 then coalesce(oa.tipo,oa.tipo_os)
                  end as invoiceCode
 ,  case c.tipo when 1 then ore.quantidade
                when 2 then ooi.quantidade
                when 3 then oap.qtde
                  end as quantity
 ,  case c.tipo when 1 then (ore.preco_venda * ore.quantidade)
                when 2 then ooi.quantidade * ooi.preco_venda
                when 3 then oap.qtde * oap.preco_venda
                  end as priceGrossValue
 ,  case c.tipo when 1 then (ore.preco_venda * ore.quantidade) - nvl(ore.valor_desconto_item,0)
                when 2 then (ooi.quantidade * ooi.preco_venda) - nvl(ooi.valor_desconto_item,0)
                when 3 then (oap.qtde * oap.preco_venda) - nvl(oap.valor_desconto_item,0)
                  end as priceNetValue
 ,  'Codigo da Peca: ' || i.cod_item as description
   from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
  left join os_requisicoes ore on o.cod_empresa = ore.cod_empresa and o.numero_os = ore.numero_os
  left join os_orcamentos_itens ooi on c.cod_empresa = ooi.cod_empresa and c.documento = ooi.numero_os and c.tipo = 2
  left join os_agenda_pecas oap on c.cod_empresa = oap.cod_empresa and c.documento = oap.cod_os_agenda and c.tipo = 3
  left join itens i on i.cod_item = case c.tipo when 1 then ore.cod_item
                                                when 2 then ooi.Cod_item
                                                else oap.cod_item end
  left join os_original ori on ori.cod_empresa = o.cod_empresa and ori.numero_os = o.numero_os and ori.item = case c.tipo when 1 then ore.item
                                                                                                                          when 2 then ooi.item
                                                                                                                            else null end
  left join os_agenda_reclamacao oar on oar.cod_empresa = oa.numero_os and oar.cod_os_agenda = oa.cod_os_agenda and oar.item = case c.tipo when 3 then oap.item
                                                                                                                            else null end
  left join reclamacao_pecas recp on recp.cod_reclamacao = case when c.tipo in (1,2) then ori.cod_reclamacao else oar.cod_reclamacao end
                                 and recp.cod_fornecedor = case c.tipo when 1 then ore.cod_fornecedor
                                                                       when 2 then ooi.cod_fornecedor
                                                                         else oap.cod_fornecedor end
                                 and recp.cod_item = case c.tipo when 1 then ore.cod_item
                                                                 when 2 then ooi.cod_item
                                                                  else oap.cod_item end
  left join reclamacoes_padroes rec on rec.cod_reclamacao = recp.cod_reclamacao
 where mei.id_mov_mb_xa = 1 --:ID_MOV
   and i.cod_item is not null
   and nvl(rec.eh_kit,'N') != 'S'`

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
func Test_Sql13(t *testing.T) {
	sql := `   select
    case c.tipo when 1 then ore.item
                when 2 then nvl(ooi.item,1)
                when 3 then oap.item
                  end as item
 ,  i.cod_item
 ,  case c.tipo when 1 then ore.cod_fornecedor
                when 2 then ooi.cod_fornecedor
                when 3 then oap.cod_fornecedor
                  end as cod_fornecedor
 ,  case c.tipo when 1 then ore.requisicao
                when 2 then 2
                when 3 then 3
                  end as requisicao
 ,  i.cod_item as codigo_peca
   from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os_requisicoes ore on c.cod_empresa = ore.cod_empresa and c.documento = ore.numero_os and c.tipo = 1
  left join os_orcamentos_itens ooi on c.cod_empresa = ooi.cod_empresa and c.documento = ooi.numero_os and c.tipo = 2
  left join os_agenda_pecas oap on c.cod_empresa = oap.cod_empresa and c.documento = oap.cod_os_agenda and c.tipo = 3
  left join itens i on i.cod_item = case c.tipo when 1 then ore.cod_item
                                                when 2 then ooi.Cod_Item
                                                else oap.cod_item end
 where mei.id_mov_mb_xa = 1 --:ID_MOV
   and i.cod_item is not null`

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
func Test_Sql14(t *testing.T) {
	sql := `   select
    case c.tipo when 3 then oar.item else ori.item end as item
  , rec.cod_reclamacao
  , rec.cod_reclamacao ||'#'|| rec.cod_kit_fabrica as numero
  , remove_caracteres_json(rec.reclamacao) as description
  , case c.tipo when 3 then oar.item else ori.item end ||'#'|| rec.cod_reclamacao as dmsExternalId
  , case c.tipo when 1 then o.tipo
                when 2 then o.tipo
                when 3 then coalesce(oa.tipo,oa.tipo_os)
                  end as invoiceCode

  , 'LOCAL' as dmsOrigin
  , 'DMS' as origin
  , 'DMS' as priceType
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
  left join os_original ori on c.cod_empresa = ori.cod_empresa and c.documento = ori.numero_os and c.tipo in (1,2)
  left join os_agenda_reclamacao oar on c.cod_empresa = oar.cod_empresa and c.documento = oar.cod_os_agenda and c.tipo = 3
  left join reclamacoes_padroes rec on rec.cod_reclamacao = case when c.tipo in (1,2) then ori.cod_reclamacao else oar.cod_reclamacao end
 where mei.id_mov_mb_xa = 1 --:ID_MOV
   and nvl(rec.cod_reclamacao,0) > 0
   and nvl(rec.eh_kit,'N') = 'S'
 order by 1`

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
func Test_Sql15(t *testing.T) {
	sql := ` select
    case c.tipo when 3 then oar.item else ori.item end as item
  , rec.cod_reclamacao
  , rece.valor_total as dmsGrossValue
  , rec.valor_total as dmsNetValue
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os_original ori on c.cod_empresa = ori.cod_empresa and c.documento = ori.numero_os and c.tipo in (1,2)
  left join os_agenda_reclamacao oar on c.cod_empresa = oar.cod_empresa and c.documento = oar.cod_os_agenda and c.tipo = 3
  left join reclamacoes_padroes rec on rec.cod_reclamacao = case when c.tipo in (1,2) then ori.cod_reclamacao else oar.cod_reclamacao end
  left join reclamacoes_padroes_empresa rece on rec.cod_reclamacao = rece.cod_reclamacao and rece.cod_empresa = c.cod_empresa
 where mei.id_mov_mb_xa = 1 --:ID_MOV
   and nvl(rec.cod_reclamacao,0) > 0
   and nvl(rec.eh_kit,'N') = 'S'`

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
func Test_Sql16(t *testing.T) {
	sql := `select
    case c.tipo when 1 then ore.item
                when 2 then nvl(ooi.item,1)
                when 3 then oap.item
                  end as item
 ,  recp.cod_reclamacao
 ,  i.cod_item
 ,  case c.tipo when 1 then ore.cod_fornecedor
                when 2 then ooi.cod_fornecedor
                when 3 then oap.cod_fornecedor
                  end as cod_fornecedor
 ,  case c.tipo when 1 then ore.requisicao
                when 2 then 2
                when 3 then 3
                  end as requisicao
 ,  'DMS' as origin
 ,  'LOCAL' as dmsOrigin
 ,  remove_caracteres_json(i.descricao) as name
 ,  case c.tipo when 1 then ore.item ||'#'||ore.cod_fornecedor
                when 2 then NVL(ooi.item,1) ||'#'||ooi.cod_fornecedor
                when 3 then oap.item ||'#'||oap.cod_fornecedor
                  end ||'#'|| i.cod_item as dmsExternalId
 ,  case c.tipo when 1 then o.tipo
                when 2 then o.tipo
                when 3 then coalesce(oa.tipo,oa.tipo_os)
                  end as invoiceCode
 ,  case c.tipo when 1 then ore.quantidade
                when 2 then ooi.quantidade
                when 3 then oap.qtde
                  end as quantity
 ,  case c.tipo when 1 then (ore.preco_venda * ore.quantidade)
                when 2 then ooi.quantidade * ooi.preco_venda
                when 3 then oap.qtde * oap.preco_venda
                  end as priceGrossValue
 ,  case c.tipo when 1 then (ore.preco_venda * ore.quantidade) - nvl(ore.valor_desconto_item,0)
                when 2 then (ooi.quantidade * ooi.preco_venda) - nvl(ooi.valor_desconto_item,0)
                when 3 then (oap.qtde * oap.preco_venda) - nvl(oap.valor_desconto_item,0)
                  end as priceNetValue
 ,  'Codigo da Peca: ' || i.cod_item as description
   from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3
  left join os_requisicoes ore on o.cod_empresa = ore.cod_empresa and o.numero_os = ore.numero_os
  left join os_orcamentos_itens ooi on c.cod_empresa = ooi.cod_empresa and c.documento = ooi.numero_os and c.tipo = 2
  left join os_agenda_pecas oap on c.cod_empresa = oap.cod_empresa and c.documento = oap.cod_os_agenda and c.tipo = 3
  left join itens i on i.cod_item = case c.tipo when 1 then ore.cod_item
                                                when 2 then ooi.Cod_item
                                                else oap.cod_item end
  left join os_original ori on ori.cod_empresa = o.cod_empresa and ori.numero_os = o.numero_os and ori.item = case c.tipo when 1 then ore.item
                                                                                                                          when 2 then ooi.item
                                                                                                                            else null end
  left join os_agenda_reclamacao oar on oar.cod_empresa = oa.numero_os and oar.cod_os_agenda = oa.cod_os_agenda and oar.item = case c.tipo when 3 then oap.item
                                                                                                                            else null end
  left join reclamacao_pecas recp on recp.cod_reclamacao = case when c.tipo in (1,2) then ori.cod_reclamacao else oar.cod_reclamacao end
                                 and recp.cod_fornecedor = case c.tipo when 1 then ore.cod_fornecedor
                                                                       when 2 then ooi.cod_fornecedor
                                                                         else oap.cod_fornecedor end
                                 and recp.cod_item = case c.tipo when 1 then ore.cod_item
                                                                 when 2 then ooi.cod_item
                                                                  else oap.cod_item end
  left join reclamacoes_padroes rec on rec.cod_reclamacao = recp.cod_reclamacao
 where mei.id_mov_mb_xa = 1 --:ID_MOV
   and i.cod_item is not null
   and recp.cod_item is not null
   and nvl(rec.eh_kit,'N') = 'S'`

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
func Test_Sql17(t *testing.T) {
	sql := `   select
    case c.tipo when 1 then ore.item
                when 2 then nvl(ooi.item,1)
                when 3 then oap.item
                  end as item
 ,  recp.cod_reclamacao
 ,  i.cod_item
 ,  case c.tipo when 1 then ore.cod_fornecedor
                when 2 then ooi.cod_fornecedor
                when 3 then oap.cod_fornecedor
                  end as cod_fornecedor
 ,  case c.tipo when 1 then ore.requisicao
                when 2 then 2
                when 3 then 3
                  end as requisicao
 ,  i.cod_item as codigo_peca
   from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os_requisicoes ore on c.cod_empresa = ore.cod_empresa and c.documento = ore.numero_os and c.tipo = 1
  left join os_orcamentos_itens ooi on c.cod_empresa = ooi.cod_empresa and c.documento = ooi.numero_os and c.tipo = 2
  left join os_agenda_pecas oap on c.cod_empresa = oap.cod_empresa and c.documento = oap.cod_os_agenda and c.tipo = 3
  left join itens i on i.cod_item = case c.tipo when 1 then ore.cod_item
                                                when 2 then ooi.cod_item
                                                else oap.cod_item end
  left join os_original ori on ori.cod_empresa = c.cod_empresa and ori.numero_os = c.documento and ori.item = case c.tipo when 1 then ore.item
                                                                                                                          when 2 then ooi.item
                                                                                                                            else null end
  left join os_agenda_reclamacao oar on oar.cod_empresa = c.cod_empresa and oar.cod_os_agenda = c.documento and oar.item = case c.tipo when 3 then oap.item
                                                                                                                            else null end
  left join reclamacao_pecas recp on recp.cod_reclamacao = case when c.tipo in (1,2) then ori.cod_reclamacao else oar.cod_reclamacao end
                                 and recp.cod_fornecedor = case c.tipo when 1 then ore.cod_fornecedor
                                                                       when 2 then ooi.cod_fornecedor
                                                                         else oap.cod_fornecedor end
                                 and recp.cod_item = case c.tipo when 1 then ore.cod_item
                                                                 when 2 then ooi.cod_item
                                                                  else oap.cod_item end
  left join reclamacoes_padroes rec on rec.cod_reclamacao = recp.cod_reclamacao
 where mei.id_mov_mb_xa = 1 --:ID_MOV
   and i.cod_item is not null
   and recp.cod_item is not null
   and nvl(rec.eh_kit,'N') = 'S'`

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
func Test_Sql18(t *testing.T) {
	sql := `   select
    case c.tipo when 1 then oss.item
                when 2 then oso.item
                when 3 then oas.item
                  end item
 ,  recs.cod_reclamacao
 ,  'DMS' as origin
 ,  s.cod_servico as operationNumber
 ,  nvl(case c.tipo when 1 then oss.tempo_padrao
                when 2 then oso.tempo
                when 3 then oas.tempo_padrao
                 end,s.tempo_padrao) as operationValue
 ,  'H' as operationUnit
 ,  case c.tipo when 1 then oss.item
                when 2 then oso.item
                when 3 then oas.item
                  end ||'#'|| s.cod_servico as dmsExternalId
 ,  case c.tipo when 1 then o.tipo
                when 2 then o.tipo
                when 3 then coalesce(oa.tipo,oa.tipo_os)
                  end as invoiceCode
 ,  remove_caracteres_json(s.descricao_servico) as descricao
 ,  case c.tipo when 1 then oss.preco_venda
                when 2 then oso.preco_venda
                when 3 then oas.preco_venda
                  end as priceGrossValue
 ,  case c.tipo when 1 then oss.preco_venda - nvl(oss.desconto_por_serv,0)
                when 2 then oso.preco_venda - nvl(oso.desconto_por_serv,0)
                when 3 then oas.preco_venda - nvl(oas.valor_desconto_serv,0)
                  end as priceNetValue
  from fab_mov_ei_cria_mb_xa mei
  left join fab_ei_cria_mb_xa ei on ei.id = mei.id_ei_cria_mb_xa
  left join xentry_controle_envio c on c.id = ei.id_controle
  left join os o on c.cod_empresa = o.cod_empresa and c.documento = o.numero_os and c.tipo in (1,2)
  left join os_agenda oa on c.cod_empresa = oa.cod_empresa and c.documento = oa.cod_os_agenda and c.tipo = 3

  left join os_servicos oss on c.cod_empresa = oss.cod_empresa and c.documento = oss.numero_os and c.tipo = 1
  left join os_serv_orc oso on c.cod_empresa = oso.cod_empresa and c.documento = oso.numero_os and c.tipo = 2
  left join os_agenda_servicos oas on c.cod_empresa = oas.cod_empresa and c.documento = oas.cod_os_agenda and c.tipo = 3
  left join servicos s on s.cod_servico = case c.tipo when 1 then oss.cod_servico
                                                      when 2 then oso.cod_servico
                                                        else oas.cod_servico end
  left join os_original ori on ori.cod_empresa = c.cod_empresa and ori.numero_os = c.documento and ori.item = case c.tipo when 1 then oss.item
                                                                                                                          when 2 then oso.item
                                                                                                                            else null end
  left join os_agenda_reclamacao oar on oar.cod_empresa = c.cod_empresa and oar.cod_os_agenda = c.documento and oar.item = case c.tipo when 3 then oar.item
                                                                                                                            else null end
  left join reclamacao_servicos recs on recs.cod_reclamacao = case when c.tipo in (1,2) then ori.cod_reclamacao else oar.cod_reclamacao end
                                    and recs.cod_servico = case c.tipo when 1 then oss.cod_servico
                                                                       when 2 then oso.cod_servico
                                                                         else oas.cod_servico end
  left join reclamacoes_padroes rec on rec.cod_reclamacao = recs.cod_reclamacao
 where mei.id_mov_mb_xa = 1 --:ID_MOV
   and s.cod_servico is not null
   and recs.cod_servico is not null
   and nvl(rec.eh_kit,'N') = 'S'`

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
