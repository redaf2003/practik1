package main

import (
	"errors"
	"fmt"
)

// Интерфейс банковского счета
type BankAccount interface {
	Deposit(amount float64)        // Пополнить баланс
	Withdraw(amount float64) error // Снять деньги (может вернуть ошибку)
	GetBalance() float64           // Получить текущий баланс
}

// Структура пользователя
type User struct {
	Id      string  // Уникальный идентификатор
	Name    string  // Имя пользователя
	Balance float64 // Текущий баланс
}

// Получить текущий баланс
func (u *User) GetBalance() float64 {
	return u.Balance
}

// Пополнить баланс
func (u *User) Deposit(amount float64) {
	u.Balance += amount // Увеличиваем баланс
}

// Снять деньги со счета
func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return errors.New("Недостаточно средств") // Проверка на достаточность средств
	}
	u.Balance -= amount // Уменьшаем баланс
	return nil
}

// Обработка операций со счетом
func processAccount(account BankAccount, depositAmt float64, withdrawAmt float64) {
	// Выводим начальный баланс
	fmt.Print("Начальный баланс: ", account.GetBalance(), "\n")

	// Пополняем счет
	account.Deposit(depositAmt)
	fmt.Print("После пополнения на ", depositAmt, ": ", account.GetBalance(), "\n")

	// Пытаемся снять деньги
	err := account.Withdraw(withdrawAmt)
	if err != nil {
		fmt.Print("Ошибка: ", err, "\n") // Выводим ошибку если не хватает средств
	} else {
		fmt.Print("После снятия ", withdrawAmt, ": ", account.GetBalance(), "\n")
	}
}

func main() {
	// Создаем двух пользователей
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

	// Приводим пользователей к интерфейсу BankAccount
	var account1 BankAccount = user1
	var account2 BankAccount = user2

	// Выполняем операции с первым пользователем
	fmt.Print("---------------------------\n")
	fmt.Print("Операция с: ", user1.Name, "\n")
	processAccount(account1, 300, 400)

	// Выполняем операции со вторым пользователем
	fmt.Print("===========================\n")
	fmt.Print("Операция с: ", user2.Name, "\n")
	processAccount(account2, 200, 300)
	fmt.Print("---------------------------\n")

	// Выводим итоговые балансы
	fmt.Print("Итоговый баланс пользователя: ", user1.Balance, " ", user2.Balance, "\n")
	fmt.Print("---------------------------\n")
}
