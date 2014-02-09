package controllers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionManager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxlifetime int64
}

type Provider interface {
	SessionInit(sid string) Session
	SessionRead(sid string) Session
	SessionDestroy(sid string)
	SessionGC()
	SessionUpdate(sid string)
	SessionExist(sid string) bool
}

type Session interface {
	Set(key, value interface{})      //set session value
	Get(key interface{}) interface{} //get session value
	Delete(key interface{})          //delete session value
	SessionID() string               //back current sessionID
}

var (
	// memory, file, database
	providers = make(map[string]Provider)
)

func InitializeSessionManager(providerName, cookieName string) (*SessionManager, error) {
	pder, exist := providers[providerName]
	if !exist || pder == nil {
		return nil, errors.New("session: unknown provider")
	}
	// 604800 => 60sec * 60min * 24h * 7days
	s_mng := &SessionManager{provider: pder,
		cookieName:  cookieName,
		maxlifetime: 604800}
	return s_mng, nil
}

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterProvider(p_name string, provider Provider) {
	if provider == nil {
		panic("session: provider is nil")
	}
	if _, dup := providers[p_name]; dup {
		panic("session: Register called twice for provide " + p_name)
	}
	providers[p_name] = provider
}

func (s_mng *SessionManager) newSessionId(r *http.Request) string {
	b := make([]byte, 32)
	io.ReadFull(rand.Reader, b)
	key := fmt.Sprintf("%s%d%s", r.RemoteAddr, time.Now().UnixNano(), b)
	hash := hmac.New(sha1.New, []byte("userKey"))
	fmt.Fprintf(hash, "%s", key)
	sid := hex.EncodeToString(hash.Sum(nil))
	return sid
}

func (s_mng *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) Session {
	cookie, err := r.Cookie(s_mng.cookieName)
	// lock session for data syncronization problems
	s_mng.lock.Lock()
	defer s_mng.lock.Unlock()

	if err != nil || cookie.Value == "" {
		sid := s_mng.newSessionId(r)
		cookie = &http.Cookie{Name: s_mng.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   0}
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
		return s_mng.provider.SessionRead(sid)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		if s_mng.provider.SessionExist(sid) {
			return s_mng.provider.SessionRead(sid)
		} else {
			sid = s_mng.newSessionId(r)
			cookie = &http.Cookie{Name: s_mng.cookieName,
				Value:    url.QueryEscape(sid),
				Path:     "/",
				HttpOnly: true,
				MaxAge:   0}
			http.SetCookie(w, cookie)
			r.AddCookie(cookie)
			return s_mng.provider.SessionInit(sid)
		}
	}
}

func (s_mng *SessionManager) GC() {
	s_mng.provider.SessionGC()
	time.AfterFunc(time.Duration(s_mng.maxlifetime)*time.Second, func() { s_mng.GC() })
}

func (s_mng *SessionManager) DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(s_mng.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		s_mng.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: s_mng.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}
