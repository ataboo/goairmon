package context

import (
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
)

func _dbTestSetup(t *testing.T) *pg.DB {
	if err := godotenv.Load("../../../.env.testing"); err != nil {
		t.Error(err)
	}

	return connectPostgres(pgOptions())
}

func TestConnection(t *testing.T) {
	_dbTestSetup(t)

	result, err := conn.Exec("SELECT * FROM users")
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}
