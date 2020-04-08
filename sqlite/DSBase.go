package sqlite

import (
	"fmt"
	"github.com/yangkaihello/go-sql-orm"
	"strings"
)

type DSBase struct {
	tableName string
	tableField []string

	fieldsString string
	fromString string
	whereString string
	offsetString string
	limitString string
	orderString string

	placeholder []string
}

func (this *DSBase) GetSql() string {
	this.fromString = fmt.Sprintf("FROM `%s`",this.tableName)

	var str = make([]string,len(this.tableField),len(this.tableField))
	for k,v := range this.tableField{
		str[k] = fmt.Sprintf("`%s`",v)
	}
	this.fieldsString = fmt.Sprintf("SELECT %s",strings.Join(str,","))

	return strings.Trim(fmt.Sprintf("%s %s %s %s %s %s",
		this.fieldsString,
		this.fromString,
		this.whereString,
		this.orderString,
		this.limitString,
		this.offsetString,
		)," ")
}

func (this *DSBase) GetPlaceholder() []string {
	return this.placeholder
}

func (this *DSBase) GetTableName() string {
	return this.tableName
}

func (this *DSBase) GetTableField() []string {
	return this.tableField
}

func (this *DSBase) Table(tableName string) databases.HandleMuster {
	this.tableName = tableName
	return this
}

func (this *DSBase) Fields(fields []string) databases.HandleMuster {
	this.tableField = fields
	return this
}

func (this *DSBase) WhereMake(where databases.Where) string {
	this.placeholder = where.Placeholder
	var whereString []string

	for key := range where.WhereString {

		if strings.Contains(where.WhereString[key],databases.DATABASE_WHERE_HANDLE_AND) ||
			strings.Contains(where.WhereString[key],databases.DATABASE_WHERE_HANDLE_OR) {
			where.WhereString[key] = fmt.Sprintf("(%s)",where.WhereString[key])
		}

		if key < 1 {
			whereString = append(whereString,where.WhereString[key])
		}else{
			whereString = append(whereString,where.Option[key-1])
			whereString = append(whereString,where.WhereString[key])
		}

	}
	return strings.Join(whereString," ")
}

func (this *DSBase) Where(where databases.Where) databases.HandleMuster {
	this.whereString = "WHERE "+this.WhereMake(where)
	return this
}

func (this *DSBase) Limit(limit int) databases.HandleMuster {
	this.limitString = fmt.Sprintf("LIMIT %d",limit)
	return this
}

func (this *DSBase) Offset(offset int) databases.HandleMuster {
	this.offsetString = fmt.Sprintf("OFFSET %d",offset)
	return this
}

func (this *DSBase) Order(field string,option string) databases.HandleMuster {
	if option != databases.DATABASE_ORDER_HANDLE_ASC &&
		option != databases.DATABASE_ORDER_HANDLE_DESC {
		option = databases.DATABASE_ORDER_HANDLE_ASC
	}
	this.orderString = fmt.Sprintf("ORDER BY `%s` %s",field,option)
	return this
}



