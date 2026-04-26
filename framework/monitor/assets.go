package monitor

import "embed"

//go:embed assets/*
var assetsFS embed.FS

// loginErrorPage é renderizado em uma falha de login. Recebe (msg, user)
// via fmt.Fprintf — manter formato compatível com escapeHTML do chamador.
const loginErrorPage = `<!doctype html>
<html lang="pt-br"><head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Login</title><link rel="stylesheet" href="/static/style.css">
</head><body class="login-bg"><form class="login-card" method="post" action="/login">
<h1>Monitor</h1>
<div class="login-error">%s</div>
<label>Usuário<input name="user" autofocus value="%s"></label>
<label>Senha<input name="pass" type="password"></label>
<button type="submit">Entrar</button>
</form></body></html>`
