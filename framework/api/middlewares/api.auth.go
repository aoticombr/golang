// Package middlewares fornece middlewares de autenticação e autorização para APIs
// Focado em JWT (JSON Web Tokens) usando a biblioteca go-chi/jwtauth
package middlewares

import (
	"net/http"

	// go-chi/jwtauth: biblioteca para autenticação JWT no framework Chi
	// Fornece funções para verificar, validar e extrair informações de tokens JWT
	"github.com/go-chi/jwtauth"
)

// DynamicVerifier cria um middleware de autenticação JWT dinâmico
// Permite que a configuração JWT seja obtida dinamicamente a cada requisição
//
// Funcionalidade:
// 1. Obtém a configuração JWT através da função authFunc
// 2. Cria um verificador JWT baseado nesta configuração
// 3. Aplica verificação e autenticação do token
// 4. Se válido: passa para o próximo handler
// 5. Se inválido: retorna erro de autenticação
//
// Parâmetros:
//   - authFunc: função que retorna *jwtauth.JWTAuth com configuração atual
//     Útil quando a chave JWT, algoritmo ou outras configurações podem mudar
//     dinamicamente (ex: rotação de chaves, multi-tenant, etc.)
//
// Retorna:
//   - Função middleware que pode ser usada na cadeia de middlewares
//
// Exemplo de uso:
//
//	func getJWTAuth() *jwtauth.JWTAuth {
//	    return jwtauth.New("HS256", []byte(getCurrentSecretKey()), nil)
//	}
//	router.Use(middlewares.DynamicVerifier(getJWTAuth))
//
// Diferença do verifier estático:
//   - Verifier estático: configuração fixada na inicialização
//   - DynamicVerifier: configuração obtida a cada requisição
func DynamicVerifier(authFunc func() *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtém a configuração JWT atual (pode variar a cada call)
			auth := authFunc()

			// Cria o verificador JWT baseado na configuração atual
			// Verifier: extrai e valida o token do header Authorization
			verifier := jwtauth.Verifier(auth)

			// Authenticator: verifica se o token é válido e não expirado
			authenticator := jwtauth.Authenticator

			// Aplica verificação e autenticação em sequência:
			// 1. verifier: extrai token do header e valida assinatura
			// 2. authenticator: verifica expiração e outras claims
			// 3. se tudo OK: chama next handler
			verifier(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authenticator(next).ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}
}

// StaticVerifier cria um middleware de autenticação JWT com configuração fixa
// Usa uma configuração JWT estática definida na inicialização da aplicação
//
// Funcionalidade:
// 1. Usa configuração JWT pré-definida (não muda durante execução)
// 2. Aplica verificação e autenticação do token
// 3. Mais performático que DynamicVerifier (sem overhead de função)
//
// Parâmetros:
//   - auth: instância *jwtauth.JWTAuth com configuração fixa
//
// Retorna:
//   - Função middleware que pode ser usada na cadeia de middlewares
//
// Exemplo de uso:
//
//	jwtAuth := jwtauth.New("HS256", []byte("secret-key"), nil)
//	router.Use(middlewares.StaticVerifier(jwtAuth))
//
// Use quando:
//   - Chave JWT é fixa durante toda execução da aplicação
//   - Performance é crítica (evita overhead de função dinâmica)
//   - Configuração simples sem necessidade de mudanças em runtime
func StaticVerifier(auth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return jwtauth.Verifier(auth)
}

// RequireAuth middleware que exige autenticação válida
// Complementa os verificadores, garantindo que o token seja válido
//
// Funcionalidade:
// 1. Verifica se o token JWT foi validado pelos middlewares anteriores
// 2. Se token é válido: permite acesso ao endpoint
// 3. Se token inválido/ausente: retorna erro 401 Unauthorized
//
// Deve ser usado APÓS um verificador (StaticVerifier ou DynamicVerifier)
//
// Exemplo de uso:
//
//	router.Use(middlewares.StaticVerifier(jwtAuth))
//	router.Use(middlewares.RequireAuth())
//
//	// Ou em grupo específico:
//	router.Route("/api/protected", func(r chi.Router) {
//	    r.Use(middlewares.StaticVerifier(jwtAuth))
//	    r.Use(middlewares.RequireAuth())
//	    r.Get("/user", getUserHandler)
//	})
//
// Retorna:
//   - Função middleware que pode ser usada na cadeia de middlewares
func RequireAuth() func(http.Handler) http.Handler {
	return jwtauth.Authenticator
}

// OptionalAuth middleware que permite acesso com ou sem autenticação
// Token JWT é verificado se presente, mas não é obrigatório
//
// Funcionalidade:
// 1. Se token presente e válido: adiciona informações do usuário ao contexto
// 2. Se token presente mas inválido: retorna erro 401
// 3. Se token ausente: permite acesso (contexto sem informações de usuário)
//
// Útil para endpoints que:
// - Funcionam para usuários anônimos e autenticados
// - Personalizam resposta baseada na autenticação
// - APIs públicas com recursos extras para usuários logados
//
// Parâmetros:
//   - auth: instância *jwtauth.JWTAuth para verificação
//
// Exemplo de uso:
//
//	router.Use(middlewares.OptionalAuth(jwtAuth))
//	router.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
//	    _, claims, _ := jwtauth.FromContext(r.Context())
//	    if claims != nil {
//	        // Usuário autenticado - mostrar posts privados também
//	    } else {
//	        // Usuário anônimo - mostrar apenas posts públicos
//	    }
//	})
func OptionalAuth(auth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Aplica verificação do token se presente
			verifier := jwtauth.Verifier(auth)
			verifier(next).ServeHTTP(w, r)
			// Note: NÃO usa Authenticator, permitindo acesso mesmo sem token válido
		})
	}
}
