package dbconnbase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aoticombr/golang/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	_ "github.com/sijms/go-ora/v2"
)

type Conn struct {
	DB           *sql.DB
	Dialect      DialectType
	DSN          string
	Log          bool
	PoolSize     int
	PoolLifetime time.Duration
	MaxOpenConns int
	ConnLifetime time.Duration
	Trace        config.Trace
	connContext  bool
	pgxTracer    *pgxFileTracer // não-nil quando Trace.Ativo em Postgres
	pgxConnStr   string         // string devolvida por stdlib.RegisterConnConfig — usada no Unregister
}

func NewConn(db config.Database) (*Conn, error) {

	fmt.Println("Criando Connection")
	conn := &Conn{
		Dialect:  DialectLowFromString(db.Db),
		DSN:      db.GetDsn(),
		PoolSize: db.PoolSize,
		Trace:    db.Trace,
	}
	fmt.Println("abrindo conexão")
	err := conn.Open()

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (co *Conn) Open() error {
	driver := co.Dialect.String()
	dsn := co.DSN

	// Postgres com trace ativo: ParseConfig + RegisterConnConfig para anexar
	// um pgx.QueryTracer ao ConnConfig (espelho do "TRACE DIR" do Oracle).
	if co.Dialect == POSTGRESQL && co.Trace.Ativo {
		cfg, err := pgx.ParseConfig(co.DSN)
		if err != nil {
			return fmt.Errorf("could not parse pgx config: %w", err)
		}
		tracer, err := newPGXFileTracer(co.Trace.Path)
		if err != nil {
			return fmt.Errorf("could not init pgx tracer: %w", err)
		}
		cfg.Tracer = tracer
		co.pgxTracer = tracer
		dsn = stdlib.RegisterConnConfig(cfg)
		co.pgxConnStr = dsn
	}

	db, err := sql.Open(driver, dsn)

	if err != nil {
		return fmt.Errorf("could not create a connection: %w", err)
	}

	db.SetMaxIdleConns(co.PoolSize)
	db.SetMaxOpenConns(co.MaxOpenConns)
	db.SetConnMaxIdleTime(co.PoolLifetime)
	db.SetConnMaxLifetime(co.ConnLifetime)

	//	if co.Dialect == ORACLE {
	if err = db.Ping(); err != nil {
		return fmt.Errorf("database is not reachable: %w", err)
	}
	//	}

	co.DB = db
	co.connContext = false

	return nil
}

// SetSizePool
// Tamanho maximo do Pool de conexão
func (co *Conn) SetSizePool(n int) {
	co.PoolSize = n
}

// SetPoolLifeTime
// Tempo de vida do Pool de conexões
func (co *Conn) SetPoolLifeTime(d time.Duration) {
	co.PoolLifetime = d
	co.DB.SetConnMaxIdleTime(d)
}

// SetMaxOpenConns
// Maximo de conexões abertas
func (co *Conn) SetMaxOpenConns(n int) {
	co.MaxOpenConns = n
	co.DB.SetMaxOpenConns(n)
}

// SetConnLifeTime
// Tempo de vida das conexões
func (co *Conn) SetConnLifeTime(d time.Duration) {
	co.ConnLifetime = d
	co.DB.SetConnMaxLifetime(d)
}

func (co *Conn) Ping() error {
	if err := co.DB.Ping(); err != nil {
		return fmt.Errorf("database is not reachable: %w", err)
	}
	return nil
}

func (co *Conn) CreateContext(ctx context.Context) (context.Context, context.CancelFunc) {
	timeout := 5 * time.Second
	return context.WithTimeout(ctx, timeout)
}

func (co *Conn) StartTransaction() (*Transaction, error) {
	tx, err := NewTransaction(co)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (co *Conn) StartTransactionContext(ctx context.Context) (*Transaction, error) {
	tx, err := NewTransactionCtx(co, ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
func (co *Conn) Exec(sql string, arg ...any) (sql.Result, error) {
	return co.DB.Exec(sql, arg)
}

func (co *Conn) Close() {
	if err := co.DB.Close(); err != nil {
		return
	}
	if co.pgxConnStr != "" {
		stdlib.UnregisterConnConfig(co.pgxConnStr)
		co.pgxConnStr = ""
	}
	if co.pgxTracer != nil {
		_ = co.pgxTracer.Close()
		co.pgxTracer = nil
	}
}
