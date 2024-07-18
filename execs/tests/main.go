package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/CustodiaJS/custodiajs-core/apiservices/httpjson"
	"github.com/CustodiaJS/custodiajs-core/apiservices/localgrpc"
	"github.com/CustodiaJS/custodiajs-core/core"
	"github.com/CustodiaJS/custodiajs-core/databaseservices"
	"github.com/CustodiaJS/custodiajs-core/filesystem"
	"github.com/CustodiaJS/custodiajs-core/identkeydatabase"
	"github.com/CustodiaJS/custodiajs-core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/kernel/external_modules"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/CustodiaJS/custodiajs-core/vmdb"
)

const spaces = "   "

var logDIR types.LOG_DIR = ""

func loadHostTlsCert() (*tls.Certificate, error) {
	// Das Host Cert wird geladen
	cert, err := os.ReadFile("/home/fluffelbuff/Schreibtisch/localhost.crt")
	if err != nil {
		panic(err)
	}

	// Der Private Schlüssel wird geladen
	key, err := os.ReadFile("/home/fluffelbuff/Schreibtisch/localhost.pem")
	if err != nil {
		panic(err)
	}

	// Erstelle ein TLS-Zertifikat aus den geladenen Dateien
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	// Das Cert wird zurückgegebn
	return &tlsCert, nil
}

func printLocalHostTlsMetaData(cert *tls.Certificate) {
	if len(cert.Certificate) == 0 {
		return
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return
	}

	fingerPrint := utils.ComputeTlsCertFingerprint(cert)
	fingerprintHex := strings.ToUpper(hex.EncodeToString(fingerPrint))

	// Extrahiere den Signaturalgorithmus als String
	sigAlgo := x509Cert.SignatureAlgorithm.String()

	// Ausgabe
	fmt.Printf("%sFingerprint (SHA3-256): %s\n   Algorithm: %s\n", spaces, fingerprintHex, sigAlgo)
}

func initLogDir() error {
	if filesystem.FolderExists(string(static.UNIX_LINUX_LOGGING_DIR)) {
		return nil
	}

	err := os.MkdirAll(string(static.UNIX_LINUX_LOGGING_DIR), os.ModePerm)
	if err != nil {
		err := os.MkdirAll(string(static.UNIX_LINUX_LOGGING_DIR_NONE_ROOT), os.ModePerm)
		if err != nil {
			return fmt.Errorf("log dir making error: " + err.Error())
		}

		logDIR = static.UNIX_LINUX_LOGGING_DIR_NONE_ROOT
		return nil
	}

	logDIR = static.UNIX_LINUX_LOGGING_DIR

	// Rückgabe
	return nil
}

func main() {
	// Compiler warnings
	if !static.CHECK_SSL_LOCALHOST_ENABLE {
		fmt.Printf("Warning: SSL verification for localhost has been completely disabled during compilation.\nThis may lead to unexpected issues, as programs or websites might not be able to communicate with the VNH1 service anymore.\nIf you have downloaded and installed VNH1 and are seeing this message, please be aware that you are not using an official build.\n\n")
	}

	// Gibt an ob das Programm in einem Linux Container ausgeführt wird
	isRunningInLinuxContainer := false

	// Die Hostinformationen werden ausgelesen
	if runtime.GOOS == "linux" {
		// Die Linux Informationen werden ausgelesen
		hostInfo, err := utils.DetectLinuxDist()
		if err != nil {
			panic(err)
		}

		// Die Host Informationen werden angezigt
		fmt.Println("Host OS:", hostInfo)

		// Es wird ermittelt ob das Programm in einem Container ausgeführt wird
		isRunningInLinuxContainer = utils.IsRunningInContainer()

		// Die Info wird angezeigt
		if isRunningInLinuxContainer {
			fmt.Println("Running in container: yes")
		} else {
			fmt.Println("Running in container: no")
		}
	}

	// Maximale Anzahl von CPU-Kernen für die Go-Runtime festlegen
	runtime.GOMAXPROCS(1)

	// Das HostCert und der Privatekey werden geladen
	fmt.Print("Loading host certificate: ")
	hostCert, err := loadHostTlsCert()
	if err != nil {
		fmt.Println("error@")
		panic(err)
	}
	fmt.Println("done")

	// Die Log Verzeichnisse werden erstellen
	fmt.Println("Prepare LOG directory...s")
	if err := initLogDir(); err != nil {
		panic(err)
	}

	// Speichert alle Abgerufen Libs ab
	extModuleLibs := make([]*external_modules.ExternalModule, 0)

	/*
		lib1, err := external_modules.LoadModuleLib("/home/fluffelbuff/Schreibtisch/lib1.so")
		if err != nil {
			panic(err)
		}
		lib2, err := external_modules.LoadModuleLib("/home/fluffelbuff/Schreibtisch/lib2.so")
		if err != nil {
			panic(err)
		}

		extModuleLibs = append(extModuleLibs, lib1)
		extModuleLibs = append(extModuleLibs, lib2)
	*/

	// Die Metadaten des Host Zertifikates werden angezeigt
	printLocalHostTlsMetaData(hostCert)

	// Die Host Ident Key database wird geladen
	fmt.Println("Loading host ident key database...")
	ikdb, err := identkeydatabase.LoadIdentKeyDatabase()
	if err != nil {
		fmt.Print("error@ ")
		panic(err)
	}

	// Die VM Datenbank wird geladen
	fmt.Println("Loading vm database...")
	vmdatabase, err := vmdb.OpenFilebasedVmDatabase()
	if err != nil {
		fmt.Print("error@ ")
		panic(err)
	}

	// Der Datenbank Hostservice wird erstellt
	dbservice := databaseservices.NewDbService()
	ipnetcon := ipnetwork.NewHostNetworkManagmentUnit()

	// Der Core wird erzeugt
	core, err := core.NewCore(hostCert, ikdb, dbservice, logDIR, ipnetcon)
	if err != nil {
		panic(err)
	}

	// Die Externen Module libs werden hinzugefügt
	for _, item := range extModuleLibs {
		// Es wird versucht das Externe Modul Lib zu laden
		if err := core.AddExternalModuleLibrary(item); err != nil {
			panic(err)
		}

		// LOG
		fmt.Printf("External module lib '%s' version %d loaded\n", item.GetName(), item.GetVersion())
	}

	// Die CLI Terminals werden erzeugt
	noneRootCLI, err := localgrpc.NewTestTCP("/home/fluffelbuff/Schreibtisch/localhost.crt", "/home/fluffelbuff/Schreibtisch/localhost.pem", static.NONE_ROOT_ADMIN)
	if err != nil {
		panic(err)
	}

	// Die CLI wird hinzuefügt
	if err := core.AddAPISocket(noneRootCLI); err != nil {
		panic(err)
	}

	// Der Localhost httpjson wird erzeugt
	localhostWebserviceV6, err := httpjson.NewLocalService("ipv6", 8080, hostCert)
	if err != nil {
		panic(err)
	}
	localhostWebserviceV4, err := httpjson.NewLocalService("ipv4", 8080, hostCert)
	if err != nil {
		panic(err)
	}

	// Der Localhost httpjson wird hinzugefügt
	if err := core.AddAPISocket(localhostWebserviceV6); err != nil {
		panic(err)
	}
	if err := core.AddAPISocket(localhostWebserviceV4); err != nil {
		panic(err)
	}

	// Es wird geprüft ob es sich um Unterstützes OS handelt
	switch runtime.GOOS {
	case "linux":
		if err := utils.VerifyLinuxSystem(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "windows":
		if err := utils.VerifyWindowsSystem(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "darwin":
		if err := utils.VerifyAppleMacOSSystem(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "freebsd", "openbsd", "netbsd":
		if err := utils.VerifyBSDSystem(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Println("It is an unsupported operating system.")
		os.Exit(1)
	}

	// Die Einzelnene Datenbank Dienste werden hinzugefügt
	for _, item := range vmdatabase.GetAllDatabaseVMBaseData() {
		if err := dbservice.AddDatabaseService(item); err != nil {
			panic(err)
		}
	}

	// Die Einzelnen VM's werden geladen und hinzugefügt
	fmt.Println("Loading JavaScript virtual machines...")
	for _, item := range vmdatabase.GetAllVirtualMachines() {
		// Die VM wird erzeugt
		newVM, err := core.AddNewVMInstance(item)
		if err != nil {
			panic(err)
		}

		// Log
		fmt.Printf("%s-> VM '%s' <-> %s loaded %d bytes [%s]\n", spaces, newVM.GetVMName(), strings.ToUpper(string(newVM.GetFingerprint())), item.GetBaseSize(), newVM.GetKId())
	}

	// Der Core wird gestartet
	fmt.Println("Starting done...")
	var waitGroupForServing sync.WaitGroup
	waitGroupForServing.Add(1)
	go func() {
		core.Serve()
		waitGroupForServing.Done()
	}()

	// Ein Channel, um Signale zu empfangen.
	sigChan := make(chan os.Signal, 1)

	// Notify sigChan, wenn ein SIGINT empfangen wird.
	signal.Notify(sigChan, syscall.SIGINT)

	// Es wird auf das Signal zum beenden gewartet
	<-sigChan
	fmt.Println("")

	// Dem Core wird Signalisert dass er beendet wird
	core.SignalShutdown()

	// Es wird gewartet bis der Core beendet wurde
	waitGroupForServing.Wait()

	// Die Externen Module Libs werden entladen
	for _, item := range extModuleLibs {
		item.Unload()
	}
}
