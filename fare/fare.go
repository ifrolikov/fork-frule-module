package fare

import (
	"context"
	"reflect"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"strings"
)

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

func (rule *FareRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return frule_module.ComparisonOrder{
		[]string{"departure_city_id", "arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "arrival_city_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "arrival_country_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "arrival_city_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "arrival_country_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_city_id", "partner", "carrier_id", "fare_spec"},
		[]string{"arrival_city_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"arrival_city_id", "partner", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"departure_country_id", "partner", "carrier_id", "fare_spec"},
		[]string{"arrival_country_id", "partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"arrival_country_id", "partner", "carrier_id", "fare_spec"},
		[]string{"partner", "connection_group", "carrier_id", "fare_spec"},
		[]string{"partner", "carrier_id", "fare_spec"},
	}
}

func (rule *FareRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return frule_module.ComparisonOperators{
		"fare_spec": func(a, b reflect.Value) bool {
			return strings.Contains(
				b.Elem().Interface().(string),
				strings.Trim(a.Elem().Interface().(string), "%"))
		},
	}
}

func (rule *FareRule) GetStrategyKeys() []string {
	return []string{
		"partner",
		"connection_group",
		"carrier_id",
		"arrival_country_id",
		"departure_country_id",
		"arrival_city_id",
		"departure_city_id",
		"fare_spec",
	}
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
