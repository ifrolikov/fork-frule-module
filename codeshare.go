package frule_module

import (
	"stash.tutu.ru/golang/resources/db"
	"time"
)

type CodeshareRule struct {
	Id               int     `gorm:"column:id"`
	Partner          *string `gorm:"column:partner"`
	ConnectionGroup  *string `gorm:"column:connection_group"`
	CarrierOperating *int64  `gorm:"column:carrier_operating"`
	CarrierMarketing *int64  `gorm:"column:carrier_marketing"`
	ServiceClass     *string `gorm:"column:service_class"`
	Result           bool    `gorm:"column:result"`
	db               *db.Database
}

func NewCodeshareFRule(db *db.Database) CodeshareRule {
	return CodeshareRule{
		db: db,
	}
}

func (a CodeshareRule) GetResultValue(testRule interface{}) interface{} {
	return a.Result
}

func (a CodeshareRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"partner", "connection_group", "carrier_operating", "carrier_marketing", "service_class"},
		[]string{"partner", "connection_group", "carrier_operating", "carrier_marketing"},
		[]string{"partner", "connection_group", "service_class"},
		[]string{"partner", "connection_group"},
		[]string{"partner", "carrier_operating", "carrier_marketing", "service_class"},
		[]string{"partner", "carrier_operating", "carrier_marketing"},
		[]string{"partner", "service_class"},
		[]string{"partner"},
	}
}

func (a CodeshareRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (a CodeshareRule) getStrategyKeys() []string {
	return []string{"partner", "connection_group", "carrier_operating", "carrier_marketing", "service_class"}
}

func (a CodeshareRule) getTableName() string {
	return "rm_frule_codeshare"
}

func (a CodeshareRule) GetDefaultValue() interface{} {
	return false
}

func (a CodeshareRule) GetLastUpdateTime() time.Time {
	return getLastUpdateTime("codeshare", a.db)
}

func (a CodeshareRule) GetDataStorage() (map[int][]FRuler, error) {
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
			return result, err
		}

		for rows.Next() {
			var rowData CodeshareRule

			if err := a.db.ScanRows(rows, &rowData); err != nil {
				return result, err
			}
			result[rank] = append(result[rank], rowData)

		}
		err = rows.Close()
		if err != nil {
			return result, err
		}
	}
	return result, nil
}
