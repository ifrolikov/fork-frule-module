package search_scheme

import (
	"context"
	frule_module "github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type SearchSchemeRule struct {
	Id                 int     `json:"id"`
	ConnectionGroup    *string `json:"connection_group"`
	DepartureCityId    *uint64 `json:"departure_city_id"`
	DepartureRegionId  *uint64 `json:"departure_region_id"`
	DepartureCountryId *uint64 `json:"departure_country_id"`
	ArrivalCityId      *uint64 `json:"arrival_city_id"`
	ArrivalRegionId    *uint64 `json:"arrival_region_id"`
	ArrivalCountryId   *uint64 `json:"arrival_country_id"`
	Result             string  `json:"result"`
	ResultParsed       []string
	repo               *frule_module.Repository
}

func NewSearchSchemeFRule(ctx context.Context, config *repository.Config) (*SearchSchemeRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &SearchSchemeRule{repo: repo}, nil
}

func (rule *SearchSchemeRule) GetResultValue(interface{}) interface{} {
	return rule.ResultParsed
}

func (rule *SearchSchemeRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

func (rule *SearchSchemeRule) GetCreateRuleHashForIndexedFieldsFunction() *frule_module.CreateRuleHashForIndexedFieldsFunction {
	return nil
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"connection_group", "departure_city_id", "arrival_city_id"},
	[]string{"connection_group", "departure_city_id", "arrival_country_id"},
	[]string{"connection_group", "departure_country_id", "arrival_city_id"},
	[]string{"connection_group", "departure_country_id", "arrival_country_id"},
	[]string{"connection_group", "departure_region_id", "arrival_country_id"},
	[]string{"connection_group", "departure_country_id", "arrival_region_id"},
	[]string{"connection_group", "departure_region_id", "arrival_region_id"},
	[]string{"connection_group", "departure_region_id"},
	[]string{"connection_group", "departure_country_id"},
	[]string{"connection_group", "departure_city_id"},
	[]string{"connection_group", "arrival_region_id"},
	[]string{"connection_group", "arrival_country_id"},
	[]string{"connection_group", "arrival_city_id"},
	[]string{"connection_group"},
}

func (rule *SearchSchemeRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *SearchSchemeRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{
	"connection_group",
	"departure_region_id",
	"departure_country_id",
	"departure_city_id",
	"arrival_region_id",
	"arrival_country_id",
	"arrival_city_id",
}

func (rule SearchSchemeRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *SearchSchemeRule) GetDefaultValue() interface{} {
	return []string{}
}

func (rule *SearchSchemeRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *SearchSchemeRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *SearchSchemeRule) GetRuleName() string {
	return "SearchScheme"
}
