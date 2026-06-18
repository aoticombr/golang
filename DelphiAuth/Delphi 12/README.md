# DelphiAuth — captura do Authorization Code (OAuth2)

App VCL (Delphi 11/12) que abre o login da Microsoft numa **janela embutida**
(TEdgeBrowser / WebView2), no estilo do Postman, intercepta o redirect para o
`CallBackUrl` e **captura o `authorization code`**.

O objetivo é apenas **capturar e exibir o `code`** — a troca por `access_token` é feita
pela API (componente `http` em Go, grant `GT_AuthorizationCode`), que consome o `code`
diretamente.

## Como funciona

1. O app monta a URL de autorização a partir dos campos (`AuthUrl`, `ClientId`,
   `CallBackUrl`, `Scope`, `State`) e navega o `TEdgeBrowser` para ela:

   ```
   {AuthUrl}?client_id=...&response_type=code&redirect_uri={CallBackUrl}&response_mode=query&scope=...&state=...
   ```

2. Você autentica com a Microsoft dentro da janela.
3. O Azure AD redireciona para `{CallBackUrl}?code=...&state=...`. No evento
   `OnNavigationStarting`, o app detecta que a navegação é para o `CallBackUrl`,
   **cancela** (a página de callback nem precisa existir) e extrai o `code` da URL.
4. O `code` aparece no campo "Code capturado" e pode ser copiado (botão **Copiar**).

> O `code` é de **uso único** e expira em poucos minutos. Entregue-o à API logo em
> seguida para a troca pelo token.

## Pré-requisitos

- **Delphi 11 (Alexandria) / 12 (Athens)** — `TEdgeBrowser` disponível.
- **Microsoft Edge WebView2 Runtime** instalado (já vem por padrão no Windows 11).
- **WebView2Loader.dll** ao lado do executável (distribuído com o Delphi, em
  `...\Embarcadero\Studio\<versão>\Redist\...\WebView2\`). Ao compilar pela IDE ele
  costuma ser copiado automaticamente; se faltar em runtime, copie-o para a pasta do `.exe`.
- No **Azure AD**, o `CallBackUrl` precisa estar registrado como *redirect URI* do app
  (tipo **Web**), exatamente igual ao informado aqui.

## Como compilar/rodar

1. Abra `DelphiAuth.dpr` na IDE do Delphi.
2. Compile e execute (F9).
3. Os campos já vêm preenchidos com os valores de exemplo (os mesmos do teste Go).
   Ajuste se necessário e clique em **Abrir login**.

## Arquivos

- `DelphiAuth.dpr` — projeto.
- `uMainForm.pas` / `uMainForm.dfm` — formulário principal (browser + captura do code).

## Relação com o componente Go

Os valores de exemplo correspondem ao teste `TestAuth2_AuthorizationCode_*` em
`testes/http/comp_test.go`. Depois de copiar o `code` aqui, use-o na API:

```go
h.Auth2.GrantType = http.GT_AuthorizationCode
h.Auth2.AccessTokenUrl = "https://login.microsoftonline.com/<tenant>/oauth2/v2.0/token"
h.Auth2.CallBackUrl    = "https://mockserver.aoti.com.br/cb" // mesmo redirect_uri
h.Auth2.ClientId       = "..."
h.Auth2.ClientSecret   = "..."
h.Auth2.Code           = "<code-capturado-no-DelphiAuth>"
token, err := h.Auth2.GetToken()
```
