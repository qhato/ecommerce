package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTPMetrics contains all HTTP-related metrics
type HTTPMetrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	RequestSize     *prometheus.HistogramVec
	ResponseSize    *prometheus.HistogramVec
	ErrorsTotal     *prometheus.CounterVec
}

// BusinessMetrics contains all business-related metrics
type BusinessMetrics struct {
	OrdersCreated      prometheus.Counter
	OrdersCompleted    prometheus.Counter
	OrdersCancelled    prometheus.Counter
	OrderValue         prometheus.Histogram
	ProductsCreated    prometheus.Counter
	ProductsArchived   prometheus.Counter
	CustomersCreated   prometheus.Counter
	PaymentsProcessed  prometheus.Counter
	PaymentsFailed     prometheus.Counter
	ShipmentsCreated   prometheus.Counter
	ShipmentsDelivered prometheus.Counter
}

// DatabaseMetrics contains all database-related metrics
type DatabaseMetrics struct {
	QueriesTotal    *prometheus.CounterVec
	QueryDuration   *prometheus.HistogramVec
	ConnectionsOpen prometheus.Gauge
	ConnectionsIdle prometheus.Gauge
}

// CacheMetrics contains all cache-related metrics
type CacheMetrics struct {
	HitsTotal   prometheus.Counter
	MissesTotal prometheus.Counter
	ErrorsTotal prometheus.Counter
	Latency     prometheus.Histogram
}

var (
	// HTTP is the singleton instance for HTTP metrics
	HTTP *HTTPMetrics

	// Business is the singleton instance for business metrics
	Business *BusinessMetrics

	// Database is the singleton instance for database metrics
	Database *DatabaseMetrics

	// Cache is the singleton instance for cache metrics
	Cache *CacheMetrics
)

// Init initializes all metrics
func Init(namespace string) {
	HTTP = initHTTPMetrics(namespace)
	Business = initBusinessMetrics(namespace)
	Database = initDatabaseMetrics(namespace)
	Cache = initCacheMetrics(namespace)
}

func initHTTPMetrics(namespace string) *HTTPMetrics {
	return &HTTPMetrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latency in seconds",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),
		RequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_size_bytes",
				Help:      "HTTP request size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 7), // 100 bytes to 100MB
			},
			[]string{"method", "path"},
		),
		ResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 7),
			},
			[]string{"method", "path"},
		),
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_errors_total",
				Help:      "Total number of HTTP errors",
			},
			[]string{"method", "path", "error_type"},
		),
	}
}

func initBusinessMetrics(namespace string) *BusinessMetrics {
	return &BusinessMetrics{
		OrdersCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "orders_created_total",
			Help:      "Total number of orders created",
		}),
		OrdersCompleted: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "orders_completed_total",
			Help:      "Total number of orders completed",
		}),
		OrdersCancelled: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "orders_cancelled_total",
			Help:      "Total number of orders cancelled",
		}),
		OrderValue: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "order_value_dollars",
			Help:      "Order value distribution in dollars",
			Buckets:   []float64{10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
		}),
		ProductsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "products_created_total",
			Help:      "Total number of products created",
		}),
		ProductsArchived: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "products_archived_total",
			Help:      "Total number of products archived",
		}),
		CustomersCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "customers_created_total",
			Help:      "Total number of customers created",
		}),
		PaymentsProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "payments_processed_total",
			Help:      "Total number of payments processed",
		}),
		PaymentsFailed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "payments_failed_total",
			Help:      "Total number of failed payments",
		}),
		ShipmentsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "shipments_created_total",
			Help:      "Total number of shipments created",
		}),
		ShipmentsDelivered: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "shipments_delivered_total",
			Help:      "Total number of shipments delivered",
		}),
	}
}

func initDatabaseMetrics(namespace string) *DatabaseMetrics {
	return &DatabaseMetrics{
		QueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "database_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "table"},
		),
		QueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "database_query_duration_seconds",
				Help:      "Database query latency in seconds",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
			},
			[]string{"operation", "table"},
		),
		ConnectionsOpen: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "database_connections_open",
			Help:      "Number of open database connections",
		}),
		ConnectionsIdle: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "database_connections_idle",
			Help:      "Number of idle database connections",
		}),
	}
}

func initCacheMetrics(namespace string) *CacheMetrics {
	return &CacheMetrics{
		HitsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "cache_hits_total",
			Help:      "Total number of cache hits",
		}),
		MissesTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "cache_misses_total",
			Help:      "Total number of cache misses",
		}),
		ErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "cache_errors_total",
			Help:      "Total number of cache errors",
		}),
		Latency: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "cache_latency_seconds",
			Help:      "Cache operation latency in seconds",
			Buckets:   []float64{.0001, .0005, .001, .0025, .005, .01, .025, .05, .1},
		}),
	}
}

// RecordHTTPRequest records an HTTP request with all its metrics
func RecordHTTPRequest(method, path, status string, duration time.Duration, requestSize, responseSize int64) {
	if HTTP == nil {
		return
	}

	HTTP.RequestsTotal.WithLabelValues(method, path, status).Inc()
	HTTP.RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	HTTP.RequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	HTTP.ResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
}

// RecordHTTPError records an HTTP error
func RecordHTTPError(method, path, errorType string) {
	if HTTP == nil {
		return
	}
	HTTP.ErrorsTotal.WithLabelValues(method, path, errorType).Inc()
}

// RecordDatabaseQuery records a database query
func RecordDatabaseQuery(operation, table string, duration time.Duration) {
	if Database == nil {
		return
	}
	Database.QueriesTotal.WithLabelValues(operation, table).Inc()
	Database.QueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// UpdateDatabaseConnections updates database connection metrics
func UpdateDatabaseConnections(open, idle int) {
	if Database == nil {
		return
	}
	Database.ConnectionsOpen.Set(float64(open))
	Database.ConnectionsIdle.Set(float64(idle))
}