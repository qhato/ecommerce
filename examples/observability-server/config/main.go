package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/qhato/ecommerce/pkg/health"
	"github.com/qhato/ecommerce/pkg/jwt"
	"github.com/qhato/ecommerce/pkg/logging"
	"github.com/qhato/ecommerce/pkg/metrics"
	"github.com/qhato/ecommerce/pkg/middleware"
	"github.com/qhato/ecommerce/pkg/ratelimit"
	"github.com/qhato/ecommerce/pkg/tracing"
)

// ServerConfig contains server configuration
type ServerConfig struct {
	Port            string
	Environment     string
	ServiceName     string
	ServiceVersion  string
	
	// Database
	DatabaseURL     string
	
	// Redis
	RedisURL        string
	
	// JWT
	JWTAccessSecret  string
	JWTRefreshSecret string
	
	// Tracing
	TracingEnabled   bool
	JaegerEndpoint   string
	
	// Rate Limiting
	RateLimitEnabled bool
	RateLimit        int
	RateLimitWindow  time.Duration
	
	// CORS
	AllowedOrigins   []string
}

// Server represents the HTTP server with observability
type Server struct {
	config        ServerConfig
	router        *chi.Mux
	httpServer    *http.Server
	db            *sql.DB
	redis         *redis.Client
	logger        logging.Logger
	jwtManager    *jwt.Manager
	healthManager *health.Manager
	tracingProvider *tracing.Provider
}

// NewServer creates a new server with observability
func NewServer(config ServerConfig) (*Server, error) {
	// Initialize logger
	var logger logging.Logger
	var err error
	if config.Environment == "production" {
		logger, err = logging.NewProductionLogger()
	} else {
		logger, err = logging.NewDevelopmentLogger()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Info("Initializing server",
		logging.String("service", config.ServiceName),
		logging.String("version", config.ServiceVersion),
		logging.String("environment", config.Environment),
	)

	// Initialize metrics
	metrics.Init(config.ServiceName)
	logger.Info("Metrics initialized")

	// Initialize tracing
	var tracingProvider *tracing.Provider
	if config.TracingEnabled {
		tracingProvider, err = tracing.Init(tracing.Config{
			ServiceName:    config.ServiceName,
			ServiceVersion: config.ServiceVersion,
			Environment:    config.Environment,
			ExporterType:   "jaeger",
			JaegerEndpoint: config.JaegerEndpoint,
			SamplingRate:   1.0, // Sample all requests in dev, adjust for prod
		})
		if err != nil {
			logger.Warn("Failed to initialize tracing, continuing without it",
				logging.Error(err))
		} else {
			logger.Info("Tracing initialized", logging.String("endpoint", config.JaegerEndpoint))
		}
	}

	// Initialize database (placeholder - implement based on your setup)
	// db, err := sql.Open("pgx", config.DatabaseURL)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to database: %w", err)
	// }
	var db *sql.DB // Replace with actual DB connection

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
	})
	
	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Warn("Failed to connect to Redis, continuing without it",
			logging.Error(err))
	} else {
		logger.Info("Redis connected")
	}

	// Initialize JWT manager
	blacklist := jwt.NewRedisBlacklist(redisClient)
	jwtManager := jwt.NewManager(jwt.Config{
		AccessSecret:  config.JWTAccessSecret,
		RefreshSecret: config.JWTRefreshSecret,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
		Issuer:        config.ServiceName,
	}, blacklist)
	logger.Info("JWT manager initialized")

	// Initialize health checks
	healthManager := health.NewManager()
	if db != nil {
		healthManager.Register("database", &health.DatabaseChecker{DB: db})
	}
	if redisClient != nil {
		healthManager.Register("redis", &health.RedisChecker{Client: redisClient})
	}
	logger.Info("Health checks configured")

	// Create server
	server := &Server{
		config:          config,
		db:              db,
		redis:           redisClient,
		logger:          logger,
		jwtManager:      jwtManager,
		healthManager:   healthManager,
		tracingProvider: tracingProvider,
	}

	// Setup routes
	server.setupRoutes()

	// Create HTTP server
	server.httpServer = &http.Server{
		Addr:         ":" + config.Port,
		Handler:      server.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server, nil
}

// setupRoutes configures all routes with middleware
func (s *Server) setupRoutes() {
	r := chi.NewRouter()

	// Global middleware (applied to all routes)
	r.Use(middleware.RequestID())
	r.Use(middleware.Security())
	
	// Logging middleware with structured logger
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(wrapped, r)
			
			duration := time.Since(start)
			logging.LogRequest(s.logger, r.Method, r.URL.Path, wrapped.statusCode, duration)
		})
	})

	// Tracing middleware
	if s.tracingProvider != nil {
		r.Use(middleware.Tracing(s.config.ServiceName))
	}

	// Metrics middleware
	r.Use(middleware.Metrics)

	// CORS
	if len(s.config.AllowedOrigins) > 0 {
		r.Use(middleware.CORS(
			s.config.AllowedOrigins,
			[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			[]string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		))
	}

	// Request size limit (10MB)
	r.Use(middleware.RequestSizeLimit(10 * 1024 * 1024))

	// Observability endpoints (no auth required)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/health", s.healthManager.Handler())
	r.Get("/health/live", health.LivenessHandler())
	r.Get("/health/ready", s.healthManager.ReadinessHandler())

	// Public API (with rate limiting)
	r.Group(func(r chi.Router) {
		if s.config.RateLimitEnabled {
			limiter := ratelimit.NewRedisLimiter(s.redis, ratelimit.Config{
				RequestsPerWindow: s.config.RateLimit,
				WindowSize:        s.config.RateLimitWindow,
			})
			r.Use(middleware.RateLimit(limiter, middleware.IPKeyFunc))
		}

		// Public endpoints
		r.Post("/api/v1/auth/login", s.handleLogin)
		r.Post("/api/v1/auth/register", s.handleRegister)
		r.Post("/api/v1/auth/refresh", s.handleRefresh)
		
		// Add your public endpoints here
	})

	// Protected API (requires authentication)
	r.Group(func(r chi.Router) {
		// JWT authentication middleware
		r.Use(middleware.JWTAuth(s.jwtManager))

		// User-based rate limiting (higher limits for authenticated users)
		if s.config.RateLimitEnabled {
			limiter := ratelimit.NewRedisLimiter(s.redis, ratelimit.Config{
				RequestsPerWindow: s.config.RateLimit * 5, // 5x limit for authenticated users
				WindowSize:        s.config.RateLimitWindow,
			})
			r.Use(middleware.RateLimit(limiter, middleware.UserKeyFunc))
		}

		// Protected endpoints
		r.Post("/api/v1/auth/logout", s.handleLogout)
		r.Get("/api/v1/profile", s.handleProfile)
		
		// Add your protected endpoints here
	})

	// Admin API (requires admin role)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(s.jwtManager))
		r.Use(middleware.RequireRole("admin", "superadmin"))

		// Admin endpoints
		r.Get("/api/v1/admin/users", s.handleAdminUsers)
		r.Post("/api/v1/admin/metrics/reset", s.handleMetricsReset)
		
		// Add your admin endpoints here
	})

	s.router = r
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting server",
		logging.String("port", s.config.Port),
		logging.String("environment", s.config.Environment),
	)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		s.logger.Info("Received shutdown signal", logging.String("signal", sig.String()))
		return s.Shutdown()
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP server shutdown error", logging.Error(err))
	}

	// Close database
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			s.logger.Error("Database close error", logging.Error(err))
		}
	}

	// Close Redis
	if s.redis != nil {
		if err := s.redis.Close(); err != nil {
			s.logger.Error("Redis close error", logging.Error(err))
		}
	}

	// Shutdown tracing
	if s.tracingProvider != nil {
		if err := s.tracingProvider.Shutdown(ctx); err != nil {
			s.logger.Error("Tracing shutdown error", logging.Error(err))
		}
	}

	s.logger.Info("Server stopped gracefully")
	return nil
}

// Handler implementations (placeholders)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Implement login logic
	ctx := r.Context()
	s.logger.WithContext(ctx).Info("Login attempt")
	
	// Example: Generate token pair
	// tokenPair, err := s.jwtManager.GenerateTokenPair(userID, email, roles)
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Login endpoint"}`))
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.WithContext(ctx).Info("Registration attempt")
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Register endpoint"}`))
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.WithContext(ctx).Info("Token refresh attempt")
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Refresh endpoint"}`))
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.WithContext(ctx).Info("Logout attempt")
	
	// Revoke token
	// authHeader := r.Header.Get("Authorization")
	// token := strings.TrimPrefix(authHeader, "Bearer ")
	// s.jwtManager.RevokeToken(ctx, token)
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logout successful"}`))
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.GetUserID(ctx)
	
	s.logger.WithContext(ctx).Info("Profile request", logging.Int64("user_id", userID))
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Profile endpoint"}`))
}

func (s *Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.WithContext(ctx).Info("Admin: list users")
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Admin users endpoint"}`))
}

func (s *Server) handleMetricsReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.WithContext(ctx).Info("Admin: reset metrics")
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Metrics reset"}`))
}

func main() {
	config := ServerConfig{
		Port:             "8080",
		Environment:      "development",
		ServiceName:      "ecommerce-api",
		ServiceVersion:   "1.0.0",
		DatabaseURL:      "postgres://localhost:5432/ecommerce",
		RedisURL:         "localhost:6379",
		JWTAccessSecret:  "your-access-secret-key-change-in-production",
		JWTRefreshSecret: "your-refresh-secret-key-change-in-production",
		TracingEnabled:   true,
		JaegerEndpoint:   "http://localhost:14268/api/traces",
		RateLimitEnabled: true,
		RateLimit:        100,
		RateLimitWindow:  time.Minute,
		AllowedOrigins:   []string{"http://localhost:3000"},
	}

	server, err := NewServer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
