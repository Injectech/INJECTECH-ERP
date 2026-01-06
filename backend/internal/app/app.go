package app

import (
	"context"
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/infrastructure/repository/postgres"
	"backend/internal/logger"
	"backend/internal/pkg/database"
	"backend/internal/transport/http/handler"
	"backend/internal/transport/http/router"
	usecaseaudit "backend/internal/usecase/audit"
	usecaseauth "backend/internal/usecase/auth"
	usecaseinventory "backend/internal/usecase/inventory"
	usecasepermission "backend/internal/usecase/permission"
	usecaseproduct "backend/internal/usecase/product"
	usecaserole "backend/internal/usecase/role"
	usecaseuser "backend/internal/usecase/user"
)

// Run bootstraps dependencies and starts the HTTP server.
func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logg, err := logger.New(cfg.Env)
	if err != nil {
		return err
	}
	defer logg.Sync() // flushes buffer, if any

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	productRepo := postgres.NewProductRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	permissionRepo := postgres.NewPermissionRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)
	auditRepo := postgres.NewAuditRepository(db)

	userUC := usecaseuser.NewService(userRepo)
	authUC := usecaseauth.NewService(userRepo, cfg.AccessTTL, cfg.RefreshTTL, cfg.JWTAccessSecret, cfg.JWTRefreshSecret)
	productUC := usecaseproduct.NewService(productRepo)
	roleUC := usecaserole.NewService(roleRepo)
	permissionUC := usecasepermission.NewService(permissionRepo)
	inventoryUC := usecaseinventory.NewService(inventoryRepo)
	auditUC := usecaseaudit.NewService(auditRepo)

	authHandler := handler.NewAuthHandler(authUC)
	userHandler := handler.NewUserHandler(userUC)
	productHandler := handler.NewProductHandler(productUC)
	roleHandler := handler.NewRoleHandler(roleUC)
	permissionHandler := handler.NewPermissionHandler(permissionUC)
	inventoryHandler := handler.NewInventoryHandler(inventoryUC)
	auditHandler := handler.NewAuditHandler(auditUC)

	engine := router.Build(logg, authUC, authHandler, userHandler, productHandler, roleHandler, permissionHandler, inventoryHandler, auditHandler)

	srv := &http.Server{
		Addr:              cfg.MustPort(),
		Handler:           engine,
		ReadHeaderTimeout: cfg.ReadTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
