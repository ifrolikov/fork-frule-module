package frule_module

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
	"stash.tutu.ru/golang/log"
	"strconv"
)

type ComparisonOrder [][]string

type ComparisonFunction func(a, b reflect.Value) bool

type ComparisonOperators map[string]ComparisonFunction

type FRuler interface {
	GetResultValue() interface{}
	GetComparisonOrder() ComparisonOrder
	GetComparisonOperators() ComparisonOperators
	getStrategyKeys() []string
	getTableName() string
	GetDefaultValue() interface{}
	GetDataStorage() (map[int][]FRuler, error)
}

type FRule struct {
	index            map[int]map[string][]FRuler
	registry         map[string]map[int]int
	primaryKeys      []string
	indexedKeys      []string
	ruleSpecificData FRuler
}

func NewFRule(ruleSpecificData FRuler) *FRule {
	definition := FRule{
		index:            make(map[int]map[string][]FRuler),
		registry:         make(map[string]map[int]int),
		ruleSpecificData: ruleSpecificData,
	}

	var indexedKeys []string
	for _, field := range definition.ruleSpecificData.getStrategyKeys() {
		if _, ok := definition.ruleSpecificData.GetComparisonOperators()[field]; !ok {
			indexedKeys = append(indexedKeys, field)
		}
	}
	definition.indexedKeys = indexedKeys

	var primaryKeys = indexedKeys

	for _, fields := range definition.ruleSpecificData.GetComparisonOrder() {
		primaryKeys = intersectSlices(primaryKeys, fields)
	}
	definition.primaryKeys = primaryKeys

	if err := definition.buildIndex(); err != nil {
		log.Logger.Error().Err(err).Msg("Building index")
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
	rulesSets, err := f.ruleSpecificData.GetDataStorage()
	if err != nil {
		return err
	}
	for rank, rulesData := range rulesSets {
		for _, rowData := range rulesData {
			hashFields := intersectSlices(f.indexedKeys, f.ruleSpecificData.GetComparisonOrder()[rank])
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
		hashFields := intersectSlices(f.indexedKeys, f.ruleSpecificData.GetComparisonOrder()[rank])
		if foundRuleSet, ok := f.index[rank][f.createRuleHash(hashFields, testRule)]; ok {
			if len(foundRuleSet) == 1 {
				return foundRuleSet[0].GetResultValue()
			} else {
			RULESET:
				for _, foundRule := range foundRuleSet {
					for fieldName, comparisonFunc := range f.ruleSpecificData.GetComparisonOperators() {
						if !comparisonFunc(getFieldValueByTag(foundRule, fieldName), getFieldValueByTag(testRule, fieldName)) {
							continue RULESET
						}
					}
					return foundRule.GetResultValue()
				}
			}
		}
	}
	return f.ruleSpecificData.GetDefaultValue()
}
