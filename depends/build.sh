#!/bin/bash

# Setze die Variablen
BUILD_DIR="./build"
MAIN1="cmd/core-service/main.go"
MAIN2="cmd/core-vm/main.go"
OUTPUT1="core-service"
OUTPUT2="core-vm"

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
    echo -e "${GREEN}Kompiliere $MAIN1...${NC}"
    go build -tags NoLocalhostSSLCheck -o $BUILD_DIR/$OUTPUT1 $MAIN1 || error_exit "Kompilierung von $MAIN1 fehlgeschlagen."
    chmod +x $BUILD_DIR/$OUTPUT1

    echo -e "${GREEN}Kompiliere $MAIN2...${NC}"
    go build -tags NoLocalhostSSLCheck -o $BUILD_DIR/$OUTPUT2 $MAIN2 || error_exit "Kompilierung von $MAIN2 fehlgeschlagen."
    chmod +x $BUILD_DIR/$OUTPUT2
    
    echo "Kompilierung abgeschlossen."
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