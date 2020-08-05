package direction

import (
	"context"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type DirectionRule struct {
	Id                 int     `json:"id"`
	Partner            *string `json:"partner"`
	ConnectionGroup    *string `json:"connection_group"`
	CarrierId          *int64  `json:"carrier_id"`
	DepartureCountryId *uint64 `json:"departure_country_id"`
	DepartureCityId    *uint64 `json:"departure_city_id"`
	ArrivalCountryId   *uint64 `json:"arrival_country_id"`
	ArrivalCityId      *uint64 `json:"arrival_city_id"`
	Result             bool    `json:"result"`
	repo               *frule_module.Repository
}

func NewDirectionFRule(ctx context.Context, config *repository.Config) (*DirectionRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &DirectionRule{repo: repo}, nil
}

func (rule *DirectionRule) GetResultValue(testRule interface{}) interface{} {
	return rule.Result
}

func (rule *DirectionRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

func (rule *DirectionRule) GetCreateRuleHashForIndexedFieldsFunction() *frule_module.CreateRuleHashForIndexedFieldsFunction {
	return nil
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "carrier_id"},
	[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "carrier_id"},
	[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "arrival_city_id", "carrier_id"},
	[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "carrier_id"},
	[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id"},
	[]string{"partner", "connection_group", "departure_country_id", "departure_city_id", "carrier_id"},
	[]string{"partner", "connection_group", "departure_country_id", "departure_city_id"},
	[]string{"partner", "connection_group", "arrival_country_id", "arrival_city_id", "carrier_id"},
	[]string{"partner", "connection_group", "arrival_country_id", "arrival_city_id"},
	[]string{"partner", "connection_group", "departure_country_id", "carrier_id"},
	[]string{"partner", "connection_group", "departure_country_id"},
	[]string{"partner", "connection_group", "arrival_country_id", "carrier_id"},
	[]string{"partner", "connection_group", "arrival_country_id"},
	[]string{"partner", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "carrier_id"},
	[]string{"partner", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"partner", "departure_country_id", "departure_city_id", "arrival_country_id", "carrier_id"},
	[]string{"partner", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"partner", "departure_country_id", "arrival_country_id", "arrival_city_id", "carrier_id"},
	[]string{"partner", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"partner", "departure_country_id", "arrival_country_id", "carrier_id"},
	[]string{"partner", "departure_country_id", "arrival_country_id"},
	[]string{"partner", "departure_country_id", "departure_city_id", "carrier_id"},
	[]string{"partner", "departure_country_id", "departure_city_id"},
	[]string{"partner", "arrival_country_id", "arrival_city_id", "carrier_id"},
	[]string{"partner", "arrival_country_id", "arrival_city_id"},
	[]string{"partner", "departure_country_id", "carrier_id"},
	[]string{"partner", "departure_country_id"},
	[]string{"partner", "arrival_country_id", "carrier_id"},
	[]string{"partner", "arrival_country_id"},
	[]string{"partner", "connection_group", "carrier_id"},
	[]string{"partner", "connection_group"},
	[]string{"partner", "carrier_id"},
	[]string{"partner"},
}

func (rule *DirectionRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *DirectionRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{"partner", "connection_group", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "carrier_id"}

func (rule *DirectionRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *DirectionRule) GetDefaultValue() interface{} {
	return false
}

func (rule *DirectionRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *DirectionRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *DirectionRule) GetRuleName() string {
	return "Direction"
}
