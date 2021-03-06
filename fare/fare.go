package fare

import (
	"context"
	"reflect"
	"regexp"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

//дефолтное значение позволяет работать по единому сценарию как в случае присутствия пользователя в какой-либо группе доступа,
//так и в том случае, если он не состоит ни в одной группе доступа
const DEFAULT_FARE_ACCESS_GROUP = "__default_fare_access_group__"

type FareRule struct {
	Id                 int     `json:"id"`
	Partner            *string `json:"partner"`
	ConnectionGroup    *string `json:"connection_group"`
	CarrierId          *int64  `json:"carrier_id"`
	DepartureCityId    *uint64 `json:"departure_city_id"`
	ArrivalCityId      *uint64 `json:"arrival_city_id"`
	DepartureCountryId *uint64 `json:"departure_country_id"`
	ArrivalCountryId   *uint64 `json:"arrival_country_id"`
	FareSpec           *string `json:"fare_spec"`
	FareAccessGroup    *string `json:"fare_access_group"`
	Result             string  `json:"result"`
	repo               *frule_module.Repository
}

func NewFareFRule(ctx context.Context, config *repository.Config) (*FareRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &FareRule{repo: repo}, nil
}

func (rule *FareRule) GetResultValue(interface{}) interface{} {
	return rule.Result
}

func (rule *FareRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

func (rule *FareRule) GetCreateRuleHashForIndexedFieldsFunction() *frule_module.CreateRuleHashForIndexedFieldsFunction {
	return nil
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"departure_city_id", "arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_city_id", "arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"departure_city_id", "arrival_city_id", "partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_city_id", "arrival_city_id", "partner", "carrier_id", "fare_spec"},
	[]string{"departure_city_id", "arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_city_id", "arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"departure_city_id", "arrival_country_id", "partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_city_id", "arrival_country_id", "partner", "carrier_id", "fare_spec"},
	[]string{"departure_country_id", "arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_country_id", "arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"departure_country_id", "arrival_city_id", "partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_country_id", "arrival_city_id", "partner", "carrier_id", "fare_spec"},
	[]string{"departure_country_id", "arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_country_id", "arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"departure_country_id", "arrival_country_id", "partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_country_id", "arrival_country_id", "partner", "carrier_id", "fare_spec"},
	[]string{"departure_city_id", "partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"departure_city_id", "partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_city_id", "partner", "carrier_id", "fare_spec"},
	[]string{"arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"arrival_city_id", "partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"arrival_city_id", "partner", "carrier_id", "fare_spec"},
	[]string{"departure_country_id", "partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"departure_country_id", "partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"departure_country_id", "partner", "carrier_id", "fare_spec"},
	[]string{"arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"arrival_country_id", "partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"arrival_country_id", "partner", "carrier_id", "fare_spec"},
	[]string{"partner", "connection_group", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"partner", "connection_group", "carrier_id", "fare_spec"},
	[]string{"partner", "carrier_id", "fare_spec", "fare_access_group"},
	[]string{"partner", "carrier_id", "fare_spec"},
}

func (rule *FareRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{
	{
		Field: "fare_spec",
		Function: func(a, b reflect.Value) bool {
			r, err := regexp.Compile(a.Elem().Interface().(string))
			if err != nil {
				return false
			}
			return r.Match([]byte(b.Elem().Interface().(string)))
		},
	},
}

func (rule *FareRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{
	"partner",
	"connection_group",
	"carrier_id",
	"arrival_country_id",
	"departure_country_id",
	"arrival_city_id",
	"departure_city_id",
	"fare_spec",
	"fare_access_group",
}

func (rule *FareRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *FareRule) GetDefaultValue() interface{} {
	return ""
}

func (rule *FareRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *FareRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *FareRule) GetRuleName() string {
	return "Fare"
}
