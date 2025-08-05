package main

import (
	"errors"
	"fmt"
)

type BankAccount interface {
	Deposit(amount float64)
	Withdraw(amount float64) error
	GetBalance() float64
}

type User struct {
	Id      string
	Name    string
	Balance float64
}

func (u *User) GetBalance() float64 {
	return u.Balance
}

func (u *User) Deposit(amount float64) {
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return errors.New("Недостаточно средств ")
	}
	u.Balance -= amount
	return nil
}

func processAccount(account BankAccount, depositAmt float64, withdrawAmt float64) {
	fmt.Printf("Начальный баланс: %.2f\n", account.GetBalance())

	account.Deposit(depositAmt)
	fmt.Printf("После пополнения на %.2f: %.2f\n", depositAmt, account.GetBalance())

	err := account.Withdraw(withdrawAmt)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("После снятия %.2f: %.2f\n", withdrawAmt, account.GetBalance())
	}
}

func main() {
	user1 := &User{
		Id:      "1",
		Name:    "Artem",
		Balance: 1000.0,
	}

	user2 := &User{
		Id:      "2",
		Name:    "Egor",
		Balance: 500.0,
	}

	var account1 BankAccount = user1
	var account2 BankAccount = user2

	fmt.Println("---------------------------")
	fmt.Println("Операция с:", user1.Name)
	processAccount(account1, 300, 400)
	fmt.Println("===========================")
	fmt.Println("Операция с:", user2.Name)
	processAccount(account2, 200, 300)
	fmt.Println("---------------------------")

	fmt.Println("Итоговый баланс пользователя:", user1.Balance, user2.Balance)
	fmt.Println("---------------------------")
}
