package dbcontext

import (
	"context"
	"database/sql"
	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	_ "github.com/lib/pq" // initialize posgresql for test
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const DSN = "postgres://127.0.0.1/go_restful?sslmode=disable&user=postgres&password=postgres"

func TestNew(t *testing.T) {
	runDBTest(t, func(db *sqlx.DB) {
		dbc := New(db)
		assert.NotNil(t, dbc)
		assert.Equal(t, db, dbc.DB())
	})
}

func runDBTest(t *testing.T, f func(db *sqlx.DB)) {
	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = DSN
	}
	db, err := dbx.MustOpen("postgres", dsn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer func() {
		_ = db.Close()
	}()

	sqls := []string{
		"CREATE TABLE IF NOT EXISTS dbcontexttest (id VARCHAR PRIMARY KEY, name VARCHAR)",
		"TRUNCATE dbcontexttest",
	}
	for _, s := range sqls {
		_, err = db.NewQuery(s).Execute()
		if err != nil {
			t.Error(err, " with SQL: ", s)
			t.FailNow()
		}
	}

	f(db)
}

func runCountQuery(t *testing.T, db *sqlx.DB) int {
	var count int
	err := db.NewQuery("SELECT COUNT(*) FROM dbcontexttest").Row(&count)
	assert.Nil(t, err)
	return count

}
