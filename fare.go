package frule_module

import (
	"reflect"
	"regexp"
	"stash.tutu.ru/golang/resources/db"
	"time"
)

type FareRule struct {
	Id                 int     `gorm:"column:id"`
	Partner            *string `gorm:"column:partner"`
	ConnectionGroup    *string `gorm:"column:connection_group"`
	CarrierId          *int64  `gorm:"column:carrier_id"`
	DepartureCityId    *int64  `gorm:"column:departure_city_id"`
	ArrivalCityId      *int64  `gorm:"column:arrival_city_id"`
	DepartureCountryId *int64  `gorm:"column:departure_country_id"`
	ArrivalCountryId   *int64  `gorm:"column:arrival_country_id"`
	FareSpec           *string `gorm:"column:fare_spec"`
	Result             string  `gorm:"column:result"`
	db                 *db.Database
}

func NewFareRule(database *db.Database) FareRule {
	return FareRule{
		db: database,
	}
}

func (f FareRule) GetResultValue(interface{}) interface{} {
	return f.Result
}

func (f FareRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"departure_city_id", "arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "arrival_city_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "arrival_country_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "arrival_city_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "arrival_country_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "partner", "carrier_id", "fare_spec"},
		[]string{"arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"arrival_city_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "partner", "carrier_id", "fare_spec"},
		[]string{"arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"arrival_country_id", "partner", "carrier_id", "fare_spec"},
		[]string{"partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"partner", "carrier_id", "fare_spec"},
	}
}

func (f FareRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{
		"fare_spec": func(a, b reflect.Value) bool {
			fareTest := regexp.MustCompile(a.Elem().Interface().(string))
			return fareTest.Match(b.Elem().Interface().([]byte))
		},
	}
}

func (f FareRule) getStrategyKeys() []string {
	return []string{
		"partner",
		"connection_group",
		"carrier_id",
		"arrival_country_id",
		"departure_country_id",
		"arrival_city_id",
		"departure_city_id",
		"fare_spec",
	}
}

func (f FareRule) getTableName() string {
	return "rm_frule_fare"
}

func (f FareRule) GetDefaultValue() interface{} {
	return ""
}

func (f FareRule) GetDataStorage() (map[int][]FRuler, error) {
	result := make(map[int][]FRuler)
	for rank, fieldList := range f.GetComparisonOrder() {
		query := f.db.Table(f.getTableName())
		for _, field := range f.getStrategyKeys() {
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
			var rowData FareRule

			if err := f.db.ScanRows(rows, &rowData); err != nil {
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

func (f FareRule) GetLastUpdateTime() time.Time {
	return getLastUpdateTime("fare", f.db)
}
