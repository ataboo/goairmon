package context

import (
	"encoding/json"
	"goairmon/business/data/models"
	"goairmon/site/helper"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func _setupMemDbContext(t *testing.T) *memDbContext {
	if err := godotenv.Load(helper.AppRoot() + "/.env.testing"); err != nil {
		t.Fatal("failed to load .env.testing")
	}

	ctx := NewMemDbContext(&MemDbConfig{
		StoragePath:      helper.MustGetEnv("STORAGE_PATH"),
		SensorPointCount: 10,
		EncodeReadible:   true,
	})

	memCtx := ctx.(*memDbContext)

	os.RemoveAll(memCtx.cfg.StoragePath)
	memCtx.loadStoredConfig()
	memCtx.loadPoints()

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

	ctx.storedConfig.Users[user1.ID] = user1.CopyTo(&models.User{})
	ctx.storedConfig.Users[user2.ID] = user2.CopyTo(&models.User{})

	user1.Username = "changed-username"

	if ctx.storedConfig.Users[user1.ID].Username != "test-username" {
		t.Error("username should not have changed")
	}

	if err := ctx.CreateOrUpdateUser(user1); err != nil {
		t.Error(err)
	}

	if len(ctx.storedConfig.Users) != 2 {
		t.Error("expected two users", len(ctx.storedConfig.Users))
	}

	if ctx.storedConfig.Users[user1.ID].Username != "changed-username" {
		t.Error("username should have changed")
	}
}

func TestSaveAndLoad(t *testing.T) {
	ctx := _setupMemDbContext(t)
	user := &models.User{
		Username: "test-username",
	}
	point := &models.SensorPoint{
		Time:     time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
		Co2Value: 1.0,
	}

	ctx.CreateOrUpdateUser(user)
	ctx.PushSensorPoint(point)

	if err := ctx.Close(); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(ctx.configFile()); err != nil {
		t.Error("failed to find user os file")
	}

	if _, err := os.Stat(ctx.pointFile()); err != nil {
		t.Error("failed to find point os file")
	}

	ctx.storedConfig.Users = nil
	ctx.sensorPoints.Clear()

	if err := ctx.loadStoredConfig(); err != nil {
		t.Error(err)
	}

	result, err := ctx.FindUser(user.ID)
	if err != nil {
		t.Error(err)
	}

	if result.Username != user.Username || result.ID != user.ID {
		t.Errorf("User mismatch: %+v, %+v", user, result)
	}

	if err := ctx.loadPoints(); err != nil {
		t.Error(err)
	}

	points, err := ctx.GetSensorPoints(1)
	if err != nil {
		t.Error(err)
	}

	if len(points) != 1 {
		t.Error("unexpected point count", 1, len(points))
	}

	if points[0].Time != point.Time || points[0].Co2Value != point.Co2Value {
		t.Error("point mismatch", points[0], point)
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

	ctx.storedConfig.Users[user1.ID] = user1
	ctx.storedConfig.Users[user2.ID] = user2

	if err := ctx.DeleteUser(user2.ID); err != nil {
		t.Error(err)
	}

	if len(ctx.storedConfig.Users) != 1 {
		t.Error("unexpected count", len(ctx.storedConfig.Users))
	}

	if _, ok := ctx.storedConfig.Users[user1.ID]; !ok {
		t.Error("user 1 should exist")
	}

	if err := ctx.DeleteUser(user2.ID); err == nil {
		t.Error("expected error")
	}
}

func TestLoadInvalidFile(t *testing.T) {
	ctx := _setupMemDbContext(t)

	os.MkdirAll(ctx.cfg.StoragePath, 0700)
	if err := ioutil.WriteFile(ctx.configFile(), []byte("garbagedata"), 0644); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(ctx.pointFile(), []byte("garbagedata"), 0644); err != nil {
		t.Error(err)
	}

	if err := ctx.loadStoredConfig(); err == nil {
		t.Error("expected error")
	}

	if err := ctx.loadPoints(); err == nil {
		t.Error("expected error")
	}

	if ctx.storedConfig.Users == nil || len(ctx.storedConfig.Users) != 0 {
		t.Error("expected empty users map set")
	}
}

func TestSaveInvalidData(t *testing.T) {
	ctx := _setupMemDbContext(t)
	ctx.storedConfig.Users = nil

	os.Remove(ctx.configFile())
	os.MkdirAll(ctx.configFile(), 0700)
	os.MkdirAll(ctx.pointFile(), 0700)
	if err := ctx.saveStoredConfig(); err == nil {
		t.Error("expected error on save")
	}

	if err := ctx.savePoints(); err == nil {
		t.Error("expected error on save")
	}

	if err := ctx.Close(); err == nil {
		t.Error("expected error")
	}

	os.Remove(ctx.configFile())
	os.Remove(ctx.pointFile())
}

func TestPushSensorPoints(t *testing.T) {
	ctx := _setupMemDbContext(t)

	point1 := &models.SensorPoint{
		Time:     time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
		Co2Value: 1.0,
	}
	point2 := &models.SensorPoint{
		Time: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if err := ctx.PushSensorPoint(point1); err != nil {
		t.Error(err)
	}

	if err := ctx.PushSensorPoint(point2); err != nil {
		t.Error(err)
	}

	peaked, err := ctx.GetSensorPoints(2)
	if err != nil {
		t.Error(err)
	}

	if peaked[0].Co2Value != point2.Co2Value || peaked[0].Time != point2.Time {
		t.Error("value mismatch", peaked[0], point2)
	}

	if peaked[1].Co2Value != point1.Co2Value || peaked[1].Time != point1.Time {
		t.Error("value mismatch", peaked[1], point1)
	}
}

func TestSaveSensorBaseline(t *testing.T) {
	ctx := _setupMemDbContext(t)

	if _, _, err := ctx.GetSensorBaseline(); err == nil {
		t.Error("expected error")
	}

	if err := ctx.SetSensorBaseline(1, 2); err != nil {
		t.Error(err)
	}

	eco2, tvoc, err := ctx.GetSensorBaseline()
	if err != nil {
		t.Error(err)
	}

	if eco2 != 1 {
		t.Error("unexpected eco2 value", 1, eco2)
	}

	if tvoc != 2 {
		t.Error("unexpected tvoc value", 2, tvoc)
	}

	if err := ctx.saveStoredConfig(); err != nil {
		t.Error(err)
	}

	ctx.storedConfig.ECO2Baseline = 0
	ctx.storedConfig.TVOCBaseline = 0

	if err := ctx.loadStoredConfig(); err != nil {
		t.Error(err)
	}

	if ctx.storedConfig.ECO2Baseline != 1 {
		t.Error("unexpected eco2", 1, ctx.storedConfig.ECO2Baseline)
	}

	if ctx.storedConfig.TVOCBaseline != 2 {
		t.Error("unexpected tvoc", 2, ctx.storedConfig.TVOCBaseline)
	}
}

func TestArchiveLastDay(t *testing.T) {
	ctx := _setupMemDbContext(t)
	ctx.cfg.SensorPointCount = 24 * 60 * 2

	ctx.sensorPoints = NewSensorPointStack(ctx.cfg.SensorPointCount)
	err := ctx.archiveLastDay()
	if err == nil {
		t.Error("expected error")
	}

	for i := 0; i < 24*60*2; i++ {
		ctx.sensorPoints.Push(&models.SensorPoint{
			Time:     time.Date(2010, 01, 01, 00, 00, 00, 00, time.UTC).Add(time.Minute * time.Duration(i)),
			Co2Value: 23,
		})
	}

	if err := ctx.PushSensorPoint(&models.SensorPoint{
		Time:     time.Date(2010, 1, 3, 0, 0, 0, 0, time.UTC),
		Co2Value: 23,
	}); err != nil {
		t.Error("unexpected error", err)
	}

	loaded, err := ioutil.ReadFile(ctx.cfg.StoragePath + "/archive_2010_01_02.json")
	if err != nil {
		t.Error("unexpected error", err)
	}

	var decoded []*models.SensorPoint
	if err := json.Unmarshal(loaded, &decoded); err != nil {
		t.Error("unexpected error", err)
	}

	if len(decoded) != 24*60-1 {
		t.Error("unexpected count", 24*60-1, len(decoded))
	}

	for _, p := range decoded {
		if p.Time.YearDay() != 2 {
			t.Error("values should be from Jan 2", p)
		}
	}
}
