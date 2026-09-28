package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dzikru07/crews-system/internal/access"
	"github.com/dzikru07/crews-system/internal/api_registry"
	"github.com/dzikru07/crews-system/internal/application"
	"github.com/dzikru07/crews-system/internal/auth"
	"github.com/dzikru07/crews-system/internal/cache"
	"github.com/dzikru07/crews-system/internal/client"
	"github.com/dzikru07/crews-system/internal/config"
	"github.com/dzikru07/crews-system/internal/credential"
	"github.com/dzikru07/crews-system/internal/database"
	"github.com/dzikru07/crews-system/internal/environment"
	"github.com/dzikru07/crews-system/internal/identity"
	"github.com/dzikru07/crews-system/internal/identity_role"
	"github.com/dzikru07/crews-system/internal/middleware"
	"github.com/dzikru07/crews-system/internal/permission"
	"github.com/dzikru07/crews-system/internal/role"
	"github.com/dzikru07/crews-system/internal/role_permission"
	"github.com/dzikru07/crews-system/internal/tenant"
	"github.com/dzikru07/crews-system/internal/user_account"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    config.Config
	db     *database.Database
	cache  *cache.Redis
	router *gin.Engine
	http   *http.Server

	tenantHandler         *tenant.Handler
	applicationHandler    *application.Handler
	environmentHandler    *environment.Handler
	clientHandler         *client.Handler
	identityHandler       *identity.Handler
	credentialHandler     *credential.Handler
	apiRegistryHandler    *api_registry.Handler
	tokenManager          *auth.TokenManager
	userAccountHandler    *user_account.Handler
	accessHandler         *access.Handler
	roleHandler           *role.Handler
	permissionHandler     *permission.Handler
	identityRoleHandler   *identity_role.Handler
	rolePermissionHandler *role_permission.Handler
}

func New(cfg config.Config, db *database.Database, cache *cache.Redis) *Server {
	router := gin.New()

	// =========================================================
	// TENANT
	// =========================================================

	tenantRepository := tenant.NewRepository(db)

	tenantService := tenant.NewService(
		tenantRepository,
	)

	tenantHandler := tenant.NewHandler(
		tenantService,
	)

	// =========================================================
	// APPLICATION
	// =========================================================

	applicationRepository := application.NewRepository(db)

	applicationService := application.NewService(
		applicationRepository,
	)

	applicationHandler := application.NewHandler(
		applicationService,
	)

	// =========================================================
	// ENVIRONMENT
	// =========================================================

	environmentRepository := environment.NewRepository(db)

	environmentService := environment.NewService(
		environmentRepository,
	)

	environmentHandler := environment.NewHandler(
		environmentService,
	)

	// =========================================================
	// CLIENT
	// =========================================================

	clientRepository := client.NewRepository(db)

	clientService := client.NewService(
		clientRepository,
	)

	clientHandler := client.NewHandler(
		clientService,
	)

	// =========================================================
	// IDENTITY
	// =========================================================

	identityRepository := identity.NewRepository(db)

	identityService := identity.NewService(
		identityRepository,
	)

	identityHandler := identity.NewHandler(
		identityService,
	)

	// =========================================================
	// CREDENTIAL
	// =========================================================

	credentialRepository := credential.NewRepository(db)

	credentialService := credential.NewService(
		credentialRepository,
	)

	credentialHandler := credential.NewHandler(
		credentialService,
		identityRepository,
	)

	// =========================================================
	// API REGISTRY
	// =========================================================

	apiRegistryRepository := api_registry.NewRepository(db)

	apiRegistryService := api_registry.NewService(
		apiRegistryRepository,
	)

	apiRegistryHandler := api_registry.NewHandler(
		apiRegistryService,
	)

	// =========================================================
	// TOKEN
	// =========================================================

	tokenManager := auth.NewTokenManager(
		cfg.JWT,
	)

	// =========================================================
	// USER ACCOUNT
	// =========================================================

	userAccountRepository := user_account.NewRepository(db)

	userAccountService := user_account.NewService(
		userAccountRepository,
		identityRepository,
	)

	userAccountHandler := user_account.NewHandler(
		userAccountService,
	)

	// =========================================================
	// AUTH
	// =========================================================

	authRepo := auth.NewRepository(
		credentialRepository,
		userAccountRepository,
		identityRepository,
	)

	authService := auth.NewService(
		authRepo,
		tokenManager,
	)

	authHandler := auth.NewHandler(
		authService,
	)

	// =========================================================
	// ACCESS
	// =========================================================

	accessRepository := access.NewRepository(
		db,
		apiRegistryRepository,
	)

	accessService := access.NewService(
		accessRepository,
	)

	accessHandler := access.NewHandler(
		accessService,
	)

	// =========================================================
	// MIDDLEWARE
	// =========================================================

	// ==
	// role permissions
	// ==

	roleRepository := role.NewRepository(db)

	roleService := role.NewService(
		roleRepository,
	)

	roleHandler := role.NewHandler(
		roleService,
	)

	permissionRepository := permission.NewRepository(db)

	permissionService := permission.NewService(
		permissionRepository,
	)

	permissionHandler := permission.NewHandler(
		permissionService,
	)

	// ===

	identityRoleRepository := identity_role.NewRepository(db)

	identityRoleService := identity_role.NewService(
		identityRoleRepository,
	)

	identityRoleHandler := identity_role.NewHandler(
		identityRoleService,
	)

	rolePermissionRepository := role_permission.NewRepository(db)

	rolePermissionService := role_permission.NewService(
		rolePermissionRepository,
	)

	rolePermissionHandler := role_permission.NewHandler(
		rolePermissionService,
	)

	router.Use(
		middleware.RequestID(),
		middleware.Recovery(),
		middleware.SecurityHeaders(),
	)

	// =========================================================
	// SERVER
	// =========================================================

	server := &Server{
		cfg:    cfg,
		db:     db,
		cache:  cache,
		router: router,

		http: &http.Server{
			Addr: fmt.Sprintf(
				"%s:%s",
				cfg.HTTP.Host,
				cfg.HTTP.Port,
			),
			Handler: router,
		},

		tenantHandler:      tenantHandler,
		applicationHandler: applicationHandler,
		environmentHandler: environmentHandler,
		clientHandler:      clientHandler,
		identityHandler:    identityHandler,
		credentialHandler:  credentialHandler,
		apiRegistryHandler: apiRegistryHandler,
		tokenManager:       tokenManager,
		userAccountHandler: userAccountHandler,
		accessHandler:      accessHandler,
		roleHandler:        roleHandler,
		permissionHandler:  permissionHandler,

		identityRoleHandler:   identityRoleHandler,
		rolePermissionHandler: rolePermissionHandler,
	}

	server.registerRoutes()

	// =========================================================
	// AUTH ROUTES
	// =========================================================

	server.router.POST(
		"/v1/oauth/token",
		authHandler.Token,
	)

	server.router.POST(
		"/v1/auth/login",
		authHandler.Login,
	)

	server.router.GET(
		"/v1/auth/me",
		auth.Middleware(tokenManager),
		func(c *gin.Context) {

			authContext, ok := auth.GetAuthContext(c)

			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "authentication_context_missing",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"tenant_id":     authContext.TenantID,
				"identity_id":   authContext.IdentityID,
				"client_id":     authContext.ClientID,
				"identity_type": authContext.IdentityType,
				"environment":   authContext.Environment,
			})
		},
	)

	return server
}

func (s *Server) registerRoutes() {

	// =========================================================
	// HEALTH
	// =========================================================

	s.router.GET("/health", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	s.router.GET("/ready", func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			2*time.Second,
		)

		defer cancel()

		if err := s.db.Ping(ctx); err != nil {

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "not_ready",
				"database": "unavailable",
			})

			return
		}

		if err := s.cache.Ping(ctx); err != nil {

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"redis":  "unavailable",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ready",
			"database": "ok",
			"redis":    "ok",
		})
	})

	// =========================================================
	// API V1
	// =========================================================

	api := s.router.Group("/v1")

	// =========================================================
	// AUTHENTICATED API
	// =========================================================

	protected := api.Group("")

	protected.Use(
		auth.Middleware(s.tokenManager),
	)

	// =========================================================
	// TENANT
	// =========================================================

	tenants := protected.Group("/tenants")

	tenants.POST(
		"",
		s.tenantHandler.Create,
	)

	tenants.GET(
		"/:tenant_id",
		s.tenantHandler.GetByID,
	)

	// =========================================================
	// APPLICATION
	// =========================================================

	tenantApplications := tenants.Group(
		"/:tenant_id/applications",
	)

	tenantApplications.POST(
		"",
		s.applicationHandler.Create,
	)

	tenantApplications.GET(
		"/:application_id",
		s.applicationHandler.GetByID,
	)

	// =========================================================
	// ENVIRONMENT
	// =========================================================

	applicationEnvironments := tenantApplications.Group(
		"/:application_id/environments",
	)

	applicationEnvironments.POST(
		"",
		s.environmentHandler.Create,
	)

	applicationEnvironments.GET(
		"/:environment_id",
		s.environmentHandler.GetByID,
	)

	// =========================================================
	// CLIENT
	// =========================================================

	environmentClients := applicationEnvironments.Group(
		"/:environment_id/clients",
	)

	environmentClients.POST(
		"",
		s.clientHandler.Create,
	)

	environmentClients.GET(
		"/:client_id",
		s.clientHandler.GetByID,
	)

	// =========================================================
	// IDENTITY
	// =========================================================

	identities := tenants.Group(
		"/:tenant_id/identities",
	)

	identities.POST(
		"",
		s.identityHandler.Create,
	)

	identities.GET(
		"/:identity_id",
		s.identityHandler.GetByID,
	)

	// =========================================================
	// MACHINE CREDENTIAL
	// =========================================================

	identities.POST(
		"/:identity_id/credentials",
		s.credentialHandler.Create,
	)

	// =========================================================
	// USER ACCOUNT
	// =========================================================

	identities.POST(
		"/:identity_id/user-account",
		s.userAccountHandler.Create,
	)

	// =========================================================
	// API REGISTRY
	// =========================================================

	protected.POST(
		"/applications/:application_id/services",
		s.apiRegistryHandler.CreateService,
	)

	protected.POST(
		"/services/:service_id/endpoints",
		s.apiRegistryHandler.CreateEndpoint,
	)

	// =========================================================
	// ENDPOINT PERMISSIONS
	// =========================================================

	protected.POST(
		"/endpoints/:endpoint_id/permissions",
		s.accessHandler.AssignPermission,
	)

	protected.GET(
		"/endpoints/:endpoint_id/permissions",
		s.accessHandler.ListPermissions,
	)

	protected.DELETE(
		"/endpoints/:endpoint_id/permissions/:permission_id",
		s.accessHandler.RemovePermission,
	)

	// == permission role
	protected.POST(
		"/tenants/:tenant_id/roles",
		s.roleHandler.Create,
	)

	protected.GET(
		"/tenants/:tenant_id/roles/:role_id",
		s.roleHandler.GetByID,
	)

	protected.POST(
		"/tenants/:tenant_id/permissions",
		s.permissionHandler.Create,
	)

	protected.GET(
		"/tenants/:tenant_id/permissions/:permission_id",
		s.permissionHandler.GetByID,
	)

	// =========================================================
	// IDENTITY ROLES
	// =========================================================

	protected.POST(
		"/identities/:identity_id/roles",
		s.identityRoleHandler.AssignRole,
	)

	protected.GET(
		"/identities/:identity_id/roles",
		s.identityRoleHandler.ListRoles,
	)

	protected.DELETE(
		"/identities/:identity_id/roles/:role_id",
		s.identityRoleHandler.RemoveRole,
	)

	// =========================================================
	// ROLE PERMISSIONS
	// =========================================================

	protected.POST(
		"/roles/:role_id/permissions",
		s.rolePermissionHandler.AssignPermission,
	)

	protected.GET(
		"/roles/:role_id/permissions",
		s.rolePermissionHandler.ListPermissions,
	)

	protected.DELETE(
		"/roles/:role_id/permissions/:permission_id",
		s.rolePermissionHandler.RemovePermission,
	)
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
