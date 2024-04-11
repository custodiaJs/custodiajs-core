package core

import extmodules "vnh1/extmodules"

// Ruft Eine Liste Externer Module anahnd ihrer Namen ab
func (o *Core) _core_util_get_list_of_extmods_by_name(names ...string) []*extmodules.ExternalModule {
	// Der Object Mutext wird verwendet
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es werden alle Externen Module herausgefiltertet
	resolve := make([]*extmodules.ExternalModule, 0)
	for _, item := range names {
		for _, mitem := range o.extModules {
			if item == mitem.GetName() {
				resolve = append(resolve, mitem)
			}
		}
	}

	// Das Objekt wird zurückgegeben
	return resolve
}
