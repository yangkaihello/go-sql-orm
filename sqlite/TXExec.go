package sqlite

import (
	"database/sql"
	"library/databases"
)

type datasetsMap map[string]databases.HandleDataset

type TXExec struct {
	db *sql.DB
	tx *sql.Tx
	datasets datasetsMap
}

func (this *TXExec) Start(config Config) *TXExec {
	this.db = config.getDb()
	this.tx,_ = this.db.Begin()
	this.datasets = make(datasetsMap)
	return this
}

func (this *TXExec) Commit()  {
	this.tx.Commit()
}

func (this *TXExec) Rollback()  {
	this.tx.Rollback()
}

func (this *TXExec) Add(table *Connect) *TXExec {
	var tableName = table.GetTableName()
	this.datasets[tableName] = table.GetDataset()
	this.datasets[tableName].SetDatasetTXExec(this)
	return this
}

func (this *TXExec) GetTx() *sql.Tx {
	return this.tx
}
