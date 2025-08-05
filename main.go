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
	amount += u.Balance
}

func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return errors.New("Недостаточно средств ")
	}
	u.Balance -= amount
	return nil
}

func processAccount(account BankAccount) {
	fmt.Println("Текущий баланс:", account.GetBalance())

	account.Deposit(500)
	fmt.Println("После пополнения 500:", account.GetBalance())

	err := account.Withdraw(200)
	if err != nil {
		fmt.Println("Ошибка при снятии:", err)
	} else {
		fmt.Println("после снятие 200:", account.GetBalance())
	}

}

func main() {
	user := &User{
		Id:      "1",
		Name:    "Artem",
		Balance: 1000.0,
	}

	var account BankAccount = user

	processAccount(account)

	fmt.Println("Итоговый баланс пользователя :", user.Balance)
}
