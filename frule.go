package frule_module

import (
	"context"
	"reflect"
	"sort"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
	"sync"
	"time"
)

type ComparisonOrder [][]string

type ComparisonFunction func(a, b reflect.Value) bool

type ComparisonOperators []ComparisonOperator

type ComparisonOperator struct {
	Field    string
	Function ComparisonFunction
}

type FRuler interface {
	GetResultValue(interface{}) interface{}
	GetComparisonOrder() ComparisonOrder
	GetComparisonOperators() ComparisonOperators
	GetStrategyKeys() []string
	GetDefaultValue() interface{}
	GetDataStorage() *RankedFRuleStorage
	GetNotificationChannel() chan repository.Notification
	GetRuleName() string
	GetCompareDynamicFieldsFunction() *CompareDynamicFieldsFunction
	GetCreateRuleHashForIndexedFieldsFunction() *CreateRuleHashForIndexedFieldsFunction
}

type CompareDynamicFieldsFunction func(testRule interface{}, foundRuleSet []FRuler) interface{}
type CreateRuleHashForIndexedFieldsFunction func(fields []string, rowSet interface{}) string

type FRule struct {
	index            map[int]map[string][]FRuler
	registry         map[string]map[int]int
	primaryKeys      []string
	indexedKeys      []string
	ruleSpecificData FRuler
	mutex            sync.Mutex
	lastUpdateTime   time.Time
}

func NewFRule(ctx context.Context, ruleSpecificData FRuler) *FRule {
	definition := FRule{
		index:            make(map[int]map[string][]FRuler),
		registry:         make(map[string]map[int]int),
		ruleSpecificData: ruleSpecificData,
	}

	var indexedKeys []string
	for _, field := range definition.ruleSpecificData.GetStrategyKeys() {
		fieldHasCustomComparisonOperator := false
		for _, comparisonOperator := range definition.ruleSpecificData.GetComparisonOperators() {
			if comparisonOperator.Field == field {
				fieldHasCustomComparisonOperator = true
				break
			}
		}
		if !fieldHasCustomComparisonOperator {
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
		log.Logger.Error().Stack().Err(err).Msg("Building index")
	}
	definition.lastUpdateTime = time.Now()

	go func(ctx context.Context, definition *FRule) {
		name := definition.ruleSpecificData.GetRuleName()
		for {
			select {
			case n := <-definition.ruleSpecificData.GetNotificationChannel():
				if n.Err != nil {
					log.Logger.Err(n.Err).Msgf("Error during FRule %s update: %s", name, n.Msg)
				} else {
					log.Logger.Info().Msgf("FRule %s update: %s", name, n.Msg)
					if n.Type == repository.NOTIFICATION_TYPE_UPDATED {
						if indexUpdateErr := definition.buildIndex(); indexUpdateErr != nil {
							log.Logger.Error().Stack().Err(indexUpdateErr).Msgf("FRule %s index update err: %v", name, indexUpdateErr)
						} else {
							log.Logger.Info().Msgf("FRule %s index updated", name)
						}
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx, &definition)

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
		case *int64, *int32:
			hashPart = strconv.FormatInt(fieldValue.Elem().Int(), 10)
		case *uint64, *uint32:
			hashPart = strconv.FormatUint(fieldValue.Elem().Uint(), 10)
		case *string:
			hashPart = fieldValue.Elem().String()
		case *bool:
			hashPart = strconv.FormatBool(fieldValue.Elem().Bool())
		}
		result += hashField + "=>" + hashPart + "|"
	}

	return result
}

func (f *FRule) buildIndex() error {
	rulesSets := f.ruleSpecificData.GetDataStorage()
	index := make(map[int]map[string][]FRuler)
	registry := make(map[string]map[int]int)

	for rank, rulesData := range *rulesSets {
		for _, rowData := range rulesData {
			var indexHash string
			rankIndexedKeys := intersectSlices(f.indexedKeys, f.ruleSpecificData.GetComparisonOrder()[rank])
			customCreateRuleHashFunc := f.ruleSpecificData.GetCreateRuleHashForIndexedFieldsFunction()
			if customCreateRuleHashFunc != nil {
				function := *customCreateRuleHashFunc
				indexHash = function(rankIndexedKeys, rowData)
			} else {
				indexHash = f.createRuleHash(rankIndexedKeys, rowData)
			}
			// indexHash будет пустой строкой в случае пустого rankIndexedKeys
			// это происходит, когда на некоторых уровнях (rank) GetComparisonOrder нет ни одного f.indexedKeys
			// все записи с этого уровня попадут в index[rank][""]
			if index[rank] == nil {
				index[rank] = make(map[string][]FRuler)
			}
			index[rank][indexHash] = append(index[rank][indexHash], rowData)

			var registryHash string
			if customCreateRuleHashFunc != nil {
				function := *customCreateRuleHashFunc
				registryHash = function(f.primaryKeys, rowData)
			} else {
				registryHash = f.createRuleHash(f.primaryKeys, rowData)
			}
			// registryHash будет пустой строкой в случае пустого f.primaryKeys
			// это происходит, когда нет полей, которые встречались бы на каждом уровне (rank) GetComparisonOrder
			// все rank попадут в один элемент регистра registry[""]
			if registry[registryHash] == nil {
				registry[registryHash] = make(map[int]int)
			}
			registry[registryHash][rank] = rank
		}
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.index = index
	f.registry = registry
	return nil
}

func (f *FRule) findRanks(testRule interface{}) []int {
	var result []int
	var registryHash string
	customCreateRuleHashFunc := f.ruleSpecificData.GetCreateRuleHashForIndexedFieldsFunction()
	if customCreateRuleHashFunc != nil {
		function := *customCreateRuleHashFunc
		registryHash = function(f.primaryKeys, testRule)
	} else {
		registryHash = f.createRuleHash(f.primaryKeys, testRule)
	}

	if indexes, ok := f.registry[registryHash]; ok {
		for _, rank := range indexes {
			result = append(result, rank)
		}
	}
	sort.Ints(result)
	return result
}

func (f *FRule) GetResult(testRule interface{}) interface{} {
	comparisonOperators := f.ruleSpecificData.GetComparisonOperators()
	for _, rank := range f.findRanks(testRule) {
		hashFields := intersectSlices(f.indexedKeys, f.ruleSpecificData.GetComparisonOrder()[rank])
		var ruleHash = ""
		if customCreateHashFunction := f.ruleSpecificData.GetCreateRuleHashForIndexedFieldsFunction(); customCreateHashFunction != nil {
			function := *customCreateHashFunction
			ruleHash = function(hashFields, testRule)
		} else {
			ruleHash = f.createRuleHash(hashFields, testRule)
		}
		if foundRuleSet, ok := f.index[rank][ruleHash]; ok {
			if comparisonFunction := f.ruleSpecificData.GetCompareDynamicFieldsFunction(); comparisonFunction != nil {
				function := *comparisonFunction
				return function(testRule, foundRuleSet)
			} else {
				return f.compareDynamicFields(
					testRule,
					foundRuleSet,
					comparisonOperators,
					f.ruleSpecificData.GetDefaultValue(),
				)
			}
		}
	}
	return f.ruleSpecificData.GetDefaultValue()
}

func (f *FRule) compareDynamicFields(testRule interface{}, foundRuleSet []FRuler, comparisonOperators ComparisonOperators, defaultValue interface{}) interface{} {
RULESET:
	for _, foundRule := range foundRuleSet {
		for _, comparisonOperator := range comparisonOperators {
			a := getFieldValueByTag(foundRule, comparisonOperator.Field)
			if !a.IsNil() {
				b := getFieldValueByTag(testRule, comparisonOperator.Field)
				if b.IsNil() || !comparisonOperator.Function(a, b) {
					continue RULESET
				}
			}
		}
		return foundRule.GetResultValue(testRule)
	}
	return defaultValue
}
