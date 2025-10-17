# Rate Limiter - Desafio GoExpert

## Descrição do Projeto

Este projeto implementa um **Rate Limiter** em Go, desenvolvido como parte do desafio GoExpert. O sistema controla o número de requisições que podem ser feitas para uma API dentro de um período específico, protegendo contra abuso e sobrecarga do servidor.

## Como Funciona o Rate Limiter

### Arquitetura

O Rate Limiter implementa um padrão de **Strategy** para persistência de dados e utiliza o **Redis** como store de dados para rastreamento de requisições. O sistema funciona de duas formas:

1. **Limitação por IP**: Controla requisições baseado no endereço IP do cliente
2. **Limitação por Token**: Controla requisições baseado em tokens de API fornecidos no header `API_KEY`

### Fluxo de Funcionamento

1. **Interceptação**: O middleware intercepta todas as requisições HTTP
2. **Identificação**: O sistema verifica se existe um token `API_KEY` no header
3. **Escolha da Estratégia**:
   - Se houver token: aplica limite baseado no token
   - Se não houver token: aplica limite baseado no IP
4. **Verificação**: Consulta o Redis para verificar o número atual de requisições
5. **Decisão**:
   - Se dentro do limite: permite a requisição e incrementa o contador
   - Se excedeu o limite: bloqueia a requisição e retorna HTTP 429
6. **Bloqueio Temporário**: Quando o limite é excedido, o sistema define um tempo de bloqueio

### Componentes Principais

- **RateLimiter**: Lógica principal do limitador
- **RateLimiterMiddleware**: Middleware HTTP que intercepta requisições
- **PersistenceStrategy**: Interface para estratégias de persistência
- **RedisPersistenceStrategy**: Implementação usando Redis
- **EnvConfig**: Gerenciamento de configurações

## Configuração

### Variáveis de Ambiente

O sistema é configurado através de variáveis de ambiente definidas no arquivo `.env`:

```env
# Janela de tempo para contagem de requisições (em segundos)
RATE_WINDOW=1

# Limite de requisições por IP por segundo
IP_LIMIT_PER_SECOND=10

# Tempo de retry após exceder o limite (em segundos)
RETRY_AFTER_SECONDS=60

# Limite para tokens desconhecidos por segundo
UNKNOWN_TOKEN_RATE_PER_SECOND=20

# Limite para token básico por segundo
BASIC_TOKEN_RATE_PER_SECOND=50

# Limite para token ultra por segundo
ULTRA_TOKEN_RATE_PER_SECOND=100

# Configurações do Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_DB=0
REDIS_USER=
REDIS_PASSWORD=
```

### Tipos de Tokens

O sistema suporta diferentes tipos de tokens com limites específicos:

- **NORMAL_TOKEN**: Limite configurado por `BASIC_TOKEN_RATE_PER_SECOND` (padrão: 50 req/seg)
- **ULTRA_TOKEN**: Limite configurado por `ULTRA_TOKEN_RATE_PER_SECOND` (padrão: 100 req/seg)
- **Tokens desconhecidos**: Limite configurado por `UNKNOWN_TOKEN_RATE_PER_SECOND` (padrão: 20 req/seg)

### Configuração dos Limites

1. **Por IP**: Se nenhum token for fornecido, aplica o limite `IP_LIMIT_PER_SECOND`
2. **Por Token**: Se um token for fornecido no header `API_KEY`, aplica o limite correspondente ao token
3. **Fallback**: Tokens não reconhecidos usam `UNKNOWN_TOKEN_RATE_PER_SECOND`

## Instalação e Execução

### Pré-requisitos

- Go 1.19 ou superior
- Redis (ou Docker para executar via container)
- Docker e Docker Compose (opcional)

### Execução Local

1. Clone o repositório:
```bash
git clone <url-do-repositorio>
cd desafio-goexpert-1
```

2. Instale as dependências:
```bash
go mod tidy
```

3. Configure o arquivo `.env` conforme necessário

4. Execute o Redis:
```bash
docker run -d -p 6379:6379 redis:alpine
```

5. Execute a aplicação:
```bash
go run cmd/main.go
```

A aplicação estará disponível em `http://localhost:8080`

### Execução com Docker

1. Execute usando Docker Compose:
```bash
docker-compose up -d
```

Isso iniciará tanto a aplicação quanto o Redis automaticamente.

## Uso

### Requisições Simples (Limitação por IP)

```bash
curl http://localhost:8080/
```

### Requisições com Token

```bash
# Token normal
curl -H "API_KEY: NORMAL_TOKEN" http://localhost:8080/

# Token ultra
curl -H "API_KEY: ULTRA_TOKEN" http://localhost:8080/

# Token personalizado
curl -H "API_KEY: meu-token-personalizado" http://localhost:8080/
```

### Testando os Limites

Para testar o rate limiting, execute múltiplas requisições rapidamente:

```bash
# Teste com limitação por IP (10 req/seg por padrão)
for i in {1..15}; do curl http://localhost:8080/; echo; done

# Teste com token (50 req/seg por padrão)
for i in {1..60}; do curl -H "API_KEY: NORMAL_TOKEN" http://localhost:8080/; echo; done
```

### Resposta de Limite Excedido

Quando o limite é excedido, o sistema retorna:

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 60

{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}
```

## Testes

Execute os testes unitários:

```bash
go test ./internal/limiter/
```

## Estrutura do Projeto

```
desafio-goexpert-1/
├── cmd/
│   └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── EnvConfig.go        # Configurações e variáveis de ambiente
│   ├── limiter/
│   │   ├── RateLimiter.go      # Implementação do rate limiter
│   │   ├── RateLimiterInterface.go  # Interface do rate limiter
│   │   └── RateLimiter_test.go # Testes unitários
│   ├── middleware/
│   │   └── RateLimiterMiddleware.go # Middleware HTTP
│   └── strategy/
│       ├── PersistenceStrategyInterface.go # Interface de persistência
│       └── impl/
│           ├── RedisPersistenceStrategy.go # Implementação Redis
│           └── RedisPersistenceFactory.go  # Factory para Redis
├── .env                        # Variáveis de ambiente
├── docker-compose.yml          # Configuração Docker Compose
├── Dockerfile                  # Imagem Docker da aplicação
├── go.mod                      # Dependências Go
└── README.md                   # Documentação
```

## Dependências

- **github.com/go-chi/chi/v5**: Router HTTP
- **github.com/go-redis/redis/v8**: Cliente Redis
- **github.com/joho/godotenv**: Carregamento de variáveis de ambiente

## Características Técnicas

- **Goroutines**: Usa goroutines para limpeza assíncrona de chaves expiradas
- **Strategy Pattern**: Permite diferentes estratégias de persistência
- **Middleware Pattern**: Integração transparente com aplicações HTTP
- **Configuração Flexível**: Todas as configurações via variáveis de ambiente
- **Fallback Gracioso**: Valores padrão para todas as configurações
- **Thread-Safe**: Operações seguras para uso concorrente

## Monitoramento e Logs

O sistema gera logs para:
- Erros de conexão com Redis
- Chaves sendo deletadas após expiração
- Erros de parsing de configuração
- Erros durante verificação de limites
