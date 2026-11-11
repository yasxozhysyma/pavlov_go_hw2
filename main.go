package main

import "fmt"

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

	for _, transaction := range ps.TransactionQueue {
		fmt.Printf(
			"Перевод от %s к %s: %.2f\n",
			transaction.FromID,
			transaction.ToID,
			transaction.Amount,
		)

		if err := ps.ProcessingTransactions(transaction); err != nil {
			fmt.Println("Ошибка:", err)
		}
	}

	fmt.Println("Итого:")
	fmt.Printf("Анна: %.2f\n", ps.Users["1"].Balance)
	fmt.Printf("Иван: %.2f\n", ps.Users["2"].Balance)
}
