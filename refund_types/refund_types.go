package refund_types

import (
	"context"
	"github.com/rs/zerolog"
	"reflect"
	"stash.tutu.ru/avia-search-common/contracts/v2/gateSearch"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
)

type RefundTypesRule struct {
	Id                      int64   `json:"id"`
	PlatingCarrierId        *int64  `json:"plating_carrier_id"`
	IssueDateFrom           *string `json:"issue_date_from"`
	IssueDateTo             *string `json:"issue_date_to"`
	DepartureDateFrom       *string `json:"departure_date_from"`
	DepartureDateTo         *string `json:"departure_date_to"`
	DepartureCountryId      *uint64 `json:"departure_country_id"`
	ArrivalCountryId        *uint64 `json:"arrival_country_id"`
	DepartureCityId         *uint64 `json:"departure_city_id"`
	ArrivalCityId           *uint64 `json:"arrival_city_id"`
	RefundType              *string `json:"refund_type"`
	repo                    *frule_module.Repository
	comparisonOrderImporter ComparisonOrderImporterInterface
	logger                  zerolog.Logger
}

func NewRefundTypesFRule(
	ctx context.Context,
	repConfig *repository.Config,
	comparisonOrderImporter ComparisonOrderImporterInterface) (*RefundTypesRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: repConfig}}, )
	if err != nil {
		return nil, err
	}

	logger := log.Logger
	logger = logger.With().Str("context.type", "refund_types_frule").Logger()

	return &RefundTypesRule{
		repo:                    repo,
		comparisonOrderImporter: comparisonOrderImporter,
		logger:                  logger,
	}, nil
}

func (rule *RefundTypesRule) GetResultValue(interface{}) interface{} {
	return gateSearch.RefundType(gateSearch.RefundType_value[*rule.RefundType])
}

func (rule *RefundTypesRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

func (rule *RefundTypesRule) GetComparisonOrder() frule_module.ComparisonOrder {
	comparisonOrder, err := rule.comparisonOrderImporter.getComparisonOrder(rule.logger)
	if err != nil {
		return frule_module.ComparisonOrder{}
	} else {
		return comparisonOrder
	}
}

var comparisonOperators = frule_module.ComparisonOperators{
	{
		Field: "issue_date_from",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) <= b.Elem().Interface().(string)
		},
	},
	{
		Field: "issue_date_to",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) > b.Elem().Interface().(string)
		},
	},
	{
		Field: "departure_date_from",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) <= b.Elem().Interface().(string)
		},
	},
	{
		Field: "departure_date_to",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) > b.Elem().Interface().(string)
		},
	},
}

func (rule *RefundTypesRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{
	"plating_carrier_id",
	"issue_date_from",
	"issue_date_to",
	"departure_date_from",
	"departure_date_to",
	"departure_country_id",
	"arrival_country_id",
	"departure_city_id",
	"arrival_city_id",
}

func (rule *RefundTypesRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *RefundTypesRule) GetDefaultValue() interface{} {
	return gateSearch.RefundType_money
}

func (rule *RefundTypesRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *RefundTypesRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *RefundTypesRule) GetRuleName() string {
	return "RefundTypesRule"
}
