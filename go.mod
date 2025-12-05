module github.com/qhato/ecommerce

go 1.24.0

require (
	github.com/elastic/go-elasticsearch/v8 v8.19.0
	github.com/expr-lang/expr v1.17.6
	github.com/gin-gonic/gin v1.11.0
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-chi/cors v1.2.2
	github.com/go-playground/validator/v10 v10.28.0
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/lib/pq v1.10.9
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.23.2
	github.com/redis/go-redis/v9 v9.17.1
	github.com/shopspring/decimal v1.4.0
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.11.1
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0
	go.opentelemetry.io/otel/sdk v1.38.0
	go.opentelemetry.io/otel/trace v1.38.0
	go.uber.org/zap v1.27.1
	golang.org/x/crypto v0.45.0
	google.golang.org/grpc v1.77.0
)

replace github.com/qhato/ecommerce/internal/catalog/application => ./internal/catalog/application

replace github.com/qhato/ecommerce/internal/catalog/application/commands => ./internal/catalog/application/commands

replace github.com/qhato/ecommerce/internal/catalog/application/queries => ./internal/catalog/application/queries

replace github.com/qhato/ecommerce/internal/catalog/infrastructure/persistence => ./internal/catalog/infrastructure/persistence

replace github.com/qhato/ecommerce/internal/catalog/ports/http => ./internal/catalog/ports/http

replace github.com/qhato/ecommerce/internal/customer/application => ./internal/customer/application

replace github.com/qhato/ecommerce/internal/customer/application/commands => ./internal/customer/application/commands

replace github.com/qhato/ecommerce/internal/customer/application/queries => ./internal/customer/application/queries

replace github.com/qhato/ecommerce/internal/customer/infrastructure/persistence => ./internal/customer/infrastructure/persistence

replace github.com/qhato/ecommerce/internal/customer/ports/http => ./internal/customer/ports/http

replace github.com/qhato/ecommerce/internal/order/application => ./internal/order/application

replace github.com/qhato/ecommerce/internal/order/application/commands => ./internal/order/application/commands

replace github.com/qhato/ecommerce/internal/order/application/queries => ./internal/order/application/queries

replace github.com/qhato/ecommerce/internal/order/infrastructure/persistence => ./internal/order/infrastructure/persistence

replace github.com/qhato/ecommerce/internal/order/ports/http => ./internal/order/ports/http

replace github.com/qhato/ecommerce/internal/offer/application => ./internal/offer/application

replace github.com/qhato/ecommerce/internal/offer/domain => ./internal/offer/domain

replace github.com/qhato/ecommerce/internal/offer/infrastructure/persistence => ./internal/offer/infrastructure/persistence

replace github.com/qhato/ecommerce/internal/inventory/application => ./internal/inventory/application

replace github.com/qhato/ecommerce/internal/inventory/domain => ./internal/inventory/domain

replace github.com/qhato/ecommerce/internal/inventory/infrastructure/persistence => ./internal/inventory/infrastructure/persistence

replace github.com/qhato/ecommerce/internal/tax/application => ./internal/tax/application

replace github.com/qhato/ecommerce/internal/payment/application/commands => ./internal/payment/application/commands

replace github.com/qhato/ecommerce/internal/payment/application/queries => ./internal/payment/application/queries

replace github.com/qhato/ecommerce/internal/payment/infrastructure/persistence => ./internal/payment/infrastructure/persistence

replace github.com/qhato/ecommerce/internal/payment/ports/http => ./internal/payment/ports/http

replace github.com/qhato/ecommerce/internal/fulfillment/application/commands => ./internal/fulfillment/application/commands

replace github.com/qhato/ecommerce/internal/fulfillment/infrastructure/persistence => ./internal/fulfillment/infrastructure/persistence

replace github.com/qhato/ecommerce/internal/fulfillment/ports/http => ./internal/fulfillment/ports/http

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/elastic-transport-go/v8 v8.7.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.10 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.1 // indirect
	go.uber.org/mock v0.5.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/arch v0.20.0 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
