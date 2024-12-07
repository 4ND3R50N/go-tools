package dbutils_test

import (
	"context"
	"errors"
	"github.com/4ND3R50N/go-tools/dbutils"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/4ND3R50N/testsetup/container"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	pgxPool *pgxpool.Pool
)

const (
	insertCount = 10
)

func TestMain(m *testing.M) {
	// Setup container.
	ctx := context.Background()
	// Initialize postgres.
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testing"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 2)
	ports, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}
	exposedPostgresPort := strings.Split(string(ports), "/")[0]
	dbURL := "postgres://test:test@" + container.AutoGuessHostname() + ":" + exposedPostgresPort + "/testing"

	// Build pgx pool
	pgxCfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		panic(err)
	}
	pPool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		panic(err)
	}
	// Migrate pgx table
	_, _ = pPool.Exec(context.Background(), `
		CREATE TABLE test (
			A varchar(255),
			B int
		);
	`)
	pgxPool = pPool
	c := m.Run()
	if err := postgresContainer.Stop(ctx, nil); err != nil {
		panic(err)
	}
	os.Exit(c)
}

func TestPGXTransaction(t *testing.T) {
	var (
		testWg sync.WaitGroup
		ctx    = context.Background()
	)
	t.Run("happy case", func(t *testing.T) {
		// Insert rows in a table randomly.
		for range insertCount {
			testWg.Add(1)
			go func() {
				conn, err := pgxPool.Acquire(ctx)
				if err != nil {
					panic(err)
				}
				defer conn.Release()
				if err := dbutils.Transaction(
					ctx,
					conn.Conn(),
					func(tx pgx.Tx) error {
						smt := "INSERT INTO test (A, B) VALUES ($1, $2);"
						if _, err := tx.Exec(ctx, smt, "test", 1); err != nil {
							return err
						}
						return nil
					},
					dbutils.WithAdvisoryLock("test1"),
					dbutils.WithAdvisoryLock("test2"),
				); err != nil {
					t.Error(err)
				}
				testWg.Done()
			}()
		}
		testWg.Wait()

		var counter int64
		if err := pgxPool.QueryRow(ctx, "SELECT count(*) FROM test").Scan(&counter); err != nil {
			t.Error(err)
		}
		assert.Equal(t, int64(insertCount), counter)
	})

	t.Run("chained transaction", func(t *testing.T) {
		// Insert rows in a table randomly.
		for range insertCount {
			testWg.Add(1)
			go func() {
				conn, err := pgxPool.Acquire(ctx)
				if err != nil {
					panic(err)
				}
				defer conn.Release()
				if err := dbutils.Transaction(
					ctx,
					conn.Conn(),
					func(tx pgx.Tx) error {
						if err := dbutils.Transaction(ctx, tx.Conn(), func(tx pgx.Tx) error {
							smt := "INSERT INTO test (A, B) VALUES ($1, $2);"
							if _, err := tx.Exec(ctx, smt, "test", 1); err != nil {
								return err
							}
							return nil
						}, dbutils.WithAdvisoryLock("test1")); err != nil {
							t.Error(err)
						}
						return nil
					},
					dbutils.WithAdvisoryLock("test1"),
					dbutils.WithAdvisoryLock("test2"),
				); err != nil {
					t.Error(err)
				}
				testWg.Done()
			}()
		}
		testWg.Wait()

		var counter int64
		if err := pgxPool.QueryRow(ctx, "SELECT count(*) FROM test").Scan(&counter); err != nil {
			t.Error(err)
		}
		assert.Equal(t, int64(insertCount*2), counter)
	})

	t.Run("check transaction timeout", func(t *testing.T) {
		var transactionWg sync.WaitGroup
		conn, err := pgxPool.Acquire(ctx)
		if err != nil {
			panic(err)
		}
		defer conn.Release()
		testWg.Add(1)
		transactionWg.Add(1)
		go func() {
			if err := dbutils.Transaction(
				ctx,
				conn.Conn(),
				func(_ pgx.Tx) error {
					transactionWg.Done()
					time.Sleep(time.Second * 3)
					return nil
				},
				dbutils.WithAdvisoryLock("test1"),
			); err != nil {
				t.Error(err)
			}
			testWg.Done()
		}()
		transactionWg.Wait()
		conn2, err2 := pgxPool.Acquire(ctx)
		if err2 != nil {
			panic(err2)
		}
		defer conn2.Release()
		if err := dbutils.Transaction(
			ctx,
			conn2.Conn(),
			func(tx pgx.Tx) error {
				smt := "INSERT INTO test (A, B) VALUES ($1, $2);"
				if _, err := tx.Exec(ctx, smt, "test", 1); err != nil {
					return err
				}
				return nil
			},
			dbutils.WithAdvisoryLock("test1"),
			dbutils.WithLockTimeout(1),
		); err != nil {
			assert.Equal(t, dbutils.ErrCouldNotAcquireLock, err)
		} else {
			assert.FailNow(t, "should throw advisory lock error")
		}
		testWg.Wait()
	})

	t.Run("rollback", func(t *testing.T) {
		expErr := errors.New("test")

		conn, err := pgxPool.Acquire(ctx)
		if err != nil {
			panic(err)
		}
		defer conn.Release()
		err = dbutils.Transaction(
			ctx,
			conn.Conn(),
			func(_ pgx.Tx) error {
				return expErr
			},
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, expErr)
		assert.Equal(t, "transaction rollback: test", err.Error())
	})
}
