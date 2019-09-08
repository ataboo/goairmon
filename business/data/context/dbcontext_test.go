package context

import (
	"goairmon/business/data/models"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
)

func _setupMemDbContext(t *testing.T) *memDbContext {
	ctx := NewMemDbContext(&MemDbConfig{
		StoragePath: "/tmp/goairmon_testing",
	})

	memCtx := ctx.(*memDbContext)

	os.Remove(memCtx.userFile())
	memCtx.load()

	return memCtx
}

func TestCreateUser(t *testing.T) {
	ctx := _setupMemDbContext(t)

	user := &models.User{
		Username: "test-username",
	}

	if user.ID != uuid.Nil {
		t.Error("expected nil starting id")
	}

	err := ctx.CreateOrUpdateUser(user)
	if err != nil {
		t.Error(err)
	}

	if user.ID == uuid.Nil {
		t.Error("ID should be set")
	}

	loaded, err := ctx.FindUser(user.ID)
	if err != nil {
		t.Error("failed to find users")
	}

	if loaded.Username != user.Username {
		t.Error("username mismatch")
	}

	if _, err := ctx.FindUser(uuid.New()); err == nil {
		t.Error("should not find user")
	}
}

func TestUpdateExistingUser(t *testing.T) {
	ctx := _setupMemDbContext(t)

	user1 := &models.User{
		ID:       uuid.New(),
		Username: "test-username",
	}
	user2 := &models.User{
		ID:       uuid.New(),
		Username: "test-username2",
	}

	ctx.users[user1.ID] = user1.CopyTo(&models.User{})
	ctx.users[user2.ID] = user2.CopyTo(&models.User{})

	user1.Username = "changed-username"

	if ctx.users[user1.ID].Username != "test-username" {
		t.Error("username should not have changed")
	}

	if err := ctx.CreateOrUpdateUser(user1); err != nil {
		t.Error(err)
	}

	if len(ctx.users) != 2 {
		t.Error("expected two users", len(ctx.users))
	}

	if ctx.users[user1.ID].Username != "changed-username" {
		t.Error("username should have changed")
	}
}

func TestSaveAndLoad(t *testing.T) {
	ctx := _setupMemDbContext(t)
	user := &models.User{
		Username: "test-username",
	}

	ctx.CreateOrUpdateUser(user)

	if err := ctx.Close(); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(ctx.userFile()); err != nil {
		t.Error("failed to find os file")
	}

	ctx.users = nil

	if err := ctx.load(); err != nil {
		t.Error(err)
	}

	result, err := ctx.FindUser(user.ID)
	if err != nil {
		t.Error(err)
	}

	if result.Username != user.Username || result.ID != user.ID {
		t.Errorf("User mismatch: %+v, %+v", user, result)
	}
}

func TestFindByName(t *testing.T) {
	ctx := _setupMemDbContext(t)
	user1 := &models.User{
		Username: "first-user",
	}
	user2 := &models.User{
		Username: "second-user",
	}

	ctx.CreateOrUpdateUser(user1)
	ctx.CreateOrUpdateUser(user2)

	_, err := ctx.FindUserByName("not-found")
	if err == nil {
		t.Error("expected error")
	}

	found1, err := ctx.FindUserByName("first-user")
	if err != nil {
		t.Error(err)
	}

	if found1.Username != user1.Username || found1.ID != user1.ID {
		t.Error("users don't match")
	}

	found2, err := ctx.FindUserByName("second-user")
	if err != nil {
		t.Error(err)
	}

	if found2.Username != user2.Username || found2.ID != user2.ID {
		t.Error("users don't match")
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := _setupMemDbContext(t)
	user1 := &models.User{
		ID:       uuid.New(),
		Username: "first-user",
	}
	user2 := &models.User{
		ID:       uuid.New(),
		Username: "second-user",
	}

	ctx.users[user1.ID] = user1
	ctx.users[user2.ID] = user2

	if err := ctx.DeleteUser(user2.ID); err != nil {
		t.Error(err)
	}

	if len(ctx.users) != 1 {
		t.Error("unexpected count", len(ctx.users))
	}

	if _, ok := ctx.users[user1.ID]; !ok {
		t.Error("user 1 should exist")
	}

	if err := ctx.DeleteUser(user2.ID); err == nil {
		t.Error("expected error")
	}
}

func TestLoadInvalidFile(t *testing.T) {
	ctx := _setupMemDbContext(t)
	if err := ioutil.WriteFile(ctx.userFile(), []byte("garbagedata"), 0644); err != nil {
		t.Error(err)
	}

	if err := ctx.load(); err == nil {
		t.Error("expected error")
	}

	if ctx.users == nil || len(ctx.users) != 0 {
		t.Error("expected empty users map set")
	}
}

func TestSaveInvalidData(t *testing.T) {
	ctx := _setupMemDbContext(t)
	ctx.users = nil

	os.Remove(ctx.userFile())
	os.MkdirAll(ctx.userFile(), 0700)
	if err := ctx.save(); err == nil {
		t.Error("expected error on save")
	}

	os.Remove(ctx.userFile())
}
