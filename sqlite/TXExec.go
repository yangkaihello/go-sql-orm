package sqlite

import (
	"database/sql"
	"library/databases"
)

type TXExec struct {
	db *sql.DB
	tx *sql.Tx
	datasets map[string]databases.HandleDataset
}

func (this *TXExec) Start(config Config) *TXExec {
	this.db = config.getDb()
	this.tx,_ = this.db.Begin()
	return this
}

func (this *TXExec) Commit()  {
	this.tx.Commit()
}

func (this *TXExec) Rollback()  {
	this.tx.Rollback()
}

func (this *TXExec) Add(tableName string,dataset databases.HandleDataset) databases.HandleTXExec {
	this.datasets[tableName].SetDatasetTXExec(new(TXExec))
	this.datasets[tableName] = dataset
	return this
}