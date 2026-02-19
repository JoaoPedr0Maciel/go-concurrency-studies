# 🔬 Go Concurrency Lab

Projeto educacional focado no estudo prático de **concorrência em Go**.

## 📚 Conceitos Estudados

### Goroutines e Worker Pools
- Múltiplos workers processando pagamentos em paralelo
- Cada worker possui seu próprio canal de entrada (`Queues`)
- Uso de `sync.WaitGroup` para sincronização

### Channels
- **Channels por worker**: Cada worker tem seu próprio canal para evitar contenção
- **Channel de resultados**: Canal compartilhado para coletar resultados do processamento
- **Buffered channels**: Uso de buffers para melhorar throughput

### Sincronização
- **Mutex (RWMutex)**: Proteção thread-safe do mapa de contas
- **Double-check locking**: Padrão para criação segura de contas
- **Mutex por conta**: Cada `Account` possui seu próprio mutex interno

### Hash Routing
- Pagamentos da mesma conta sempre vão para o mesmo worker
- Garante processamento sequencial por conta (evita race conditions)
- Implementado usando `hash/fnv` para distribuição consistente

### Context e Cancelamento
- Uso de `context.Context` para controle de ciclo de vida dos workers
- Cancelamento graceful ao parar o processador

### Prevenção de Deadlocks
- Consumo de resultados **antes** de enfileirar pagamentos
- Evita que o buffer de resultados encha e trave os workers

## 🏗️ Arquitetura

```
main.go
  └── Processor
      ├── Workers (goroutines)
      │   ├── Queue[0] (channel)
      │   ├── Queue[1] (channel)
      │   └── ...
      ├── Result (channel compartilhado)
      └── Accounts (map[string]*Account)
          └── Account (com mutex interno)
```

## 🚀 Como Executar

```bash
go run cmd/api/main.go
```

## 💡 Pontos Importantes

1. **Hash Routing**: Pagamentos da mesma `AccountID` sempre vão para o mesmo worker, garantindo ordem sequencial
2. **Consumo Antecipado**: Resultados são consumidos antes de enfileirar para evitar deadlock
3. **Double-Check Locking**: Padrão usado em `getOrCreateAccount` para criação thread-safe
4. **Graceful Shutdown**: Fechamento ordenado de canais e espera de workers finalizarem

## 📝 Estrutura do Projeto

```
go-projects/
├── cmd/
│   └── api/
│       └── main.go          # Ponto de entrada
├── internal/
│   ├── payment/
│   │   ├── processor.go     # Lógica de processamento concorrente
│   │   ├── payment.go       # Estruturas de pagamento
│   │   └── generate-payments.go
│   ├── account/
│   │   └── account.go       # Conta com mutex interno
│   └── infra/
│       └── db/
│           └── db.go        # DB fake para simulação
└── go.mod
```
