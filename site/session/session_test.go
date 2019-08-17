package session

import (
	"testing"
	"time"
)

func _assertSessionMatch(session *Session, id string, startTime time.Time, userId string, t *testing.T) {
	if session.Id != id {
		t.Errorf("unnexpected session id: %s, %s", id, session.Id)
	}

	if session.StartTime != startTime {
		t.Errorf("unnexpected session start time: %s, %s", startTime, session.StartTime)
	}

	if session.Values["user_id"] != userId {
		t.Errorf("unnexpected user id: %s, %s", userId, session.Values["user_id"])
	}
}

func _assertSessionStore(store *SessionStore, sessionCount int, hasSessionId string, t *testing.T) {
	if len(store.sessions) != sessionCount {
		t.Errorf("unnexpected session count: %d, %d", sessionCount, len(store.sessions))
	}

	if store.idStack.Count() != sessionCount {
		t.Errorf("unnexpected id stack count: %d, %d", sessionCount, store.idStack.Count())
	}

	idx := store.idStack.IndexOf(hasSessionId)
	if idx < 0 {
		t.Errorf("id missing from idstack")
	}

	_, ok := store.sessions[hasSessionId]
	if !ok {
		t.Errorf("id missing from sessions")
	}
}

func TestFindSession(t *testing.T) {
	store := _sessionSetup()
	timePro := (store.timeProvider).(*mockTimeProvider)
	startNow := time.Unix(503539200, 0)
	timePro.NowCallback = func() time.Time {
		return startNow
	}

	_, err := store.Find("not_found")
	if err == nil {
		t.Error("expected error")
	}

	store.sessions["first_id"] = &Session{
		Id:        "first_id",
		StartTime: startNow,
		Values:    map[string]string{"user_id": "first_user"},
	}
	store.idStack.PushBack("first_id")
	store.sessions["second_id"] = &Session{
		Id:        "second_id",
		StartTime: startNow,
		Values:    map[string]string{"user_id": "second_user"},
	}
	store.idStack.PushBack("second_id")

	sess, err := store.Find("first_id")
	if err != nil {
		t.Error("error when finding session", err)
	}

	_assertSessionMatch(sess, "first_id", startNow, "first_user", t)
	_assertSessionStore(store, 2, "first_id", t)

	_, err = store.Find("not_a_match")
	if err == nil {
		t.Error("expected error finding session")
	}
}

func TestRemoveExpiredSessionOnFind(t *testing.T) {
	store := _sessionSetup()
	timePro := (store.timeProvider).(*mockTimeProvider)
	startNow := time.Unix(503539200, 0)
	timePro.NowCallback = func() time.Time {
		return startNow
	}

	store.sessions["first_id"] = &Session{
		Id:        "expired_session",
		StartTime: startNow.Add(-store.Config.expirationDuration()),
	}

	store.idStack.PushBack("first_id")

	sess, err := store.Find("first_id")
	if err != nil || sess == nil {
		t.Error("expected to find session")
	}

	if sess.StartTime != startNow {
		t.Error("expected start time to be set to now on successful find")
	}

	store.sessions["first_id"].StartTime = sess.StartTime.Add(-store.Config.expirationDuration() - time.Minute)

	sess, err = store.Find("first_id")
	if err == nil || sess != nil {
		t.Error("expected not to find session")
	}

	if len(store.sessions) != 0 || store.idStack.Count() != 0 {
		t.Error("expected session to be deleted")
	}
}

func TestNewOrExisting(t *testing.T) {
	store := _sessionSetup()
	timePro := (store.timeProvider).(*mockTimeProvider)
	startNow := time.Unix(503539200, 0)
	timePro.NowCallback = func() time.Time {
		return startNow
	}

	sess, err := store.NewOrExisting("first_session", "first_user")
	if err != nil {
		t.Error("unnexpected err", err)
	}

	_assertSessionMatch(sess, "first_session", startNow, "first_user", t)
	_assertSessionStore(store, 1, "first_session", t)

	secondSess, err := store.NewOrExisting("second_session", "second_user")
	if err != nil {
		t.Error("unnexpected err", err)
	}

	_assertSessionMatch(secondSess, "second_session", startNow, "second_user", t)
	_assertSessionStore(store, 2, "second_session", t)
	if store.idStack.Peak() != "first_session" {
		t.Error("first session should be ontop of stack")
	}

	sess, err = store.NewOrExisting("first_session", "first_user")
	if err != nil {
		t.Error("unnexpected err", err)
	}

	if store.idStack.Peak() != "second_session" {
		t.Error("second session should be ontop of stack")
	}

	_assertSessionMatch(sess, "first_session", startNow, "first_user", t)
	_assertSessionStore(store, 2, "first_session", t)
}

func TestRemoveSession(t *testing.T) {
	store := _sessionSetup()

	err := store.Remove("not_found")
	if err == nil {
		t.Error("expected error")
	}

	store.sessions["first_id"] = &Session{
		Id:        "first_id",
		StartTime: time.Now(),
		Values:    map[string]string{"user_id": "first_user"},
	}
	store.idStack.PushBack("first_id")
	store.sessions["second_id"] = &Session{
		Id:        "second_id",
		StartTime: time.Now(),
		Values:    map[string]string{"user_id": "second_user"},
	}
	store.idStack.PushBack("second_id")

	err = store.Remove("first_id")
	if err != nil {
		t.Error("unexpected error")
	}

	_assertSessionStore(store, 1, "second_id", t)
}

func TestRemoveExpiredSessions(t *testing.T) {
	store := _sessionSetup()
	timePro := (store.timeProvider).(*mockTimeProvider)
	startNow := time.Unix(503539200, 0)
	timePro.NowCallback = func() time.Time {
		return startNow
	}

	store.removeExpiredSessions()

	store.idStack.PushBack("this_has_no_session")

	store.removeExpiredSessions()

	if store.idStack.Count() > 0 {
		t.Error("id stack should have been cleared")
	}

	if _, err := store.NewOrExisting("first_id", "first_user"); err != nil {
		t.Error(err)
	}

	startNow = startNow.Add(store.Config.expirationDuration() + time.Minute)

	if _, err := store.NewOrExisting("second_id", "second_user"); err != nil {
		t.Error(err)
	}

	_assertSessionStore(store, 2, "first_id", t)

	store.removeExpiredSessions()

	_assertSessionStore(store, 1, "second_id", t)
}

func TestGCRemovesSessions(t *testing.T) {
	store := _sessionSetup()
	timePro := (store.timeProvider).(*mockTimeProvider)
	startNow := time.Unix(503539200, 0)
	timePro.NowCallback = func() time.Time {
		return startNow
	}

	if _, err := store.NewOrExisting("first_id", "first_user"); err != nil {
		t.Error(err)
	}

	startNow = startNow.Add(store.Config.expirationDuration() + time.Minute)

	if _, err := store.NewOrExisting("second_id", "second_user"); err != nil {
		t.Error(err)
	}

	err := store.StartGC()
	if err != nil {
		t.Error("unexpected error")
	}

	err = store.StartGC()
	if err == nil {
		t.Error("expected error")
	}

	_assertSessionStore(store, 2, "first_id", t)

	timePro.TickChannel <- startNow

	_assertSessionStore(store, 1, "second_id", t)

	startNow = startNow.Add(store.Config.expirationDuration() + time.Minute)
}

func TestSystemTime(t *testing.T) {
	sysTime := systemTime{}
	if time.Until(sysTime.Now()) > time.Second {
		t.Errorf("now should match %s, %s", sysTime.Now(), time.Now())
	}

	tickChan := sysTime.Tick(time.Microsecond)
	timeout := time.After(time.Second)
	select {
	case <-tickChan:
		break
	case <-timeout:
		t.Error("timed out of tick")
	}
}

func _sessionSetup() *SessionStore {
	tickChan := make(chan time.Time)
	timeProvider := &mockTimeProvider{
		TickChannel: tickChan,
		NowCallback: nil,
	}
	cfg := Config{
		SessionKey:     "goairmon_session",
		ExpirationSecs: 3600,
		GCDelaySeconds: 600,
	}
	store := NewSessionStore(cfg)
	store.timeProvider = timeProvider

	return store
}

type mockTimeProvider struct {
	NowCallback func() time.Time
	TickChannel chan time.Time
}

func (t mockTimeProvider) Now() time.Time {
	if t.NowCallback == nil {
		return time.Now()
	}

	return t.NowCallback()
}

func (t mockTimeProvider) Tick(duration time.Duration) <-chan time.Time {
	return t.TickChannel
}
