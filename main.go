package main

import (
	"errors"
	"fmt"
	"sync"
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
	TransactionQueue []Transaction    // очередь транзакции
	// добавили канал для транзакции
}

func (pp *PaymentProcessor) ProcessingTransactions() error {
	// Создаем канал и WaitGroup
	ch := make(chan Transaction, len(pp.TransactionQueue))
	var wg sync.WaitGroup

	// Запускаем 3 воркера
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go pp.Worker(ch, &wg)
	}

	// Отправляем транзакции в канал
	for _, tx := range pp.TransactionQueue {
		ch <- tx
	}
	close(ch)

	// Ждем завершения всех горутин
	wg.Wait()

	// Очищаем очередь
	pp.TransactionQueue = nil
	return nil
}

// Методы пользователя
func (u *User) Deposit(amount float64) {
	u.Balance += amount // пополнение баланса
}

func (pp *PaymentProcessor) Worker(ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	for t := range ch {
		fromUser, exists := pp.Users[t.GetFromID()]
		if !exists {
			fmt.Printf("Ошибка: отправитель не найден (tx %v -> %v)\n", t.GetFromID(), t.GetToID())
			continue
		}

		toUser, exists := pp.Users[t.GetToID()]
		if !exists {
			fmt.Printf("Ошибка: получатель не найден (tx %v -> %v)\n", t.GetFromID(), t.GetToID())
			continue
		}

		if err := fromUser.Withdraw(t.GetAmount()); err != nil {
			fmt.Printf("Ошибка транзакции: %v\n", err)
			continue
		}

		toUser.Deposit(t.GetAmount())
	}
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

func (pp *PaymentProcessor) GetUserBalance(id string) (float64, error) {
	user, exists := pp.Users[id]
	if !exists {
		return 0, fmt.Errorf("user not found: %s", id) // проверка существования пользователя
	}
	return user.Balance, nil // возврат баланса
}

func main() {
	var ps PaymentSystem = &PaymentProcessor{
		Users: make(map[string]*User),
	}

	// Добавляем пользователей
	fmt.Println("Создаю UserID: 1 с балансом 1000")
	user1 := &User{ID: "1", Name: "Егор", Balance: 1000}
	ps.AddUser(user1)

	fmt.Println("Создаю UserID: 2 с балансом 500")
	user2 := &User{ID: "2", Name: "Артем", Balance: 500}
	ps.AddUser(user2)

	// Добавляем транзакции
	fmt.Println("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200")
	ps.AddTransaction(&PaymentTransaction{FromID: "1", ToID: "2", Amount: 200})

	fmt.Println("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50")
	ps.AddTransaction(&PaymentTransaction{FromID: "2", ToID: "1", Amount: 50})

	// Обрабатываем транзакции
	if err := ps.ProcessingTransactions(); err != nil {
		fmt.Println("Ошибка обработки транзакций:", err)
		return
	}

	fmt.Println("Итого")
	if balance, err := ps.GetUserBalance("1"); err == nil {
		fmt.Printf("У первого пользователя получилось: %.2f\n", balance)
	}
	if balance, err := ps.GetUserBalance("2"); err == nil {
		fmt.Printf("У второго пользователя получилось: %.2f\n", balance)
	}

}
