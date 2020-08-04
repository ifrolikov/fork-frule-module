package manual_exchange_refund

import (
	"context"
	"github.com/rs/zerolog"
	"reflect"
	"regexp"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
)

type ManualExchangeRefundRule struct {
	Id                       int64                    `json:"id"`
	ServiceClass             *string                  `json:"service_class"`
	CarrierId                *int64                   `json:"carrier_id"`
	Fare                     *string                  `json:"fare"`
	HoursBeforeDeparture     *int64                   `json:"hours_before_departure"`
	PenaltyStrategy          *string                  `json:"penalty_strategy"`
	PassengerType            *string                  `json:"passenger_type"`
	IsTransit                *bool                    `json:"is_transit"`
	UsedType                 *string                  `json:"used_type"`
	DepartureCityId          *uint64                  `json:"departure_city_id"`
	ArrivalCityId            *uint64                  `json:"arrival_city_id"`
	FlightType               *string                  `json:"flight_type"`
	DepartureCountryId       *uint64                  `json:"departure_country_id"`
	ArrivalCountryId         *uint64                  `json:"arrival_country_id"`
	MaxExchangeCount         *int64                   `json:"max_exchange_count"`
	DaysAfterTariffStart     *int64                   `json:"days_after_tariff_start"`
	TariffStartType          *string                  `json:"tariff_start_type"`
	SegmentNumberInRoute     *int64                   `json:"segment_number_in_route"`
	SegmentNumberInItinerary *int64                   `json:"segment_number_in_itinerary"`
	Context                  *Context                 `json:"context"`
	IssueDateFrom            *string                  `json:"issue_date_from"`
	IssueDateTo              *string                  `json:"issue_date_to"`
	DepartureDateFrom        *string                  `json:"departure_date_from"`
	DepartureDateTo          *string                  `json:"departure_date_to"`
	Destination              *FeeDestination          `json:"destination"`
	ApplyStrategy            *ApplyStrategy           `json:"apply_strategy"`
	FarePercent              *float64                 `json:"fare_percent"`
	Amount                   *int64                   `json:"amount"`
	Currency                 *string                  `json:"currency"`
	FarePercentDestination   *FarePercentDestination  `json:"fare_percent_destination"`
	IsChangeable             *bool                    `json:"is_changeable"`
	IsRefundable             *bool                    `json:"is_refundable"`
	CalculationUnit          *CalculationUnit         `json:"calculation_unit"`
	Brand                    *string                  `json:"brand"`
	TariffCalculateFor       *TariffCalculationSource `json:"tariff_calculate_for"`
	repo                     *frule_module.Repository
	comparisonOrderImporter  ComparisonOrderImporterInterface
	logger                   zerolog.Logger
}

func NewManualExchangeRefundFRule(
	ctx context.Context,
	repConfig *repository.Config,
	comparisonOrderImporter ComparisonOrderImporterInterface) (*ManualExchangeRefundRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: repConfig}}, )
	if err != nil {
		return nil, err
	}

	logger := log.Logger
	logger = logger.With().Str("context.type", "manual_exchange_refund_frule").Logger()

	return &ManualExchangeRefundRule{
		repo:                    repo,
		comparisonOrderImporter: comparisonOrderImporter,
		logger:                  logger,
	}, nil
}

func (rule *ManualExchangeRefundRule) GetResultValue(interface{}) interface{} {
	var isAvailable = true
	switch *rule.Context {
	case ContextExchange:
		isAvailable = *rule.IsChangeable
		break
	case ContextRefund:
		isAvailable = *rule.IsRefundable
	}

	return NewManualExchangeRefundResult(
		rule.Id,
		rule.Destination,
		rule.ApplyStrategy,
		rule.FarePercent,
		rule.FarePercentDestination,
		rule.Amount,
		rule.Currency,
		rule.CalculationUnit,
		rule.Brand,
		rule.TariffCalculateFor,
		isAvailable,
		rule.HoursBeforeDeparture,
	)
}

func (rule *ManualExchangeRefundRule) GetComparisonOrder() frule_module.ComparisonOrder {
	comparisonOrder, err := rule.comparisonOrderImporter.getComparisonOrder(rule.logger)
	if err != nil {
		return frule_module.ComparisonOrder{}
	} else {
		return comparisonOrder
	}
}

var comparisonOperators = frule_module.ComparisonOperators{
	{
		Field: "hours_before_departure",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(int64) <= b.Elem().Interface().(int64)
		},
	},
	{
		Field: "days_after_tariff_start",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(int64) <= b.Elem().Interface().(int64)
		},
	},
	{
		Field: "max_exchange_count",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(int64) >= b.Elem().Interface().(int64)
		},
	},
	{
		Field: "issue_date_from",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) <= b.Elem().Interface().(string)
		},
	},
	{
		Field: "issue_date_to",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) >= b.Elem().Interface().(string)
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
	{
		Field: "fare",
		Function: func(a, b reflect.Value) bool {
			r, err := regexp.Compile(a.Elem().Interface().(string))
			if err != nil {
				return false
			}
			return r.Match([]byte(b.Elem().Interface().(string)))
		},
	},
}

func (*ManualExchangeRefundRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{
	"service_class",
	"carrier_id",
	"fare",
	"hours_before_departure",
	"penalty_strategy",
	"passenger_type",
	"is_transit",
	"used_type",
	"departure_city_id",
	"arrival_city_id",
	"flight_type",
	"departure_country_id",
	"arrival_country_id",
	"max_exchange_count",
	"days_after_tariff_start",
	"tariff_start_type",
	"segment_number_in_route",
	"segment_number_in_itinerary",
	"context",
	"issue_date_from",
	"issue_date_to",
	"departure_date_from",
	"departure_date_to",
}

func (*ManualExchangeRefundRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (*ManualExchangeRefundRule) GetDefaultValue() interface{} {
	return ManualExchangeRefundResult{}
}

func (rule *ManualExchangeRefundRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *ManualExchangeRefundRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (*ManualExchangeRefundRule) GetRuleName() string {
	return "ManualExchangeRefundRule"
}
