package vmdb

import (
	"fmt"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/utils"
)

func (o *VMEntryBaseData) GetDatabaseFingerprint() DatabaseFingerprint {
	fprintHash, err := utils.BuildStringHashChain(o.Type, o.Host, fmt.Sprintf("%d", o.Port), o.Username, o.Password, o.Database, o.Alias)
	if err != nil {
		panic("VMEntryBaseData->GetDatabaseFingerprint: " + err.Error())
	}
	return DatabaseFingerprint(strings.ToLower(fprintHash))
}
