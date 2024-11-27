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

type DocumentRowRecord struct {
	ID          uint   `json:"id"`
	Field       string `json:"field"`
	DocumentKey string `json:"documentKey"`
	TypeField   string `json:"typeField"`
	Value       any    `json:"value"`
}

type EditableDoumentRowRecord struct {
	ID         uint `json:"id"`
	DocumentID uint `json:"documentID"`
	Records    any  `json:"records"`
}

type TableByID struct {
	Document ReportData               `json:"document"`
	Table    []map[string]interface{} `json:"table"`
	Columns  []TableViewDefinition    `json:"columns"`
}

type ReportByID struct {
	Report  ReportData               `json:"report"`
	Table   []map[string]interface{} `json:"table"`
	Columns []TableViewDefinition    `json:"columns"`
	Charts  []ChartReport            `json:"charts"`
}

type UploadDocument struct {
	DocumentID uint   `json:"documentID"`
	File       string `json:"file"`
}
