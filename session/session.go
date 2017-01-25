package session

import (
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

//TODO secure needs to be true for https
const (
	//8 hours in seconds
	DEFAULT_SESSION_AGE    int    = 28800
	DEFAULT_SESSION_SECURE bool   = false
	SESSION_KEY            string = "session"
)

type ISession interface {
	Set(string, string)
	Get(string) string
	Remove()
	GetCookie() *http.Cookie
}

type Session struct {
	expires time.Time
	values  map[string]string
	cookie  *http.Cookie
}

var sessionMap map[string]*Session
var lock *sync.RWMutex

func init() {
	lock = new(sync.RWMutex)
	sessionMap = make(map[string]*Session)
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		select {
		case <-ticker.C:
			cleanSessionMap()
		}
	}()
}

func cleanSessionMap() {
	var keys []string = make([]string, 0)
	lock.RLock()
	for k, v := range sessionMap {
		if v.isExpired() {
			keys = append(keys, k)
		}
	}
	lock.RUnlock()
	lock.Lock()
	defer lock.Unlock()
	for _, k := range keys {
		delete(sessionMap, k)
	}
}

func (s *Session) Remove() {
	key := s.cookie.Value
	if key == "" {
		return
	}
	lock.RLock()
	if _, found := sessionMap[key]; !found {
		return
	}
	lock.RUnlock()
	lock.Lock()
	defer lock.Unlock()
	delete(sessionMap, key)
}

//age in seconds
func NewSession(cookies []*http.Cookie, age int) (*Session, error) {
	var found *http.Cookie
	for _, c := range cookies {
		if c.Name == SESSION_KEY {
			if found == nil {
				found = c
			} else {
				return nil, fmt.Errorf("More than one cookie found with name '%s'.", SESSION_KEY)
			}
		}
	}
	var key string
	var oldValues map[string]string = make(map[string]string)
	if found != nil {
		if _, ok := sessionMap[found.Value]; ok {
			key = found.Value
			lock.RLock()
			oldValues = sessionMap[key].values
			lock.RUnlock()
		}
	}
	toRet := newSession(age, key)
	toRet.values = oldValues
	return toRet, nil
}

func (s *Session) GetCookie() *http.Cookie {
	return s.cookie
}

func (s *Session) Get(k string) string {
	if s.isExpired() {
		lock.Lock()
		delete(sessionMap, s.cookie.Name)
		lock.Unlock()
		return ""
	}
	return s.values[k]
}

func (s *Session) Set(k, v string) {
	//Don't do anything if expired
	if s.isExpired() {
		lock.Lock()
		delete(sessionMap, s.cookie.Name)
		lock.Unlock()
		return
	}
	s.values[k] = v
}

func (s *Session) isExpired() bool {
	return time.Now().UTC().After(s.expires)
}

func newSession(age int, name string) *Session {
	var key string
	if name == "" {
		key = getRandKey()
	} else {
		key = name
	}
	expires := getExpires(age)
	toRet := &Session{
		expires: expires,
		values:  make(map[string]string),
		cookie:  createCookie(key, age, expires),
	}
	lock.Lock()
	sessionMap[key] = toRet
	lock.Unlock()
	return toRet
}

func createCookie(value string, age int, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Expires:    expires,
		RawExpires: getRawExpires(age),
		MaxAge:     age,
		Secure:     DEFAULT_SESSION_SECURE,
		HttpOnly:   true,
		Name:       SESSION_KEY,
		Value:      value,
	}
}

func getExpires(age int) time.Time {
	if age <= 0 {
		return time.Now().UTC()
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", age))
	if err != nil {
		return time.Now().UTC()
	}
	return time.Now().UTC().Add(duration)
}

func getRawExpires(age int) string {
	return getExpires(age).Format(time.RFC1123)
}

func getRandKey() string {
	lock.RLock()
	defer lock.RUnlock()
	for {
		s := make([]byte, 64)
		rand.Read(s)
		key := b64.StdEncoding.EncodeToString(s)
		if _, ok := sessionMap[key]; !ok {
			return key
		}
	}
}
