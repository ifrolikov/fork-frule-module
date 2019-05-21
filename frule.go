package frule_module

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

type ComparisonOrder [][]string

type ComparisonOperators map[string]string

type FRuler interface {
	GetResultValue() interface{}
	GetComparisonOrder() ComparisonOrder
	GetComparisonOperators() ComparisonOperators
	getStrategyKeys() []string
	GetIndexedKeys() []string
	getTableName() string
	GetDefaultValue() interface{}
	GetDataStorage() map[int][]FRuler
}

type FRule struct {
	index            map[int]map[string][]FRuler
	registry         map[string]map[int]int
	primaryKeys      []string
	ruleSpecificData FRuler
}

func NewFRule(ruleSpecificData FRuler) *FRule {
	definition := FRule{
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
	for rank, rulesData := range f.ruleSpecificData.GetDataStorage() {
		for _, rowData := range rulesData {
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
