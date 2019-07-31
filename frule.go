package frule_module

import (
	"context"
	"reflect"
	"sort"
	"stash.tutu.ru/golang/log"
	"strconv"
	"sync"
	"time"
)

type ComparisonOrder [][]string

type ComparisonFunction func(a, b reflect.Value) bool

type ComparisonOperators map[string]ComparisonFunction

type FRuler interface {
	GetResultValue(interface{}) interface{}
	GetComparisonOrder() ComparisonOrder
	GetComparisonOperators() ComparisonOperators
	GetStrategyKeys() []string
	GetDefaultValue() interface{}
	GetDataStorage() *RankedFRuleStorage
	GetNotificationChannel() chan error
}

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
	definition.lastUpdateTime = time.Now()

	go func(ctx context.Context, definition *FRule) {
		for {
			select {
			case err := <- definition.ruleSpecificData.GetNotificationChannel():
				log.Logger.Info().Msgf("repository update message received: %v", err)
				if err != nil {
					log.Logger.Error().Err(err).Msgf("repository update error: %v", err)
				} else {
					log.Logger.Info().Msg("start index update")
					if indexUpdateErr := definition.buildIndex(); indexUpdateErr != nil {
						log.Logger.Error().Err(err).Msgf("index update err: %v", err)
					} else {
						log.Logger.Info().Msg("index updated")
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
			hashFields := intersectSlices(f.indexedKeys, f.ruleSpecificData.GetComparisonOrder()[rank])
			hash := f.createRuleHash(hashFields, rowData)
			if hash != "" {
				if index[rank] == nil {
					index[rank] = make(map[string][]FRuler)
				}
				index[rank][hash] = append(index[rank][hash], rowData)
			}
			registryHash := f.createRuleHash(f.primaryKeys, rowData)
			if registryHash != "" {
				if registry[registryHash] == nil {
					registry[registryHash] = make(map[int]int)
				}
				registry[registryHash][rank] = rank
			}
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
	registryHash := f.createRuleHash(f.primaryKeys, testRule)
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
		if foundRuleSet, ok := f.index[rank][f.createRuleHash(hashFields, testRule)]; ok {
		RULESET:
			for _, foundRule := range foundRuleSet {
				for fieldName, comparisonFunc := range comparisonOperators {
					a := getFieldValueByTag(foundRule, fieldName)
					if !a.IsNil() {
						b := getFieldValueByTag(testRule, fieldName)
						if b.IsNil() || !comparisonFunc(a, b) {
							continue RULESET
						}
					}
				}
				return foundRule.GetResultValue(testRule)
			}
		}
	}
	return f.ruleSpecificData.GetDefaultValue()
}
