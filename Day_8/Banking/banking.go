package main

import (
	"fmt"
)

type BankAcc struct {
	Owner   string
	balance float64
}

func (b BankAcc) DispBalance() {
	fmt.Printf("Owner: %s ,Balance : %f\n", b.Owner, b.balance)
}

func (b *BankAcc) Deposit(amount float64) {
	if amount > 0 {
		b.balance += amount
		fmt.Printf("Deposit amount : %f . New Balance : %f \n", amount, b.balance)
	} else {
		fmt.Println("Deposit must be positive")
	}
}
func (b *BankAcc) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Println("withdraw amount must be positive")
		return
	}
	if amount <= b.balance {
		b.balance -= amount
		fmt.Printf("Withdraw amount : %f,New Balance: %f \n", amount, b.balance)
	} else {
		fmt.Println("No Sufficient balance")
	}
}

func main() {
	account := BankAcc{Owner: "Pankaj", balance: 10000}
	account.DispBalance()
	account.Deposit(1000)
	account.Withdraw(500)
	account.DispBalance()

	CpAccount := account
	CpAccount.Deposit(500)

	fmt.Println("\n Final balance with original account:")
	account.DispBalance()
	fmt.Println("Balance of copied account after deposite:")
	CpAccount.DispBalance()
}
