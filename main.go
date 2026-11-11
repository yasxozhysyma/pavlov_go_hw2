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

func main() {
	user1 := &User{ID: "1", Name: "Анна", Balance: 1000}
	user2 := &User{ID: "2", Name: "Иван", Balance: 500}

	user1.Deposit(200)
	fmt.Printf("%s после пополнения: %.2f\n", user1.Name, user1.Balance)

	user2.Deposit(100)
	fmt.Printf("%s после пополнения: %.2f\n", user2.Name, user2.Balance)

	if err := user1.Withdraw(300); err != nil {
		fmt.Println("Ошибка:", err)
	}
	fmt.Printf("%s после снятия: %.2f\n", user1.Name, user1.Balance)

	if err := user2.Withdraw(1000); err != nil {
		fmt.Println("Ошибка:", err)
	}
	fmt.Printf("%s после попытки снятия: %.2f\n", user2.Name, user2.Balance)
}
