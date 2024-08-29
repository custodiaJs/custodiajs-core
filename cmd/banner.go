package cmd

import (
	"fmt"
	"runtime"

	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/utils"
)

// Welcome Banner
func ShowBanner() {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	isAdmin := utils.CheckAdmin()

	banner := fmt.Sprintf(`
 ██████╗██╗   ██╗███████╗████████╗ ██████╗ ██████╗ ██╗ █████╗      ██╗███████╗
██╔════╝██║   ██║██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗██║██╔══██╗     ██║██╔════╝
██║     ██║   ██║███████╗   ██║   ██║   ██║██║  ██║██║███████║     ██║███████╗
██║     ██║   ██║╚════██║   ██║   ██║   ██║██║  ██║██║██╔══██║██   ██║╚════██║
╚██████╗╚██████╔╝███████║   ██║   ╚██████╔╝██████╔╝██║██║  ██║╚█████╔╝███████║
 ╚═════╝ ╚═════╝ ╚══════╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝ ╚════╝ ╚══════╝
                                                                             
	JavaScript VM Interpreter                  
	Version: %s
	Author: %s
	OS: %s
	Architecture: %s
	User is Admin: %t
----------------------------------------------------------------------------------`, utils.FormatNumberWithDots(int(static.C_VESION)), "fluffelpuff", osName, arch, isAdmin)
	fmt.Println(banner)

	if !static.CHECK_SSL_LOCALHOST_ENABLE {
		fmt.Printf("Warning: SSL verification for localhost has been completely disabled during compilation.\nThis may lead to unexpected issues, as programs or websites might not be able to communicate with the VNH1 service anymore.\nIf you have downloaded and installed VNH1 and are seeing this message, please be aware that you are not using an official build.\n\n")
	}
}
