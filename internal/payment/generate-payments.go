package payment

import (
	"fmt"
	"math/rand"
)

func GeneratePayments(total int, accounts int) []Payment {

	payments := make([]Payment, 0, total)

	for i := range total {
		accID := fmt.Sprintf("acc-%d", rand.Intn(accounts))

		var pType string
		var amount int64

		if rand.Intn(2) == 0 {
			pType = "deposit"
			amount = int64(rand.Intn(500) + 1)
			fmt.Println("deposit", amount)
		} else {
			pType = "withdraw"
			amount = int64(rand.Intn(300) + 1)
			fmt.Println("withdraw", amount)
		}

		payments = append(payments, Payment{
			ID:        fmt.Sprintf("pay-%d", i),
			AccountID: accID,
			Amount:    amount,
			Type:      pType,
		})
	}

	return payments
}
