package vmdb

import (
	"fmt"
	"strings"
	"vnh1/utils"
)

func (o *VMEntryBaseData) GetDatabaseFingerprint() DatabaseFingerprint {
	fprintHash, err := utils.BuildStringHashChain(o.Type, o.Host, fmt.Sprintf("%d", o.Port), o.Username, o.Password, o.Database, o.Alias)
	if err != nil {
		panic("VMEntryBaseData->GetDatabaseFingerprint: " + err.Error())
	}
	return DatabaseFingerprint(strings.ToLower(fprintHash))
}
