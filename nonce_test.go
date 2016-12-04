package main

import (
	"testing"
	"time"
)

//Test is skipped usually, requires manually
//modifying values
func TestCleanUp(t *testing.T) {
	t.Skip()
	nonce := NewNonce(1)
	time.Sleep(2 * time.Second)
	if nonce.CheckNonce() {
		t.Logf("CheckNonce() should return false")
		t.FailNow()
	}
}

func TestNewNonce(t *testing.T) {
	nonce := NewNonce(60)
	if !nonce.CheckNonce() {
		t.Logf("CheckNonce() should return true")
		t.FailNow()
	}
	if nonce.CheckNonce() {
		t.Logf("CheckNonce() should return false")
		t.FailNow()
	}
	nonce = NewNonce(1)
	time.Sleep(2 * time.Second)
	if nonce.CheckNonce() {
		t.Logf("CheckNonce() should return false")
		t.FailNow()
	}
}
