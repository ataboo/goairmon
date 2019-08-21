package context

import (
	"goairmon/site/helper"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func NewPostgresContext() DbContext {
	db := connectPostgres(pgOptions())

	return &postgresContext{
		db: db,
	}
}

func connectPostgres(options *pg.Options) *pg.DB {
	return pg.Connect(options)
}

func pgOptions() *pg.Options {
	return &pg.Options{
		Addr:     helper.MustGetEnv("DB_ADDR"),
		User:     helper.MustGetEnv("DB_USERNAME"),
		Password: helper.MustGetEnv("DB_PASSWORD"),
		Database: helper.MustGetEnv("DB_DATABASE"),
	}
}

type DbContext interface {
	Close() error
	// AddUser(user models.User)
	// FindUser(id int)
	// ExecRaw(rawSql string, args ...interface{}) (sql.Result, error)
	// Query(rawSql string, args ...interface{}) (*sql.Rows, error)
}

type postgresContext struct {
	db *pg.DB
}

func (d *postgresContext) Close() error {
	return d.db.Close()
}

func (d *postgresContext) DbExists() bool {
	result := d.mustExec("SELECT to_regclass('$1.users')", helper.MustGetEnv("DB_DATABASE"))

	return result.RowsReturned() > 0
}

func (d *postgresContext) mustExec(query string, params ...interface{}) orm.Result {
	result, err := d.db.Exec(query, params)
	if err != nil {
		panic(err)
	}

	return result
}

// func (d *postgresContext) GetUsers()
