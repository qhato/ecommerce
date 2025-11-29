#!/bin/bash

echo "🔧 Script de Corrección de Errores de Compilación"
echo "=================================================="
echo ""

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Función para imprimir con color
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# 1. Instalar herramientas necesarias
echo "1. Instalando herramientas..."
go install golang.org/x/tools/cmd/goimports@latest
print_status "goimports instalado"

# 2. Fix imports automáticamente
echo ""
echo "2. Corrigiendo imports automáticamente..."
~/go/bin/goimports -w internal/ pkg/ cmd/ 2>/dev/null && print_status "Imports corregidos" || print_warning "Algunos imports no se pudieron corregir automáticamente"

# 3. go mod tidy
echo ""
echo "3. Limpiando dependencias..."
go mod tidy && print_status "Dependencias limpias"

# 4. Compilar y contar errores
echo ""
echo "4. Compilando proyecto..."
ERRORS=$(go build ./... 2>&1)
ERROR_COUNT=$(echo "$ERRORS" | grep -E "undefined|too many|not enough|redeclared|declared and not used|imported and not used" | wc -l | tr -d ' ')

if [ "$ERROR_COUNT" -eq "0" ]; then
    print_status "¡Compilación exitosa! Sin errores."
    exit 0
else
    print_warning "Quedan $ERROR_COUNT errores de compilación"
    echo ""
    echo "Guardando reporte de errores en compilation-errors.log..."
    echo "$ERRORS" > compilation-errors.log
    print_status "Reporte guardado en compilation-errors.log"
fi

# 5. Mostrar resumen de errores
echo ""
echo "📊 Resumen de Errores por Tipo:"
echo "================================"
echo "$ERRORS" | grep -oE "(undefined|too many arguments|not enough arguments|redeclared|declared and not used|imported and not used)" | sort | uniq -c | sort -rn

# 6. Mostrar paquetes con errores
echo ""
echo "📦 Paquetes con Errores:"
echo "======================="
echo "$ERRORS" | grep "^#" | sort -u

# 7. Instrucciones finales
echo ""
echo "📝 Próximos Pasos:"
echo "=================="
echo "1. Revisa el archivo: COMPILATION_ERRORS_REPORT.md"
echo "2. Revisa los errores detallados en: compilation-errors.log"
echo "3. Ejecuta: go build ./... para ver errores en tiempo real"
echo ""
echo "💡 Tips:"
echo "  - Para ver errores de un paquete específico: go build ./internal/customer/..."
echo "  - Para fix rápido de imports: goimports -w <archivo>.go"
echo ""

exit $ERROR_COUNT
