package main

import (
	"fmt"
	"sync"

	"github.com/JoaoPedr0Maciel/go-studies/internal/payment"
)

func main() {
	processor := payment.InitializeProcessor(10, 10)

	payments := payment.GeneratePayments(1, 10)

	// consumir resultados ANTES de enfileirar.
	// se enfileirar muitos pagamentos primeiro, o buffer de Result enche e trava os workers,
	// os workers travam no send de Result e o producer trava enfileirando resultando em um deadlock.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range processor.Result {
			if result.Success {
				fmt.Printf("Payment %s processed successfully. Balance: %d\n", result.PaymentID, result.Balance)
			} else {
				fmt.Printf("Payment %s failed: %s. Balance: %d\n", result.PaymentID, result.Error, result.Balance)
			}
		}
	}()

	for i := range payments {
		processor.Queue(payments[i])
	}

	processor.Stop()
	wg.Wait()

	if acc, exists := processor.Accounts["acc-1"]; exists {
		fmt.Printf("Final balance: %d\n", acc.Balance)
	}
}
