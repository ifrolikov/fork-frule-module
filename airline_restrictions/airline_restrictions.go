package airline_restrictions

import (
	"context"
	"reflect"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type AirlineRestrictionsRule struct {
	Id                  int     `json:"id"`
	DepartureDateFrom   *string `json:"departure_date_from"`
	DepartureDateTo     *string `json:"departure_date_to"`
	DeparturePeriodFrom *int64  `json:"departure_period_from"`
	DeparturePeriodTo   *int64  `json:"departure_period_to"`
	PlatingCarrierId    *int64  `json:"plating_carrier_id"`
	MarketingCarrierId  *int64  `json:"marketing_carrier_id"`
	OperatingCarrierId  *int64  `json:"operating_carrier_id"`
	Partner             *string `json:"partner"`
	Gds                 *string `json:"gds"`
	DepartureCountryId  *uint64 `json:"departure_country_id"`
	DepartureCityId     *uint64 `json:"departure_city_id"`
	ArrivalCountryId    *uint64 `json:"arrival_country_id"`
	ArrivalCityId       *uint64 `json:"arrival_city_id"`
	Result              bool    `json:"result"`
	repo                *frule_module.Repository
}

func NewAirlineRestrictionsFRule(ctx context.Context, config *repository.Config) (*AirlineRestrictionsRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &AirlineRestrictionsRule{repo: repo}, nil
}

func (rule *AirlineRestrictionsRule) GetResultValue(testRule interface{}) interface{} {
	return rule.Result
}

var comparisonOrder = frule_module.ComparisonOrder{
	//пирамида 576 строк это 8 блоков (разные сочетания по перевозчикам) по 72 строки (самый нижкий базовый блок)

	// все перевозчики 'plating_carrier_id', 'marketing_carrier_id', 'operating_carrier_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "operating_carrier_id"},

	//сочетание только 'plating_carrier_id', 'marketing_carrier_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "marketing_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "marketing_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "marketing_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "marketing_carrier_id"},

	//сочетание только 'plating_carrier_id', 'operating_carrier_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "operating_carrier_id"},

	//сочетание только 'marketing_carrier_id', 'operating_carrier_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "operating_carrier_id"},

	// только plating_carrier_id
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "plating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "plating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "plating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "plating_carrier_id"},

	// только operating_carrier_id
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "operating_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "operating_carrier_id"},

	//только marketing_carrier_id
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "arrival_country_id"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "marketing_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "marketing_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "marketing_carrier_id"},
	[]string{"departure_date_from", "departure_date_to", "marketing_carrier_id"},

	// без перевозчиков, 72 строки (базовая пирамида, к ней выше добавляются перевозчики и их сочетания)
	//добавляем сочетание 'departure_country_id', 'departure_city_id', 'arrival_country_id', 'arrival_city_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	//добавляем сочетание 'departure_country_id', 'departure_city_id', 'arrival_country_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "departure_city_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "departure_city_id", "arrival_country_id"},

	//добавляем сочетание 'departure_country_id', 'arrival_country_id', 'arrival_city_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	//добавляем сочетание 'departure_country_id', 'arrival_country_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "arrival_country_id"},

	//добавляем сочетание 'departure_country_id', 'departure_city_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "departure_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "departure_city_id"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "departure_city_id"},

	//добавляем сочетание 'arrival_country_id', 'arrival_city_id'
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "arrival_country_id", "arrival_city_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "arrival_country_id", "arrival_city_id"},
	[]string{"departure_date_from", "departure_date_to", "arrival_country_id", "arrival_city_id"},

	//добавляем только departure_country_id
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_country_id"},
	[]string{"departure_date_from", "departure_date_to", "departure_country_id"},

	//добавляем только arrival_country_id
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "arrival_country_id", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "arrival_country_id", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "partner", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "gds", "arrival_country_id"},
	[]string{"departure_date_from", "departure_date_to", "arrival_country_id"},

	//условный интуитивно понятный базовый блок правил 8 строк
	[]string{"departure_date_from", "departure_date_to", "partner", "gds", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "partner", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "gds", "departure_period_from", "departure_period_to"},
	[]string{"departure_date_from", "departure_date_to", "departure_period_from", "departure_period_to"},

	[]string{"departure_date_from", "departure_date_to", "partner", "gds"},
	[]string{"departure_date_from", "departure_date_to", "partner"},
	[]string{"departure_date_from", "departure_date_to", "gds"},
	[]string{"departure_date_from", "departure_date_to"},
}

func (rule *AirlineRestrictionsRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{
	"departure_date_from": func(a, b reflect.Value) bool {
		return a.Elem().Interface().(string) <= b.Elem().Interface().(string)
	},
	"departure_date_to": func(a, b reflect.Value) bool {
		return a.Elem().Interface().(string) > b.Elem().Interface().(string)
	},
	"departure_period_from": func(a, b reflect.Value) bool {
		return a.Elem().Interface().(int64) <= b.Elem().Interface().(int64)
	},
	"departure_period_to": func(a, b reflect.Value) bool {
		return a.Elem().Interface().(int64) > b.Elem().Interface().(int64)
	},
}

func (rule *AirlineRestrictionsRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{
	"departure_date_from",
	"departure_date_to",
	"partner",
	"gds",
	"plating_carrier_id",
	"marketing_carrier_id",
	"operating_carrier_id",
	"departure_country_id",
	"departure_city_id",
	"arrival_country_id",
	"arrival_city_id",
	"departure_period_from",
	"departure_period_to"}

func (rule *AirlineRestrictionsRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *AirlineRestrictionsRule) GetDefaultValue() interface{} {
	return false
}

func (rule *AirlineRestrictionsRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *AirlineRestrictionsRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *AirlineRestrictionsRule) GetRuleName() string {
	return "AirlineRestrictions"
}
