#!/bin/bash

# Script para corregir llamadas a logger con demasiados argumentos

echo "Corrigiendo llamadas a logger..."

# Patrón: h.log.Debug("msg", "key", value) -> h.log.Debugf("msg key=%v", value)
# Patrón: h.log.Error("msg", "key", err) -> h.log.WithError(err).Error("msg")
# Patrón: h.logger.Error("msg", "key", err) -> h.logger.WithError(err).Error("msg")

# Archivos a procesar
FILES=$(grep -l "too many arguments in call to.*log" compilation-errors.log | sed 's/:.*//g' | sort -u)

for file in $FILES; do
    echo "Procesando: $file"

    # Backup
    cp "$file" "$file.bak"

    # No podemos hacer esto automáticamente porque cada caso es diferente
    # Lo marcaremos para revisión manual
done

echo "Archivos que necesitan corrección manual:"
echo "$FILES"
