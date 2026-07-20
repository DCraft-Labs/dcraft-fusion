package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dcraft-fusion/control-plane-kernel/internal/audit"
	"github.com/dcraft-fusion/control-plane-kernel/internal/auth"
	"github.com/dcraft-fusion/control-plane-kernel/internal/httpapi"
	"github.com/dcraft-fusion/control-plane-kernel/internal/observability"
	"github.com/dcraft-fusion/control-plane-kernel/internal/phase1"
	"github.com/dcraft-fusion/control-plane-kernel/internal/postgres"
	"github.com/dcraft-fusion/control-plane-kernel/internal/tenancy"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	auditLog, orgRepo, tenantRepo, projectRepo, phase1Store := repositories()
	server := httpapi.NewServerWithPhase1(orgRepo, tenantRepo, projectRepo, auditLog, phase1Store)
	metrics := observability.NewMetrics()
	authConfig := auth.FromEnv()
	authMiddleware := auth.NewMiddleware(authConfig)
	provider, err := auth.NewProvider(authConfig)
	if err != nil {
		slog.Error("failed to initialize oidc provider", "error", err)
		os.Exit(1)
	}
	routes := http.NewServeMux()
	routes.HandleFunc("GET /.well-known/openid-configuration", provider.OpenIDConfiguration)
	routes.HandleFunc("GET /oidc/jwks", provider.JWKS)
	routes.HandleFunc("GET /oidc/authorize", provider.Authorize)
	routes.HandleFunc("POST /oidc/token", provider.Token)
	routes.HandleFunc("GET /api/v1/auth/users", provider.Users)
	routes.HandleFunc("POST /api/v1/auth/login", provider.Login)
	routes.HandleFunc("POST /api/v1/auth/logout", provider.Logout)
	routes.Handle("/metrics", metrics.Handler())
	routes.Handle("/", metrics.Wrap(authMiddleware.Wrap(server.Routes())))

	addr := ":8080"
	slog.Info("starting Fusion control-plane kernel", "addr", addr)
	if err := http.ListenAndServe(addr, routes); err != nil {
		slog.Error("kernel stopped", "error", err)
		os.Exit(1)
	}
}

func repositories() (audit.Log, tenancy.OrganizationRepository, tenancy.TenantRepository, tenancy.ProjectRepository, *phase1.Store) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		slog.Warn("POSTGRES_DSN not set; using in-memory repositories")
		auditLog := audit.NewMemoryLog()
		orgRepo := tenancy.NewMemoryOrganizationRepository()
		tenantRepo := tenancy.NewMemoryTenantRepository()
		projectRepo := tenancy.NewMemoryProjectRepository()
		phase1.SeedTenancyCatalog(orgRepo, tenantRepo, projectRepo)
		return auditLog, orgRepo, tenantRepo, projectRepo, phase1.NewSeedStore(auditLog)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("failed to open postgres", "error", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("failed to ping postgres", "error", err)
		os.Exit(1)
	}
	if err := postgres.Migrate(ctx, db); err != nil {
		slog.Error("failed to migrate postgres", "error", err)
		os.Exit(1)
	}

	auditLog := postgres.NewAuditLog(db)
	phase1Store, err := phase1.NewPersistentStore(db, auditLog)
	if err != nil {
		slog.Error("failed to initialize phase1 store", "error", err)
		os.Exit(1)
	}
	orgRepo := postgres.NewOrganizationRepository(db)
	tenantRepo := postgres.NewTenantRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	phase1.SeedTenancyCatalog(orgRepo, tenantRepo, projectRepo)

	return auditLog, orgRepo, tenantRepo, projectRepo, phase1Store
}
