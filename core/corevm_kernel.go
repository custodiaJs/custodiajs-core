package core

func (o *Core) _init_vm_kernel_base(vm *CoreVM) error {
	// Die Runtime wird abgerufen
	vmRuntime := vm.GetRuntime()

	// Die JS Exports werden bereitgestellt
	vmRuntime.Set("exports", vm.vmExports)

	// Der Vorgang ist ohne Fehler durchgeführt wurden
	return nil
}
