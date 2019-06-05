package frule_module

import (
	"reflect"
	"stash.tutu.ru/golang/resources/db"
	"time"
)

type PartnerPercentRule struct {
	Id                 int     `gorm:"column:id"`
	CarrierId          *int    `gorm:"column:carrier_id"`
	Partner            *string `gorm:"column:partner"`
	ConnectionGroup    *string `gorm:"column:connection_group"`
	DateOfPurchaseFrom *string `gorm:"column:date_of_purchase_from"`
	DateOfPurchaseTo   *string `gorm:"column:date_of_purchase_to"`
	CarrierCountryId   *int    `gorm:"column:carrier_country_id"`
	FareType           *string `gorm:"column:fare_type"`
	Result             float64 `gorm:"column:result"`
	db                 *db.Database
}

func NewPartnerPercentFRule(db *db.Database) PartnerPercentRule {
	return PartnerPercentRule{
		db: db,
	}
}

func (a PartnerPercentRule) GetResultValue(testRule interface{}) interface{} {
	return float64(a.Result)
}

func (a PartnerPercentRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_id", "fare_type", "connection_group"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_id", "fare_type"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_id", "connection_group"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_id"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_country_id", "connection_group"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_country_id"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "connection_group"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to"},
	}
}

func (a PartnerPercentRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{
		"date_of_purchase_from": func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) <= b.Elem().Interface().(string)
		},
		"date_of_purchase_to": func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) > b.Elem().Interface().(string)
		},
	}
}

func (a PartnerPercentRule) getStrategyKeys() []string {
	return []string{"partner", "date_of_purchase_from", "date_of_purchase_to", "connection_group", "carrier_country_id",
		"carrier_id", "fare_type"}
}

func (a PartnerPercentRule) getTableName() string {
	return "rm_frule_partner_percent"
}

func (a PartnerPercentRule) GetDefaultValue() interface{} {
	return 0.0
}

func (a PartnerPercentRule) GetLastUpdateTime() time.Time {
	return getLastUpdateTime("partner_percent", a.db)
}

func (a PartnerPercentRule) GetDataStorage() (map[int][]FRuler, error) {
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
			var rowData PartnerPercentRule

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
