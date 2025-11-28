# Multi-Platform Build Script

Script para compilar binarios de la plataforma e-commerce para múltiples sistemas operativos y arquitecturas.

## Características

- ✅ Compilación para **Linux** (amd64, arm64)
- ✅ Compilación para **macOS** (amd64, arm64/M1)
- ✅ Compilación para **Windows** (amd64)
- ✅ Binarios optimizados con flags `-ldflags="-s -w"`
- ✅ Información de versión embebida
- ✅ Generación automática de archivos `.tar.gz` para releases
- ✅ Soporte para compilación selectiva (solo admin o storefront)

## Uso

### Desde el script directamente

```bash
# Compilar para todas las plataformas
./scripts/build.sh all

# Compilar solo para Linux
./scripts/build.sh linux

# Compilar solo para macOS
./scripts/build.sh macos

# Compilar solo para Windows
./scripts/build.sh windows

# Crear archivos de release (compila + genera .tar.gz)
./scripts/build.sh release

# Limpiar directorio de build
./scripts/build.sh clean

# Ver ayuda
./scripts/build.sh help
```

### Desde Makefile (recomendado)

```bash
# Compilar para todas las plataformas
make build-all-platforms

# Compilar para Linux
make build-linux

# Compilar para macOS
make build-macos

# Compilar para Windows
make build-windows

# Crear release completo
make build-release

# Limpiar
make build-clean

# Solo compilar admin para todas las plataformas
make build-admin-only

# Solo compilar storefront para todas las plataformas
make build-storefront-only
```

## Opciones Avanzadas

### Compilar solo un binario

```bash
# Solo admin para todas las plataformas
./scripts/build.sh all --admin-only

# Solo storefront para todas las plataformas
./scripts/build.sh all --storefront-only

# Solo admin para Linux
./scripts/build.sh linux --admin-only
```

### Especificar versión

```bash
# Con variable de entorno
VERSION=2.0.0 ./scripts/build.sh all

# Con flag
./scripts/build.sh all --version 2.0.0

# Desde Makefile
VERSION=2.0.0 make build-all-platforms
```

## Estructura de Salida

Los binarios compilados se colocan en el directorio `build/`:

```
build/
├── linux-amd64/
│   ├── admin
│   └── storefront
├── linux-arm64/
│   ├── admin
│   └── storefront
├── darwin-amd64/
│   ├── admin
│   └── storefront
├── darwin-arm64/
│   ├── admin
│   └── storefront
└── windows-amd64/
    ├── admin.exe
    └── storefront.exe
```

### Con archivos de release

Al ejecutar `./scripts/build.sh release`, también se generan archivos comprimidos:

```
build/
├── ecommerce-linux-amd64-1.0.0.tar.gz
├── ecommerce-linux-arm64-1.0.0.tar.gz
├── ecommerce-darwin-amd64-1.0.0.tar.gz
├── ecommerce-darwin-arm64-1.0.0.tar.gz
└── ecommerce-windows-amd64-1.0.0.tar.gz
```

## Información Embebida en Binarios

Cada binario incluye información de versión embebida:

- **Version**: Número de versión (default: 1.0.0)
- **BuildTime**: Fecha y hora de compilación
- **GitCommit**: Hash corto del commit Git actual

Esta información se puede acceder desde el código con:

```go
var (
    Version   string
    BuildTime string
    GitCommit string
)

func main() {
    fmt.Printf("Version: %s\n", Version)
    fmt.Printf("Build Time: %s\n", BuildTime)
    fmt.Printf("Git Commit: %s\n", GitCommit)
}
```

## Plataformas Soportadas

| OS      | Arquitectura | Soporte | Notas                          |
|---------|-------------|---------|--------------------------------|
| Linux   | amd64       | ✅      | Servidores, desktops          |
| Linux   | arm64       | ✅      | ARM servers, Raspberry Pi     |
| macOS   | amd64       | ✅      | Intel Macs                    |
| macOS   | arm64       | ✅      | Apple Silicon (M1/M2/M3)      |
| Windows | amd64       | ✅      | Windows 64-bit                |

## Optimizaciones de Compilación

El script utiliza los siguientes flags de optimización:

```bash
CGO_ENABLED=0          # Compilación estática sin CGO
-ldflags="-s -w"       # Reducir tamaño del binario
-s                     # Omitir tabla de símbolos
-w                     # Omitir información de debug DWARF
```

Esto resulta en binarios más pequeños y portables.

## Ejemplos Prácticos

### 1. Compilar versión de producción

```bash
VERSION=1.0.0 ./scripts/build.sh release
```

### 2. Compilar solo para servidores Linux

```bash
./scripts/build.sh linux
```

### 3. Compilar para desarrollo en macOS M1

```bash
./scripts/build.sh macos
# Los binarios estarán en: build/darwin-arm64/
```

### 4. Crear release para distribución

```bash
# Limpiar builds anteriores
make build-clean

# Compilar con versión específica
VERSION=2.1.0 make build-release

# Los archivos .tar.gz estarán listos para distribución
ls -lh build/*.tar.gz
```

### 5. Compilar solo el admin API para Windows

```bash
./scripts/build.sh windows --admin-only
# El ejecutable estará en: build/windows-amd64/admin.exe
```

## Troubleshooting

### Error: Permission denied

```bash
chmod +x scripts/build.sh
```

### Error: GOOS not supported

Verifica que tu versión de Go soporte la plataforma objetivo:

```bash
go tool dist list
```

### Binario muy grande

Los binarios incluyen toda la aplicación. Tamaños típicos:
- Admin: ~15-25 MB
- Storefront: ~15-25 MB

Para reducir más, considera usar `upx`:

```bash
upx --best build/linux-amd64/admin
```

## Integración con CI/CD

### GitHub Actions

```yaml
name: Build

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build all platforms
        run: make build-all-platforms
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: binaries
          path: build/
```

### GitLab CI

```yaml
build:
  stage: build
  image: golang:1.21
  script:
    - make build-all-platforms
  artifacts:
    paths:
      - build/
```

## Notas Adicionales

- Los binarios de Windows tienen extensión `.exe`
- Los binarios de Linux y macOS no tienen extensión
- Todos los binarios son completamente estáticos (sin dependencias externas)
- Compatibles con Docker multi-stage builds
- Optimizados para tamaño y rendimiento
