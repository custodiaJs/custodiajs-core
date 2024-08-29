#!/bin/bash

# Setze die Variablen
BUILD_DIR="./build"
CORE_MAIN="cmd/core-service/main.go"
VM_MAIN="cmd/vm/"
OUTPUT1="coresrvce"
OUTPUT2="vm"

# Farbe für Ausgaben
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Funktion zur Fehlerbehandlung
error_exit() {
    echo -e "${RED}Fehler: $1${NC}" 1>&2
    exit 1
}

# Bereinigen
clean() {
    echo -e "${GREEN}Bereinigen...${NC}"
    rm -rf $BUILD_DIR
    mkdir -p $BUILD_DIR || error_exit "Konnte Build-Verzeichnis nicht erstellen."
    echo "Bereinigung abgeschlossen."
}

# Kompilieren der Go-Programme
build() {
    echo -e "${GREEN}Kompiliere $CORE_MAIN für macOS...${NC}"
    if ! go build -tags NoLocalhostSSLCheck -o $BUILD_DIR/$OUTPUT1 $CORE_MAIN; then
        error_exit "Fehler bei der Kompilierung von $CORE_MAIN für macOS: $(go build -tags NoLocalhostSSLCheck -o $BUILD_DIR/$OUTPUT1 $CORE_MAIN 2>&1)"
    fi
    chmod +x $BUILD_DIR/$OUTPUT1

    cd $VM_MAIN
    echo -e "${GREEN}Kompiliere das Verzeichnis $VM_MAIN für macOS...${NC}"
    if ! go build -tags NoLocalhostSSLCheck -o ../../$BUILD_DIR/$OUTPUT2; then
        error_exit "Fehler bei der Kompilierung des Verzeichnisses $VM_MAIN für macOS: $(go build -tags NoLocalhostSSLCheck -o $BUILD_DIR/$OUTPUT2 $VM_MAIN 2>&1)"
    fi
    chmod +x ../../$BUILD_DIR/$OUTPUT2
    cd ..
    
    echo "Kompilierung für macOS abgeschlossen."
}

error_exit() {
    echo "$1" 1>&2
    exit 1
}

# Hilfe anzeigen
usage() {
    echo "Usage: $0 {clean|build|all}"
    exit 1
}

# Hauptlogik des Skripts
case "$1" in
    clean)
        clean
        ;;
    build)
        clean
        build
        ;;
    all)
        clean
        build
        ;;
    *)
        usage
        ;;
esac

exit 0