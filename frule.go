package frule_module

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

type ComparisonOrder [][]string

type ComparisonOperators map[string]string

type FRuler interface {
	GetResultValue() interface{}
	GetContainer() FRuler
	GetComparisonOrder() ComparisonOrder
	GetComparisonOperators() ComparisonOperators
	GetStrategyKeys() []string
	GetIndexedKeys() []string
	GetTableName() string
	GetDefaultValue() interface{}
}

type FRule struct {
	index            map[int]map[string][]FRuler
	registry         map[string]map[int]int
	db               *gorm.DB
	primaryKeys      []string
	ruleSpecificData FRuler
}

func NewFRule(db *gorm.DB, ruleSpecificData FRuler) *FRule {
	definition := FRule{
		db:               db,
		index:            make(map[int]map[string][]FRuler),
		registry:         make(map[string]map[int]int),
		ruleSpecificData: ruleSpecificData,
	}

	var primaryKeys = definition.ruleSpecificData.GetIndexedKeys()

	for _, fields := range definition.ruleSpecificData.GetComparisonOrder() {
		primaryKeys = intersectSlices(primaryKeys, fields)
	}
	definition.primaryKeys = primaryKeys

	if err := definition.buildIndex(); err != nil {
		fmt.Println(err)
	}
	return &definition
}

func (f *FRule) createRuleHash(hashFields []string, rule interface{}) string {
	var result string

	for _, hashField := range hashFields {
		fieldValue := getFieldValueByTag(rule, hashField)
		if fieldValue.IsNil() {
			continue
		}
		var hashPart string
		switch fieldValue.Interface().(type) {
		case *int:
			hashPart = strconv.Itoa(int(fieldValue.Elem().Int()))
		case *string:
			hashPart = fieldValue.Elem().String()
		}
		result += hashField + "=>" + hashPart + "|"
	}
	return result
}

func (f *FRule) buildIndex() error {
	for rank, fieldList := range f.ruleSpecificData.GetComparisonOrder() {
		query := f.db.Table(f.ruleSpecificData.GetTableName())
		for _, field := range f.ruleSpecificData.GetStrategyKeys() {
			if inSlice(field, fieldList) {
				query = query.Where(field + " IS NOT NULL")
			} else {
				query = query.Where(field + " IS NULL")
			}
		}
		rows, err := query.Rows()
		if err != nil {
			return err
		}

		for rows.Next() {
			var rowData = f.ruleSpecificData.GetContainer()

			if err := f.db.ScanRows(rows, &rowData); err != nil {
				fmt.Println(err)
				return err
			}
			hashFields := intersectSlices(f.ruleSpecificData.GetIndexedKeys(), f.ruleSpecificData.GetComparisonOrder()[rank])
			hash := f.createRuleHash(hashFields, rowData)
			if hash != "" {
				if f.index[rank] == nil {
					f.index[rank] = make(map[string][]FRuler)
				}
				f.index[rank][hash] = append(f.index[rank][hash], rowData)
			}
			registryHash := f.createRuleHash(f.primaryKeys, rowData)
			if registryHash != "" {
				if f.registry[registryHash] == nil {
					f.registry[registryHash] = make(map[int]int)
				}
				f.registry[registryHash][rank] = rank
			}

		}
		err = rows.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (f *FRule) findRanks(testRule interface{}) []int {
	var result []int
	registryHash := f.createRuleHash(f.primaryKeys, testRule)
	if indexes, ok := f.registry[registryHash]; ok {
		for _, rank := range indexes {
			result = append(result, rank)
		}
	}
	return result
}

func (f *FRule) GetResult(testRule interface{}) interface{} {
	for _, rank := range f.findRanks(testRule) {
		hashFields := intersectSlices(f.ruleSpecificData.GetIndexedKeys(), f.ruleSpecificData.GetComparisonOrder()[rank])
		if foundRuleSet, ok := f.index[rank][f.createRuleHash(hashFields, testRule)]; ok {
			if len(foundRuleSet) == 1 {
				return foundRuleSet[0].GetResultValue()
			}
		}
	}
	return f.ruleSpecificData.GetDefaultValue()
}
