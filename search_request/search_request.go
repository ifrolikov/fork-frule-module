package search_request

import (
	"context"
	"encoding/json"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
	"time"
)

type SearchRequestRule struct {
	Id                 int     `json:"id"`
	ConnectionGroup    *string `json:"connection_group"`
	DepartureCityId    *uint64 `json:"departure_city_id"`
	ArrivalCityId      *uint64 `json:"arrival_city_id"`
	DepartureCountryId *uint64 `json:"departure_country_id"`
	ArrivalCountryId   *uint64 `json:"arrival_country_id"`
	ServiceClass       *string `json:"service_class"`
	Result             string  `json:"result"`
	ResultParsed       []frule_module.CronStructBool
	repo               *frule_module.Repository
}

func NewSearchRequestFRule(ctx context.Context, config *repository.Config) (*SearchRequestRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &SearchRequestRule{repo: repo}, nil
}

func (rule *SearchRequestRule) GetResultValue(interface{}) interface{} {
	for i := range rule.ResultParsed {
		if frule_module.CronSpec(&rule.ResultParsed[i].Spec, time.Now()) {
			return rule.ResultParsed[i].Value
		}
	}
	return false
}

func (rule *SearchRequestRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

var comparisonOrder = frule_module.ComparisonOrder{
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

func (rule *SearchRequestRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comaprisonOperators = frule_module.ComparisonOperators{}

func (rule *SearchRequestRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comaprisonOperators
}

var strategyKeys = []string{
	"connection_group",
	"arrival_country_id",
	"departure_country_id",
	"arrival_city_id",
	"departure_city_id",
	"service_class",
}

func (rule SearchRequestRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *SearchRequestRule) GetDefaultValue() interface{} {
	return false
}

func (rule *SearchRequestRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

/*
func (rule SearchRequestRule) GetDataStorage() (map[int][]frule_module.FRuler, error) {
	result := make(map[int][]frule_module.FRuler)
	repo := createRepository(rule.config)
	for rank, ruleList := range repo.GetStorage() {
		for _, ruleItem := range ruleList {
			var unserialized map[interface{}]interface{}

			err := phpserialize.Unmarshal([]byte(ruleItem.Result), &unserialized)
			if err != nil {
				return nil, err
			}

			var resultParsed []frule_module.CronStructBool

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

				resultParsed = append(resultParsed, frule_module.CronStructBool{Spec: key.(string), Value: val})
			}
			ruleItem.ResultParsed = resultParsed
			result[rank] = append(result[rank], ruleItem)
		}
	}
	return result, nil
}*/

func (rule *SearchRequestRule) parseCronSpecField(value string) []frule_module.CronStructBool {
	var resultParsed []frule_module.CronStructString

	err := json.Unmarshal([]byte(value), &resultParsed)

	if err != nil {
		log.Logger.Error().Stack().Err(err).Msg("Unmarshal")
	}

	result := make([]frule_module.CronStructBool, 0, len(resultParsed))
	for _, item := range resultParsed {
		val, _ := strconv.ParseBool(item.Value)
		result = append(result, frule_module.CronStructBool{Spec: item.Spec, Value: val})
	}
	return result
}

func (rule *SearchRequestRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *SearchRequestRule) GetRuleName() string {
	return "SearchRequest"
}
