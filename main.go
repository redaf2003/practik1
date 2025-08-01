package main

import (
	"errors"
	"fmt"
)

// Transaction - информация о денежном переводе.
// Содержит:
//
//	FromID - ID отправителя (кто переводит)
//	ToID - ID получателя (кому переводят)
//	Amount - сумма перевода (в рублях/долларах)
type Transaction struct {
	FromID string  // ID отправителя
	ToID   string  // ID получателя
	Amount float64 // Сумма перевода (например: 150.50)
}

// PaymentSystem - центральная система переводов.
// Содержит:
//
//	Users - всех зарегистрированных пользователей (ключ: ID, значение: данные пользователя)
//	Transactions - полную историю всех платежей (список операций)
type PaymentSystem struct {
	Users        map[string]*User // База пользователей (ID → данные)
	Transactions []Transaction    // История всех платежных операций
}

// Добавления пользователей в систему (берет ID пользователя и добавляет его в мапу).
func (ps *PaymentSystem) AddUser(u User) {
	ps.Users[u.Id] = &u // Добавили нового пользователя в мапу
}

// Добавления транзакций в очередь (добавлие элемента в слайс)
func (ps *PaymentSystem) AddTransaction(t Transaction) {
	ps.Transactions = append(ps.Transactions, t) // Добавили транзакцию в слайс
}

func (ps *PaymentSystem) ProcessingTransactions() []error {
	var errorsList []error

	// Проходим по всем транзакциям
	for _, transaction := range ps.Transactions {
		fromUser, fromExists := ps.Users[transaction.FromID]
		toUser, toExists := ps.Users[transaction.ToID]

		// Проверка на существование отправителя
		if !fromExists {
			errorsList = append(errorsList, errors.New("пользователь с ID "+transaction.FromID+" не найден"))
			continue
		}

		// Проверка на существование получателя
		if !toExists {
			errorsList = append(errorsList, errors.New("пользователь с ID "+transaction.ToID+" не найден"))
			continue
		}

		// Проверка на успешность снятия средств
		if err := fromUser.Withdraw(transaction.Amount); err != nil {
			errorsList = append(errorsList, err)
			continue
		}

		// Пополнение счета получателя
		toUser.Deposit(transaction.Amount)
	}

	// Очищаем список транзакций после обработки
	ps.Transactions = nil
	return errorsList
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
func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return errors.New("Баланс меньше чем сумма перевода")
	}
	u.Balance -= amount
	return nil
}

func main() {
	// Создаем экземпляр PaymentSystem
	paymentSystem := PaymentSystem{
		Users:        make(map[string]*User),
		Transactions: []Transaction{},
	}

	// Создаем пользователей
	user1 := User{Id: "1", Name: "Тёмочка", Balance: 100.0}
	user2 := User{Id: "2", Name: "Бро", Balance: 50.0}
	user3 := User{Id: "3", Name: "Маша", Balance: 200.0}

	// Добавляем пользователей в систему
	paymentSystem.AddUser(user1)
	paymentSystem.AddUser(user2)
	paymentSystem.AddUser(user3)

	// Добавляем несколько транзакций
	paymentSystem.AddTransaction(Transaction{FromID: "1", ToID: "2", Amount: 30})
	paymentSystem.AddTransaction(Transaction{FromID: "2", ToID: "3", Amount: 20})
	paymentSystem.AddTransaction(Transaction{FromID: "3", ToID: "1", Amount: 50})

	// Обрабатываем транзакции
	if errList := paymentSystem.ProcessingTransactions(); len(errList) > 0 {
		for _, err := range errList {
			fmt.Println("Ошибка:", err)
		}
	} else {
		fmt.Println("Транзакции успешно обработаны!")
	}

	// Проверяем балансы
	fmt.Println("Баланс Тёмочки:", paymentSystem.Users["1"].Balance)
	fmt.Println("Баланс Бро:", paymentSystem.Users["2"].Balance)
	fmt.Println("Баланс Маши:", paymentSystem.Users["3"].Balance)
}
