# Implementação de Proxy - Arquitetura Modular

## Filosofia da Implementação

Esta implementação segue o princípio de **partes independentes**, onde cada componente (proxy, certificados, transport customizados) é tratado de forma modular e só é ativado quando realmente necessário.

## Características Principais

### 1. Transport Inteligente
- O `GetTransport()` **só cria um transport quando necessário**
- Se nenhuma configuração especial é necessária (sem proxy, sem certificados, sem transport customizado), retorna `nil`
- O cliente HTTP usa o transport padrão do Go quando `nil` é passado

### 2. Configuração Independente
- Cada parte pode ser configurada independentemente
- Proxy pode ser ativado/desativado sem afetar outras configurações
- Certificados funcionam independentemente do proxy
- Transport customizado funciona independentemente de proxy e certificados

### 3. Padrão Set/Get Consistente
- Todos os métodos seguem o padrão `Set/Get` estabelecido no código
- Propriedades podem ser configuradas individualmente
- Métodos convenientes para configuração rápida

## Estrutura de Código

### Transport Logic (GetTransport)
```go
func (H *THttp) GetTransport() *http.Transport {
    var needTransport bool
    var transport *http.Transport

    // Verificar se precisamos de um transport customizado
    if H.Proxy != nil && H.Proxy.Ativo {
        needTransport = true
    }
    if H.Certificate.PathCrt != "" && H.Certificate.PathPriv != "" {
        needTransport = true
    }
    if H.TransportType != TNenhum {
        needTransport = true
    }

    // Só criar transport se realmente precisar
    if needTransport {
        transport = &http.Transport{}
        
        // Configurar proxy se ativo
        if H.Proxy != nil && H.Proxy.Ativo {
            H.Proxy.SetProxy(transport) // Modifica o transport existente
        }
        
        return transport
    }

    return nil // Usa transport padrão do Go
}
```

### Proxy Logic (SetProxy)
```go
func (p *proxy) SetProxy(transport *http.Transport) error {
    if !p.Ativo {
        return nil // Proxy desabilitado, não configura nada
    }
    
    // Validações
    if p.Host == "" {
        return fmt.Errorf("host do proxy não pode estar vazio")
    }
    
    if p.Port <= 0 || p.Port > 65535 {
        return fmt.Errorf("porta do proxy deve estar entre 1 e 65535")
    }

    // Configurar o transport existente
    proxyURL, err := url.Parse(p.getUrl())
    if err != nil {
        return fmt.Errorf("erro ao fazer parse da URL do proxy: %v", err)
    }
    
    transport.Proxy = http.ProxyURL(proxyURL)
    return nil
}
```

## API de Uso

### Configuração Básica
```go
client := http.NewHttp()
defer client.Free()

// Configuração completa de uma vez
client.SetProxyConfig("proxy.empresa.com", 8080, "usuario", "senha")

// Verificar status
fmt.Printf("Proxy ativo: %t\n", client.GetProxyAtivo())
```

### Configuração Individual
```go
// Configurar cada propriedade individualmente
client.SetProxyHost("meu-proxy.com")
client.SetProxyPort(3128)
client.SetProxyUserName("usuario")
client.SetProxyPassword("senha")
client.SetProxyAtivo(true)

// Ler configurações
host := client.GetProxyHost()
port := client.GetProxyPort()
ativo := client.GetProxyAtivo()
```

### Ativação/Desativação
```go
// Desabilitar proxy temporariamente
client.SetProxyAtivo(false)

// Reabilitar
client.SetProxyAtivo(true)
```

## Comportamento do Sistema

### Sem Configurações Especiais
- `GetTransport()` retorna `nil`
- Cliente HTTP usa transport padrão do Go
- Performance otimizada, sem overhead

### Com Proxy Ativo
- `GetTransport()` cria novo transport
- Configura proxy no transport
- Cliente HTTP usa transport customizado

### Com Certificados
- `GetTransport()` cria novo transport
- Configura TLS no transport
- Funciona independentemente do proxy

### Com Proxy + Certificados
- `GetTransport()` cria um transport
- Configura proxy E TLS no mesmo transport
- Ambos funcionam em conjunto

## Vantagens da Implementação

1. **Performance**: Só cria transport quando necessário
2. **Modularidade**: Cada componente é independente
3. **Flexibilidade**: Configuração granular ou em bloco
4. **Consistência**: Segue padrões estabelecidos no código
5. **Robustez**: Validações e tratamento de erros adequados
6. **Manutenibilidade**: Código claro e bem estruturado

## Exemplos de Uso

Veja o arquivo `exemplos/proxy_example.go` para exemplos práticos de todas as funcionalidades.
