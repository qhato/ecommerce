module github.com/qhato/ecommerce

go 1.24.0

require (
	github.com/expr-lang/expr v1.17.6
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-chi/cors v1.2.2
	github.com/go-playground/validator/v10 v10.28.0
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/redis/go-redis/v9 v9.17.1
	github.com/shopspring/decimal v1.4.0
	github.com/spf13/viper v1.21.0
	go.uber.org/zap v1.27.1
	golang.org/x/crypto v0.45.0
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
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.10 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)
