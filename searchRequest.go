package frule_module

import (
	"github.com/elliotchance/phpserialize"
	"stash.tutu.ru/golang/resources/db"
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
	ResultParsed       []cronStrucBool
	db                 *db.Database
}

func NewSearchRequest(database *db.Database) SearchRequest {
	return SearchRequest{
		db: database,
	}
}

func (sr SearchRequest) GetResultValue(interface{}) interface{} {
	for i := range sr.ResultParsed {
		if cronSpec(&sr.ResultParsed[i].spec, time.Now()) {
			return sr.ResultParsed[i].value
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

			var unserialized map[interface{}]interface{}

			err := phpserialize.Unmarshal([]byte(rowData.Result), &unserialized)
			if err != nil {
				return nil, err
			}

			var resultParsed []cronStrucBool

			for key, value := range unserialized {
				var val bool
				switch value.(type) {
				case string:
					if value.(string) == "1" {
						val = true
					} else {
						val = false
					}
				case int64:
					if value.(int64) == int64(1) {
						val = true
					} else {
						val = false
					}
				}

				resultParsed = append(resultParsed, cronStrucBool{key.(string), val})
			}
			rowData.ResultParsed = resultParsed
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
