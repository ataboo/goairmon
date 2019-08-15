package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type Config struct {
	SessionCookieName string
	SessionKey        string
	ExpirationSecs    int
}

type Session struct {
	Id        string
	StartTime time.Time
	Values    map[string]string
}

type SessionStore struct {
	Config      Config
	CookieStore sessions.CookieStore
	startedGc   bool
	sessions    map[string]Session
}

func (s *SessionStore) StartGC() error {
	if s.startedGc {
		return fmt.Errorf("GC already started")
	}

	s.startedGc = true

	go func() {
		for {
			select {
			case <-time.Tick(time.Minute):
				s.runGc()
			}
		}
	}()

	return nil
}

func (s *SessionStore) New(userId string, r *http.Request) Session {
	cookieSession, _ := s.CookieStore.Get(r, s.Config.SessionCookieName)
	newSession := Session{
		Id:        cookieSession.ID,
		StartTime: time.Now(),
	}
	newSession.Values["user_id"] = userId

	s.sessions[cookieSession.ID] = newSession

	return newSession
}

func (s *SessionStore) Current(r *http.Request) (Session, error) {
	cookieSession, err := s.CookieStore.Get(r, s.Config.SessionCookieName)
	if err != nil {
		return Session{}, err
	}

	sessionId := cookieSession.ID

	currentSession, ok := s.sessions[sessionId]
	if !ok {
		return Session{}, fmt.Errorf("session not found")
	}

	return currentSession, nil
}

func (s *SessionStore) runGc() {
	// for(i:=0; len(i<s.sessions); i++) {
	// }
}

func MyHandler(r *http.Request) {
	store := sessions.NewCookieStore()

	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(c.Request(), "goairmon_session")
	// Set some session values.
	session.Values["session_id"] = uuid.New()
	// Save it before we write to the response/return from the handler.
	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		c.Logger().Fatal("failed to save session", err)
	}
}
