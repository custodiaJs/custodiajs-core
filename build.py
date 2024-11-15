import os
import sys
import platform
import subprocess


# Farbe für Ausgaben
GREEN = "\033[0;32m"
RED = "\033[0;31m"
NC = "\033[0m"  # No Color


def error_exit(message):
    """Beendet das Skript bei einem Fehler."""
    print(f"{RED}Fehler: {message}{NC}", file=sys.stderr)
    sys.exit(1)


def check_go():
    """Prüft, ob Go installiert ist."""
    print(f"{GREEN}Prüfe, ob Go installiert ist...{NC}")
    try:
        result = subprocess.run(["go", "version"], check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        print(f"{GREEN}Go gefunden: {result.stdout.decode().strip()}{NC}")
    except FileNotFoundError:
        error_exit("Go ist nicht installiert. Bitte installiere Go von https://go.dev/dl/")
    except subprocess.CalledProcessError as e:
        error_exit(f"Fehler beim Überprüfen von Go: {e.stderr.decode().strip()}")


def check_v8go_dependencies():
    """Prüft, ob die notwendigen Abhängigkeiten für v8go installiert sind."""
    print(f"{GREEN}Prüfe Abhängigkeiten für v8go...{NC}")
    system = platform.system().lower()
    dependencies = []

    if system == "linux":
        dependencies = ["build-essential", "libssl-dev", "pkg-config", "clang", "python3"]
    elif system == "darwin":  # macOS
        dependencies = ["llvm", "pkg-config", "python3"]
    elif system == "freebsd":
        dependencies = ["llvm", "pkgconf", "python3"]
    elif system == "windows":
        print(f"{RED}Prüfung für v8go auf Windows wird nicht unterstützt. Bitte installieren Sie die notwendigen Abhängigkeiten manuell.{NC}")
        return
    else:
        error_exit(f"Nicht unterstütztes Betriebssystem: {system}")

    for dep in dependencies:
        print(f"Prüfe {dep}...")
        try:
            result = subprocess.run(
                ["which", dep] if system != "windows" else ["where", dep],
                check=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            print(f"{GREEN}{dep} gefunden: {result.stdout.decode().strip()}{NC}")
        except FileNotFoundError:
            error_exit(f"{dep} ist nicht installiert. Bitte installieren.")
        except subprocess.CalledProcessError as e:
            error_exit(f"Fehler beim Überprüfen von {dep}: {e.stderr.decode().strip()}")

    print(f"{GREEN}Alle Abhängigkeiten für v8go sind vorhanden.{NC}")


def main():
    """Hauptlogik des Skripts."""
    print(f"{GREEN}Prüfe System und Abhängigkeiten...{NC}")
    check_go()
    check_v8go_dependencies()
    print(f"{GREEN}Alle erforderlichen Tools und Abhängigkeiten sind vorhanden!{NC}")


if __name__ == "__main__":
    main()
