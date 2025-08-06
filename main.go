package main

import (
	"errors"
	"fmt"
)

// Интерфейсы для банковского аккаунта, транзакции и платежной системы
type BankAccount interface {
	Deposit(amount float64)
	Withdraw(amount float64) error
	GetBalance() float64
}

type Transaction interface {
	GetFromID() string  // ID отправителя
	GetToID() string    // ID получателя
	GetAmount() float64 // сумма перевода
}

type PaymentSystem interface {
	AddUser(user *User)                        // добавить пользователя
	AddTransaction(t Transaction)              // добавить транзакцию в очередь
	ProcessingTransactions() error             // обработать все транзакции
	GetUserBalance(id string) (float64, error) // получить баланс пользователя
}

// Структура пользователя
type User struct {
	ID      string  // уникальный идентификатор
	Name    string  // имя пользователя
	Balance float64 // текущий баланс
}

// Структура транзакции
type PaymentTransaction struct {
	FromID string  // ID отправителя
	ToID   string  // ID получателя
	Amount float64 // сумма перевода
}

// Платежный процессор (основная логика)
type PaymentProcessor struct {
	Users            map[string]*User // хранилище пользователей
	TransactionQueue []Transaction    // очередь транзакций
}

// Методы пользователя
func (u *User) Deposit(amount float64) {
	u.Balance += amount // пополнение баланса
}

func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return errors.New("insufficient funds") // ошибка если недостаточно средств
	}
	u.Balance -= amount // списание средств
	return nil
}

func (u *User) GetBalance() float64 {
	return u.Balance // получение текущего баланса
}

// Методы транзакции
func (t *PaymentTransaction) GetFromID() string {
	return t.FromID
}

func (t *PaymentTransaction) GetToID() string {
	return t.ToID
}

func (t *PaymentTransaction) GetAmount() float64 {
	return t.Amount
}

// Методы платежного процессора
func (pp *PaymentProcessor) AddUser(user *User) {
	if pp.Users == nil {
		pp.Users = make(map[string]*User) // инициализация мапы при первом добавлении
	}
	pp.Users[user.ID] = user // добавление пользователя
}

func (pp *PaymentProcessor) AddTransaction(t Transaction) {
	pp.TransactionQueue = append(pp.TransactionQueue, t) // добавление транзакции в очередь
}

func (pp *PaymentProcessor) ProcessingTransactions() error {
	// Обработка всех транзакций в очереди
	for _, t := range pp.TransactionQueue {
		fromUser, exists := pp.Users[t.GetFromID()]
		if !exists {
			return fmt.Errorf("sender not found: %s", t.GetFromID()) // проверка отправителя
		}

		toUser, exists := pp.Users[t.GetToID()]
		if !exists {
			return fmt.Errorf("recipient not found: %s", t.GetToID()) // проверка получателя
		}

		if err := fromUser.Withdraw(t.GetAmount()); err != nil {
			return fmt.Errorf("transaction failed: %v", err) // списание средств
		}

		toUser.Deposit(t.GetAmount()) // зачисление средств
	}

	pp.TransactionQueue = nil // очистка очереди после обработки
	return nil
}

func (pp *PaymentProcessor) GetUserBalance(id string) (float64, error) {
	user, exists := pp.Users[id]
	if !exists {
		return 0, fmt.Errorf("user not found: %s", id) // проверка существования пользователя
	}
	return user.Balance, nil // возврат баланса
}

func main() {
	var ps PaymentSystem = &PaymentProcessor{} // создание платежной системы

	// Демонстрация работы системы
	fmt.Println("Создаю UserID: 1 с балансом 1000")
	user1 := &User{ID: "1", Name: "User1", Balance: 1000}
	ps.AddUser(user1)

	fmt.Println("Создаю UserID: 2 с балансом 500")
	user2 := &User{ID: "2", Name: "User2", Balance: 500}
	ps.AddUser(user2)

	fmt.Println("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200")
	ps.AddTransaction(&PaymentTransaction{FromID: "1", ToID: "2", Amount: 200})

	fmt.Println("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50")
	ps.AddTransaction(&PaymentTransaction{FromID: "2", ToID: "1", Amount: 50})

	// Обработка транзакций
	if err := ps.ProcessingTransactions(); err != nil {
		fmt.Println("Ошибка обработки транзакций:", err)
		return
	}

	// Вывод результатов
	fmt.Println("Итого")
	if balance, err := ps.GetUserBalance("1"); err == nil {
		fmt.Printf("У первого пользователя получилось: %.2f\n", balance)
	}
	if balance, err := ps.GetUserBalance("2"); err == nil {
		fmt.Printf("У второго пользователя получилось: %.2f\n", balance)
	}
}
