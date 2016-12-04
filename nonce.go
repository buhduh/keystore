package main

import (
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"sync"
	"time"
)

const (
	CLEANUP_INT    int = 5
	NONCE_LIFETIME     = 3600 //1 hour
)

type NonceChecker interface {
	CheckNonce() bool
}

type Nonce string

//value is expiration
var nonceMap map[Nonce]time.Time

var nonceLock *sync.Mutex

func init() {
	nonceLock = new(sync.Mutex)
	nonceMap = make(map[Nonce]time.Time)
	ticker := time.NewTicker(time.Duration(CLEANUP_INT) * time.Minute)
	//swap these two lines to run TestCleanUp
	//ticker := time.NewTicker(time.Duration(1) * time.Second)
	go func() {
		select {
		case <-ticker.C:
			cleanNonceMap()
		}
	}()
}

//Keep the nonce map from growing without bounds
func cleanNonceMap() {
	nonceLock.Lock()
	defer nonceLock.Unlock()
	now := time.Now().UTC()
	toDelete := make([]Nonce, 0)
	for k, v := range nonceMap {
		if v.Before(now) {
			toDelete = append(toDelete, k)
		}
	}
	for _, k := range toDelete {
		delete(nonceMap, k)
	}
}

//lifetime in seconds
func NewNonce(lifetime int) Nonce {
	nonceLock.Lock()
	defer nonceLock.Unlock()
	var s []byte
	var toRet Nonce
	for {
		s = make([]byte, 128)
		//keep spinning if a failure, unsure of the consequences
		if _, err := rand.Read(s); err != nil {
			continue
		}
		toRet = Nonce(b64.StdEncoding.EncodeToString(s))
		//keep spinning until nonce is unique
		if _, found := nonceMap[toRet]; found {
			continue
		}
		break
	}
	durr, _ := time.ParseDuration(fmt.Sprintf("%ds", lifetime))
	nonceMap[toRet] = time.Now().UTC().Add(durr)
	return toRet
}

func (n Nonce) CheckNonce() bool {
	nonceLock.Lock()
	defer nonceLock.Unlock()
	var expires time.Time
	var found bool
	if expires, found = nonceMap[n]; !found {
		return false
	}
	if expires.Before(time.Now().UTC()) {
		return false
	}
	delete(nonceMap, n)
	return true
}
