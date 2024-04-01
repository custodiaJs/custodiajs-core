package httpapi

func (o *FunctionParameterCapsle) GetType() string {
	return o.CType
}

func (o *FunctionParameterCapsle) GetValue() interface{} {
	return o.Value
}
