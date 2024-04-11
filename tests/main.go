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
	"vnh1/apiservices/httpjson"
	"vnh1/apiservices/localgrpc"
	"vnh1/core"
	"vnh1/core/databaseservices"
	"vnh1/core/identkeydatabase"
	"vnh1/core/vmdb"
	"vnh1/extmodules"
	"vnh1/types"
	"vnh1/utils"
)

const spaces = "   "

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

func main() {
	// Compiler warnings
	if !utils.CHECK_SSL_LOCALHOST_ENABLE {
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

	// Anzahl der verfügbaren CPU-Kerne ermitteln
	numCPU := runtime.NumCPU()
	fmt.Println("Total CPU-Cores avail:", numCPU-2)

	// Maximale Anzahl von CPU-Kernen für die Go-Runtime festlegen
	runtime.GOMAXPROCS(numCPU - 2)

	// Das HostCert und der Privatekey werden geladen
	fmt.Print("Loading host certificate: ")
	hostCert, err := loadHostTlsCert()
	if err != nil {
		fmt.Println("error@")
		panic(err)
	}
	fmt.Println("done")

	// Speichert alle Abgerufen Libs ab
	extModuleLibs := make([]*extmodules.ExternalModule, 0)

	lib1, err := extmodules.LoadModuleLib("/home/fluffelbuff/Schreibtisch/lib1.so")
	if err != nil {
		panic(err)
	}
	extModuleLibs = append(extModuleLibs, lib1)

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

	// Der Core wird erzeugt
	core, err := core.NewCore(hostCert, ikdb, dbservice)
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
	noneRootCLI, err := localgrpc.NewTestTCP("/home/fluffelbuff/Schreibtisch/localhost.crt", "/home/fluffelbuff/Schreibtisch/localhost.pem", types.NONE_ROOT_ADMIN)
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

	// Die Einzelnene Datenbank Dienste werden hinzugefügt
	for _, item := range vmdatabase.GetAllDatabaseConfigurations() {
		if err := dbservice.AddDatabaseService(item); err != nil {
			panic(err)
		}
	}

	// Die Einzelnen VM's werden geladen und hinzugefügt
	fmt.Println("Loading JavaScript virtual machines...")
	for _, item := range vmdatabase.GetAllVirtualMachines() {
		// Die VM wird erzeugt
		newVM, err := core.AddScriptContainer(item)
		if err != nil {
			panic(err)
		}

		// Log
		fmt.Printf("%s-> VM '%s' <-> %s loaded %d bytes\n", spaces, newVM.GetVMName(), strings.ToUpper(string(newVM.GetFingerprint())), item.GetBaseSize())
	}

	// Der Core wird gestartet
	fmt.Println()
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

	// Dem Core wird Signalisert dass er beendet wird
	core.SignalShutdown()

	// Es wird gewartet bis der Core beendet wurde
	waitGroupForServing.Wait()

	// Die Externen Module Libs werden entladen
	for _, item := range extModuleLibs {
		item.Unload()
	}
}
