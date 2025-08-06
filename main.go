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
		return errors.New("Недостаточно средств ") // Проверка на достаточность средств
	}
	u.Balance -= amount // Уменьшаем баланс
	return nil
}

// Обработка операций со счетом
func processAccount(account BankAccount, depositAmt float64, withdrawAmt float64) {
	// Выводим начальный баланс
	fmt.Printf("Начальный баланс: %.2f\n", account.GetBalance())

	// Пополняем счет
	account.Deposit(depositAmt)
	fmt.Printf("После пополнения на %.2f: %.2f\n", depositAmt, account.GetBalance())

	// Пытаемся снять деньги
	err := account.Withdraw(withdrawAmt)
	if err != nil {
		fmt.Println("Ошибка:", err) // Выводим ошибку если не хватает средств
	} else {
		fmt.Printf("После снятия %.2f: %.2f\n", withdrawAmt, account.GetBalance())
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
	fmt.Println("---------------------------")
	fmt.Println("Операция с:", user1.Name)
	processAccount(account1, 300, 400)

	// Выполняем операции со вторым пользователем
	fmt.Println("===========================")
	fmt.Println("Операция с:", user2.Name)
	processAccount(account2, 200, 300)
	fmt.Println("---------------------------")

	// Выводим итоговые балансы
	fmt.Println("Итоговый баланс пользователя:", user1.Balance, user2.Balance)
	fmt.Println("---------------------------")
}
