package partner_percent

import (
	"context"
	"reflect"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type PartnerPercentRule struct {
	Id                 int     `json:"id"`
	CarrierId          *int64  `json:"carrier_id"`
	Partner            *string `json:"partner"`
	ConnectionGroup    *string `json:"connection_group"`
	DateOfPurchaseFrom *string `json:"date_of_purchase_from"`
	DateOfPurchaseTo   *string `json:"date_of_purchase_to"`
	CarrierCountryId   *int64  `json:"carrier_country_id"`
	FareType           *string `json:"fare_type"`
	Result             float64 `json:"result"`
	repo               *frule_module.Repository
}

func NewPartnerPercentFRule(ctx context.Context, config *repository.Config) (*PartnerPercentRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &PartnerPercentRule{repo: repo}, nil
}

func (rule *PartnerPercentRule) GetResultValue(testRule interface{}) interface{} {
	return float64(rule.Result)
}

func (rule *PartnerPercentRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return frule_module.ComparisonOrder{
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_id", "fare_type", "connection_group"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_id", "fare_type"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_id", "connection_group"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_id"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_country_id", "connection_group"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "carrier_country_id"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to", "connection_group"},
		[]string{"partner", "date_of_purchase_from", "date_of_purchase_to"},
	}
}

func (rule *PartnerPercentRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return frule_module.ComparisonOperators{
		"date_of_purchase_from": func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) <= b.Elem().Interface().(string)
		},
		"date_of_purchase_to": func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) > b.Elem().Interface().(string)
		},
	}
}

func (rule *PartnerPercentRule) GetStrategyKeys() []string {
	return []string{"partner", "date_of_purchase_from", "date_of_purchase_to", "connection_group", "carrier_country_id",
		"carrier_id", "fare_type"}
}

func (rule *PartnerPercentRule) GetDefaultValue() interface{} {
	return 0.0
}

func (rule *PartnerPercentRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *PartnerPercentRule) GetNotificationChannel() chan error {
	return rule.repo.NotificationChannel
}