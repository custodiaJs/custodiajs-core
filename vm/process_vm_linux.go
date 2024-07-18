package vm

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func hiddenF() {
	// Hardcoded Command zum Starten eines Bash-Prozesses
	command := "/bin/bash"

	// Konfiguriert die Namespaces und Isolierungsoptionen
	cmd := exec.Command(command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Setzt die benötigten Namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER,
	}

	// Setzt das eigene Root-Verzeichnis für das Dateisystem im Container
	cmd.SysProcAttr.Chroot = "/var/lib/container"
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: 0, Gid: 0}

	// Startet den isolierten Prozess
	if err := cmd.Start(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// Wartet auf das Ende des Prozesses
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command: ", err)
		os.Exit(1)
	}
}
