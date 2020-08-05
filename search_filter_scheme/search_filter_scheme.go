package search_filter_scheme

import (
	"context"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type SearchFilterSchemeRule struct {
	Id                 int     `json:"id"`
	DepartureCountryId *uint64 `json:"departure_country_id"`
	DepartureCityId    *uint64 `json:"departure_city_id"`
	ArrivalCountryId   *uint64 `json:"arrival_country_id"`
	ArrivalCityId      *uint64 `json:"arrival_city_id"`
	JourneyType        *string `json:"journey_type"`
	Result             *string `json:"result"`
	repo               *frule_module.Repository
}

type RuleResult struct {
	Id     int
	Result *string
}

func NewSearchFilterSchemeFRule(ctx context.Context, config *repository.Config) (*SearchFilterSchemeRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &SearchFilterSchemeRule{repo: repo}, nil
}

func (rule *SearchFilterSchemeRule) GetResultValue(testRule interface{}) interface{} {
	return &RuleResult{Id: rule.Id, Result: rule.Result}
}

func (rule *SearchFilterSchemeRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

func (rule *SearchFilterSchemeRule) GetCreateRuleHashForIndexedFieldsFunction() *frule_module.CreateRuleHashForIndexedFieldsFunction {
	return nil
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"departure_city_id",    "arrival_city_id",    "journey_type"},
	[]string{"departure_city_id",    "arrival_country_id", "journey_type"},
	[]string{"departure_country_id", "arrival_city_id",    "journey_type"},
	[]string{"departure_country_id", "arrival_country_id", "journey_type"},
	[]string{"departure_city_id",    "journey_type"},
	[]string{"arrival_city_id",      "journey_type"},
	[]string{"departure_country_id", "journey_type"},
	[]string{"arrival_country_id",   "journey_type"},
	[]string{"journey_type"},
}

func (rule *SearchFilterSchemeRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *SearchFilterSchemeRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{"departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "journey_type"}

func (rule *SearchFilterSchemeRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *SearchFilterSchemeRule) GetDefaultValue() interface{} {
	return nil
}

func (rule *SearchFilterSchemeRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *SearchFilterSchemeRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *SearchFilterSchemeRule) GetRuleName() string {
	return "SearchFilterScheme"
}
