#HTTP Package
1. [THttp Structure](#THttpStructure)
   
 The http package provides a set of functions and structures for making HTTP requests and handling WebSocket connections. This package is designed to simplify working with HTTP and WebSocket protocols in your Go applications.

3. [Fields](#Fields)
   
#req: Represents an HTTP request.
#ws: Represents a WebSocket connection.
#Auth2: An instance of the auth2 type for handling authentication.
#Request: An instance of the Request type to configure the HTTP request.
#Response: An instance of the Response type to store the HTTP response.
#Metodo: An enum representing the HTTP request method (e.g., GET, POST).
#AuthorizationType: The type of authorization (e.g., AutoDetect, Bearer, Basic).
#Authorization: The authorization token.
#Password: The user's password (for Basic authentication).
#UserName: The user's username (for Basic authentication).
#url: The full URL for the HTTP request.
#Protocolo: The protocol (http or https).
#Host: The host (e.g., www.example.com).
#Path: The path (e.g., /product).
#Varibles: A collection of variables.
#Params: A collection of query parameters.
#Proxy: An instance of the proxy type for configuring proxy settings.
#EncType: The type of request content (e.g., Form data, Raw, Binary).
#Timeout: The request timeout in seconds.
#OnSend: An interface for handling WebSocket events.

5. [Methods](#Methods)

#SetMetodoStr(value string) error: Set the HTTP request method using a string.
#GetMetodoStr() string: Get the HTTP request method as a string.
#SetMetodo(value TMethod) error: Set the HTTP request method using an enum.
#GetMetodo() TMethod: Get the HTTP request method.
#SetAuthorizationType(value AuthorizationType) error: Set the authorization type.
#GetAuthorizationType() AuthorizationType: Get the authorization type.
#SetUrl(value string) error: Set the URL for the HTTP request.
#GetFullURL() (string, error): Get the full URL.
#GetUrl() string: Get the URL with query parameters.
#completHeader(): Complete the request headers.
#completAutorization(req *http.Request) error: Complete the authorization for the request.
#Send() (*Response, error): Send the HTTP request and receive the response.
#websocketClient() error: Establish a WebSocket connection.
#Conectar() error: Connect to a WebSocket server.
#IsConect() bool: Check if the WebSocket connection is established.
#Desconectar() error: Close the WebSocket connection.
#EnviarBinario(messageType int, data []byte) error: Send binary data over WebSocket.
#EnviarTexto(messageType int, data string) error: Send text data over WebSocket.
#EnviarTextTypeTextMessage(data []byte) error: Send text data with a TextMessage type.
#EnviarBinarioTypeBinaryMessage(data []byte) error: Send binary data with a BinaryMessage type.
#NewHttp() *THttp: Create a new THttp instance.
This package simplifies making HTTP requests and handling WebSocket connections in your Go applications. It provides methods and structures to configure and send requests, handle responses, and establish WebSocket connections.

For more details, refer to the Go source code in the http package.
