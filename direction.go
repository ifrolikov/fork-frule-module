package frule_module

import (
	"stash.tutu.ru/golang/resources/db"
	"time"
)

type DirectionRule struct {
	Id                 int     `gorm:"column:id"`
	Partner            *string `gorm:"column:partner"`
	ConnectionGroup    *string `gorm:"column:connection_group"`
	CarrierId          *int64  `gorm:"column:carrier_id"`
	DepartureCountryId *uint64 `gorm:"column:departure_country_id"`
	DepartureCityId    *uint64 `gorm:"column:departure_city_id"`
	ArrivalCountryId   *uint64 `gorm:"column:arrival_country_id"`
	ArrivalCityId      *uint64 `gorm:"column:arrival_city_id"`
	Result             bool    `gorm:"column:result"`
	db                 *db.Database
}

func NewDirectionFRule(db *db.Database) DirectionRule {
	return DirectionRule{
		db: db,
	}
}

func (a DirectionRule) GetResultValue(testRule interface{}) interface{} {
	return a.Result
}

func (a DirectionRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "carrier_id"},
		[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
		[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "carrier_id"},
		[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id"},
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "arrival_city_id", "carrier_id"},
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "arrival_city_id"},
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "carrier_id"},
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id"},
		[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "carrier_id"},
		[]string{"partner", "connection_group", "departure_country_id", "departure_city_id"},
		[]string{"partner", "connection_group", "arrival_country_id", "arrival_city_id", "carrier_id"},
		[]string{"partner", "connection_group", "arrival_country_id", "arrival_city_id"},
		[]string{"partner", "connection_group", "departure_country_id", "carrier_id"},
		[]string{"partner", "connection_group", "departure_country_id"},
		[]string{"partner", "connection_group", "arrival_country_id", "carrier_id"},
		[]string{"partner", "connection_group", "arrival_country_id"},
		[]string{"partner", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "carrier_id"},
		[]string{"partner", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
		[]string{"partner", "departure_country_id", "departure_city_id", "arrival_country_id", "carrier_id"},
		[]string{"partner", "departure_country_id", "departure_city_id", "arrival_country_id"},
		[]string{"partner", "departure_country_id", "arrival_country_id", "arrival_city_id", "carrier_id"},
		[]string{"partner", "departure_country_id", "arrival_country_id", "arrival_city_id"},
		[]string{"partner", "departure_country_id", "arrival_country_id", "carrier_id"},
		[]string{"partner", "departure_country_id", "arrival_country_id"},
		[]string{"partner", "departure_country_id", "departure_city_id", "carrier_id"},
		[]string{"partner", "departure_country_id", "departure_city_id"},
		[]string{"partner", "arrival_country_id", "arrival_city_id", "carrier_id"},
		[]string{"partner", "arrival_country_id", "arrival_city_id"},
		[]string{"partner", "departure_country_id", "carrier_id"},
		[]string{"partner", "departure_country_id"},
		[]string{"partner", "arrival_country_id", "carrier_id"},
		[]string{"partner", "arrival_country_id"},
		[]string{"partner", "connection_group", "carrier_id"},
		[]string{"partner", "connection_group"},
		[]string{"partner", "carrier_id"},
		[]string{"partner"},
	}
}

func (a DirectionRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (a DirectionRule) getStrategyKeys() []string {
	return []string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "carrier_id"}
}

func (a DirectionRule) getTableName() string {
	return "rm_frule_direction"
}

func (a DirectionRule) GetDefaultValue() interface{} {
	return false
}

func (a DirectionRule) GetLastUpdateTime() time.Time {
	return getLastUpdateTime("direction", a.db)
}

func (a DirectionRule) GetDataStorage() (map[int][]FRuler, error) {
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
			var rowData DirectionRule

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
