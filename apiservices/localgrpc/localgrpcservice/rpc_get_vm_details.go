package localgrpcservice

import (
	"context"
	"fmt"
	"strings"
	"vnh1/localgrpcproto"
	"vnh1/static"
	"vnh1/types"
	"vnh1/utils"
)

func (s *CliGrpcServer) GetVMDetails(ctx context.Context, vmDetailParms *localgrpcproto.VmDetailsParms) (*localgrpcproto.VmDetailsResponse, error) {
	// Die VM wird ermittelt
	var foundedVM types.CoreVMInterface
	var foundVM bool
	var err error
	switch vmDetailParms.Value.(type) {
	case *localgrpcproto.VmDetailsParms_Id:
		// Es wird geprüft ob es sich um eine gültige VM Id handelt
		if !utils.ValidateVMIdString(vmDetailParms.GetId()) {
			return nil, fmt.Errorf("invalid vm id")
		}

		// Es wird versucht die VM abzurufen
		foundedVM, foundVM, err = s.core.GetScriptContainerVMByID(vmDetailParms.GetId())
	case *localgrpcproto.VmDetailsParms_Name:
		// Es wird geprüft ob es sich um einen zulässigen VM Namen handelt
		if !utils.ValidateVMName(vmDetailParms.GetName()) {
			return nil, fmt.Errorf("invalid vm name")
		}

		// Es wird versucht die VM mittels ihres Namens zu Extrahieren
		foundedVM, err = s.core.GetScriptContainerByVMName(vmDetailParms.GetName())
		if foundedVM != nil {
			foundVM = true
		}
	default:
		return nil, fmt.Errorf("invalid 'get vm details' parameter vm id/name")
	}

	// Es wird geprüft ob die VM gefunden wurde
	if err != nil {
		return nil, err
	}

	// Es wird geprüft ob die VM gefunden wurde
	if !foundVM {
		return nil, fmt.Errorf("unkown vm")
	}

	// Die Whitelist wird extrahiert
	extractedWhitelist := make([]*localgrpcproto.VmDetailWhitelistEntry, 0)
	for _, item := range foundedVM.GetWhitelist() {
		extractedWhitelist = append(extractedWhitelist, &localgrpcproto.VmDetailWhitelistEntry{
			WildCardDomains: item.WildCardDomains,
			ExactDomains:    item.ExactDomains,
			Methods:         item.Methods,
			IPv4List:        item.IPv4List,
			Ipv6List:        item.Ipv6List,
		})
	}

	// Die HostCA Members werden abgerufen
	extractedHostCAList := make([]*localgrpcproto.VmDetailHostCAMemberEntry, 0)
	for _, item := range foundedVM.GetRootMemberIDS() {
		extractedHostCAList = append(extractedHostCAList, &localgrpcproto.VmDetailHostCAMemberEntry{
			Type:        1,
			Fingerprint: strings.ToUpper(item.Fingerprint),
		})
	}

	// Die Datenbanken werden extrahiert
	extractedDBList := make([]*localgrpcproto.VmDetailDatabaseEntry, 0)
	for _, item := range foundedVM.GetDatabaseServices() {
		extractedDBList = append(extractedDBList, &localgrpcproto.VmDetailDatabaseEntry{
			Type:     item.Type,
			Host:     item.Host,
			Port:     uint32(item.Port),
			Username: item.Username,
			Database: item.Database,
			Alias:    item.Alias,
		})
	}

	// Die geteilten Funktionen werden abgerufen
	sharedFunctions := make([]*localgrpcproto.VmDetailSharedFunctionEntry, 0)
	for _, item := range foundedVM.GetAllSharedFunctions() {
		// Die Parameter Typen werden extrahiert
		extractedParmsList := make([]uint32, 0)
		for _, item := range item.GetParmTypes() {
			switch item {
			case "boolean":
				extractedParmsList = append(extractedParmsList, 0)
			case "number":
				extractedParmsList = append(extractedParmsList, 1)
			case "string":
				extractedParmsList = append(extractedParmsList, 2)
			case "array":
				extractedParmsList = append(extractedParmsList, 3)
			case "object":
				extractedParmsList = append(extractedParmsList, 4)
			case "bytes":
				extractedParmsList = append(extractedParmsList, 5)
			default:
				return nil, fmt.Errorf("GetVMDetails: unsuported parameter datatype")
			}
		}

		// Der Eintrag wird erzeugt und zwischengspeichert
		sharedFunctions = append(sharedFunctions, &localgrpcproto.VmDetailSharedFunctionEntry{
			Name:      item.GetName(),
			Mode:      "unkown",
			ParmTypes: extractedParmsList,
		})
	}

	// Der Stauts wird als String ermittelt
	var stateStr string
	switch foundedVM.GetState() {
	case static.Closed:
		stateStr = "closed"
	case static.Running:
		stateStr = "running"
	case static.Starting:
		stateStr = "starting"
	case static.StillWait:
		stateStr = "still wait"
	default:
		stateStr = "unkown"
	}

	// Das Rückgabe Objekt wird erstellt
	responseObj := &localgrpcproto.VmDetailsResponse{
		Name:            foundedVM.GetVMName(),
		Version:         10000000000000000000,
		Owner:           foundedVM.GetOwner(),
		Repourl:         foundedVM.GetRepoURL(),
		Mode:            foundedVM.GetMode(),
		State:           stateStr,
		StartTimestamp:  foundedVM.GetStartingTimestamp(),
		Whitelist:       extractedWhitelist,
		Hostcamember:    extractedHostCAList,
		Databases:       extractedDBList,
		SharedFunctions: sharedFunctions,
	}

	// Die Daten werden zurückgesendet
	return responseObj, nil
}
