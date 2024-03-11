package types

type RBlockCallStatement struct {
	URI        string // Die URI als String
	Config     map[string]interface{}
	Params     []*Param    // Eine Liste von Parametern
	CatchBlock *CatchBlock // Optionaler Catch-Block
	FinalBlock *FinalBlock // Optionaler Final-Block
}

type Param struct {
	Key   string     // Der Schlüssel des Parameters
	Value Expression // Der Wert des Parameters, kann ein einfacher Wert oder ein komplexerer Ausdruck sein
}

type CatchBlock struct {
	Parameter string          // Der Parameter des Catch-Blocks, typischerweise eine Fehlervariable
	Body      *BlockStatement // Der Körper des Catch-Blocks
}

type FinalBlock struct {
	Parameter string          // Der Parameter des Final-Blocks, typischerweise die Resultatvariable
	Body      *BlockStatement // Der Körper des Final-Blocks
}

type BlockStatement struct {
	Statements []Statement // Eine Liste von Anweisungen im Block
}

type RBlockCallOptions struct {
}

type RBlockCallPassParms struct {
	Key   string
	Value Param
}
