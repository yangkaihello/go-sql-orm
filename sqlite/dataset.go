package sqlite

import (
	"database/sql"
	"github.com/yangkaihello/go-sql-orm"
)

type Dataset struct {
	db *sql.DB
	TemplateTable interface{}
	TXExec databases.HandleTXExec
	StructTables []databases.DataTemplate
}

func (this *Dataset) SetDatasetDb(db *sql.DB)  {
	this.db = db
}

func (this *Dataset) SetDatasetTXExec(TXExec databases.HandleTXExec)  {
	this.TXExec = TXExec
}

func (this *Dataset) SetDatasetTemplateTable(table interface{})  {
	this.TemplateTable = table
}

func (this *Dataset) SetDatasetStructTablesReset()  {
	this.StructTables = []databases.DataTemplate{}
}

func (this *Dataset) GetDataset() databases.HandleDataset {
	return this
}

func (this *Dataset) GetOne(sqlString string,placeholder []string) (databases.DataTemplate,error) {
	if len(this.StructTables) != 0 {
		return this.StructTables[0],nil
	}
	var rows *sql.Rows
	var err error
	var p = make([]interface{},len(placeholder),len(placeholder))

	for k,v := range placeholder {
		p[k] = v
	}

	if rows,err = this.db.Query(sqlString+" LIMIT 1",p...); err != nil {
		return nil,err
	}
	columns,_ := rows.Columns()
	var templateInterface = make([]interface{},len(columns))
	var templateScan = make([]interface{},len(columns))

	for k := range templateInterface {
		templateScan[k] = &templateInterface[k]
	}
	for rows.Next() {
		var DataTemplate = make(databases.DataTemplate)

		rows.Scan(templateScan...)
		for k,v := range templateInterface {
			DataTemplate[columns[k]] = v
		}
		this.StructTables = append(this.StructTables, DataTemplate)
	}
	rows.Close()

	var table databases.DataTemplate
	if len(this.StructTables) != 0 {
		table = this.StructTables[0]
	}
	return table,nil
}

func (this *Dataset) GetAll(sqlString string,placeholder []string) ([]databases.DataTemplate,error) {
	if len(this.StructTables) != 0 {
		return this.StructTables,nil
	}
	var rows *sql.Rows
	var err error
	var p = make([]interface{},len(placeholder),len(placeholder))

	for k,v := range placeholder {
		p[k] = v
	}

	if rows,err = this.db.Query(sqlString,p...); err != nil {
		return nil,err
	}
	columns,_ := rows.Columns()
	var templateInterface = make([]interface{},len(columns))
	var templateScan = make([]interface{},len(columns))

	for k := range templateInterface {
		templateScan[k] = &templateInterface[k]
	}
	for rows.Next() {
		var DataTemplate = make(databases.DataTemplate)

		rows.Scan(templateScan...)
		for k,v := range templateInterface {
			DataTemplate[columns[k]] = v
		}
		this.StructTables = append(this.StructTables, DataTemplate)
	}
	rows.Close()

	return this.StructTables,nil
}

func (this *Dataset) SetExec(sqlString string,placeholder []string) (result sql.Result,err error) {
	var p = make([]interface{},len(placeholder),len(placeholder))
	for k,v := range placeholder {
		p[k] = v
	}
	if this.TXExec != nil {
		begin := this.TXExec.GetTx()
		result,err = begin.Exec(sqlString,p...)
	}else{
		result,err = this.db.Exec(sqlString,p...)
	}

	return result,err
}


