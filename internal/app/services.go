package app

import (
	"context"
	"database/sql"

	"github.com/imnulhaqueruman/opencode-poc/internal/config"
	"github.com/imnulhaqueruman/opencode-poc/internal/db"
	"github.com/imnulhaqueruman/opencode-poc/internal/logging"
	"github.com/imnulhaqueruman/opencode-poc/internal/message"
	"github.com/imnulhaqueruman/opencode-poc/internal/permission"
	"github.com/imnulhaqueruman/opencode-poc/internal/session"
)

type App struct {
	Context context.Context

	Sessions    session.Service
	Messages    message.Service
	Permissions permission.Service

	Logger logging.Interface
}

func New(ctx context.Context, conn *sql.DB) *App {
	q := db.New(conn)
	log := logging.NewLogger(logging.Options{
		Level: config.Get().Log.Level,
	})
	sessions := session.NewService(ctx, q)
	messages := message.NewService(ctx, q)

	return &App{
		Context:     ctx,
		Sessions:    sessions,
		Messages:    messages,
		Permissions: permission.Default,
		Logger:      log,
	}
}