package frule_module

import (
	"stash.tutu.ru/golang/resources/db"
	"strconv"
	"strings"
	"time"
)

type InterlineRule struct {
	Id                   int     `gorm:"column:id"`
	Partner              *string `gorm:"column:partner"`
	ConnectionGroup      *string `gorm:"column:connection_group"`
	CarrierPlating       *int64  `gorm:"column:carrier_plating"`
	PureInterline        *bool   `gorm:"column:pure_interline"`
	CarriersForbid       string  `gorm:"column:carriers_forbid"`
	CarriersForbidParsed []int64
	CarriersNeed         string `gorm:"column:carrier_need"`
	CarriersNeedParsed   []int64
	Carriers             []int64
	Result               bool `gorm:"column:result"`
	db                   *db.Database
}

func NewInterlineFRule(db *db.Database) InterlineRule {
	return InterlineRule{
		db: db,
	}
}

func (ir InterlineRule) GetResultValue(testRule interface{}) interface{} {
	params := testRule.(InterlineRule)

	if len(params.Carriers) > 0 {
		if len(ir.CarriersNeedParsed) > 0 {
			s := true
			for _, carrierId := range ir.CarriersNeedParsed {
				if !inSliceInt64(carrierId, params.Carriers) {
					s = false
				}
			}
			if !s {
				return false
			}
		}

		if len(ir.CarriersForbidParsed) > 0 {
			for _, carrierId := range ir.CarriersForbidParsed {
				if inSliceInt64(carrierId, params.Carriers) {
					return false
				}
			}
		}
	}

	return ir.Result
}

func (ir InterlineRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"partner", "connection_group", "carrier_plating", "pure_interline"},
		[]string{"partner", "connection_group", "carrier_plating"},
		[]string{"partner", "connection_group", "pure_interline"},
		[]string{"partner", "connection_group"},
		[]string{"partner", "carrier_plating", "pure_interline"},
		[]string{"partner", "carrier_plating"},
		[]string{"partner", "pure_interline"},
		[]string{"partner"},
	}
}

func (ir InterlineRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (ir InterlineRule) getStrategyKeys() []string {
	return []string{"partner", "connection_group", "carrier_plating", "pure_interline"}
}

func (ir InterlineRule) getTableName() string {
	return "rm_frule_interline"
}

func (ir InterlineRule) GetDefaultValue() interface{} {
	return false
}

func (ir InterlineRule) GetLastUpdateTime() time.Time {
	return getLastUpdateTime("interline", ir.db)
}

func (ir InterlineRule) GetDataStorage() (map[int][]FRuler, error) {
	result := make(map[int][]FRuler)

	for rank, fieldList := range ir.GetComparisonOrder() {
		query := ir.db.Table(ir.getTableName())
		for _, field := range ir.getStrategyKeys() {
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
			var rowData InterlineRule

			if err := ir.db.ScanRows(rows, &rowData); err != nil {
				return result, err
			}

			rowData.CarriersForbidParsed, err = ir.splitCarriersString(rowData.CarriersForbid)
			if err != nil {
				return result, err
			}

			rowData.CarriersNeedParsed, err = ir.splitCarriersString(rowData.CarriersNeed)
			if err != nil {
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

func (ir InterlineRule) splitCarriersString(carriersString string) ([]int64, error) {
	var carriersStringParsed []int64

	if carriersString != "" {
		for _, s := range strings.Split(carriersString, "|") {
			d, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, err
			}
			carriersStringParsed = append(carriersStringParsed, d)
		}
	}

	return carriersStringParsed, nil
}
