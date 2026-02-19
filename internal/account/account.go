package account

import (
	"errors"
	"sync"
)

type Account struct {
	ID      string
	Balance int64
	mutex   sync.Mutex
}

func (a *Account) Deposit(amount int64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.Balance += amount
}

func (a *Account) Withdraw(amount int64) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.Balance < amount {
		return errors.New("insufficient funds")
	}

	a.Balance -= amount
	return nil
}
