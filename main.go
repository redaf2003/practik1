package main

import (
	"errors"
	"fmt"
)

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

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
	//Создайте несколько объектов пользователей с разными значениями баланса и именами.
	user1 := &User{
		Id:      "1",
		Name:    "Artem",
		Balance: 500.0,
	}

	user2 := &User{
		Id:      "2",
		Name:    "Egor",
		Balance: 300.0,
	}

	//Провели операции пополнения и снятия средств. После каждой операции выводите баланс.
	fmt.Println("\nНачальные балансы:")
	fmt.Printf("User1: %+v\n", user1)
	fmt.Printf("User2: %+v\n", user2)

	fmt.Println("\nПосле пополнения:")
	user1.Deposit(200)
	user2.Deposit(300)
	fmt.Printf("User1: %+v\n", user1)
	fmt.Printf("User2: %+v\n", user2)

	fmt.Println("\nПосле снятия:")
	if err := user1.Withdraw(1000); err != nil {
		fmt.Println("Ошибка:", err)
	}
	if err := user2.Withdraw(200); err != nil {
		fmt.Println("Ошибка:", err)
	}
	fmt.Printf("User1: %+v\n", user1)
	fmt.Printf("User2: %+v\n", user2)
}
