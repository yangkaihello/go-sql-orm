package databases

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	DATABASE_WHERE_HANDLE_AND  = "AND"
	DATABASE_WHERE_HANDLE_OR   = "OR"
	DATABASE_ORDER_HANDLE_ASC  = "ASC"
	DATABASE_ORDER_HANDLE_DESC = "DESC"
	TAG_NAME = "model"
	TAG_IGNORE = "-"
)

type DataTemplate map[string]interface{}

type WhereOperation struct {
	Handle string
	Key string
	Value string
}

type Where struct {
	Placeholder []string	//最后的参数值
	WhereString []string	//最后需要整合的where条件
	Option []string			//最后需要整合的where操作符
}

func (this *Where) Where(where []WhereOperation,option string) *Where {

	if option != DATABASE_WHERE_HANDLE_AND &&
		option != DATABASE_WHERE_HANDLE_OR {
		option = DATABASE_WHERE_HANDLE_AND
	}

	var whereString []string
	for _,value := range where {
		whereString = append(whereString,fmt.Sprintf("`%s` %s ?",strings.Trim(value.Key,""),strings.Trim(value.Handle,"")))
		this.Placeholder = append(this.Placeholder,value.Value)
	}
	this.WhereString = append(this.WhereString,strings.Join(whereString," "+option+" "))
	return this
}

func (this *Where) WhereOr(where []WhereOperation) *Where {
	var option = DATABASE_WHERE_HANDLE_OR
	var whereString []string
	for _,value := range where {
		whereString = append(whereString,fmt.Sprintf("`%s` %s ?",strings.Trim(value.Key,""),strings.Trim(value.Handle,"")))
		this.Placeholder = append(this.Placeholder,value.Value)
	}
	this.WhereString = append(this.WhereString,strings.Join(whereString," "+option+" "))
	return this
}

func (this *Where) Clean(option []string) (*Where , error) {
	if len(this.WhereString)-1 != len(option) {
		return nil,errors.New("option length incorrect")
	}
	for _,value := range option{
		if value != DATABASE_WHERE_HANDLE_AND &&
			value != DATABASE_WHERE_HANDLE_OR {
			return nil,errors.New("option describe the error")
		}
	}
	this.Option = option
	return this,nil
}

type HandleTXExec interface {
	Add(tableName string,dataset HandleDataset) HandleTXExec
}

type HandleDataset interface {
	GetDataset() HandleDataset
	GetOne(sqlString string,placeholder []string) (DataTemplate,error)
	GetAll(sqlString string,placeholder []string) ([]DataTemplate,error)
	SetExec(sqlString string,placeholder []string) (result sql.Result,err error)
	SetDatasetDb(db *sql.DB)
	SetDatasetTXExec(TXExec HandleTXExec)
	SetDatasetTemplateTable(table interface{})
	SetDatasetStructTablesReset()
}

//构造sql的结构体
type HandleMuster interface {
	Table(table string) HandleMuster
	Fields(fields []string) HandleMuster
	Where(where Where) HandleMuster
	Limit(limit int) HandleMuster
	Offset(offset int) HandleMuster
	Order(field string,option string) HandleMuster
	WhereMake(where Where) string
	GetSql() string
	GetTableName() string
	GetTableField() []string
	GetPlaceholder() []string
}

func SingleMuster(DS HandleMuster) HandleMuster {
	return DS
}

func SingleDataset(data HandleDataset) HandleDataset {
	return data
}

