# E-Commerce Platform - Go Edition

Plataforma de comercio electrónico construida en Go con arquitectura hexagonal y Domain-Driven Design (DDD), basada en el schema de Broadleaf Commerce v7 Community.

## 🏗️ Arquitectura

Este proyecto implementa una arquitectura hexagonal (Ports & Adapters) con DDD, separando claramente las responsabilidades:

```
ecommerce/
├── cmd/                          # Entry points
│   ├── admin/                    # Admin API server
│   └── storefront/               # Storefront API server (public)
├── internal/                     # Bounded contexts
│   └── catalog/                  # Catalog bounded context
│       ├── domain/               # Domain layer (entities, value objects, interfaces)
│       ├── application/          # Application layer (use cases, DTOs)
│       │   ├── commands/         # Write operations (CQRS)
│       │   └── queries/          # Read operations (CQRS)
│       ├── infrastructure/       # Infrastructure layer (repositories, external services)
│       │   └── persistence/      # Database implementations
│       └── ports/                # Ports layer (HTTP handlers, gRPC, etc.)
│           └── http/             # HTTP handlers
├── pkg/                          # Shared kernel (utilities compartidas)
│   ├── database/                 # Database connection pool
│   ├── cache/                    # Cache implementations
│   ├── event/                    # Event bus
│   ├── logger/                   # Structured logging
│   ├── errors/                   # Error handling
│   ├── middleware/               # HTTP middleware
│   └── validator/                # Request validation
├── config/                       # Configuration management
└── scripts/                      # Utility scripts
```

## 🚀 Inicio Rápido con Docker

La forma más rápida de ejecutar la plataforma completa:

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd ecommerce

# 2. Crear config.yaml desde el ejemplo
cp config.example.yaml config.yaml

# 3. Iniciar todos los servicios
docker-compose up -d

# 4. Verificar que los servicios estén corriendo
docker-compose ps
```

Los servicios estarán disponibles en:
- **Admin API**: http://localhost:8080
- **Storefront API**: http://localhost:8081
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## 📦 Instalación Local

### Prerrequisitos

- Go 1.21 o superior
- PostgreSQL 12+
- Redis (opcional)
- Make

### Configuración

1. **Instalar dependencias**:
```bash
make install
```

2. **Configurar la base de datos**:
```bash
# Opción 1: Usando Make
make db-create
make db-migrate

# Opción 2: Usando el script de migraciones
chmod +x scripts/migrate.sh
./scripts/migrate.sh create
./scripts/migrate.sh migrate

# Opción 3: Reset completo (drop + create + migrate)
make db-reset
# o
./scripts/migrate.sh reset
```

3. **Configurar variables de entorno**:
```bash
cp config.example.yaml config.yaml
# Editar config.yaml con tus configuraciones
```

## 🏃 Ejecución

### Modo Desarrollo

```bash
# Iniciar Admin API
make run-admin
# o
go run cmd/admin/main.go

# Iniciar Storefront API (en otra terminal)
make run-storefront
# o
go run cmd/storefront/main.go
```

### Compilar Binarios

#### Compilación Local (Desarrollo)

```bash
# Compilar ambos binarios
make build

# Compilar solo Admin
make build-admin

# Compilar solo Storefront
make build-storefront

# Ejecutar binarios compilados
./bin/admin
./bin/storefront
```

#### Compilación Multi-Plataforma

Compilar para Linux, macOS y Windows:

```bash
# Compilar para todas las plataformas (Linux, macOS, Windows)
make build-all-platforms

# Compilar solo para Linux (amd64 + arm64)
make build-linux

# Compilar solo para macOS (amd64 + arm64/M1)
make build-macos

# Compilar solo para Windows (amd64)
make build-windows

# Crear archivos de release (.tar.gz)
make build-release

# Limpiar directorio de build
make build-clean
```

Los binarios se generan en `build/`:
```
build/
├── linux-amd64/
├── linux-arm64/
├── darwin-amd64/
├── darwin-arm64/
└── windows-amd64/
```

Ver [scripts/BUILD.md](scripts/BUILD.md) para más detalles.

#### Producción

```bash
# Compilar con optimizaciones para producción (Linux)
make build-prod

# Compilar con versión específica para todas las plataformas
VERSION=2.0.0 make build-release

# Los binarios estarán en bin/ o build/ según el comando
```

## 🐳 Docker

### Comandos útiles

```bash
# Construir imágenes Docker
make docker-build

# Iniciar servicios
make docker-up

# Detener servicios
make docker-down

# Ver logs
make docker-logs

# Ver estado de contenedores
make docker-ps
```

## 🔧 Makefile Targets

```bash
make help              # Mostrar todos los comandos disponibles
make install           # Instalar dependencias
make build             # Compilar todos los binarios
make run-admin         # Ejecutar Admin API
make run-storefront    # Ejecutar Storefront API
make test              # Ejecutar tests
make test-coverage     # Ejecutar tests con reporte de cobertura
make clean             # Limpiar binarios y archivos temporales
make fmt               # Formatear código
make lint              # Ejecutar linter

# Compilación Multi-Plataforma
make build-all-platforms    # Compilar para Linux, macOS, Windows
make build-linux            # Compilar solo para Linux
make build-macos            # Compilar solo para macOS
make build-windows          # Compilar solo para Windows
make build-release          # Crear archivos de release
make build-clean            # Limpiar directorio build/
make build-admin-only       # Solo admin para todas las plataformas
make build-storefront-only  # Solo storefront para todas las plataformas

# Docker
make docker-build      # Construir imágenes Docker
make docker-up         # Iniciar servicios con Docker Compose
make docker-down       # Detener servicios
make docker-logs       # Ver logs de Docker

# Base de datos
make db-create         # Crear base de datos
make db-drop           # Eliminar base de datos
make db-migrate        # Ejecutar migraciones
make db-reset          # Reset completo (drop + create + migrate)
make db-shell          # Abrir shell de PostgreSQL
```

## 📡 API Endpoints

### Admin API (Puerto 8080) - CRUD Completo

#### Productos

```
POST   /admin/products              # Crear producto
GET    /admin/products              # Listar productos (paginado)
GET    /admin/products/{id}         # Obtener producto por ID
PUT    /admin/products/{id}         # Actualizar producto
DELETE /admin/products/{id}         # Eliminar producto (soft delete)
POST   /admin/products/{id}/archive # Archivar producto
GET    /admin/products/search       # Buscar productos (?q=query)
```

#### Categorías

```
POST   /admin/categories                # Crear categoría
GET    /admin/categories                # Listar categorías (paginado)
GET    /admin/categories/root           # Listar categorías raíz
GET    /admin/categories/{id}           # Obtener categoría por ID
PUT    /admin/categories/{id}           # Actualizar categoría
DELETE /admin/categories/{id}           # Eliminar categoría
GET    /admin/categories/{id}/children  # Listar subcategorías
GET    /admin/categories/{id}/path      # Obtener ruta completa
```

#### SKUs

```
POST   /admin/skus                     # Crear SKU
GET    /admin/skus                     # Listar SKUs (paginado)
GET    /admin/skus/{id}                # Obtener SKU por ID
PUT    /admin/skus/{id}                # Actualizar SKU
DELETE /admin/skus/{id}                # Eliminar SKU
PUT    /admin/skus/{id}/pricing        # Actualizar pricing
PUT    /admin/skus/{id}/availability   # Actualizar disponibilidad
GET    /admin/skus/upc/{upc}           # Buscar SKU por UPC
GET    /admin/skus/product/{product_id} # Listar SKUs de un producto
```

### Storefront API (Puerto 8081) - Solo Lectura

#### Productos

```
GET /catalog/products                 # Listar productos activos
GET /catalog/products/{id}            # Obtener producto por ID
GET /catalog/products/url/{url}       # Obtener producto por URL
GET /catalog/products/search          # Buscar productos (?q=query)
```

#### Categorías

```
GET /catalog/categories               # Listar categorías raíz activas
GET /catalog/categories/{id}          # Obtener categoría por ID
GET /catalog/categories/url/{url}     # Obtener categoría por URL
GET /catalog/categories/{id}/children # Listar subcategorías
GET /catalog/categories/{id}/products # Listar productos de categoría
GET /catalog/categories/{id}/path     # Obtener ruta completa
```

#### SKUs

```
GET /catalog/skus                      # Listar SKUs activos y disponibles
GET /catalog/skus/{id}                 # Obtener SKU por ID
GET /catalog/skus/upc/{upc}            # Buscar SKU por UPC
GET /catalog/skus/product/{product_id} # Listar SKUs de un producto
```

## 📝 Ejemplos de Uso

### Crear un Producto

```bash
curl -X POST http://localhost:8080/admin/products \
  -H "Content-Type: application/json" \
  -d '{
    "manufacture": "Apple",
    "model": "iPhone 15 Pro",
    "url": "/products/iphone-15-pro",
    "url_key": "iphone-15-pro",
    "can_sell_without_options": false,
    "enable_default_sku": true,
    "meta_title": "iPhone 15 Pro - The Ultimate Smartphone",
    "meta_description": "Experience the power of iPhone 15 Pro",
    "attributes": {
      "color": "Natural Titanium",
      "storage": "256GB"
    }
  }'
```

### Crear una Categoría

```bash
curl -X POST http://localhost:8080/admin/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smartphones",
    "description": "Latest smartphones and mobile devices",
    "url": "/categories/smartphones",
    "url_key": "smartphones",
    "root_display_order": 1.0,
    "meta_title": "Smartphones - Buy Latest Mobile Phones",
    "meta_description": "Browse our collection of smartphones"
  }'
```

### Crear un SKU

```bash
curl -X POST http://localhost:8080/admin/skus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro - 256GB - Natural Titanium",
    "description": "iPhone 15 Pro with 256GB storage",
    "upc": "195949038123",
    "currency_code": "USD",
    "price": 999.00,
    "retail_price": 1099.00,
    "sale_price": 949.00,
    "available": true,
    "discountable": true,
    "taxable": true,
    "attributes": {
      "color": "Natural Titanium",
      "storage": "256GB"
    }
  }'
```

### Listar Productos (Storefront)

```bash
# Listar productos con paginación
curl "http://localhost:8081/catalog/products?page=1&page_size=20"

# Buscar productos
curl "http://localhost:8081/catalog/products/search?q=iphone"

# Obtener producto por URL
curl "http://localhost:8081/catalog/products/url/iphone-15-pro"
```

## ✅ Características Implementadas

### Catalog Bounded Context (100% Completo)

#### Domain Layer ✅
- ✅ Product entity con métodos de dominio
- ✅ Category entity con jerarquía y activación por fechas
- ✅ SKU entity con pricing y disponibilidad
- ✅ Domain events (8 tipos)
- ✅ Repository interfaces
- ✅ Filters para paginación

#### Application Layer (CQRS) ✅
- ✅ **Product**: 4 Commands + 5 Queries
- ✅ **Category**: 3 Commands + 5 Queries
- ✅ **SKU**: 5 Commands + 4 Queries
- ✅ DTOs con conversores automáticos
- ✅ Validación de comandos
- ✅ Publicación de domain events

#### Infrastructure Layer ✅
- ✅ PostgreSQL Product Repository
- ✅ PostgreSQL Category Repository
- ✅ PostgreSQL SKU Repository
- ✅ Cache integration (Redis/Memory)
- ✅ Queries optimizadas con paginación

#### Ports Layer ✅
- ✅ Admin HTTP Handlers (CRUD completo)
- ✅ Storefront HTTP Handlers (read-only)
- ✅ Paginación en todos los endpoints
- ✅ Búsqueda y filtrado

### Shared Kernel ✅
- ✅ Configuration management
- ✅ Database connection pool con transacciones
- ✅ Cache abstraction (Redis + In-Memory)
- ✅ Event bus (In-Memory)
- ✅ Structured logging (Zap)
- ✅ Error handling
- ✅ HTTP utilities
- ✅ Middleware (Auth, Logging, Recovery, CORS)
- ✅ Validator

### DevOps & Tooling ✅
- ✅ Makefile con 20+ targets
- ✅ Docker Compose con PostgreSQL + Redis
- ✅ Dockerfiles multi-stage optimizados
- ✅ Migration scripts
- ✅ .gitignore y .dockerignore
- ✅ Documentación completa

## 📊 Estadísticas del Proyecto

- **Archivos Go**: ~30 archivos
- **Líneas de código**: ~7,000 líneas
- **Bounded Contexts**: 1/10 (Catalog completo)
- **Entidades de dominio**: 3 (Product, Category, SKU)
- **Repositorios**: 3 (Product, Category, SKU)
- **Commands**: 12 (Product: 4, Category: 3, SKU: 5)
- **Queries**: 14 (Product: 5, Category: 5, SKU: 4)
- **HTTP Handlers**: 4 (Admin Product, Admin Category, Admin SKU, Storefront)
- **Domain Events**: 8 tipos
- **API Endpoints**: 40+ endpoints

## 🔜 Pendiente

- [ ] Tests unitarios e integración
- [ ] Customer bounded context
- [ ] Order bounded context
- [ ] Payment bounded context
- [ ] Fulfillment bounded context
- [ ] Authentication & Authorization JWT completo
- [ ] Documentación OpenAPI/Swagger
- [ ] CI/CD pipeline
- [ ] Métricas y observabilidad (Prometheus/Grafana)
- [ ] Rate limiting
- [ ] API versioning

## 🎯 Patrones y Principios

- **Hexagonal Architecture**: Separación clara entre dominio, aplicación, infraestructura y puertos
- **Domain-Driven Design (DDD)**: Bounded contexts, entities, value objects, domain events
- **CQRS**: Separación de comandos (escritura) y queries (lectura)
- **Repository Pattern**: Abstracción del acceso a datos
- **Dependency Injection**: Inyección manual en entry points
- **Clean Code**: Código limpio y mantenible
- **SOLID Principles**: Single Responsibility, Open/Closed, etc.

## 🚀 Escalabilidad

El proyecto está diseñado para escalar horizontalmente:

- **Stateless APIs**: Sin estado en los servidores
- **Cache distribuido**: Redis para compartir cache entre instancias
- **Database pooling**: Connection pooling optimizado
- **Bounded Contexts**: Cada contexto puede ser microservicio independiente
- **Event-Driven**: Comunicación asíncrona vía eventos
- **Docker**: Fácil deployment en cualquier plataforma

## ⚡ Performance

- Structured logging con Zap (alto rendimiento)
- Cache de queries frecuentes (Redis/Memory)
- Connection pooling de PostgreSQL
- Queries optimizadas con índices
- Graceful shutdown
- Multi-stage Docker builds (imágenes pequeñas)

## 🧪 Testing

```bash
# Ejecutar tests
make test

# Tests con cobertura
make test-coverage

# Ver reporte HTML de cobertura
open coverage.html
```

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una branch para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la branch (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

[MIT License](LICENSE)

## 👥 Contacto

- GitHub: [@qhato](https://github.com/qhato)

---

**⭐ Si este proyecto te resulta útil, considera darle una estrella en GitHub!**
