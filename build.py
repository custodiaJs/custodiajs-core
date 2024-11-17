import os
import sys
import subprocess
import platform
import shutil
import argparse


# Farben für Ausgaben
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


def detect_host_platform():
    """Ermittelt die Host-Plattform und -Architektur."""
    system = platform.system().lower()
    architecture = platform.machine().lower()

    if system == "windows":
        goos = "windows"
    elif system == "linux":
        goos = "linux"
    elif system == "darwin":
        goos = "darwin"
    else:
        error_exit(f"Unbekannte Host-Plattform: {system}")

    # Mapping der Architektur
    if architecture in ["x86_64", "amd64"]:
        goarch = "amd64"
    elif architecture in ["arm64", "aarch64"]:
        goarch = "arm64"
    else:
        error_exit(f"Unbekannte Host-Architektur: {architecture}")

    return goos, goarch


def clean_or_create_build_dir(build_dir):
    """Bereinigt das Build-Verzeichnis oder erstellt es, falls es nicht existiert."""
    build_dir = os.path.abspath(build_dir)  # Absoluter Pfad sicherstellen
    if os.path.exists(build_dir):
        print(f"{GREEN}Bereinige vorhandenes Build-Verzeichnis: {build_dir}{NC}")
        shutil.rmtree(build_dir)
    os.makedirs(build_dir)
    print(f"{GREEN}Erstelle Build-Verzeichnis: {build_dir}{NC}")


def build(core_service_path, build_dir, output_name, platforms):
    """Kompiliert plattformspezifische Dateien aus dem Core-Service-Verzeichnis."""
    clean_or_create_build_dir(build_dir)

    core_service_path = os.path.abspath(core_service_path)  # Absoluter Pfad zum Core-Service-Ordner

    for goos, goarch in platforms:
        print(f"{GREEN}Kompiliere Core-Service für {goos}/{goarch}...{NC}")
        env = os.environ.copy()
        env["GOOS"] = goos
        env["GOARCH"] = goarch

        # Ausgabe-Dateiname anpassen
        output_file = f"{output_name}_{goos}_{goarch}" + (".exe" if goos == "windows" else "")
        output_path = os.path.join(build_dir, output_file)

        try:
            # Setze das Working Directory auf den Core-Service-Ordner und kompiliere
            subprocess.run(
                ["go", "build", "-tags", "NoLocalhostSSLCheck", "-o", output_path],
                check=True,
                env=env,
                cwd=core_service_path  # Setze das Working Directory
            )
            os.chmod(output_path, 0o755)
            print(f"{GREEN}Erfolgreich kompiliert: {output_path}{NC}")
        except subprocess.CalledProcessError as e:
            error_exit(f"Fehler bei der Kompilierung von Core-Service für {goos}/{goarch}: {e.stderr.decode() if e.stderr else 'Keine Fehlermeldung verfügbar'}")

    print(f"{GREEN}Kompilierung abgeschlossen.{NC}")


def parse_arguments():
    """Parst die Kommandozeilenargumente."""
    parser = argparse.ArgumentParser(description="Go-Build-Skript für Core-Service mit plattformspezifischen Dateien.")
    parser.add_argument(
        "--platform", "-p", nargs="*", 
        help="Plattform und Architektur, z.B. linux/amd64 windows/arm64. Standard: Host-Plattform."
    )
    parser.add_argument(
        "--output", "-o", default="core_service", 
        help="Basisname der Ausgabedateien. Standard: core_service."
    )
    parser.add_argument(
        "--build-dir", "-b", default="build", 
        help="Verzeichnis für die kompilierten Dateien. Standard: build."
    )
    return parser.parse_args()


def main():
    """Hauptlogik des Skripts."""
    args = parse_arguments()
    
    print(f"{GREEN}Prüfe System und Abhängigkeiten...{NC}")
    check_go()

    # Stammverzeichnis des Skripts
    script_dir = os.path.dirname(os.path.abspath(__file__))

    # Neuer Pfad zum Core-Service und Build-Verzeichnis
    core_service_path = os.path.join(script_dir, "src/cmd/core-service")
    build_dir = os.path.join(script_dir, args.build_dir)

    # Plattform- und Architekturverarbeitung
    if args.platform:
        platforms = []
        for plat in args.platform:
            try:
                goos, goarch = plat.split("/")
                platforms.append((goos, goarch))
            except ValueError:
                error_exit(f"Ungültiges Plattformformat: {plat}. Erwartet: OS/ARCH (z.B. linux/amd64)")
    else:
        # Host-Plattform verwenden, wenn keine Plattform angegeben
        platforms = [detect_host_platform()]

    build(core_service_path, build_dir, args.output, platforms)
    print(f"{GREEN}Alle Schritte erfolgreich abgeschlossen!{NC}")


if __name__ == "__main__":
    main()
