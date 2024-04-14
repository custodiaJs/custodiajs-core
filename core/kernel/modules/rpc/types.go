package kmodulerpc

import (
	v8 "rogchap.com/v8go"
)

type SharedLocalFunction struct {
	v8VM         *v8.Context
	callFunction *v8.Function
	name         string
	parmTypes    []string
	returnType   string
}

type SharedPublicFunction struct {
	v8VM         *v8.Context
	callFunction *v8.Function
	name         string
	parmTypes    []string
	returnType   string
}

type FunctionCallReturn struct {
	CType string
	Value interface{}
}
