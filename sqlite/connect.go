package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yangkaihello/go-sql-orm"
	"github.com/yangkaihello/go-sql-orm/factory"
	"reflect"
	"strings"
)

type Config struct {
	Path string
}

func (this *Config) getDb() *sql.DB {
	db,_ := sql.Open("sqlite3",this.Path)
	return db
}

type Connect struct {
	TemplateTable interface{}
	Config Config
	db            *sql.DB
	databases.HandleDataset
	databases.HandleMuster
}

func (this *Connect) Start(config Config,table interface{}) *Connect {
	switch reflect.TypeOf(table).Kind() {
	case reflect.Struct:
		table = &table
	}
	db := config.getDb()
	ds,data := factory.SingleSqlIte(new(DSBase),new(Dataset))

	//对象的映射
	var templateTableType = reflect.TypeOf(table)
	var templateTableValueFun = reflect.ValueOf(table)

	if templateTableType.Kind() == reflect.Ptr {
		templateTableType = templateTableType.Elem()
	}
	//没有设置table表名的时候默认表名
	if ds.GetTableName() == "" {
		ds.Table(strings.ToLower(templateTableType.Name()))
	}

	//func 定义表名 (验证是否存在func,返回值是否存在，返回值是否是string)
	if tn := templateTableValueFun.MethodByName("TableName");
		tn.IsValid() &&
			tn.Type().NumOut() > 0 &&
			tn.Type().Out(0).Kind() == reflect.String {
		ds.Table(tn.Call(nil)[0].String())
	}

	return &Connect{table,config,db,data,ds}
}

func (this *Connect) Select() *Connect {
	//对象的映射
	var templateTableValue = reflect.Indirect(reflect.ValueOf(this.TemplateTable))
	var templateTableType = reflect.TypeOf(this.TemplateTable)

	if templateTableType.Kind() == reflect.Ptr {
		templateTableType = templateTableType.Elem()
	}

	//判断是否需要设置结构体的字段
	if len(this.GetTableField()) < 1 {
		var f []string
		for i := 0; i < templateTableValue.NumField() ; i++ {
			var field string
			if	templateTableType.Field(i).Tag.Get(databases.TAG_NAME) != "" &&
				templateTableType.Field(i).Tag.Get(databases.TAG_NAME) == databases.TAG_IGNORE {
				continue
			}
			if templateTableType.Field(i).Tag.Get(databases.TAG_NAME) != "" {
				field = templateTableType.Field(i).Tag.Get(databases.TAG_NAME)
			}else{
				field = strings.ToLower(templateTableType.Field(i).Name)
			}
			f = append(f, field)
		}
		this.Fields(f)
	}

	//操作结果集的配置
	this.setDataset()
	this.SetDatasetStructTablesReset()

	return this
}

func (this *Connect) setDataset()  {
	if err := this.db.Ping(); err != nil {
		this.db = this.Config.getDb()
	}

	this.SetDatasetDb(this.db)
	this.SetDatasetTemplateTable(this.TemplateTable)
}

func (this *Connect) One() databases.DataTemplate {
	data,_ := this.GetOne(this.GetSql(),this.GetPlaceholder())
	return data
}

func (this *Connect) All() []databases.DataTemplate {
	data,_ := this.GetAll(this.GetSql(),this.GetPlaceholder())
	return data
}

func (this *Connect) Insert(template databases.DataTemplate) error {
	this.setDataset()

	if len(template) == 0 {
		return errors.New("databases.DataTemplate not")
	}

	var sqlString string
	var fields []string
	var places []string
	var values []string

	for k,v := range template{

		var value string
		switch v.(type) {
		case int:
			value = fmt.Sprintf("%d",v)
		default:
			value = fmt.Sprintf("%s",v)
		}

		fields = append(fields, fmt.Sprintf("`%s`",k))
		places = append(places, "?")
		values = append(values, value)

	}

	sqlString = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",this.GetTableName(),strings.Join(fields,","),strings.Join(places,","))
	_,err := this.SetExec(sqlString,values)
	return err
}

func (this *Connect) Update(where databases.Where,template databases.DataTemplate) error {
	this.setDataset()

	var sqlString string
	var whereString string
	var fields []string
	var values []string

	if whereString = this.WhereMake(where); whereString == "" || len(template) == 0 {
		return errors.New("databases.Where not OR databases.DataTemplate not")
	}

	for k,v := range template{
		fields = append(fields, fmt.Sprintf("`%s` = ?",k))
		values = append(values, reflect.ValueOf(v).String())
	}

	sqlString = fmt.Sprintf("UPDATE %s SET %s WHERE %s",this.GetTableName(),strings.Join(fields,","),whereString)
	_,err := this.SetExec(sqlString,append(values,this.GetPlaceholder()...))
	return err
}

func (this *Connect) Delete(where databases.Where) error {
	this.setDataset()

	var sqlString string
	var whereString string
	if whereString = this.WhereMake(where); whereString == "" {
		return errors.New("databases.Where not OR databases.DataTemplate not")
	}

	sqlString = fmt.Sprintf("DELETE FROM %s WHERE %s",this.GetTableName(),whereString)
	_,err := this.SetExec(sqlString,this.GetPlaceholder())
	return err
}

func (this *Connect) Exec(sqlString string,placeholder... string) error {
	this.setDataset()

	_,err := this.SetExec(sqlString,placeholder)
	return err
}

func (this *Connect) Query(sqlString string,placeholder... string) []databases.DataTemplate {
	this.setDataset()

	data,_ := this.GetAll(sqlString,placeholder)
	return data
}

func (this *Connect) Close() {
	this.db.Close()
}





