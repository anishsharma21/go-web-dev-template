package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal/handlers"
	"github.com/anishsharma21/go-web-dev-template/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

//go:embed templates/*.html
var templateFS embed.FS

var (
	dbConnStr string

	dbPool    *pgxpool.Pool
	templates *template.Template
)

func init() {
	dbConnStr = os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		slog.Error("DATABASE_URL environment variable not set")
		os.Exit(1)
	}

	// Set up slog as default logger across the application
	// Default logger is JSON logger with request_id and user_id fields added if they exist in the context
	defaultHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(&middleware.CustomLogHandler{Handler: defaultHandler}))

	// Parse html templates
	var err error
	templates, err = template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		os.Exit(1)
	}
	slog.Info("Templates parsed successfully")
}

func main() {
	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup database connection pool
	dbPool, err := setupDBPool(ctx)
	if err != nil {
		slog.Error("Failed to initialise database connection pool", "error", err)
		return
	}
	defer dbPool.Close()

	// Run database migrations if environment variable is set for it
	if os.Getenv("RUN_MIGRATION") == "true" {
		slog.Info("Attempting to run database migrations...")
		err := runMigrations()
		if err != nil {
			slog.Error("Failed to run database migrations", "error", err)
			return
		}
		slog.Info("Database migrations complete.")
	} else {
		slog.Info("Database migrations skipped.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: setupRoutes(dbPool),
		BaseContext: func(l net.Listener) context.Context {
			url := "http://" + l.Addr().String()
			slog.Info(fmt.Sprintf("Server started on %s", url))
			return ctx
		},
	}

	shutdownChan := make(chan bool, 1)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server closed early", "error", err)
		}
		slog.Info("Stopped server new connections.")
		shutdownChan <- true
	}()

	// Listen for OS signals (SIGINT, SIGTERM) to shutdown server gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Warn("Received signal", "signal", sig.String())

	// Shutdown server gracefully within 10 seconds
	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error occurred", "error", err)
	}
	<-shutdownChan
	close(shutdownChan)

	slog.Info("Graceful server shutdown complete.")
}

func setupDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	// Parse database connection string into pgxpool config
	config, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse database connection string: %v", err)
	}

	// Set connection pool configurations
	// Sets the maximum time an idle connection can remain in the pool before being closed
	config.MaxConnIdleTime = 1 * time.Minute
	// To prevent database and backend from ever sleeping, uncomment the following line
	config.MinConns = 1

	// Try to initialise database connection pool 5 times with exponential backoff
	var dbPool *pgxpool.Pool
	for i := 1; i <= 5; i++ {
		dbPool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil && dbPool != nil {
			break
		}
		slog.Warn("Failed to initialise database connection pool", "error", err)
		slog.Info(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if dbPool == nil {
		return nil, fmt.Errorf("Failed to initialise database connection pool after 5 attempts")
	}

	// Try to ping database connection pool 5 times with exponential backoff
	for i := 1; i <= 5; i++ {
		err = dbPool.Ping(ctx)
		if err == nil && dbPool != nil {
			break
		}
		slog.Warn("Failed to ping database connection pool", "error", err)
		slog.Info(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to ping database connection pool after 5 attempts")
	}

	return dbPool, nil
}

type routeConfig struct {
	Handler      http.Handler
	ApplyLogging bool
	ApplyJWT     bool
}

func setupRoutes(dbPool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	routes := map[string]routeConfig{
		// auth
		"POST /api/signup": {
			Handler:      handlers.SignUp(dbPool),
			ApplyLogging: true,
			ApplyJWT:     false,
		},
		"POST /api/login": {
			Handler:      handlers.Login(dbPool),
			ApplyLogging: true,
			ApplyJWT:     false,
		},
		"POST /api/refresh-token": {
			Handler:      handlers.RefreshToken(),
			ApplyLogging: false,
			ApplyJWT:     false,
		},
		"GET /login": {
			Handler:      handlers.RenderLoginView(templates),
			ApplyLogging: true,
			ApplyJWT:     false,
		},

		// users
		"GET /users": {
			Handler:      handlers.RenderBaseUserView(dbPool, templates),
			ApplyLogging: true,
			ApplyJWT:     true,
		},
		"GET /": {
			Handler:      handlers.RenderBaseView(templates),
			ApplyLogging: true,
			ApplyJWT:     false,
		},

		// other
		"GET /static/": {
			Handler:      http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
			ApplyLogging: false,
			ApplyJWT:     false,
		},
	}

	for pattern, config := range routes {
		handler := config.Handler
		if config.ApplyLogging {
			handler = middleware.LoggingMiddleware(handler)
		}
		if config.ApplyJWT {
			handler = middleware.JWTMiddleware(handler)
		}
		mux.Handle(pattern, handler)
	}

	return mux
}

func chainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func runMigrations() error {
	if gooseDriver := os.Getenv("GOOSE_DRIVER"); gooseDriver == "" {
		return fmt.Errorf("Goose driver not set: GOOSE_DRIVER=?")
	}

	if gooseDbString := os.Getenv("GOOSE_DBSTRING"); gooseDbString == "" {
		return fmt.Errorf("Goose db string not set: GOOSE_DBSTRING=?")
	}

	if gooseMigrationDir := os.Getenv("GOOSE_MIGRATION_DIR"); gooseMigrationDir == "" {
		return fmt.Errorf("Goose migration dir not set: GOOSE_MIGRATION_DIR=?")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return fmt.Errorf("Failed to open database connection for *sql.DB: %v\n", err)
	}
	defer db.Close()

	if err = goose.Status(db, "migrations"); err != nil {
		return fmt.Errorf("Failed to retrieve status of migrations: %v\n", err)
	}

	if err = goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("Failed to run `goose up` command: %v\n", err)
	}

	return nil
}
