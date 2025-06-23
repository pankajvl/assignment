package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func captureOutput(f func()) string {

	old := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		panic(fmt.Sprintf("failed to create pipe: %v", err))
	}
	os.Stdout = w

	f()

	err = w.Close()
	if err != nil {
		panic(fmt.Sprintf("failed to close pipe writer: %v", err))
	}
	os.Stdout = old

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		panic(fmt.Sprintf("failed to read from pipe reader: %v", err))
	}
	return buf.String()
}

func TestDeposit(t *testing.T) {
	acc := &BankAcc{Owner: "Alice", balance: 100}
	acc.Deposit(50)

	if acc.balance != 150 {
		t.Errorf("expected balance 150, got %f", acc.balance)
	}

	acc.Deposit(-10)
	if acc.balance != 150 {
		t.Errorf("balance should remain unchanged after negative deposit, got %f", acc.balance)
	}
}

func TestWithdraw(t *testing.T) {
	acc := &BankAcc{Owner: "Bob", balance: 200}

	acc.Withdraw(50)
	if acc.balance != 150 {
		t.Errorf("expected balance 150, got %f", acc.balance)
	}

	acc.Withdraw(300)
	if acc.balance != 150 {
		t.Errorf("balance should remain unchanged after overdraft attempt, got %f", acc.balance)
	}

	acc.Withdraw(-10)
	if acc.balance != 150 {
		t.Errorf("balance should remain unchanged after negative withdrawal, got %f", acc.balance)
	}
}

func TestDispBalance(t *testing.T) {
	acc := &BankAcc{Owner: "Charlie", balance: 300}
	output := captureOutput(func() {
		acc.DispBalance()
	})

	expected := fmt.Sprintf("Owner: %s ,Balance : %f\n", acc.Owner, acc.balance)
	if output != expected {
		t.Errorf("expected output %q, got %q", expected, output)
	}
}

func TestCopyBehavior(t *testing.T) {
	original := BankAcc{Owner: "David", balance: 500}
	copyAcc := original

	copyAcc.Deposit(100)

	if original.balance != 500 {
		t.Errorf("original balance should remain unchanged, got %f", original.balance)
	}

	if copyAcc.balance != 600 {
		t.Errorf("copy balance expected 600, got %f", copyAcc.balance)
	}
}
