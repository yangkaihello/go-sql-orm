package factory

import (
	"library/databases"
)

func SingleSqlIte(base databases.HandleMuster,dataset databases.HandleDataset) (databases.HandleMuster,databases.HandleDataset) {
	return databases.SingleMuster(base),databases.SingleDataset(dataset)
}
