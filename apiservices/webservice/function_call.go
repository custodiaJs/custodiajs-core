package webservice

type FunctionParameterCapsle struct {
	Value interface{}
	CType string
}

func (o *FunctionParameterCapsle) GetType() string {
	return o.CType
}

func (o *FunctionParameterCapsle) GetValue() interface{} {
	return o.Value
}
