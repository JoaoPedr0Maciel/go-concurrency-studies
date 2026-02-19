package payment

type Payment struct {
	ID        string
	AccountID string // ID da conta para roteamento
	Amount    int64
	Type      string
}

type PaymentResult struct {
	PaymentID string
	Balance   int64
	Success   bool
	Error     string
}

func (p *Processor) handlePayment(payment Payment) *PaymentResult {

	p.DB.Create(payment)

	// Obtém ou cria a conta para este AccountID
	acc := p.getOrCreateAccount(payment.AccountID)

	var err error
	switch payment.Type {
	case "deposit":
		acc.Deposit(payment.Amount)
	case "withdraw":
		err = acc.Withdraw(payment.Amount)
	}

	result := &PaymentResult{
		PaymentID: payment.ID,
		Balance:   acc.Balance,
		Success:   err == nil,
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result
}
