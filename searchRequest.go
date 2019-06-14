package frule_module

import (
	"encoding/json"
	"stash.tutu.ru/golang/resources/db"
	"strconv"
	"time"
)

type SearchRequest struct {
	Id                 int     `gorm:"column:id"`
	ConnectionGroup    *string `gorm:"column:connection_group"`
	DepartureCityId    *uint64 `gorm:"column:departure_city_id"`
	ArrivalCityId      *uint64 `gorm:"column:arrival_city_id"`
	DepartureCountryId *uint64 `gorm:"column:departure_country_id"`
	ArrivalCountryId   *uint64 `gorm:"column:arrival_country_id"`
	ServiceClass       *string `gorm:"column:service_class"`
	Result             string  `gorm:"column:result"`
	db                 *db.Database
}

func NewSearchRequest(database *db.Database) SearchRequest {
	return SearchRequest{
		db: database,
	}
}

func (sr SearchRequest) GetResultValue(interface{}) interface{} {
	var result map[string]string
	err := json.Unmarshal([]byte(sr.Result), &result)
	if err != nil {
		return false
	}
	for key, value := range result {
		if cronSpec(&key, time.Now()) {
			val, err := strconv.ParseBool(value)
			if err != nil {
				return false
			}
			return val
		}
	}
	return false
}

func (sr SearchRequest) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "service_class"},
		[]string{"connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
		[]string{"connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "service_class"},
		[]string{"connection_group", "departure_country_id", "departure_city_id", "arrival_country_id"},
		[]string{"connection_group", "departure_country_id", "arrival_country_id", "arrival_city_id", "service_class"},
		[]string{"connection_group", "departure_country_id", "arrival_country_id", "arrival_city_id"},
		[]string{"connection_group", "departure_country_id", "departure_city_id", "service_class"},
		[]string{"connection_group", "departure_country_id", "departure_city_id"},
		[]string{"connection_group", "arrival_country_id", "arrival_city_id", "service_class"},
		[]string{"connection_group", "arrival_country_id", "arrival_city_id"},
		[]string{"connection_group", "departure_country_id", "arrival_country_id", "service_class"},
		[]string{"connection_group", "departure_country_id", "arrival_country_id"},
		[]string{"connection_group", "departure_country_id", "service_class"},
		[]string{"connection_group", "departure_country_id"},
		[]string{"connection_group", "arrival_country_id", "service_class"},
		[]string{"connection_group", "arrival_country_id"},
		[]string{"connection_group", "service_class"},
		[]string{"connection_group"},
	}
}

func (sr SearchRequest) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (sr SearchRequest) getStrategyKeys() []string {
	return []string{
		"connection_group",
		"arrival_country_id",
		"departure_country_id",
		"arrival_city_id",
		"departure_city_id",
		"service_class",
	}
}

func (sr SearchRequest) getTableName() string {
	return "rm_frule_search_request"
}

func (sr SearchRequest) GetDefaultValue() interface{} {
	return false
}

func (sr SearchRequest) GetDataStorage() (map[int][]FRuler, error) {
	result := make(map[int][]FRuler)
	for rank, fieldList := range sr.GetComparisonOrder() {
		query := sr.db.Table(sr.getTableName())
		for _, field := range sr.getStrategyKeys() {
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
			var rowData SearchRequest

			if err := sr.db.ScanRows(rows, &rowData); err != nil {
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

func (sr SearchRequest) GetLastUpdateTime() time.Time {
	return getLastUpdateTime("search_request", sr.db)
}
