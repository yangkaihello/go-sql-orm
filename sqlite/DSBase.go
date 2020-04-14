package sqlite

import (
	"fmt"
	"github.com/yangkaihello/go-sql-orm"
	"strings"
)

type DSBase struct {
	tableName   string
	tableField  []string
	originField []string

	fieldsString string
	fromString   string
	whereString  string
	offsetString string
	limitString  string
	orderString  string

	placeholder []string
}

func (this *DSBase) GetSql() string {
	this.fromString = fmt.Sprintf("FROM `%s`", this.tableName)

	var str = make([]string, 0, 0)
	for _, v := range this.tableField {
		str = append(str, fmt.Sprintf("`%s`", v))
	}
	for _, v := range this.originField {
		str = append(str, fmt.Sprintf("%s", v))
	}
	this.fieldsString = fmt.Sprintf("SELECT %s", strings.Join(str, ","))

	return strings.Trim(fmt.Sprintf("%s %s %s %s %s %s",
		this.fieldsString,
		this.fromString,
		this.whereString,
		this.orderString,
		this.limitString,
		this.offsetString,
	), " ")
}

func (this *DSBase) GetPlaceholder() []string {
	return this.placeholder
}

func (this *DSBase) GetTableName() string {
	return this.tableName
}

func (this *DSBase) GetField() []string {
	return this.tableField
}

func (this *DSBase) GetOriginField() []string {
	return this.originField
}

func (this *DSBase) GetTableField() []string {
	var str = append(this.tableField, this.originField...)
	return str
}

func (this *DSBase) Table(tableName string) databases.HandleMuster {
	this.tableName = tableName
	return this
}

func (this *DSBase) Fields(fields []string) databases.HandleMuster {
	this.tableField = fields
	return this
}

func (this *DSBase) OriginFields(fields []string) databases.HandleMuster {
	this.originField = fields
	return this
}

func (this *DSBase) Where(where databases.Where) databases.HandleMuster {
	if where.GetWhereString() == "" || len(where.GetPlaceholder()) == 0 {
		return this
	}
	this.whereString = "WHERE " + where.GetWhereString()
	this.placeholder = where.GetPlaceholder()
	return this
}

func (this *DSBase) Limit(limit int) databases.HandleMuster {
	this.limitString = fmt.Sprintf("LIMIT %d", limit)
	return this
}

func (this *DSBase) Offset(offset int) databases.HandleMuster {
	this.offsetString = fmt.Sprintf("OFFSET %d", offset)
	return this
}

func (this *DSBase) Order(field string, option string) databases.HandleMuster {
	if option != databases.DATABASE_ORDER_HANDLE_ASC &&
		option != databases.DATABASE_ORDER_HANDLE_DESC {
		option = databases.DATABASE_ORDER_HANDLE_ASC
	}
	this.orderString = fmt.Sprintf("ORDER BY `%s` %s", field, option)
	return this
}
