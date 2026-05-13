# github.com/aoticombr/golang

Coleção de pacotes Go para construção de aplicações de backend: acesso a banco de dados (Oracle, PostgreSQL, MySQL, SQL Server, Firebird, SQLite), cliente HTTP/WebSocket, OAuth2, SMTP, criptografia PGP, utilitários de configuração, logging e mais.

## Instalação

```bash
go get github.com/aoticombr/golang
```

E importe individualmente os pacotes que precisar:

```go
import (
    "github.com/aoticombr/golang/config"
    "github.com/aoticombr/golang/dbconndataset"
    "github.com/aoticombr/golang/http"
)
```

## Visão geral dos pacotes

| Pacote | Função |
|---|---|
| [config](#config) | Carregamento e gerenciamento de configuração via JSON |
| [dbconnbase](#dbconnbase) | Conexão e transação com banco (multi-dialeto) |
| [dbconndataset](#dbconndataset) | Wrapper conveniente combinando `dbconnbase` + `dbdataset` |
| [dbdataset](#dbdataset) | DataSet rico: query, parâmetros, mapeamento para struct |
| [orm](#orm) | Geração de SQL INSERT/UPDATE/DELETE a partir de tags em struct |
| [preprocesssql](#preprocesssql) | Pré-processador SQL: extrai parâmetros/macros, trata escapes |
| [http](#http) | Cliente HTTP/WebSocket com OAuth2, proxy, TLS, multipart |
| [mail](#mail) | Envio de e-mail SMTP com anexos, MIME e TLS |
| [pgp](#pgp) | Criptografia PGP (OpenPGP) com chave + passphrase |
| [m3u8](#m3u8) | Download de playlists HLS (M3U8) |
| [restmonitorserver](#restmonitorserver) | Servidor HTTP de debug para inspecionar requisições + JWT |
| [mail](#mail), [file](#file), [log](#log) | Utilitários de I/O e logging |
| [variant](#variant) | Tipo genérico com conversões seguras (`AsInt`, `AsString`, etc.) |
| [stringlist](#stringlist) | Lista de strings com delimitador configurável |
| [jsonconfig](#jsonconfig) | Carregamento de configuração JSON estruturada |
| [lib](#lib) | Utilitários gerais (arquivos, crypto, strings, GUIDs) |
| [framework](#framework) | Bootstrap de aplicação com pattern de Options |
| [memoryleak](#memoryleak) | Servidor pprof para profiling de memória |

---

## config

Carregamento, manipulação e persistência de configurações de aplicação a partir de JSON. Inclui suporte a múltiplos bancos, APIs, serviços, bots, JWT, e flag de tracing por banco.

**Tipos principais:** `Config`, `Database`, `Trace`, `Api`, `Service`, `Bot`

```go
cfg := config.NewConfig()

found, db := cfg.GetDB("principal")
if !found {
    cfg.AddDb(config.Database{
        Name: "principal",
        Db:   "PG",
        Host: "localhost",
        Port: 5432,
        User: "postgres",
        Pass: "secret",
        Schema: "app",
        Trace: config.Trace{Ativo: true, Path: "C:/logs/pgx"},
    })
    cfg.Save()
}
```

---

## dbconnbase

Camada base de conexão SQL multi-dialeto sobre `database/sql`. Faz pool, transação, ping, e suporte a query tracing nativo do `pgx` (Postgres) via diretório de log — espelhando o comportamento de `TRACE DIR` do go-ora (Oracle).

**Tipos principais:** `Conn`, `Transaction`, `DialectType`

**Dialetos suportados:** Oracle, PostgreSQL, MySQL, SQL Server, Firebird, SQLite, Interbase.

```go
conn, err := dbconnbase.NewConn(dbConfig)
if err != nil { log.Fatal(err) }
defer conn.Close()

tx, _ := conn.StartTransaction()
defer tx.Rollback()
// ... usa tx
tx.Commit()
```

**Drivers internos:** `github.com/jackc/pgx/v5/stdlib`, `github.com/sijms/go-ora/v2`.

---

## dbconndataset

Combina `Conn` com `DataSet` para o caso comum de "abrir conexão e executar queries". Apenas adiciona açúcar sintático.

**Tipos principais:** `ConnDataSet`

```go
cds, err := dbconndataset.NewConn(dbConfig)
defer cds.Close()

ds := cds.NewDataSetName("users")
ds.AddSql("SELECT * FROM users WHERE status = :status")
ds.SetInputParam("status", "active")
ds.Open()
```

---

## dbdataset

DataSet rico com parâmetros (`:nome`), macros, cursor, mapeamento para struct, suporte a CLOB/BLOB e BindMode posicional/nominal. Usa `preprocesssql` para extrair os parâmetros da SQL e `replaceParamPG` para converter `:name` em `$N` quando o dialeto é Postgres.

**Tipos principais:** `DataSet`, `Field`, `Fields`, `Param`, `Params`, `Macro`

**Métodos chave:** `Open`, `Exec`, `SetInputParam`, `SetInputDateISO8601`, `FieldByName`, `First`, `Next`, `Eof`, `ToStruct`.

```go
ds := dbdataset.NewDataSet(conn)
ds.AddSql(`
    update t set payload = :payload::jsonb
     where id = :id
`)
ds.SetInputParam("payload", `{"foo":"bar"}`)
ds.SetInputParam("id", 42)
if err := ds.Exec(); err != nil { log.Fatal(err) }
```

---

## orm

Gera SQL `INSERT`, `UPDATE` e `DELETE` a partir de tags `column` em structs. Suporta primary key, required, auto-guid, timestamps automáticos, transformações (`md5`, `upper`, `lower`), `omitempty` e `nullempty`.

```go
type User struct {
    ID    *string `column:"id,#primarykey,#autoguid"`
    Name  string  `column:"name,#insert,#update,#required"`
    Email string  `column:"email,#insert,#update"`
}

u := &User{Name: "João", Email: "joao@x.com"}
tb := orm.NewTable(u)
sqlIns, _ := tb.SqlInsert()
// insert into user (id, name, email) values (uuid_generate_v4(), :name, :email)
```

`isEmpty` detecta corretamente: `*string` nil, `*string` apontando para `""`, slices/maps vazios, e interfaces segurando typed-nil.

---

## preprocesssql

Pré-processador de SQL portado de FireDAC. Extrai parâmetros nomeados (`:nome`), macros (`!nome`, `&nome`), escapes (`{d 'data'}`, `{fn func()}`) e respeita strings literais, comentários e casts (`::jsonb` no Postgres).

```go
sql := "select * from users where id = :id and status = :status"
params, _, _, err := preprocesssql.PreprocessSQL(sql, true, true, true, true, true)
for _, p := range params.Items {
    fmt.Println("param:", p.Name) // id, status
}
```

---

## http

Cliente HTTP/HTTPS e WebSocket multi-uso com:

- Vários `EncType`: `RAW`, `BINARY`, `FORM_DATA` (multipart), `X_WWW_FORM_URLENCODED`.
- Autenticação Basic, Bearer, OAuth2 client_credentials (auto-renovação de token).
- Proxy com escape correto de credenciais.
- Certificados TLS (PEM e PFX), `InsecureSkipVerify`.
- Trace e logging via `OnSend` (interface `IWebsocket`).
- WebSocket com reconexão automática e thread-safety (`sync.RWMutex`).

**Tipos principais:** `THttp`, `Request`, `Response`, `WebSocket`, `Auth2`, `TCert`

**Métodos:** `Send`, `Conectar`, `Desconectar`, `EnviarTexto`, `EnviarBinario`, `SetUrl`, `SetMetodo`, `Free`.

```go
h := http.NewHttp()
h.SetUrl("https://api.exemplo.com/recurso")
h.SetMetodo(http.M_POST)
h.AuthorizationType = http.AT_Auth2
h.Auth2.AuthUrl = "https://idp/oauth/token"
h.Auth2.ClientId = "..."
h.Auth2.ClientSecret = "..."

resp, err := h.Send()
if err != nil { log.Fatal(err) }
fmt.Println(resp.StatusCode, string(resp.Body))
```

WebSocket:

```go
h := http.NewHttp()
h.SetUrl("wss://exemplo.com/ws")
h.OnSend = meuHandler // implementa http.IWebsocket
h.Conectar()
h.EnviarTextTypeTextMessage([]byte(`{"hello":"world"}`))
```

**Deps:** `github.com/gorilla/websocket`, `golang.org/x/crypto/pkcs12`.

---

## mail

Composição de e-mails MIME e envio via SMTP, com suporte a HTML, anexos, embeds (imagens inline) e TLS.

```go
m := mail.NewMessage()
m.SetHeader("From", "remetente@exemplo.com")
m.SetHeader("To", "destino@exemplo.com")
m.SetHeader("Subject", "Olá")
m.SetBody("text/html", "<h1>Olá mundo</h1>")
m.Attach("relatorio.pdf")

d := mail.NewDialer("smtp.gmail.com", 587, "user", "pass")
if err := d.DialAndSend(m); err != nil { log.Fatal(err) }
```

---

## pgp

Encriptação e decriptação PGP/OpenPGP usando chave + passphrase. Funções utilitárias para conteúdo curto codificado em base64 e para arquivos.

```go
chaveSecreta := os.Getenv("PGP_PRIVATE_KEY")
texto := pgp.DecodePGPKeyPass(cifrado_base64, chaveSecreta, "senha")
```

**Deps:** `golang.org/x/crypto/openpgp`.

---

## m3u8

Download de vídeo segmentado HLS — recebe a URL de uma playlist `.m3u8`, baixa todos os segments e concatena em arquivo único.

```go
m := m3u8.NewM3u8()
bytes, err := m.GetVideoByte("https://exemplo.com/play.m3u8")
m.SaveByteToFile("video.mp4")
```

---

## restmonitorserver

Servidor HTTP standalone para debug — loga toda requisição recebida (corpo + headers + URL) em `./data/<timestamp>.json` e emite tokens JWT em `POST /token`. Útil pra inspecionar callbacks de OAuth2 e webhooks.

```bash
go run ./restmonitorserver
# Escuta em http://localhost:3003
```

**Deps:** `github.com/dgrijalva/jwt-go`, `github.com/google/uuid`.

---

## file

Logger simples por arquivo: escreve `.log` diário em pasta padrão com timestamp, opcionalmente espelhando em stdout.

```go
log := file.GetLog()
log.SendMsg("aplicação iniciada")
log.SendErro("falha ao conectar", err)
```

---

## log

Logger com níveis (`DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL`) e writer plugável.

```go
lg, _ := log.NewLogger("INFO", os.Stderr, "[APP] ", "")
lg.Info("servidor pronto")
lg.Warning("memória baixa")
lg.Error("erro de conexão")
```

---

## variant

Tipo genérico (`Variant`) que aceita qualquer valor (`Value any`) e expõe conversões seguras — útil pra mapear célula de banco sem se preocupar com o tipo original.

```go
v := &variant.Variant{Value: "42"}
n := v.AsInt()        // 42
s := v.AsString()     // "42"
ok := v.IsNull()      // false
```

**Deps:** `github.com/araddon/dateparse`, `github.com/sijms/go-ora/v2`.

---

## stringlist

Lista de strings com delimitador configurável (`,`, `;`, `\n`, etc.) — útil pra construir IN-lists, payloads, etc.

```go
sl := stringlist.NewStrings(stringlist.WithDelimiter(","))
sl.Append("a").Append("b").Append("c")
fmt.Println(sl.Text()) // "a,b,c"
```

---

## jsonconfig

Outro carregador de configuração via JSON — focado em estruturas de boot/services/APIs/schemas. (Coexiste com `config`; verifique qual atende seu caso.)

```go
cfg := jsonconfig.NewJsonConfig()
cfg.Name = "config.json"
cfg.Load()

api := cfg.GetConfig().GetApiByName("checkout")
```

---

## lib

Bag de utilitários: existência de arquivo, leitura/escrita JSON, GUIDs, validação de e-mail, CPF, geração de string aleatória, encrypt/decrypt, gzip, remoção de acentos, etc.

```go
exists := lib.FileExists("config.json")
lib.ReadJsonFileToStruct("config.json", &cfg)

id := lib.GuidString()
ok := lib.IsValidEmail("foo@bar.com")
plain := lib.RemoveAcentos("João José")
```

**Deps:** `github.com/google/uuid`, `github.com/kardianos/service`.

---

## framework

Bootstrap de aplicação com pattern Options — registra APIs, bots, monitor e serviços de uma vez. Integra com `github.com/kardianos/service` para rodar como serviço Windows/systemd.

```go
app := framework.NewApp()
app.Execute(
    framework.AddApi(),
    framework.AddMonitor(),
)
```

---

## memoryleak

Sobe servidor pprof em `:6060` pra investigação de memory leak.

```go
import "github.com/aoticombr/golang/memoryleak"

memoryleak.MemoryLeak()
// → http://localhost:6060/debug/pprof/heap
```

---

## Compatibilidade

- **Go:** 1.22+
- **Postgres driver:** `pgx/v5` (via `database/sql` stdlib adapter)
- **Oracle driver:** `go-ora/v2`
- **OS:** Windows / Linux / macOS

## Licença

Veja [LICENCE](./LICENCE).

## Contribuindo

Veja [CONTRIBUTING.MD](./CONTRIBUTING.MD) e [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md).
