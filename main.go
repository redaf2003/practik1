package main

import (
	"errors"
)

// Структура User: ID (строка), Name (строка), Balance (число).
type User struct {
	Id      string
	Name    string
	Balance float64
}

// Реализуйте метод пополнения счета для пользователя, увеличивающий баланс на указанную сумму.
func (u *User) Deposit(amount float64) {
	u.Balance += amount
}

// Реализуйте метод Withdraw для структуры User, который уменьшает баланс на указанную сумму.
// Проверяйте достаточность средств, возвращая ошибку при нехватке.
func (u *User) Withdraw(amounts float64) error {
	if amounts <= 0 {
		return errors.New("Cумма должна быть положительной")
	}
	if amounts > u.Balance {
		return errors.New("недостаточно средств")
	}
	u.Balance -= amounts

	return nil
}

func main() {

}
