package controllers

import (
	"sync"
	"time"
)

type MemProvider struct {
	lock        sync.RWMutex // ReadWrite Mutex
	s_map       map[string]interface{}
	maxlifetime int64
}

type MemSession struct {
	sid          string                      // unique session id
	timeAccessed time.Time                   // last access time
	sessions     map[interface{}]interface{} // session value stored inside
	lock         sync.RWMutex                // ReadWrite Mutex
}

//set session value
func (session *MemSession) Set(key, value interface{}) {
	session.lock.Lock() // Lock locks for writing.
	defer session.lock.Unlock()
	session.sessions[key] = value
}

//get session value
func (session *MemSession) Get(key interface{}) interface{} {
	session.lock.RLock()
	defer session.lock.RUnlock()
	if v, exist := session.sessions[key]; exist {
		return v
	}
	return ""
}

//delete session value
func (session *MemSession) Delete(key interface{}) {
	session.lock.Lock() // Lock locks for writing.
	defer session.lock.Unlock()
	delete(session.sessions, key)
}

//back current sessionID
func (session *MemSession) SessionID() string {
	return session.sid
}

// implements initialization of session, it returns new session
// variable if it succeed.
func (p *MemProvider) SessionInit(s_id string) Session {
	p.lock.Lock() // Lock locks for writing.
	defer p.lock.Unlock()

	s_map := make(map[interface{}]interface{})

	newSession := &MemSession{sid: s_id,
		timeAccessed: time.Now(),
		sessions:     s_map}

	p.s_map[s_id] = newSession

	return newSession
}

// returns session variable that is represented by corresponding
// session_id, it creates a new session variable and return if it does not exist.
func (p *MemProvider) SessionRead(sid string) Session {
	p.lock.RLock()
	if _, exist := p.s_map[sid]; exist {
		go p.SessionUpdate(sid)
		p.lock.RUnlock()
		return p.s_map[sid].(Session)
	}
	p.lock.RUnlock()
	return p.SessionInit(sid)
}

func (p *MemProvider) SessionExist(sid string) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if _, exist := p.s_map[sid]; exist {
		return true
	} else {
		return false
	}
}

// deletes session variable by corresponding session_id.
func (p *MemProvider) SessionDestroy(sid string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exist := p.s_map[sid]; exist {
		delete(p.s_map, sid)
	}
}

// deletes expired session variables according to maxLifeTime
func (p *MemProvider) SessionGC() {
	p.lock.RLock()
	for key, value := range p.s_map {
		if (value.(*MemSession).timeAccessed.Unix() + p.maxlifetime) < time.Now().Unix() {
			p.lock.RUnlock()
			p.lock.Lock()
			delete(p.s_map, key)
			p.lock.Unlock()
			p.lock.RLock()
		}
	}
	p.lock.RUnlock()
}

func (p *MemProvider) SessionUpdate(sid string) {
	p.lock.Lock() // Lock locks for writing.
	defer p.lock.Unlock()
	if session, exist := p.s_map[sid]; exist {
		session.(*MemSession).timeAccessed = time.Now()
	}
}

func NewMemProvider() *MemProvider {
	return &MemProvider{s_map: make(map[string]interface{}),
		maxlifetime: 604800}
}
