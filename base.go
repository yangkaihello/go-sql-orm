package databases

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const (
	DATABASE_WHERE_HANDLE_AND  = "AND"
	DATABASE_WHERE_HANDLE_OR   = "OR"
	DATABASE_ORDER_HANDLE_ASC  = "ASC"
	DATABASE_ORDER_HANDLE_DESC = "DESC"
	TAG_NAME                   = "model"
	TAG_IGNORE                 = "-"
)

type DataTemplate map[string]interface{}

type WhereOperation struct {
	Key         string
	Handle      string
	Value       string
	WhereOption *WhereOption
}

type WhereOption struct {
	Option    string
	Operation []WhereOperation
}

type Where struct {
	placeholder    []string //最后的参数值
	wheres         string
	WhereInterface []interface{}
}

func (this *Where) Add(wheres WhereOption) *Where {
	this.WhereInterface = append(this.WhereInterface, append(this.whereMap(wheres), wheres.Option))
	//var whereString []string
	//for _,value := range where {
	//	whereString = append(whereString,fmt.Sprintf("`%s` %s ?",strings.Trim(value.Key,""),strings.Trim(value.Handle,"")))
	//	this.Placeholder = append(this.Placeholder,value.Value)
	//}
	//this.WhereString = append(this.WhereString,strings.Join(whereString," "+option+" "))
	return this
}

func (this *Where) GetWhereString() string {
	this.wheres = ""
	for _, value := range this.WhereInterface {
		this.whereString(value)
	}
	if this.wheres != "" {
		byteWhere := []byte(this.wheres)
		this.wheres = strings.Trim(string(byteWhere[0:len(byteWhere)-4]), " ")
	}
	return this.wheres
}

func (this *Where) GetPlaceholder() []string {
	return this.placeholder
}

func (this *Where) whereString(s interface{}) {

	for _, value := range s.([]interface{}) {

		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			this.wheres += value.(string) + " "
		case reflect.Slice:
			this.whereString(value)
		}

	}
}

func (this *Where) whereMap(wheres WhereOption) []interface{} {

	var option = wheres.Option

	if option != DATABASE_WHERE_HANDLE_AND &&
		option != DATABASE_WHERE_HANDLE_OR {
		option = DATABASE_WHERE_HANDLE_AND
	}
	var sliceInterface = []interface{}{}
	for _, value := range wheres.Operation {
		if value.WhereOption == nil {
			sliceInterface = append(sliceInterface, fmt.Sprintf("`%s` %s ?", value.Key, value.Handle), option)
			this.placeholder = append(this.placeholder, value.Value)
		} else {
			sliceInterface = append(sliceInterface, "(", this.whereMap(*value.WhereOption), ")", option)
		}
	}

	if len(sliceInterface) != 0 {
		sliceInterface = sliceInterface[0 : len(sliceInterface)-1]
	}

	return sliceInterface
}

//func (this *Where) WhereAnd(where []WhereOperation) *Where {
//	return this.Where(where,DATABASE_WHERE_HANDLE_AND)
//}
//
//func (this *Where) WhereOr(where []WhereOperation) *Where {
//	return this.Where(where,DATABASE_WHERE_HANDLE_OR)
//}

//func (this *Where) Clean(option []string) (*Where , error) {
//	if len(this.WhereString)-1 != len(option) {
//		return nil,errors.New("option length incorrect")
//	}
//	for _,value := range option{
//		if value != DATABASE_WHERE_HANDLE_AND &&
//			value != DATABASE_WHERE_HANDLE_OR {
//			return nil,errors.New("option describe the error")
//		}
//	}
//	this.Option = option
//	return this,nil
//}

type HandleTXExec interface {
	GetTx() *sql.Tx
}

type HandleDataset interface {
	GetDataset() HandleDataset
	GetOne(sqlString string, placeholder []string) (DataTemplate, error)
	GetAll(sqlString string, placeholder []string) ([]DataTemplate, error)
	SetExec(sqlString string, placeholder []string) (result sql.Result, err error)
	SetDatasetDb(db *sql.DB)
	SetDatasetTXExec(TXExec HandleTXExec)
	SetDatasetTemplateTable(table interface{})
	SetDatasetStructTablesReset()
}

//构造sql的结构体
type HandleMuster interface {
	Table(table string) HandleMuster
	Fields(fields []string) HandleMuster
	OriginFields(fields []string) HandleMuster
	Where(where Where) HandleMuster
	Limit(limit int) HandleMuster
	Offset(offset int) HandleMuster
	Order(field string, option string) HandleMuster
	GetSql() string
	GetField() []string
	GetOriginField() []string
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
