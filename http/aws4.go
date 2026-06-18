package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// aws4 guarda as credenciais/contexto para assinar requisições com
// AWS Signature Version 4 (SigV4).
//
// O header gerado tem a forma:
//
//	Authorization: AWS4-HMAC-SHA256 Credential=<AccessKey>/<data>/<regiao>/<servico>/aws4_request,
//	               SignedHeaders=host;x-amz-date, Signature=<assinatura>
type aws4 struct {
	AccessKey    string
	SecretKey    string
	Region       string // ex.: us-east-1
	Service      string // ex.: sqs, s3, execute-api
	SessionToken string // opcional - credenciais temporarias (STS)
}

func NewAws4() *aws4 {
	return &aws4{}
}

// sign calcula a assinatura SigV4 da requisição req (com o corpo payload) e
// define os headers X-Amz-Date, Authorization (e X-Amz-Security-Token, se houver
// SessionToken). Assina o conjunto minimo de headers (host;x-amz-date[;x-amz-security-token]),
// como faz o Postman para servicos do tipo Query (SQS, STS, etc.).
func (A *aws4) sign(req *http.Request, payload []byte) error {
	return A.signAt(req, payload, time.Now().UTC())
}

// signAt é o núcleo de sign() com o instante injetável (facilita testes determinísticos).
func (A *aws4) signAt(req *http.Request, payload []byte, now time.Time) error {
	if A.AccessKey == "" || A.SecretKey == "" || A.Region == "" || A.Service == "" {
		return fmt.Errorf("AWS SigV4: AccessKey, SecretKey, Region e Service são obrigatórios")
	}

	now = now.UTC()
	amzDate := now.Format("20060102T150405Z") // ISO8601 basico
	dateStamp := now.Format("20060102")       // YYYYMMDD

	req.Header.Set("X-Amz-Date", amzDate)
	if A.SessionToken != "" {
		req.Header.Set("X-Amz-Security-Token", A.SessionToken)
	}

	host := req.URL.Host

	// ----- 1) Canonical Request -----
	canonicalURI := awsCanonicalURI(req.URL.EscapedPath())
	canonicalQuery := awsCanonicalQueryString(req.URL.Query())
	payloadHash := awsSha256Hex(payload)

	// headers assinados (ordenados por nome, em minusculas)
	type kv struct{ k, v string }
	headers := []kv{
		{"host", host},
		{"x-amz-date", amzDate},
	}
	if A.SessionToken != "" {
		headers = append(headers, kv{"x-amz-security-token", A.SessionToken})
	}
	sort.Slice(headers, func(i, j int) bool { return headers[i].k < headers[j].k })

	var canonicalHeaders strings.Builder
	names := make([]string, 0, len(headers))
	for _, h := range headers {
		canonicalHeaders.WriteString(h.k)
		canonicalHeaders.WriteString(":")
		canonicalHeaders.WriteString(strings.TrimSpace(h.v))
		canonicalHeaders.WriteString("\n")
		names = append(names, h.k)
	}
	signedHeaders := strings.Join(names, ";")

	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders.String(),
		signedHeaders,
		payloadHash,
	}, "\n")

	// ----- 2) String to Sign -----
	const algorithm = "AWS4-HMAC-SHA256"
	credentialScope := strings.Join([]string{dateStamp, A.Region, A.Service, "aws4_request"}, "/")
	stringToSign := strings.Join([]string{
		algorithm,
		amzDate,
		credentialScope,
		awsSha256Hex([]byte(canonicalRequest)),
	}, "\n")

	// ----- 3) Signing key (HMAC encadeado) -----
	kDate := awsHmacSHA256([]byte("AWS4"+A.SecretKey), dateStamp)
	kRegion := awsHmacSHA256(kDate, A.Region)
	kService := awsHmacSHA256(kRegion, A.Service)
	kSigning := awsHmacSHA256(kService, "aws4_request")

	// ----- 4) Signature -----
	signature := hex.EncodeToString(awsHmacSHA256(kSigning, stringToSign))

	req.Header.Set("Authorization", fmt.Sprintf(
		"%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, A.AccessKey, credentialScope, signedHeaders, signature))

	return nil
}

func awsSha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func awsHmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

// awsCanonicalURI devolve o caminho canonico. Para serviços diferentes de S3 a AWS
// recomenda URI-encode (duplo) de cada segmento; aqui usamos o EscapedPath, que cobre
// os casos comuns (raiz "/" e caminhos simples como em SQS/STS/execute-api).
func awsCanonicalURI(path string) string {
	if path == "" {
		return "/"
	}
	return path
}

// awsCanonicalQueryString monta a query canonica: chaves ordenadas e cada chave/valor
// URI-encoded no estilo AWS (RFC 3986, espaço = %20, ~ não codificado).
func awsCanonicalQueryString(values url.Values) string {
	if len(values) == 0 {
		return ""
	}
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		vals := append([]string(nil), values[k]...)
		sort.Strings(vals)
		ek := awsURIEncode(k, true)
		for _, v := range vals {
			parts = append(parts, ek+"="+awsURIEncode(v, true))
		}
	}
	return strings.Join(parts, "&")
}

// awsURIEncode codifica conforme as regras do SigV4 (RFC 3986). Mantém sem codificar
// apenas os "unreserved": A-Z a-z 0-9 - _ . ~ (e '/' quando encodeSlash=false).
func awsURIEncode(s string, encodeSlash bool) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~':
			b.WriteByte(c)
		case c == '/' && !encodeSlash:
			b.WriteByte('/')
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}
