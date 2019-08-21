package migration

import (
	"fmt"
	"goairmon/business/data/context"
	"goairmon/site/helper"
)

func Migrate(dbContext context.DbContext) {
	existsSql := "SELECT datname FROM pg_catalog.pg_database WHERE datname = $1;"
	rows, err := dbContext.Query(existsSql, helper.MustGetEnv("DB_DATABASE"))
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		fmt.Printf("Found it!")
	} else {
		fmt.Printf("Not Fount it!")
	}

}
