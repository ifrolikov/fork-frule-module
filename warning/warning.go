package warning

import (
	"context"
	"reflect"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"time"
)

type Warning struct {
	Type    *string `json:"type"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Name    *string `json:"name"`
}

type WarningRule struct {
	Id                 int     `json:"id"`
	CarrierId          *int64  `json:"carrier_id"`
	DepartureCountryId *int64  `json:"departure_country_id"`
	DepartureCityId    *int64  `json:"departure_city_id"`
	ArrivalCountryId   *int64  `json:"arrival_country_id"`
	ArrivalCityId      *int64  `json:"arrival_city_id"`
	DepartureDateFrom  *string `json:"departure_date_from"`
	DepartureDateTo    *string `json:"departure_date_to"`
	StartDate          *string `json:"start_date"`
	ParsedStartDate    *time.Time
	FinishDate         *string `json:"finish_date"`
	ParsedFinishDate   *time.Time
	ConnectionGroup    *string   `json:"connection_group"`
	Lang               *string   `json:"lang"`
	Result             []Warning `json:"result"`
	repo               *frule_module.Repository
}

type RuleResult struct {
	Id     int
	Result []Warning
}

func NewWarningFRule(ctx context.Context, config *repository.Config) (*WarningRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &WarningRule{repo: repo}, nil
}

func (rule *WarningRule) GetResultValue(testRule interface{}) interface{} {
	if rule.isActual(*testRule.(WarningRule).ParsedStartDate, rule.ParsedStartDate, rule.ParsedFinishDate) {
		return &RuleResult{Id: rule.Id, Result: rule.Result}
	} else {
		return nil
	}
}

func (rule *WarningRule) isActual(departureDate time.Time, startDate *time.Time, finishDate *time.Time) bool {
	if startDate == nil && finishDate == nil {
		return true
	} else if startDate == nil && finishDate != nil {
		return departureDate.Before(*finishDate)
	} else if startDate != nil && finishDate == nil {
		return departureDate.After(*startDate) || departureDate.Equal(*startDate)
	} else {
		return (departureDate.After(*startDate) || departureDate.Equal(*startDate)) && departureDate.Before(*finishDate)
	}
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"lang", "carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"lang", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"lang", "carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"lang", "carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"lang", "carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"lang", "carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"lang", "carrier_id", "departure_date_from", "departure_date_to"},
	[]string{"lang", "carrier_id", "departure_date_to"},
	[]string{"lang", "carrier_id", "departure_date_from"},
	[]string{"lang", "departure_country_id", "arrival_country_id", "departure_city_id"},
	[]string{"lang", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"lang", "connection_group", "carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"lang", "carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"lang", "connection_group", "departure_country_id", "arrival_country_id"},
	[]string{"lang", "departure_country_id", "arrival_country_id"},
	[]string{"lang", "departure_country_id", "departure_city_id"},
	[]string{"lang", "arrival_country_id", "arrival_city_id"},
	[]string{"lang", "connection_group", "carrier_id", "departure_country_id"},
	[]string{"lang", "connection_group", "carrier_id", "arrival_country_id"},
	[]string{"lang", "carrier_id", "departure_country_id"},
	[]string{"lang", "carrier_id", "arrival_country_id"},
	[]string{"lang", "connection_group", "departure_country_id"},
	[]string{"lang", "connection_group", "arrival_country_id"},
	[]string{"lang", "departure_country_id"},
	[]string{"lang", "arrival_country_id"},
	[]string{"lang", "connection_group", "carrier_id"},
	[]string{"lang", "carrier_id"},
	[]string{"lang", "connection_group"},
}

func (rule *WarningRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{
	{
		Field: "departure_date_from",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) < b.Elem().Interface().(string)
		},
	},
	{
		Field: "departure_date_to",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) > b.Elem().Interface().(string)
		},
	},
}

func (rule *WarningRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{
	"lang",
	"carrier_id",
	"departure_date_from",
	"departure_date_to",
	"departure_country_id",
	"departure_city_id",
	"arrival_country_id",
	"arrival_city_id",
	"connection_group",
}

func (rule *WarningRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *WarningRule) GetDefaultValue() interface{} {
	return nil
}

func (rule *WarningRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *WarningRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *WarningRule) GetRuleName() string {
	return "Warning"
}
