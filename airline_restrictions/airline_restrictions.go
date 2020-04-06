package airline_restrictions

import (
"context"
"stash.tutu.ru/avia-search-common/frule-module"
"stash.tutu.ru/avia-search-common/repository"
)

type AirlineRestrictionsRule struct {
	Id                       int     `json:"id"`
	PurchaseDateFrom         string `json:"purchase_date_from"`
	PurchaseDateTo           string `json:"purchase_date_to"`
	PurchasePeriodFrom       int `json:"purchase_period_from"`
	PurchasePeriodTo         int `json:"purchase_period_to"`
	ValidatingCarrierId      *int64  `json:"validating_carrier_id"`
	MarketingCarrierId       *int64  `json:"marketing_carrier_id"`
	OperatingCarrierId       *int64  `json:"operating_carrier_id"`
	Partner                  *string `json:"partner"`
	Gds                      *string `json:"gds"`
	DepartureCountryId       *uint64 `json:"departure_country_id"`
	DepartureCityId          *uint64 `json:"departure_city_id"`
	ArrivalCountryId         *uint64 `json:"arrival_country_id"`
	ArrivalCityId            *uint64 `json:"arrival_city_id"`
	Result                   bool    `json:"result"`
	repo                     *frule_module.Repository
}

func NewAirlineRestrictionFRule(ctx context.Context, config *repository.Config) (*AirlineRestrictionsRule, error) {
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

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},


	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "departure_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "validating_carrier_id", "arrival_country_id"},


	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner","validating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "departure_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "operating_carrier_id", "arrival_country_id"},


	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "departure_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "marketing_carrier_id", "operating_carrier_id", "arrival_country_id"},


	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "departure_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "validating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "validating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "validating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "validating_carrier_id", "arrival_country_id"},


	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "departure_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "departure_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "operating_carrier_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "operating_carrier_id", "arrival_country_id"},


	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "departure_city_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "departure_city_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "departure_city_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "departure_city_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "arrival_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "departure_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "departure_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "departure_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "arrival_country_id", "arrival_city_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "arrival_country_id", "arrival_city_id"},
	[]string{"purchase_date_from", "purchase_date_to", "arrival_country_id", "arrival_city_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "departure_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "departure_country_id"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "arrival_country_id", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "arrival_country_id", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "arrival_country_id"},
	[]string{"purchase_date_from", "purchase_date_to", "arrival_country_id"},


	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "partner", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "gds", "purchase_period_from", "purchase_period_to"},
	[]string{"purchase_date_from", "purchase_date_to", "purchase_period_from", "purchase_period_to"},

	[]string{"purchase_date_from", "purchase_date_to", "partner", "gds"},
	[]string{"purchase_date_from", "purchase_date_to", "partner"},
	[]string{"purchase_date_from", "purchase_date_to", "gds"},
	[]string{"purchase_date_from", "purchase_date_to"},
}

func (rule *AirlineRestrictionsRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *AirlineRestrictionsRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{
	"purchase_date_from",
	"purchase_date_to",
	"partner",
	"gds",
	"validating_carrier_id",
	"marketing_carrier_id",
	"operating_carrier_id",
	"departure_country_id",
	"departure_city_id",
	"arrival_country_id",
	"arrival_city_id",
	"purchase_period_from",
	"purchase_period_to"}

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
	return "AirlineRestriction"
}

