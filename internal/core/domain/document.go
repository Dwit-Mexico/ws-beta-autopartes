package domain

type DocumentWithDetails struct {
	Document
	Details []DetailDocument `json:"fields"`
}

type EditableDocument struct {
	ID      uint             `json:"id"`
	Name    string           `json:"name"`
	Details []DetailDocument `json:"fields"`
}

type Align string

const (
	AlignStart  Align = "start"
	AlignEnd    Align = "end"
	AlignCenter Align = "center"
)

type DocumentDataTypes string

const (
	NVARCHARThreeHundred DocumentDataTypes = "NVARCHAR(300)"
	NVARCHAROneHundred   DocumentDataTypes = "NVARCHAR(100)"
	NVARCHARFifty        DocumentDataTypes = "NVARCHAR(50)"
	NVARCHARTwenty       DocumentDataTypes = "NVARCHAR(20)"
	INT                  DocumentDataTypes = "INT"
	DECIMALTenTwo        DocumentDataTypes = "DECIMAL(10, 2)"
	DATETIME             DocumentDataTypes = "DATETIME"
)

type DinamicTable struct {
	ID    uint   `json:"uid"`
	Name  string `json:"name"`
	Align Align  `json:"align"`
}

type TableViewDefinition struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Align Align  `json:"align"`
}

type ReportData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TableByID struct {
	Document ReportData               `json:"document"`
	Table    []map[string]interface{} `json:"table"`
	Columns  []TableViewDefinition    `json:"columns"`
}

type UploadDocument struct {
	DocumentID uint   `json:"documentID"`
	File       string `json:"file"`
}
