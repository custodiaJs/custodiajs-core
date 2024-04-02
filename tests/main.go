package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"vnh1/apiservices/cligrpc"
	"vnh1/apiservices/httpapi"
	"vnh1/apiservices/webgrpc"
	"vnh1/core"
	"vnh1/core/identkeydatabase"
	"vnh1/core/vmdb"
	"vnh1/types"
	"vnh1/utils"

	"golang.org/x/crypto/sha3"
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

func loadHostIdentKeyDatabase() (*identkeydatabase.IdenKeyDatabase, error) {
	return &identkeydatabase.IdenKeyDatabase{}, nil
}

func loadVMDatabase() (*vmdb.VmDatabase, error) {
	return vmdb.OpenFilebasedVmDatabase()
}

func printLocalHostTlsMetaData(cert *tls.Certificate) {
	if len(cert.Certificate) == 0 {
		return
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return
	}

	// Berechne den Fingerprint des Zertifikats (hier weiterhin SHA-256)
	hash := sha3.New256()
	_, err = hash.Write(x509Cert.Raw)
	if err != nil {
		return
	}
	fingerprintBytes := hash.Sum(nil)
	fingerprint := hex.EncodeToString(fingerprintBytes)

	// Extrahiere den Signaturalgorithmus als String
	sigAlgo := x509Cert.SignatureAlgorithm.String()

	// Ausgabe
	fmt.Printf("%sFingerprint (SHA3-256): %s\n   Algorithm: %s\n", spaces, fingerprint, sigAlgo)
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

	// Das HostCert und der Privatekey werden geladen
	fmt.Print("Loading host certificate: ")
	hostCert, err := loadHostTlsCert()
	if err != nil {
		fmt.Println("error@")
		panic(err)
	}
	fmt.Println("done")

	// Die Metadaten des Host Zertifikates werden angezeigt
	printLocalHostTlsMetaData(hostCert)

	// Die Host Ident Key database wird geladen
	fmt.Print("Loading host ident key database: ")
	ikdb, err := loadHostIdentKeyDatabase()
	if err != nil {
		fmt.Print("error@ ")
		panic(err)
	}
	fmt.Println("done")

	// Die VM Datenbank wird geladen
	fmt.Print("Loading vm database: ")
	vmdatabase, err := loadVMDatabase()
	if err != nil {
		fmt.Print("error@ ")
		panic(err)
	}
	fmt.Println("done")

	// Der Core wird erzeugt
	core, err := core.NewCore(hostCert, ikdb)
	if err != nil {
		panic(err)
	}

	// Die CLI Terminals werden erzeugt
	fmt.Println("cligrpc: enabled")
	noneRootCLI, err := cligrpc.New("/tmp/vnh1_none_root", types.NONE_ROOT_ADMIN)
	if err != nil {
		panic(err)
	}

	// Die CLI wird hinzuefügt
	if err := core.AddAPISocket(noneRootCLI); err != nil {
		panic(err)
	}

	// Der Localhost httpapi wird erzeugt
	fmt.Println("httpapi (localhost): enabled")
	localhostWebserviceV6, err := httpapi.NewLocalService("ipv6", 8080, hostCert)
	if err != nil {
		panic(err)
	}
	localhostWebserviceV4, err := httpapi.NewLocalService("ipv4", 8080, hostCert)
	if err != nil {
		panic(err)
	}

	// Der Localhost httpapi wird hinzugefügt
	if err := core.AddAPISocket(localhostWebserviceV6); err != nil {
		panic(err)
	}
	if err := core.AddAPISocket(localhostWebserviceV4); err != nil {
		panic(err)
	}

	// Der grpcservice wird erzeugt
	fmt.Println("grpcapi (localhost): enabled")
	localhostGrpcServiceV6, err := webgrpc.NewLocalService("ipv6", 8081, hostCert)
	if err != nil {
		panic(err)
	}
	localhostGrpcServiceV4, err := webgrpc.NewLocalService("ipv4", 8081, hostCert)
	if err != nil {
		panic(err)
	}

	// Der Localhost grpcservice wird hinzugefügt
	if err := core.AddAPISocket(localhostGrpcServiceV6); err != nil {
		panic(err)
	}
	if err := core.AddAPISocket(localhostGrpcServiceV4); err != nil {
		panic(err)
	}

	// Die Einzelnen VM's werden geladen
	fmt.Println("Loading JavaScript virtual machines...")
	vms, err := vmdatabase.LoadAllVirtualMachines()
	if err != nil {
		panic(err)
	}

	// Die Einzelnen VM's werden gestartet
	for _, item := range vms {
		// Die VM wird erzeugt
		newVM, err := core.AddScriptContainer(item)
		if err != nil {
			fmt.Print("error@ ")
			panic(err)
		}

		// Log
		fmt.Printf("%s-> VM '%s' <-> %s loaded %d bytes\n%s%s--> Total NodeJS submodules: %d\n", spaces, newVM.GetVMName(), newVM.GetFingerprint(), item.GetBaseSize(), spaces, spaces, item.GetTotalNodeJsModules())
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
}
