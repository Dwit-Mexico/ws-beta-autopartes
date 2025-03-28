package domain

type ListString []string

type ObjectString map[string][]string
type Align string

const (
	AlignStart  Align = "start"
	AlignEnd    Align = "end"
	AlignCenter Align = "center"
)

type TableViewDefinition struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Align Align  `json:"align"`
}
