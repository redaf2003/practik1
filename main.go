package main

import (
	"errors"
	"fmt"
	"sync"
)

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

type User struct {
	ID      string
	Name    string
	Balance float64
	mu      sync.Mutex
}

type PaymentTransaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentProcessor struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

func (u *User) Deposit(amount float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.Balance < amount {
		return errors.New("insufficient funds")
	}
	u.Balance -= amount
	return nil
}

func (u *User) GetBalance() float64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.Balance
}

func (t *PaymentTransaction) GetFromID() string  { return t.FromID }
func (t *PaymentTransaction) GetToID() string    { return t.ToID }
func (t *PaymentTransaction) GetAmount() float64 { return t.Amount }

func (pp *PaymentProcessor) AddUser(user *User) {
	if pp.Users == nil {
		pp.Users = make(map[string]*User)
	}
	pp.Users[user.ID] = user
}

func (pp *PaymentProcessor) AddTransaction(t Transaction) {
	pp.TransactionQueue = append(pp.TransactionQueue, t)
}

func (pp *PaymentProcessor) ProcessingTransactions() error {
	ch := make(chan Transaction, len(pp.TransactionQueue))
	var wg sync.WaitGroup

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

	for _, tx := range pp.TransactionQueue {
		ch <- tx
	}
	close(ch)

	wg.Wait()
	pp.TransactionQueue = nil
	return nil
}

func (pp *PaymentProcessor) GetUserBalance(id string) (float64, error) {
	user, exists := pp.Users[id]
	if !exists {
		return 0, fmt.Errorf("user not found: %s", id)
	}
	return user.GetBalance(), nil
}

func main() {
	ps := &PaymentProcessor{
		Users: make(map[string]*User),
	}

	fmt.Print("Создаю UserID: 1 с балансом 1000\n")
	ps.AddUser(&User{ID: "1", Name: "Артем", Balance: 1000})

	fmt.Print("Создаю UserID: 2 с балансом 500\n")
	ps.AddUser(&User{ID: "2", Name: "Егор", Balance: 500})

	fmt.Print("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200\n")
	ps.AddTransaction(&PaymentTransaction{FromID: "1", ToID: "2", Amount: 200})

	fmt.Print("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50\n")
	ps.AddTransaction(&PaymentTransaction{FromID: "2", ToID: "1", Amount: 50})

	if err := ps.ProcessingTransactions(); err != nil {
		fmt.Print("Ошибка обработки транзакций:", err, "\n")
		return
	}

	fmt.Print("Итого\n")
	if balance, err := ps.GetUserBalance("1"); err == nil {
		fmt.Print("У Артема получилось: ", balance, "\n")
	}
	if balance, err := ps.GetUserBalance("2"); err == nil {
		fmt.Print("У Егора получилось: ", balance, "\n")
	}
}
