package main

import (
	"fmt"
	"sync"
)

type User struct {
	ID      string
	Name    string
	Balance float64
}

func (u *User) Deposit(amount float64) {
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	if u.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	u.Balance -= amount
	return nil
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentSystem struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

func (ps *PaymentSystem) AddUser(user *User) {
	ps.Users[user.ID] = user
}

func (ps *PaymentSystem) AddTransaction(transaction Transaction) {
	ps.TransactionQueue = append(ps.TransactionQueue, transaction)
}

func (ps *PaymentSystem) ProcessingTransactions(transaction Transaction) error {
	from, exists := ps.Users[transaction.FromID]
	if !exists {
		return fmt.Errorf("sender not found: %s", transaction.FromID)
	}

	to, exists := ps.Users[transaction.ToID]
	if !exists {
		return fmt.Errorf("recipient not found: %s", transaction.ToID)
	}

	if err := from.Withdraw(transaction.Amount); err != nil {
		return err
	}

	to.Deposit(transaction.Amount)
	return nil
}

func (ps *PaymentSystem) Worker(ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	for transaction := range ch {
		if err := ps.ProcessingTransactions(transaction); err != nil {
			fmt.Println("Ошибка:", err)
		}
	}
}

func main() {
	ps := &PaymentSystem{
		Users: make(map[string]*User),
	}

	ps.AddUser(&User{ID: "1", Name: "Анна", Balance: 1000})
	ps.AddUser(&User{ID: "2", Name: "Иван", Balance: 500})

	fmt.Println("Создан пользователь 1 с балансом 1000")
	fmt.Println("Создан пользователь 2 с балансом 500")

	ps.AddTransaction(Transaction{FromID: "1", ToID: "2", Amount: 200})
	ps.AddTransaction(Transaction{FromID: "2", ToID: "1", Amount: 50})

	ch := make(chan Transaction, len(ps.TransactionQueue))
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go ps.Worker(ch, &wg)
	}

	for _, transaction := range ps.TransactionQueue {
		ch <- transaction
	}

	close(ch)
	wg.Wait()

	fmt.Println("Итого:")
	fmt.Printf("Анна: %.2f (ожидается 850)\n", ps.Users["1"].Balance)
	fmt.Printf("Иван: %.2f (ожидается 650)\n", ps.Users["2"].Balance)
}
