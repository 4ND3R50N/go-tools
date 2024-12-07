package dbutils_test

import (
	"context"
	"errors"

	"testing"

	"github.com/4ND3R50N/go-tools/dbutils"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewPGXLocks(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mock.Close()
		mock.
			ExpectExec("SELECT pg_advisory_xact_lock(?)").
			WithArgs(int64(-4578387130389545126)).
			WillReturnResult(pgxmock.NewResult("test", 1))
		mock.
			ExpectExec("SELECT pg_advisory_xact_lock(?)").
			WithArgs(int64(-4578387130389545127)).
			WillReturnResult(pgxmock.NewResult("test", 1))
		assert.NoError(t, dbutils.NewPGXLocks(context.Background(), mock, "test1", "test2"))
	})

	t.Run("error cases", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mock.Close()
		someError := errors.New("some error")
		mock.
			ExpectExec("SELECT pg_advisory_xact_lock(?)").
			WithArgs(int64(-4578387130389545126)).
			WillReturnError(err)
		assert.ErrorAs(t,
			dbutils.NewPGXLocks(context.Background(), mock, "test1", "test2"), &someError)
	})
}
