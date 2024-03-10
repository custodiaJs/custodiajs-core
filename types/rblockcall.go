package types

// RBlockCallStatement repräsentiert ein 'rblockcall' Statement in deiner Sprache.
type RBlockCallStatement struct {
	URI        string         // Die URI als String
	Config     *ObjectLiteral // Optionales Konfigurationsobjekt
	Params     []*Param       // Eine Liste von Parametern
	CatchBlock *CatchBlock    // Optionaler Catch-Block
	FinalBlock *FinalBlock    // Optionaler Final-Block
}

// Param repräsentiert einen Parameter im 'rblockcall' Statement.
type Param struct {
	Key   string     // Der Schlüssel des Parameters
	Value Expression // Der Wert des Parameters, kann ein einfacher Wert oder ein komplexerer Ausdruck sein
}

// CatchBlock repräsentiert den 'catch' Block eines 'rblockcall' Statements.
type CatchBlock struct {
	Parameter string          // Der Parameter des Catch-Blocks, typischerweise eine Fehlervariable
	Body      *BlockStatement // Der Körper des Catch-Blocks
}

// FinalBlock repräsentiert den 'final' Block eines 'rblockcall' Statements.
type FinalBlock struct {
	Parameter string          // Der Parameter des Final-Blocks, typischerweise die Resultatvariable
	Body      *BlockStatement // Der Körper des Final-Blocks
}

// BlockStatement repräsentiert einen Block von Anweisungen.
type BlockStatement struct {
	Statements []Statement // Eine Liste von Anweisungen im Block
}

// ObjectLiteral könnte ein Konfigurationsobjekt für 'rblockcall' darstellen.
type ObjectLiteral struct {
	// Felder je nach Bedarf
}

type RBlockCallOptions struct {
}

type RBlockCallPassParms struct {
	Key   string
	Value Param
}
