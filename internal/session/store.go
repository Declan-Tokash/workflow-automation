package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type Session struct {
	UserID       int
	GitHubLogin  string
	AccessToken  string
}

type Store struct {
	sessions map[string]Session
	mu       sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		sessions: make(map[string]Session),
	}
}

func (s *Store) Create(session Session) string {
	sessionIDBytes := make([]byte, 32)

	if _, err := rand.Read(sessionIDBytes); err != nil {
		panic(err)
	}

	sessionID := hex.EncodeToString(sessionIDBytes)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[sessionID] = session

	return sessionID
}

func (s *Store) Get(sessionID string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionID]

	return session, ok
}