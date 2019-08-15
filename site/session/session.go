package session

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

type Config struct {
	SessionCookieName string
	SessionKey        string
	ExpirationSecs    int
	GCDelaySeconds    int
}

func NewSessionStore(cfg Config) *SessionStore {
	store := SessionStore{
		Config:       cfg,
		startedGc:    false,
		sessions:     make(map[string]*Session),
		idStack:      NewIdStack(),
		lock:         sync.Mutex{},
		timeProvider: systemTime{},
	}

	return &store
}

type SessionStore struct {
	Config       Config
	startedGc    bool
	sessions     map[string]*Session
	idStack      *IdStack
	lock         sync.Mutex
	timeProvider TimeProvider
}

func (c *Config) expiration() time.Duration {
	return time.Duration(c.ExpirationSecs) * time.Second
}

func (c *Config) gcDelay() time.Duration {
	return time.Duration(c.GCDelaySeconds) * time.Second
}

// Allows injection of mock timing for tests but intended to use traditional system time functions.
type TimeProvider interface {
	Now() time.Time
	Tick(delay time.Duration) <-chan time.Time
}

type systemTime struct {
}

func (p systemTime) Now() time.Time {
	return time.Now()
}

func (p systemTime) Tick(duration time.Duration) <-chan time.Time {
	return time.NewTicker(duration).C
}

type Session struct {
	Id        string
	StartTime time.Time
	Values    map[string]string
}

func (s *SessionStore) StartGC() error {
	if s.startedGc {
		return fmt.Errorf("GC already started")
	}

	s.startedGc = true

	go func() {
		for range s.timeProvider.Tick(s.Config.gcDelay()) {
			s.removeExpiredSessions()
		}
	}()

	return nil
}

func (s *SessionStore) NewOrExisting(sessionId string, userId string) (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	sess, ok := s.sessions[sessionId]
	if ok {
		_ = s.idStack.Remove(sessionId)
	} else {
		sess = &Session{
			Id:     sessionId,
			Values: make(map[string]string),
		}
	}

	sess.StartTime = s.timeProvider.Now()
	s.sessions[sessionId] = sess
	sess.Values["user_id"] = userId
	s.idStack.PushBack(sessionId)

	return sess, nil
}

func (s *SessionStore) Find(sessionId string) (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	existing, ok := s.sessions[sessionId]
	if !ok {
		return existing, fmt.Errorf("session not found")
	}

	_ = s.idStack.Remove(sessionId)
	s.idStack.PushBack(sessionId)

	return existing, nil
}

func (s *SessionStore) Remove(sessionId string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	errors := []string{}
	_, ok := s.sessions[sessionId]
	if ok {
		delete(s.sessions, sessionId)
	} else {
		errors = append(errors, "no session found")
	}

	err := s.idStack.Remove(sessionId)
	if err != nil {
		errors = append(errors, "no id in stack")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}

func (s *SessionStore) removeExpiredSessions() {
	s.lock.Lock()
	defer s.lock.Unlock()

	deadline := s.timeProvider.Now().Add(-s.Config.expiration())

	for {
		if s.idStack.Count() == 0 {
			return
		}

		id := s.idStack.Peak()
		session, ok := s.sessions[id]
		if !ok {
			log.Errorf("stacked id %s has no corresponding session", id)
			_, _ = s.idStack.Pop()
			continue
		}

		if session.StartTime.After(deadline) {
			return
		}

		delete(s.sessions, id)
		_, _ = s.idStack.Pop()
	}
}
