package dbconnbase

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

// pgxFileTracer implementa pgx.QueryTracer escrevendo cada query (start/end)
// em um arquivo dentro do diretório informado. Espelha o comportamento do
// "TRACE DIR" usado no Oracle (go-ora) para que o setup de trace seja
// equivalente entre os dois bancos.
//
// Se dir for vazio, escreve em os.Stdout.
// Se dir não existir, é criado (MkdirAll). Se existir como arquivo, falha.
type pgxFileTracer struct {
	mu sync.Mutex
	w  io.Writer
	f  *os.File // nil quando saída é stdout
}

func newPGXFileTracer(dir string) (*pgxFileTracer, error) {
	if dir == "" {
		return &pgxFileTracer{w: os.Stdout}, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("pgx tracer: criar diretório %q: %w", dir, err)
	}
	name := fmt.Sprintf("pgx_trace_%s.log", time.Now().Format("20060102_150405"))
	path := filepath.Join(dir, name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("pgx tracer: abrir %q: %w", path, err)
	}
	return &pgxFileTracer{w: f, f: f}, nil
}

func (t *pgxFileTracer) Close() error {
	if t.f == nil {
		return nil
	}
	return t.f.Close()
}

type pgxTraceCtxKey struct{}

type pgxTraceCtx struct {
	start time.Time
	sql   string
	args  []any
}

func (t *pgxFileTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, d pgx.TraceQueryStartData) context.Context {
	return context.WithValue(ctx, pgxTraceCtxKey{}, &pgxTraceCtx{
		start: time.Now(),
		sql:   d.SQL,
		args:  d.Args,
	})
}

func (t *pgxFileTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, d pgx.TraceQueryEndData) {
	tc, _ := ctx.Value(pgxTraceCtxKey{}).(*pgxTraceCtx)
	if tc == nil {
		return
	}
	dur := time.Since(tc.start)

	t.mu.Lock()
	defer t.mu.Unlock()

	stamp := tc.start.Format("2006-01-02 15:04:05.000")
	if d.Err != nil {
		fmt.Fprintf(t.w, "[%s] dur=%s ERR=%v\n  sql=%s\n  args=%v\n",
			stamp, dur, d.Err, tc.sql, tc.args)
	} else {
		fmt.Fprintf(t.w, "[%s] dur=%s rows=%d\n  sql=%s\n  args=%v\n",
			stamp, dur, d.CommandTag.RowsAffected(), tc.sql, tc.args)
	}
}
