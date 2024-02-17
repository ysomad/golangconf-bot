package postgres

import (
	"context"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns     = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Client struct {
	maxConns     int32
	connAttempts uint8
	connTimeout  time.Duration
	tracer       pgx.QueryTracer

	Builder sq.StatementBuilderType
	Pool    *pgxpool.Pool
}

func NewClient(connString string, opts ...Option) (*Client, error) {
	c := &Client{
		maxConns:     defaultMaxConns,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	poolConf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	if c.tracer != nil {
		poolConf.ConnConfig.Tracer = c.tracer
	}

	poolConf.MaxConns = c.maxConns

	ctx := context.Background()

	c.Pool, err = pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, err
	}

	for c.connAttempts > 0 {
		if err = c.Pool.Ping(ctx); err == nil {
			break
		}

		slog.Info("trying connecting to postgres", "attempts_left", c.connAttempts)
		time.Sleep(c.connTimeout)
		c.connAttempts--
	}

	if err != nil {
		return nil, err
	}

	c.Builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	slog.Info("connected to postgres")

	return c, nil
}

func (p *Client) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

type Option func(*Client)

func WithMaxConns(conns int32) Option {
	return func(c *Client) {
		c.maxConns = conns
	}
}

func WithConnAttempts(attempts uint8) Option {
	return func(c *Client) {
		c.connAttempts = attempts
	}
}

func WithConnTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.connTimeout = timeout
	}
}

func WithQueryTracer(tracer pgx.QueryTracer) Option {
	return func(c *Client) {
		c.tracer = tracer
	}
}
