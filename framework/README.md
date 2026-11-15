# framework — guia de uso

Bootstrap de aplicações Go com **pattern de Options**: um único executável que pode subir
**API HTTP**, **Bots** (ciclos periódicos), **Srv** (serviço de 1 min) e **Monitor**, e que pode
ser **instalado como serviço do sistema operacional** (Windows Service, systemd, launchd) via
[`kardianos/service`](https://github.com/kardianos/service).

> Pacote Go: `github.com/aoticombr/golang/framework` — o nome do pacote é **`core`**.
> Importe com alias: `core "github.com/aoticombr/golang/framework"`.

Exemplos reais: `Embravor-Admin/back.go` (API + JWT), `Aoti-Site-Admin/back.go/bot` (bots),
`Embravor-Studio-Backup` (API + bot agendador, sem banco de tenant).

---

## 1. Esqueleto mínimo

```go
package main

import (
    core "github.com/aoticombr/golang/framework"

    _ "meuapp/api" // init() registra as rotas
    _ "meuapp/bot" // init() registra os bots
)

func main() {
    // NewAppDev(nome) → o nome é a chave usada no config.json (apis[].name, bots[].name)
    // NewApp()        → usa o nome do executável
    core.NewAppDev("meuapp").Execute(
        core.AddApi(),
        core.AddBot(),
        // core.AddSrv(), core.AddMonitor()
    )
}
```

`Execute` faz:
1. No Windows, se o diretório atual contém `SYSTEM32` (rodando como serviço), muda para a pasta do exe.
2. Cria o `service.Service` (kardianos).
3. Se houver verbo na linha de comando (`install`, `uninstall`, `start`, `stop`, `restart`, `status`) executa
   `service.Control` e sai (ver seção 6).
4. Senão chama `s.Run()` → `Start` → `Run`: lê `config.json`, configura o log e executa cada Option
   em goroutine.
5. No `Stop` (Ctrl+C ou parada do serviço) cancela `app.Ctx`, faz shutdown gracioso da API
   (`timeouts.shutdown_timeout`, padrão 10s) e do Monitor (5s).

Singletons expostos: `core.CoreApi`, `core.CoreBot`, `core.CoreSrv`, `core.CoreMonitor`.

---

## 2. config.json

Lido da **pasta atual** (como serviço: pasta do exe). Se não existir, um padrão é gravado.
Todas as buscas por nome são case-insensitive e exigem `"ativo": true` (exceto `jwt` e `parans`).

```jsonc
{
  "dbs": [                       // conexões (PG ou ORA) usadas como tenants da API
    { "name": "EMB", "db": "PG", "host": "localhost", "port": 5432, "user": "u", "pass": "p",
      "schema": "public", "sid": "", "poolsize": 30, "maxconn": 100, "ativo": true,
      "trace": { "path": "tracepgx", "ativo": false } }
  ],
  "apis": [
    { "name": "meuapp",           // == nome do app
      "host": "127.0.0.1", "port": 4101, "ativo": true,
      "dbs": ["EMB"],             // cada nome vira um tenant /{db}; pode ficar vazio
      "cors":  { "ativo": true, "allowOrigins": ["*"], "allowMethods": ["GET","POST","PUT","DELETE","OPTIONS","PATCH"],
                 "allowHeaders": ["Authorization","Content-Type"], "exposedHeaders": [], "allowCredentials": false, "maxAge": 300 },
      "https": { "ativo": false, "cert": "", "key": "" },
      "timeouts": { "read_header_timeout": 30, "read_timeout": 120, "write_timeout": 120,
                    "idle_timeout": 120, "shutdown_timeout": 30 } }
  ],
  "bots":     [ { "name": "meuapp", "ativo": true, "dbs": [] } ],   // liga o runner de bots
  "services": [ { "name": "meuapp", "ativo": true } ],              // liga o runner Srv
  "jwt":      [ { "name": "admin", "expirationTime": 30, "secret": "64-hex..." } ],
  "certs":    [ { "name": "pgp", "ativo": true, "path_crt": "...", "path_private": "...", "pass": "..." } ],
  "parans":   [],
  "log":      { "nivel": 2, "screen": true },
  "monitor":  { "ativo": false }
}
```

Acesso no código: `config.NewConfig().GetApi(nome)`, `GetDB`, `GetBot`, `GetService`, `GetJwt`,
`GetCert`, `Parans.GetParam(marcador, nome)`; `Database.GetDsn()` monta a URL do driver.

---

## 3. Options

| Option | Liga quando | O que faz |
|---|---|---|
| `AddApi()` | `apis[name==app].ativo` | Abre uma conexão por nome em `apis[].dbs` e registra como tenant; sobe o router chi com todos os `RouteGroup`s registrados. |
| `AddBot()` | `bots[name==app].ativo` | Para cada bot registrado, uma goroutine que executa `Executar()` imediatamente e depois a cada N minutos. |
| `AddSrv()` | `services[name==app].ativo` | Igual ao bot, com ciclo fixo de 1 minuto (`srv.RegistrarSrv`). **Bloqueante.** |
| `AddMonitor()` | `monitor.ativo` | Sobe a UI web de monitoramento (login próprio). |
| `AddApiPrefixo(p)` | `apis[name==app].ativo` | API com prefixo fixo, sem exigir tenant ([seção 8](#8-api-sem-tenant--addapiprefixo)). |
| `AddServiceConfig(o)` | sempre | Personaliza o registro no gerenciador de serviços ([seção 6](#personalizar-o-serviço--addserviceconfig)). |

Criar uma Option própria: qualquer `func(*core.App)` serve — `type OptionsApp func(*App)`.

```go
func AddMinhaCoisa() core.OptionsApp {
    return func(app *core.App) {
        go func() { <-app.Ctx.Done(); /* liberar recursos */ }()
    }
}
```

---

## 4. API HTTP

### Registrar rotas

```go
package api

import (
    "net/http"
    fapi "github.com/aoticombr/golang/framework/api"
    mw   "github.com/aoticombr/golang/framework/api/middlewares"
    "github.com/go-chi/jwtauth"
)

var tokenAdmin = jwtauth.New("HS256", []byte(segredo), nil)

func init() {
    pub := fapi.RegisterRouterGroup("Publico")
    pub.RegisterRoute("/health", fapi.GET, health)            // GET /{db}/health

    priv := fapi.RegisterRouterGroup("Privado")
    priv.Use(mw.StaticVerifier(tokenAdmin))                   // JWT obrigatório
    priv.RegisterRoute("/api/itens", fapi.GET, listar)
    priv.RegisterRoute("/api/itens/{id}", fapi.PUT, salvar)
    priv.RegisterRoute("/ws/eventos", fapi.WEBSOCKET, ws)     // HandleFunc puro
}

func health(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) }
```

Métodos: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `WEBSOCKET`. Parâmetros: `chi.URLParam(r, "id")`.
Erros: `fapi.ErrorHandler(w, r, status, body, err)`.

### O prefixo `/{db}`

Toda rota recebe `/{db}` na frente. Ele existe para **multi-tenant**: o handler resolve a conexão com

```go
ts := core.CoreApi.Processo.Tenants.FindByName(chi.URLParam(r, "db"))
ds := ts.Connection.NewDataSetName("X"); defer ds.Free()
```

**Se a aplicação não usa banco de tenant**, deixe `apis[].dbs` vazio e use qualquer nome fixo
no primeiro segmento — por convenção o nome do produto. Ex.: o Embravor Studio Backup usa
`/backup/api/projetos`. O segmento é aceito por qualquer valor; valide-o no middleware se quiser
restringir.

### Middlewares prontos (`framework/api/middlewares`)
- `StaticVerifier(auth)` / `DynamicVerifier(func() *jwtauth.JWTAuth)` — JWT obrigatório.
- `RequireAuth`, `OptionalAuth` — variações.
- `api.ratelimit.go` (memória) e `api.ratelimit.redis.go`.

---

## 5. Bots (ciclos periódicos)

```go
package bot

import (
    "context"
    "time"

    fbot "github.com/aoticombr/golang/framework/bot"
)

type Limpeza struct{}

// método no PONTEIRO; o registro recebe um VALOR (reflect.New cria a instância)
func (b *Limpeza) Executar() { /* trabalho do ciclo */ }

// Bot que recebe o contexto da aplicação (cancelado quando o serviço para)
type Agendador struct{}

func (a *Agendador) ExecutarCtx(ctx context.Context) { /* ... */ }

func init() {
    fbot.RegistrarBot("Limpeza", Limpeza{}, 5)                         // a cada 5 minutos
    fbot.RegistrarBotIntervalo("Agendador", Agendador{}, 30*time.Second) // intervalo livre
}
```

Requisitos: `core.AddBot()` no `Execute`, import em branco do pacote e `"bots": [{"name":"<app>","ativo":true}]`.

Comportamento:
- Executa **imediatamente** ao subir, depois a cada intervalo (`time.Ticker`).
- Um bot pode implementar `Executar()` (clássico) **ou** `ExecutarCtx(ctx)` (`IControllerBotCtx`); se tiver os dois, `ExecutarCtx` é usado. O `ctx` é o `app.Ctx`, cancelado na parada do serviço.
- Ciclos são **sequenciais**: o próximo só começa quando o atual termina (sem sobreposição). Se um ciclo demorar mais que o intervalo, o próximo começa logo após.
- Um **panic** dentro do ciclo é registrado no log (com stack) e o bot continua no próximo ciclo.
- `RegistrarBot(key, clazz, minutos)` — intervalo em minutos (0 vira 1 minuto).
- `RegistrarBotIntervalo(key, clazz, time.Duration)` — qualquer intervalo.
- Struct sem `Executar`/`ExecutarCtx` no ponteiro: `AddBot` registra erro e não sobe os bots.
- Para agendas de calendário (dia/hora), use um bot de ciclo curto que avalia o que está vencido e despacha em goroutines (padrão do Embravor Studio Backup).

---

## 6. Instalar como serviço do SO

O mesmo executável funciona em console e como serviço:

```bash
meuapp            # roda em console (Ctrl+C para parar)
meuapp install    # instala o serviço (nome = app.Name)
meuapp start | stop | restart
meuapp status     # imprime: executando | parado | desconhecido
meuapp uninstall
```

Regras:
- O verbo precisa ser o **primeiro** argumento (`os.Args[1]`); qualquer outro argumento (ex.: `-migrate-only`) roda o app normalmente.
- Verbos só são tratados quando o processo é iniciado **por um usuário no terminal** (`service.Interactive()`). Quando o gerenciador de serviços inicia o exe, ele sempre roda — units antigas com `ExecStart=/opt/app/app start` continuam funcionando.
- Execuções de teste (`-test.*`) nunca são tratadas como verbo.

| SO | Backend | Observações |
|---|---|---|
| Windows | Service Control Manager | Executar como Administrador. Pasta de trabalho = pasta do exe (ajuste automático). |
| Linux | systemd (ou SysV/upstart) | Executar como root. Unit gerada em `/etc/systemd/system/<nome>.service`. |
| macOS | launchd | `sudo`. Plist em `/Library/LaunchDaemons`. |

### Personalizar o serviço — `AddServiceConfig`

```go
core.NewAppDev("embravor-backup").Execute(
    core.AddServiceConfig(core.ServiceOptions{
        DisplayName:  "Embravor Studio Backup",
        Description:  "Agendador e executor de backups",
        Dependencias: []string{"After=network-online.target"}, // Linux; no Windows: nomes de serviços
        Usuario:      "",       // vazio = LocalSystem / root
        Restart:      "always", // systemd
        InicioManual: false,    // Windows: true = início manual
    }),
    core.AddApi(),
)
// ou encadeado: core.NewAppDev("x").ServiceConfig(core.ServiceOptions{...}).Execute(...)
```

Sem a option, `DisplayName`/`Description` continuam `"<exe> Go lang application AOTI"`.

---

## 7. Log

`lib.NewLog()` (singleton): `.Debug/.Info/.Warning/.Error/.Critical/.Fatal(app, v...)` e `.Screen`.
Arquivo `log_<App>_<dd-mm-yyyy-HH>.log` na pasta atual, rotação por hora. Nível em `config.log.nivel`
(0 nenhum … 5 crítico); `log.screen` também imprime no console.

---

## 8. API sem tenant — `AddApiPrefixo`

Para apps que não usam banco de tenant (ex.: o catálogo é SQLite próprio):

```go
core.NewAppDev("embravor-backup").Execute(
    core.AddApiPrefixo("backup"), // rotas em /backup/... ; "" = sem prefixo
)
```

Diferenças em relação a `AddApi()` (que continua igual):
- Prefixo fixo no lugar de `/{db}`; `apis[].dbs` pode ficar vazio (se houver, continuam abertos como tenants).
- Escuta em `host:port` do config (`host` vazio = todas as interfaces; `127.0.0.1` = só local).
- Timeouts do `http.Server` aplicados **somente quando informados** (> 0) no config — com 0, não há limite de leitura/escrita (WebSocket e downloads longos funcionam).
- `middleware.Recoverer`: panic em handler vira HTTP 500 em vez de derrubar a conexão.
- `api.RegistrarHandlerRaiz("/*", handler)` monta um handler **fora** do prefixo (ex.: SPA/UI web em `/`).
- `app.OnApiErro = func(err error){...}` é chamado se o servidor terminar com erro (ex.: porta ocupada).
- Rota `api.WEBSOCKET` registra o handler para **todos os métodos** — use também para handlers que tratam GET/POST/DELETE juntos (ex.: endpoint MCP streamable).

---

## 9. Histórico de novidades

| Versão | Novidade |
|---|---|
| v1.0.196 | Correção de `install/uninstall/start/stop/restart` + novo verbo `status`; `AddServiceConfig`/`App.ServiceConfig`; `AddApiPrefixo` + `RegistrarHandlerRaiz` + `App.OnApiErro`; bots com `IControllerBotCtx`, `RegistrarBotIntervalo`, `recover` por ciclo, `Tempo 0` não causa mais panic. |

Em avaliação: `AddAgendador(...)` genérico de calendário (diário hh:mm, dias da semana, mensal), se o modelo do Embravor Studio Backup se provar reutilizável.
