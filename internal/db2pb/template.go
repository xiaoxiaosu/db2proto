package db2pb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

const structTpl =
`syntax = "proto3"

package {{.TableName}}

service {{.TableName}} {
rpc Add{{camel $.TableName}}(Add{{camel $.TableName}}Request) returns (Add{{camel $.TableName}}Reply)
rpc Update{{camel $.TableName}}(Add{{camel $.TableName}}Request) returns (Update{{camel $.TableName}}Reply)
}

message Add{{camel $.TableName}}Request {
{{range $index, $value := .Columns}}{{.Type}} {{.Name}} = {{inc $index}}
{{end}}
}

message Update{{camel $.TableName}}Request {
{{range $index, $value := .Columns}}{{.Type}} {{.Name}} = {{inc $index}}
{{end}}
}

message Add{{camel $.TableName}}Reply {
int32 code = 1
string msg = 2
}

message Update{{camel $.TableName}}Reply {
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
		"camel": func(s string) string {
			s = strings.Replace(s, "_", " ", -1)
			s = strings.Title(s)
			return strings.Replace(s, " ", "", -1)
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