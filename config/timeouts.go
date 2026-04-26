package config

import "time"

// Timeouts define os tempos limite do servidor HTTP da API.
// Todos os valores são em SEGUNDOS. Quando 0 (ou ausente no JSON), o
// helper correspondente aplica o default seguro indicado em cada campo.
type Timeouts struct {
	// ReadHeaderTimeout — tempo máximo para o servidor ler o cabeçalho
	// HTTP da requisição (do início da conexão até o fim dos headers).
	//
	// Para que serve: protege contra ataques tipo "slowloris", em que
	// um cliente abre muitas conexões e envia bytes lentamente para
	// esgotar as conexões disponíveis.
	//
	// Aumentar: clientes em redes muito lentas conseguem completar
	// o handshake/headers. Pode ser necessário em redes móveis ruins.
	// Diminuir: rejeita conexões lentas/maliciosas mais cedo,
	// liberando recursos do servidor.
	//
	// Default quando 0: 10 segundos.
	ReadHeaderTimeout int `json:"read_header_timeout"`

	// ReadTimeout — tempo máximo para o servidor ler a requisição
	// inteira (cabeçalho + body) do cliente.
	//
	// Para que serve: limita por quanto tempo uma requisição pode
	// ficar "em andamento" do lado de leitura.
	//
	// Aumentar: necessário se sua API recebe uploads grandes ou
	// clientes lentos enviando POSTs/PUTs longos.
	// Diminuir: abandona requisições demoradas mais cedo, libera
	// goroutines presas. Cuidado: se for menor que o tempo real
	// de upload típico, uploads serão cortados no meio.
	//
	// Default quando 0: 30 segundos.
	ReadTimeout int `json:"read_timeout"`

	// WriteTimeout — tempo máximo para o servidor escrever a resposta
	// completa para o cliente (a partir do fim da leitura do header).
	//
	// Para que serve: limita por quanto tempo o handler pode demorar
	// para gerar e enviar a resposta.
	//
	// Aumentar: necessário se você tem endpoints SÍNCRONOS que
	// processam por muito tempo (ex: relatórios pesados, exportações,
	// queries longas). Caso contrário o cliente recebe a conexão
	// fechada no meio da resposta.
	// Diminuir: corta handlers que demoram demais, evita "pendurar"
	// o servidor com handlers lentos.
	//
	// Default quando 0: 60 segundos.
	WriteTimeout int `json:"write_timeout"`

	// IdleTimeout — tempo máximo que uma conexão keep-alive ociosa
	// fica aberta aguardando a próxima requisição do mesmo cliente.
	//
	// Para que serve: controla reuso de conexões TCP entre múltiplas
	// requisições do mesmo cliente.
	//
	// Aumentar: reduz o overhead de reabrir conexões para clientes
	// que fazem requisições frequentes (ex: SPAs, microserviços
	// chamando uns aos outros). Mantém mais conexões abertas.
	// Diminuir: libera recursos mais rápido em clientes esporádicos.
	// Pode aumentar latência se o cliente sempre precisa renegociar.
	//
	// Default quando 0: 120 segundos.
	IdleTimeout int `json:"idle_timeout"`

	// ShutdownTimeout — tempo máximo que o servidor espera as
	// requisições em andamento terminarem ao receber sinal de parada
	// (graceful shutdown).
	//
	// Para que serve: ao parar o serviço, aguarda os handlers ativos
	// completarem suas respostas em vez de cortar conexões abruptamente.
	//
	// Aumentar: dá mais tempo para handlers longos completarem antes
	// do shutdown forçado. Necessário se você tem handlers que podem
	// demorar para terminar.
	// Diminuir: o serviço para mais rápido, mas pode cortar respostas
	// no meio. Útil em ambientes onde reinício rápido é prioridade.
	//
	// Default quando 0: 10 segundos.
	ShutdownTimeout int `json:"shutdown_timeout"`
}

func (t Timeouts) ReadHeader() time.Duration {
	if t.ReadHeaderTimeout <= 0 {
		return 10 * time.Second
	}
	return time.Duration(t.ReadHeaderTimeout) * time.Second
}

func (t Timeouts) Read() time.Duration {
	if t.ReadTimeout <= 0 {
		return 30 * time.Second
	}
	return time.Duration(t.ReadTimeout) * time.Second
}

func (t Timeouts) Write() time.Duration {
	if t.WriteTimeout <= 0 {
		return 60 * time.Second
	}
	return time.Duration(t.WriteTimeout) * time.Second
}

func (t Timeouts) Idle() time.Duration {
	if t.IdleTimeout <= 0 {
		return 120 * time.Second
	}
	return time.Duration(t.IdleTimeout) * time.Second
}

func (t Timeouts) Shutdown() time.Duration {
	if t.ShutdownTimeout <= 0 {
		return 10 * time.Second
	}
	return time.Duration(t.ShutdownTimeout) * time.Second
}
