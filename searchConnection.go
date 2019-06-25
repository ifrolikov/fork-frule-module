package frule_module

import (
	"github.com/elliotchance/phpserialize"
	"regexp"
	"stash.tutu.ru/golang/resources/db"
	"strconv"
	"time"
)

type SearchConnectionRule struct {
	Id                     int     `gorm:"column:id"`
	Partner                *string `gorm:"column:partner"`
	ConnectionGroup        *string `gorm:"column:connection_group"`
	DepartureDate          time.Time
	MinDepartureDate       *string `gorm:"column:min_departure_date"`
	MinDepartureDateParsed []cronStrucString
	MaxDepartureDate       *string `gorm:"column:max_departure_date"`
	MaxDepartureDateParsed []cronStrucString
	db                     *db.Database
}

var specSearchConnectionRegexp = regexp.MustCompile(`\+(\d+)([dwmy])`)

func NewSearchConnectionFRule(db *db.Database) SearchConnectionRule {
	return SearchConnectionRule{
		db: db,
	}
}

func (sc SearchConnectionRule) GetResultValue(testRule interface{}) interface{} {
	nowLocal := time.Now()
	nowUtc := time.Now().UTC()

	if minDepartureDateBorder := sc.getSpecInterval(sc.MinDepartureDateParsed, nowLocal); minDepartureDateBorder != "" {
		if minDepartureDate := sc.getDateToCompare(minDepartureDateBorder, nowUtc); minDepartureDate != nil {
			if !testRule.(SearchConnectionRule).DepartureDate.After(*minDepartureDate) {
				return false
			}
		}

	}

	if maxDepartureDateBorder := sc.getSpecInterval(sc.MaxDepartureDateParsed, nowLocal); maxDepartureDateBorder != "" {
		if maxDepartureDate := sc.getDateToCompare(maxDepartureDateBorder, nowUtc); maxDepartureDate != nil {
			if !testRule.(SearchConnectionRule).DepartureDate.Before(*maxDepartureDate) {
				return false
			}
		}
	}

	return true
}

func (sc SearchConnectionRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"partner", "connection_group"},
	}
}

func (sc SearchConnectionRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (sc SearchConnectionRule) getStrategyKeys() []string {
	return []string{"partner", "connection_group"}
}

func (sc SearchConnectionRule) getTableName() string {
	return "rm_frule_search_connection"
}

func (sc SearchConnectionRule) GetDefaultValue() interface{} {
	return false
}

func (sc SearchConnectionRule) GetLastUpdateTime() time.Time {
	return getLastUpdateTime("search_connection", sc.db)
}

func (sc SearchConnectionRule) GetDataStorage() (map[int][]FRuler, error) {
	result := make(map[int][]FRuler)

	for rank, fieldList := range sc.GetComparisonOrder() {
		query := sc.db.Table(sc.getTableName())
		for _, field := range sc.getStrategyKeys() {
			if inSlice(field, fieldList) {
				query = query.Where(field + " IS NOT NULL")
			} else {
				query = query.Where(field + " IS NULL")
			}
		}
		rows, err := query.Rows()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var rowData SearchConnectionRule

			if err := sc.db.ScanRows(rows, &rowData); err != nil {
				return nil, err
			}

			rowData.MinDepartureDateParsed, err = sc.parseCronSpecField(*rowData.MinDepartureDate)
			if err != nil {
				return nil, err
			}

			rowData.MaxDepartureDateParsed, err = sc.parseCronSpecField(*rowData.MaxDepartureDate)
			if err != nil {
				return nil, err
			}

			result[rank] = append(result[rank], rowData)

		}
		err = rows.Close()
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (sc SearchConnectionRule) parseCronSpecField(value string) ([]cronStrucString, error) {
	var unserialized map[interface{}]interface{}

	err := phpserialize.Unmarshal([]byte(value), &unserialized)
	if err != nil {
		return nil, err
	}

	var resultParsed []cronStrucString

	for key, value := range unserialized {
		var val string
		switch value.(type) {
		case string:
			val = value.(string)
			if val == "0" {
				val = ""
			}
		case int64:
			val = ""
		}

		resultParsed = append(resultParsed, cronStrucString{key.(string), val})
	}

	return resultParsed, nil
}

func (sc SearchConnectionRule) getSpecInterval(specs []cronStrucString, t time.Time) string {
	for i := range specs {
		if cronSpec(&specs[i].spec, t) {
			return specs[i].value
		}
	}
	return ""
}

func (sc SearchConnectionRule) getDateToCompare(intervalString string, t time.Time) *time.Time {
	intervalParsed := specSearchConnectionRegexp.FindStringSubmatch(intervalString)

	if len(intervalParsed) >= 3 {
		intervalInt, err := strconv.Atoi(intervalParsed[1])
		if err != nil {
			return nil
		}

		var years, months, days int

		switch intervalParsed[2] {
		case "d":
			days = intervalInt
		case "w":
			days = intervalInt * 7
		case "m":
			months = intervalInt
		case "y":
			years = intervalInt
		default:
			return nil
		}

		minDepartureDate := t.AddDate(years, months, days)

		return &minDepartureDate
	}

	return nil
}
