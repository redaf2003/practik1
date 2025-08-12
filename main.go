package main

import (
	"errors"
	"fmt"
	"sync"
)

// Интерфейсы остаются без изменений
type BankAccount interface {
	Deposit(amount float64)
	Withdraw(amount float64) error
	GetBalance() float64
}

type Transaction interface {
	GetFromID() string
	GetToID() string
	GetAmount() float64
}

type PaymentSystem interface {
	AddUser(user *User)
	AddTransaction(t Transaction)
	ProcessingTransactions() error
	GetUserBalance(id string) (float64, error)
}

// Структура пользователя
type User struct {
	ID      string
	Name    string
	Balance float64
}

// Структура транзакции
type PaymentTransaction struct {
	FromID string
	ToID   string
	Amount float64
}

// Платежный процессор
type PaymentProcessor struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

func (pp *PaymentProcessor) ProcessingTransactions() error {
	ch := make(chan Transaction, len(pp.TransactionQueue))
	var wg sync.WaitGroup

	// Запускаем 3 воркера
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range ch {
				fromUser := pp.Users[t.GetFromID()]
				toUser := pp.Users[t.GetToID()]

				if fromUser == nil || toUser == nil {
					fmt.Print("Ошибка: пользователь не найден (tx ", t.GetFromID(), " -> ", t.GetToID(), ")\n")
					continue
				}

				if err := fromUser.Withdraw(t.GetAmount()); err != nil {
					fmt.Print("Ошибка транзакции: ", err, "\n")
					continue
				}

				toUser.Deposit(t.GetAmount())
			}
		}()
	}

	// Отправляем транзакции
	for _, tx := range pp.TransactionQueue {
		ch <- tx
	}
	close(ch)

	wg.Wait()
	pp.TransactionQueue = nil
	return nil
}

// Методы пользователя
func (u *User) Deposit(amount float64) {
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return errors.New("insufficient funds")
	}
	u.Balance -= amount
	return nil
}

func (u *User) GetBalance() float64 {
	return u.Balance
}

// Методы транзакции
func (t *PaymentTransaction) GetFromID() string  { return t.FromID }
func (t *PaymentTransaction) GetToID() string    { return t.ToID }
func (t *PaymentTransaction) GetAmount() float64 { return t.Amount }

// Методы платежного процессора
func (pp *PaymentProcessor) AddUser(user *User) {
	if pp.Users == nil {
		pp.Users = make(map[string]*User)
	}
	pp.Users[user.ID] = user
}

func (pp *PaymentProcessor) AddTransaction(t Transaction) {
	pp.TransactionQueue = append(pp.TransactionQueue, t)
}

func (pp *PaymentProcessor) GetUserBalance(id string) (float64, error) {
	user, exists := pp.Users[id]
	if !exists {
		return 0, fmt.Errorf("user not found: %s", id)
	}
	return user.Balance, nil
}

func main() {
	var ps PaymentSystem = &PaymentProcessor{
		Users: make(map[string]*User),
	}

	// Добавляем пользователей
	fmt.Print("Создаю UserID: 1 с балансом 1000\n")
	ps.AddUser(&User{ID: "1", Name: "Артем", Balance: 300})

	fmt.Print("Создаю UserID: 2 с балансом 500\n")
	ps.AddUser(&User{ID: "2", Name: "Егор", Balance: 200})

	// Добавляем транзакции
	fmt.Print("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200\n")
	ps.AddTransaction(&PaymentTransaction{FromID: "1", ToID: "2", Amount: 200})

	fmt.Print("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50\n")
	ps.AddTransaction(&PaymentTransaction{FromID: "2", ToID: "1", Amount: 50})

	// Обработка транзакций
	if err := ps.ProcessingTransactions(); err != nil {
		fmt.Print("Ошибка обработки транзакций: ", err, "\n")
		return
	}

	// Вывод результатов
	fmt.Print("Итого\n")
	if balance, err := ps.GetUserBalance("1"); err == nil {
		fmt.Print("У Артема получилось: ", balance, "\n")
	}
	if balance, err := ps.GetUserBalance("2"); err == nil {
		fmt.Print("У Егора получилось: ", balance, "\n")
	}
}
