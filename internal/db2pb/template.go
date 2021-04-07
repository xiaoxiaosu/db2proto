package db2pb

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

const structTpl =
`syntax = "proto3"

package {{.TableName}}

service {{.TableName}} {
rpc Add{{$.TableName}}(Add{{$.TableName}}Request) returns (Add{{$.TableName}}Reply)
rpc Update{{$.TableName}}(Add{{$.TableName}}Request) returns (Update{{$.TableName}}Reply)
}

message Add{{$.TableName}}Request {
{{range $index, $value := .Columns}}
{{.Type}} {{.Name}} = {{inc $index}}
{{end}}
}

message Update{{$.TableName}}Request {
{{range $index, $value := .Columns}}
{{.Type}} {{.Name}} = {{inc $index}}
{{end}}
}

message Add{{$.TableName}}Reply {
int32 code = 1
string msg = 2
}

message Update{{$.TableName}}Reply {
int32 code = 1
string msg = 2
}
`

type StructTemplate struct {
	structTpl string
}

type StructColumn struct {
	Name string
	Type string
	Comment string
}

type StructTemplateDb struct {
	TableName string
	Columns []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{structTpl: structTpl}
}

func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column :=range tbColumns {
		tplColumns = append(tplColumns, &StructColumn{
			Name: column.ColumnName,
			Type: DBTypeToPbType[column.DataType],
			Comment: column.ColumnComment,
		})
	}

	return tplColumns
}

func (t *StructTemplate) Generate(tableName string, tplColumns []*StructColumn) error {
	tpl,  _ := template.Must(template.New("db2pb"), nil).Funcs(template.FuncMap{
		"inc": func(i int) int{
			return i+1
		},
	}).Parse(t.structTpl)
	tplDB := StructTemplateDb{
		TableName: tableName,
		Columns: tplColumns,
	}
	f, err := os.Create("./"+tableName+".proto")
	if err != nil {
		log.Fatalf("create file err:%v", err)
	}

	err = tpl.Execute(f, tplDB)
	if err != nil {
		return err
	}

	fmt.Printf("%v generate succ",f.Name())

	return nil
}