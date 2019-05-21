package frule_module

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type AirlineRule struct {
	Id              int     `sql:"id"`
	CarrierId       *int    `sql:"carrier_id"`
	Partner         *string `sql:"partner"`
	ConnectionGroup *string `sql:"connection_group"`
	Result          bool    `sql:"result"`
	db              *gorm.DB
}

func NewAirlineFRule(db *gorm.DB) AirlineRule {
	return AirlineRule{
		db: db,
	}
}

func (a AirlineRule) GetResultValue() interface{} {
	return a.Result
}

func (a AirlineRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"carrier_id", "partner", "connection_group"},
		[]string{"partner", "connection_group"},
		[]string{"carrier_id", "partner"},
		[]string{"partner"},
	}
}

func (a AirlineRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (a AirlineRule) getStrategyKeys() []string {
	return []string{"carrier_id", "partner", "connection_group"}
}

func (a AirlineRule) GetIndexedKeys() []string {
	return []string{"carrier_id", "partner", "connection_group"}
}

func (a AirlineRule) getTableName() string {
	return "rm_frule_airline"
}

func (a AirlineRule) GetDefaultValue() interface{} {
	return false
}

func (a AirlineRule) GetDataStorage() map[int][]FRuler {
	result := make(map[int][]FRuler)
	for rank, fieldList := range a.GetComparisonOrder() {
		query := a.db.Table(a.getTableName())
		for _, field := range a.getStrategyKeys() {
			if inSlice(field, fieldList) {
				query = query.Where(field + " IS NOT NULL")
			} else {
				query = query.Where(field + " IS NULL")
			}
		}
		rows, err := query.Rows()
		if err != nil {
			fmt.Println(err)
		}

		for rows.Next() {
			var rowData AirlineRule

			if err := a.db.ScanRows(rows, &rowData); err != nil {
				fmt.Println(err)
			}
			result[rank] = append(result[rank], rowData)

		}
		err = rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	return result
}
