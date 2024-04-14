package kernel

import (
	"vnh1/types"
)

func (o *Kernel) GetLocalSharedFunctions() []types.SharedLocalFunctionInterface {
	extracted := make([]types.SharedLocalFunctionInterface, 0)
	for _, item := range o.sharedLocalFunctions {
		extracted = append(extracted, item)
	}
	return extracted
}

func (o *Kernel) GetPublicSharedFunctions() []types.SharedPublicFunctionInterface {
	extracted := make([]types.SharedPublicFunctionInterface, 0)
	for _, item := range o.sharedPublicFunctions {
		extracted = append(extracted, item)
	}
	return extracted
}

func (o *Kernel) GetAllSharedFunctions() []types.SharedFunctionInterface {
	vat := make([]types.SharedFunctionInterface, 0)
	for _, item := range o.GetLocalSharedFunctions() {
		vat = append(vat, item)
	}
	for _, item := range o.GetPublicSharedFunctions() {
		vat = append(vat, item)
	}
	return vat
}
