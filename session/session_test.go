package session

import (
	"net/http"
	"testing"
	"time"
)

func TestSessionExpiration(t *testing.T) {
	//The clock is ticking!!!!Will test expiration in one second
	var key, val string = "foo", "bar"
	session, err := NewSession(make([]*http.Cookie, 0), 1)
	if err != nil {
		t.Logf("Should not get err, got '%s'.", err)
		t.FailNow()
	}
	if session == nil {
		t.Logf("session should not be nil.")
		t.FailNow()
	}
	if session.isExpired() {
		t.Logf("session should not be expired.")
		t.Fail()
	}
	session.Set(key, val)
	if session.Get(key) != val {
		t.Logf("Expected '%s', got '%s'.", val, session.Get(key))
		t.Fail()
	}
	if session.Get("blarg") != "" {
		t.Logf("Get should return empty string.")
		t.Fail()
	}
	time.Sleep(1 * time.Second)
	if !session.isExpired() {
		t.Logf("session should be expired.")
		t.Fail()
	}
	if session.Get(key) != "" {
		t.Logf("Expected empty string, got '%s'.", session.Get(key))
		t.Fail()
	}
	if _, ok := sessionMap[session.cookie.Name]; ok {
		t.Logf("Session should have been removed from map from get or set.")
		t.Fail()
	}
}

func Test2Sessions(t *testing.T) {
	var key, val1, val2 string = "foo", "bar", "baz"
	sess1, err := NewSession(make([]*http.Cookie, 0), 1)
	sess1.Set(key, val1)
	if err != nil {
		t.Logf("Should not get err, got '%s'.", err)
		t.FailNow()
	}
	if sess1 == nil {
		t.Logf("session should not be nil.")
		t.FailNow()
	}
	cName := "yoyoyo"
	expires := time.Now().UTC().Add(1 * time.Second)
	cookie := createCookie(cName, 1, expires)
	sessionMap[cName] = &Session{
		expires: expires,
		values:  map[string]string{key: val2},
		cookie:  cookie,
	}
	sess2, err := NewSession([]*http.Cookie{cookie}, 1)
	if err != nil {
		t.Logf("Should not get err, got '%s'.", err)
		t.FailNow()
	}
	if sess2 == nil {
		t.Logf("session should not be nil.")
		t.FailNow()
	}
	if sess2.Get(key) != val2 {
		t.Logf("Expected '%s', got '%s'.", val2, sess2.Get(key))
		t.Fail()
	}
	if sess1.Get(key) != val1 {
		t.Logf("Expected '%s', got '%s'.", val1, sess1.Get(key))
		t.Fail()
	}
	if sess2.cookie.Value != cName {
		t.Logf("Expected '%s', got '%s'.", cName, sess2.cookie.Value)
		t.Fail()
	}
}
