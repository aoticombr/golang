package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const secretKey = "123456789"

func main() {
	dataDir := path.Join(".", "data")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, os.ModeDir)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		fmt.Println(r.Method)

		requestData := map[string]interface{}{
			"method":  r.Method,
			"url":     r.URL.Path,
			"headers": r.Header,
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse the URL-encoded data using querystring package
		parsedData, err := urlEncodedToMap(string(data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		requestData["body"] = parsedData

		now := time.Now()
		guid := uuid.New().String()
		filename := fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d-%s.json", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), guid)
		filePath := path.Join(dataDir, filename)

		requestDataJSON, err := json.MarshalIndent(requestData, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile(filePath, requestDataJSON, 0644)
		if err != nil {
			fmt.Println("Erro ao salvar dados da solicitação no arquivo:", err)
			http.Error(w, "Erro ao salvar dados da solicitação no arquivo.", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Dados da solicitação salvos em %s\n", filename)

		if r.URL.Path == "/token" {
			fmt.Println("Recebendo solicitação de token...")

			// Extract the scope from the request (if provided)
			requestedScope := parsedData["scope"]
			if requestedScope == "" {
				requestedScope = "default_scope"
			}

			// Build the token payload with the necessary information
			payload := jwt.MapClaims{
				"user_id": 123, // Replace this with the appropriate user ID
				"scope":   requestedScope,
			}

			// Generate the JWT token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
			tokenString, err := token.SignedString([]byte(secretKey))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Build the response including the JWT token
			tokenResponse := map[string]interface{}{
				"access_token":       tokenString,
				"expires_in":         300,
				"refresh_expires_in": 0,
				"token_type":         "Bearer",
				"not-before-policy":  0,
				"scope":              requestedScope,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse)
		} else {
			// Here you can add JWT token validation
			authorizationHeader := r.Header.Get("Authorization")

			if strings.HasPrefix(authorizationHeader, "Bearer ") {
				tokenString := authorizationHeader[7:]

				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(secretKey), nil
				})

				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Header().Set("Content-Type", "text/plain")
					w.Write([]byte("Token JWT inválido"))
				} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					fmt.Println("Token JWT válido", claims)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(requestData)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Header().Set("Content-Type", "text/plain")
					w.Write([]byte("Token JWT inválido"))
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(requestData)
			}
		}
	})

	port := 3003
	fmt.Printf("O servidor está em execução em http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func urlEncodedToMap(data string) (map[string]string, error) {
	parsedData, err := url.ParseQuery(data)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for key, values := range parsedData {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result, nil
}
