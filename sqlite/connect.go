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
	db, _ := sql.Open("sqlite3", this.Path)
	return db
}

type Connect struct {
	TemplateTable interface{}
	Config        Config
	db            *sql.DB
	modelField map[string]string
	databases.HandleDataset
	databases.HandleMuster
}

func (this *Connect) Start(config Config, table interface{}) *Connect {
	switch reflect.TypeOf(table).Kind() {
	case reflect.Struct:
		table = &table
	}
	db := config.getDb()
	ds, data := factory.SingleSqlIte(new(DSBase), new(Dataset))

	//对象的映射
	var templateTableType = reflect.TypeOf(table)
	var templateTableValue = reflect.ValueOf(table)
	var templateTableValueElem = templateTableValue.Elem()

	if templateTableType.Kind() == reflect.Ptr {
		templateTableType = templateTableType.Elem()
	}
	//没有设置table表名的时候默认表名
	if ds.GetTableName() == "" {
		ds.Table(strings.ToLower(templateTableType.Name()))
	}

	//func 定义表名 (验证是否存在func,返回值是否存在，返回值是否是string)
	if tn := templateTableValue.MethodByName("TableName");
		tn.IsValid() &&
			tn.Type().NumOut() > 0 &&
			tn.Type().Out(0).Kind() == reflect.String {
		ds.Table(tn.Call(nil)[0].String())
	}

	//func 预定义加入时间
	if tn := templateTableValue.MethodByName("CreatedAt");
		tn.IsValid() &&
			tn.Type().NumOut() > 0 &&
			tn.Type().Out(0).Kind() == reflect.String {
		//this.created_at = tn.Call(nil)[0].String()
	}

	var f  = map[string]string{}
	for i := 0; i < templateTableValueElem.NumField(); i++ {
		var field string

		if d, ok := templateTableType.Field(i).Tag.Lookup(databases.TAG_NAME); !ok || d == databases.TAG_IGNORE {
			continue
		}

		if templateTableType.Field(i).Tag.Get(databases.TAG_NAME) != "" {
			field = templateTableType.Field(i).Tag.Get(databases.TAG_NAME)
		} else {
			field = strings.ToLower(templateTableType.Field(i).Name)
		}
		f[field] = ""
	}

	return &Connect{table, config, db,f, data, ds}
}

func (this *Connect) Select() *Connect {
	//对象的映射
	var templateTableType = reflect.TypeOf(this.TemplateTable)

	if templateTableType.Kind() == reflect.Ptr {
		templateTableType = templateTableType.Elem()
	}

	//判断是否需要设置结构体的字段
	if len(this.GetTableField()) < 1 {
		var field []string
		for k,_ := range this.modelField {
			field = append(field, k)
		}
		this.Fields(field)
	}

	//操作结果集的配置
	this.setDataset()
	this.SetDatasetStructTablesReset()

	return this
}

func (this *Connect) setDataset() {
	if err := this.db.Ping(); err != nil {
		this.db = this.Config.getDb()
	}

	this.SetDatasetDb(this.db)
	this.SetDatasetTemplateTable(this.TemplateTable)
}

func (this *Connect) One() databases.DataTemplate {
	data, _ := this.GetOne(this.GetSql(), this.GetPlaceholder())
	return data
}

func (this *Connect) All() []databases.DataTemplate {
	data, _ := this.GetAll(this.GetSql(), this.GetPlaceholder())
	return data
}

func (this *Connect) Count() int64 {
	this.setDataset()
	this.SetDatasetStructTablesReset()
	var fields = this.GetField()
	var originField = this.GetOriginField()
	//重新设置field
	this.OriginFields([]string{"count(*) as count"}).Fields([]string{})
	data, _ := this.GetOne(this.GetSql(), this.GetPlaceholder())
	//还原field
	this.Fields(fields).OriginFields(originField)
	//log.Println(data,err)
	if d, ok := data["count"]; ok {
		return d.(int64)
	} else {
		return 0
	}

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

	//对象的映射
	var templateTableValue = reflect.ValueOf(this.TemplateTable)

	//自动创建时间数据
	_,CreatedAtOkOne := template["created_at"]
	_,CreatedAtOkTwo := this.modelField["created_at"]
	if !CreatedAtOkOne && CreatedAtOkTwo {
		if tn := templateTableValue.MethodByName("CreatedAt");
			tn.IsValid() &&
				tn.Type().NumOut() > 0 &&
				tn.Type().Out(0).Kind() == reflect.String {
			template["created_at"] = tn.Call(nil)[0].String()
		}
	}

	//自动创建时间数据
	_,UpdatedAtOkOne := template["updated_at"]
	_,UpdatedAtOkTwo := this.modelField["updated_at"]
	if !UpdatedAtOkOne && UpdatedAtOkTwo {
		if tn := templateTableValue.MethodByName("UpdatedAt");
			tn.IsValid() &&
				tn.Type().NumOut() > 0 &&
				tn.Type().Out(0).Kind() == reflect.String {
			template["updated_at"] = tn.Call(nil)[0].String()
		}
	}

	for k, v := range template {

		var value string
		switch v.(type) {
		case int:
			value = fmt.Sprintf("%d", v)
		default:
			value = fmt.Sprintf("%s", v)
		}
		if _,ok := this.modelField[k]; !ok {
			continue
		}

		fields = append(fields, fmt.Sprintf("`%s`", k))
		places = append(places, "?")
		values = append(values, value)

	}

	sqlString = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", this.GetTableName(), strings.Join(fields, ","), strings.Join(places, ","))
	_, err := this.SetExec(sqlString, values)
	return err
}

func (this *Connect) Update(where databases.Where, template databases.DataTemplate) error {
	this.setDataset()

	var sqlString string
	var whereString string
	var fields []string
	var values []string

	if whereString = where.GetWhereString(); whereString == "" || len(template) == 0 {
		return errors.New("databases.Where not OR databases.DataTemplate not")
	}

	//对象的映射
	var templateTableValue = reflect.ValueOf(this.TemplateTable)

	//自动创建时间数据
	_,UpdatedAtOkOne := template["updated_at"]
	_,UpdatedAtOkTwo := this.modelField["updated_at"]
	if !UpdatedAtOkOne && UpdatedAtOkTwo {
		if tn := templateTableValue.MethodByName("UpdatedAt");
			tn.IsValid() &&
				tn.Type().NumOut() > 0 &&
				tn.Type().Out(0).Kind() == reflect.String {
			template["updated_at"] = tn.Call(nil)[0].String()
		}
	}

	for k, v := range template {
		if _,ok := this.modelField[k]; !ok {
			continue
		}
		fields = append(fields, fmt.Sprintf("`%s` = ?", k))
		values = append(values, reflect.ValueOf(v).String())
	}

	sqlString = fmt.Sprintf("UPDATE %s SET %s WHERE %s", this.GetTableName(), strings.Join(fields, ","), whereString)
	_, err := this.SetExec(sqlString, append(values, this.GetPlaceholder()...))
	return err
}

func (this *Connect) Delete(where databases.Where) error {
	this.setDataset()

	var sqlString string
	var whereString string
	if whereString = where.GetWhereString(); whereString == "" {
		return errors.New("databases.Where not")
	}

	sqlString = fmt.Sprintf("DELETE FROM %s WHERE %s", this.GetTableName(), whereString)
	_, err := this.SetExec(sqlString, this.GetPlaceholder())
	return err
}

func (this *Connect) Exec(sqlString string, placeholder ...string) error {
	this.setDataset()

	_, err := this.SetExec(sqlString, placeholder)
	return err
}

func (this *Connect) Query(sqlString string, placeholder ...string) []databases.DataTemplate {
	this.setDataset()

	data, _ := this.GetAll(sqlString, placeholder)
	return data
}

func (this *Connect) Close() {
	this.db.Close()
}
