package payment

import (
	"context"
	"hash/fnv"
	"sync"

	"github.com/JoaoPedr0Maciel/go-studies/internal/account"
	infra "github.com/JoaoPedr0Maciel/go-studies/internal/infra/db"
)

type Processor struct {
	DB       *infra.FakeDB               // representa um db fake
	Accounts map[string]*account.Account // mapa de contas por AccountID
	Queues   []chan Payment              // um canal por worker
	Workers  int
	Wg       *sync.WaitGroup
	Result   chan *PaymentResult
	Ctx      context.Context
	Cancel   context.CancelFunc
	mu       sync.RWMutex // mutex para proteger o mapa de Accounts
}

func InitializeProcessor(workers int, queueSize int) *Processor {
	ctx, cancel := context.WithCancel(context.Background())

	fakeDB := &infra.FakeDB{}

	// Cria um canal por worker
	queues := make([]chan Payment, workers)
	for i := range queues {
		queues[i] = make(chan Payment, queueSize)
	}

	processor := &Processor{
		DB:       fakeDB,
		Accounts: make(map[string]*account.Account),
		Queues:   queues,
		Workers:  workers,
		Wg:       &sync.WaitGroup{},
		Result:   make(chan *PaymentResult, queueSize*workers),
		Ctx:      ctx,
		Cancel:   cancel,
	}

	processor.Start()

	return processor
}

func (p *Processor) Start() {
	// Inicia N workers, cada um com seu próprio canal
	// Pagamentos da mesma conta sempre vão para o mesmo worker (via hash)
	for i := range p.Workers {
		p.Wg.Add(1)
		go p.Worker(i)
	}
}

// getOrCreateAccount retorna a conta para o AccountID, criando se não existir
func (p *Processor) getOrCreateAccount(accountID string) *account.Account {
	p.mu.RLock()
	if acc, exists := p.Accounts[accountID]; exists {
		p.mu.RUnlock()
		return acc
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check: pode ter sido criado por outro worker enquanto esperava
	if acc, exists := p.Accounts[accountID]; exists {
		return acc
	}

	acc := &account.Account{ID: accountID, Balance: 1000}
	p.Accounts[accountID] = acc
	return acc
}

// routePayment determina qual worker deve processar o pagamento baseado no AccountID
func (p *Processor) routePayment(payment Payment) int {
	accountID := payment.AccountID
	if accountID == "" {
		accountID = "acc-1" // default
	}

	// Usa hash do AccountID para garantir que a mesma conta sempre vai para o mesmo worker
	h := fnv.New32a()
	h.Write([]byte(accountID))
	return int(h.Sum32()) % p.Workers
}

func (p *Processor) Worker(id int) {
	defer p.Wg.Done()

	for {
		select {
		case <-p.Ctx.Done():
			return

		case payment, ok := <-p.Queues[id]:
			if !ok {
				return
			}
			result := p.handlePayment(payment)

			p.Result <- result
		}
	}
}

func (p *Processor) Queue(payment Payment) {
	workerID := p.routePayment(payment)
	p.Queues[workerID] <- payment
}

func (p *Processor) Stop() {
	for i := range p.Queues {
		close(p.Queues[i])
	}
	p.Wg.Wait()
	close(p.Result)
}
