package kmodulerpc

import (
	"vnh1/types"
)

type FunctionType string

const (
	Local  FunctionType = "local"
	Public FunctionType = "public"
)

// Stellt eine Geteilte Funktion dar
type SharedFunction struct {
	kernel             types.KernelInterface
	name               string   // Speichert den Namen der Funktion ab
	parmTypes          []string // Speichert die Parameterdatentypen welche die Funktion erwartet ab
	returnType         string   // Speichert den Rückgabetypen welche die Funktion zurück gibt ab
	functionSourceCode string   // Speichert den Quellcode der Funktion ab
}

// Stellt eine Lokale geteilte Funktion dar
type SharedLocalFunction struct {
	*SharedFunction
}

// Stellt eine Öffentliche geteilte Funktion dar
type SharedPublicFunction struct {
	*SharedFunction
}

// Stellt eine Funktionsanfrage dar
type SharedFunctionRequest struct {
	resolveChan chan *types.FunctionCallState
	parms       *types.RpcRequest
}

type CallFunctionSignature struct {
	*types.FunctionSignature
}

type RpcRequest struct {
	parms []*types.FunctionParameterCapsle
}
