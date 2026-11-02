package service

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]Session
}
type Session struct {
	UserID    int
	ExpiresAt time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]Session),
	}
}

func (s *SessionStore) Create(userID int) (sessionID string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}
	sessionID = hex.EncodeToString(b)
	var session Session
	session.UserID = userID
	session.ExpiresAt = time.Now().Add(10 * time.Minute)
	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()
	return sessionID, nil
}

func (s *SessionStore) Delete(sessionID string) error {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
	return nil
}

func (s *SessionStore) Get(sessionID string) (userID int, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[sessionID]
	if ok && !time.Now().Before(session.ExpiresAt) {
		ok = false
		delete(s.sessions, sessionID)
	}
	if !ok {
		return 0, false
	}
	return session.UserID, ok
}
