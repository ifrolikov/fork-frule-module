package manual_exchange_refund

import (
	"context"
	"github.com/ifrolikov/fork-frule-module"
	"github.com/rs/zerolog"
	"regexp"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
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

func (rule *ManualExchangeRefundRule) GetResultValue(resultRule interface{}) interface{} {
	frule := resultRule.(*ManualExchangeRefundRule)
	var isAvailable = true
	switch *frule.Context {
	case ContextExchange:
		isAvailable = *frule.IsChangeable
		break
	case ContextRefund:
		isAvailable = *frule.IsRefundable
	}

	return NewManualExchangeRefundResult(
		frule.Id,
		frule.Destination,
		frule.ApplyStrategy,
		frule.FarePercent,
		frule.FarePercentDestination,
		frule.Amount,
		frule.Currency,
		frule.CalculationUnit,
		frule.Brand,
		frule.TariffCalculateFor,
		isAvailable,
		frule.HoursBeforeDeparture,
	)
}

func (rule *ManualExchangeRefundRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	var result frule_module.CompareDynamicFieldsFunction = func(testRule interface{}, foundRuleSet []frule_module.FRuler) interface{} {
	RULESET:
		for _, foundRule := range foundRuleSet {
			frule := foundRule.(*ManualExchangeRefundRule)
			if frule == nil {
				continue RULESET
			}
			tRule := testRule.(*ManualExchangeRefundRule)
			if !rule.compareHoursBeforeDeparture(frule.HoursBeforeDeparture, tRule.HoursBeforeDeparture) ||
				!rule.compareDaysAfterTariffStart(frule.DaysAfterTariffStart, tRule.DaysAfterTariffStart) ||
				!rule.compareMaxExchangeCount(frule.MaxExchangeCount, tRule.MaxExchangeCount) ||
				!rule.compareIssueDateFrom(frule.IssueDateFrom, tRule.IssueDateFrom) ||
				!rule.compareIssueDateTo(frule.IssueDateTo, tRule.IssueDateTo) ||
				!rule.compareDepartureDateFrom(frule.DepartureDateFrom, tRule.DepartureDateFrom) ||
				!rule.compareDepartureDateTo(frule.DepartureDateTo, tRule.DepartureDateTo) ||
				!rule.compareFare(frule.Fare, tRule.Fare) {
				continue RULESET
			}
			return rule.GetResultValue(frule)
		}
		return rule.GetDefaultValue()
	}
	return &result
}

func (rule *ManualExchangeRefundRule) GetCreateRuleHashForIndexedFieldsFunction() *frule_module.CreateRuleHashForIndexedFieldsFunction {
	var result frule_module.CreateRuleHashForIndexedFieldsFunction = func(fields []string, rowSet interface{}) string {
		var inFields = func(field string, fields []string) bool {
			for _, fieldToCompare := range fields {
				if fieldToCompare == field {
					return true
				}
			}
			return false
		}
		var result = ""
		frule := rowSet.(*ManualExchangeRefundRule)
		if inFields("service_class", fields) && frule.ServiceClass != nil {
			result = result + "service_class" + "=>" + *frule.ServiceClass + "|"
		}
		if inFields("carrier_id", fields) && frule.CarrierId != nil {
			result = result + "carrier_id" + "=>" + strconv.FormatInt(*frule.CarrierId, 10) + "|"
		}
		if inFields("fare", fields) && frule.Fare != nil {
			result = result + "fare" + "=>" + *frule.Fare + "|"
		}
		if inFields("hours_before_departure", fields) && frule.HoursBeforeDeparture != nil {
			result = result + "hours_before_departure" + "=>" + strconv.FormatInt(*frule.HoursBeforeDeparture, 10) + "|"
		}
		if inFields("penalty_strategy", fields) && frule.PenaltyStrategy != nil {
			result = result + "penalty_strategy" + "=>" + *frule.PenaltyStrategy + "|"
		}
		if inFields("passenger_type", fields) && frule.PassengerType != nil {
			result = result + "passenger_type" + "=>" + *frule.PassengerType + "|"
		}
		if inFields("is_transit", fields) && frule.IsTransit != nil {
			result = result + "is_transit" + "=>" + strconv.FormatBool(*frule.IsTransit) + "|"
		}
		if inFields("used_type", fields) && frule.UsedType != nil {
			result = result + "used_type" + "=>" + *frule.UsedType + "|"
		}
		if inFields("departure_city_id", fields) && frule.DepartureCityId != nil {
			result = result + "departure_city_id" + "=>" + strconv.FormatUint(*frule.DepartureCityId, 10) + "|"
		}
		if inFields("arrival_city_id", fields) && frule.ArrivalCityId != nil {
			result = result + "arrival_city_id" + "=>" + strconv.FormatUint(*frule.ArrivalCityId, 10) + "|"
		}
		if inFields("flight_type", fields) && frule.FlightType != nil {
			result = result + "flight_type" + "=>" + *frule.FlightType + "|"
		}
		if inFields("departure_country_id", fields) && frule.DepartureCountryId != nil {
			result = result + "departure_country_id" + "=>" + strconv.FormatUint(*frule.DepartureCountryId, 10) + "|"
		}
		if inFields("arrival_country_id", fields) && frule.ArrivalCountryId != nil {
			result = result + "arrival_country_id" + "=>" + strconv.FormatUint(*frule.ArrivalCountryId, 10) + "|"
		}
		if inFields("max_exchange_count", fields) && frule.MaxExchangeCount != nil {
			result = result + "max_exchange_count" + "=>" + strconv.FormatInt(*frule.MaxExchangeCount, 10) + "|"
		}
		if inFields("days_after_tariff_start", fields) && frule.DaysAfterTariffStart != nil {
			result = result + "days_after_tariff_start" + "=>" + strconv.FormatInt(*frule.DaysAfterTariffStart, 10) + "|"
		}
		if inFields("tariff_start_type", fields) && frule.TariffStartType != nil {
			result = result + "tariff_start_type" + "=>" + *frule.TariffStartType + "|"
		}
		if inFields("segment_number_in_route", fields) && frule.SegmentNumberInRoute != nil {
			result = result + "segment_number_in_route" + "=>" + strconv.FormatInt(*frule.SegmentNumberInRoute, 10) + "|"
		}
		if inFields("segment_number_in_itinerary", fields) && frule.SegmentNumberInItinerary != nil {
			result = result + "segment_number_in_itinerary" + "=>" + strconv.FormatInt(*frule.SegmentNumberInItinerary, 10) + "|"
		}
		if inFields("context", fields) && frule.Context != nil {
			result = result + "context" + "=>" + string(*frule.Context) + "|"
		}
		if inFields("issue_date_from", fields) && frule.IssueDateFrom != nil {
			result = result + "issue_date_from" + "=>" + *frule.IssueDateFrom + "|"
		}
		if inFields("issue_date_to", fields) && frule.IssueDateTo != nil {
			result = result + "issue_date_to" + "=>" + *frule.IssueDateTo + "|"
		}
		if inFields("departure_date_from", fields) && frule.DepartureDateFrom != nil {
			result = result + "departure_date_from" + "=>" + *frule.DepartureDateFrom + "|"
		}
		if inFields("departure_date_to", fields) && frule.DepartureDateTo != nil {
			result = result + "departure_date_to" + "=>" + *frule.DepartureDateTo + "|"
		}

		return result
	}

	return &result
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
	},
	{
		Field: "days_after_tariff_start",
	},
	{
		Field: "max_exchange_count",
	},
	{
		Field: "issue_date_from",
	},
	{
		Field: "issue_date_to",
	},
	{
		Field: "departure_date_from",
	},
	{
		Field: "departure_date_to",
	},
	{
		Field: "fare",
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

func (*ManualExchangeRefundRule) compareHoursBeforeDeparture(a *int64, b *int64) bool {
	if a != nil {
		if b == nil || !(*a <= *b) {
			return false
		}
	}
	return true
}

func (*ManualExchangeRefundRule) compareDaysAfterTariffStart(a *int64, b *int64) bool {
	if a != nil {
		if b == nil || !(*a <= *b) {
			return false
		}
	}
	return true
}

func (*ManualExchangeRefundRule) compareMaxExchangeCount(a *int64, b *int64) bool {
	if a != nil {
		if b == nil || !(*a >= *b) {
			return false
		}
	}
	return true
}

func (*ManualExchangeRefundRule) compareIssueDateFrom(a *string, b *string) bool {
	if a != nil {
		if b == nil || !(*a <= *b) {
			return false
		}
	}
	return true
}

func (*ManualExchangeRefundRule) compareIssueDateTo(a *string, b *string) bool {
	if a != nil {
		if b == nil || !(*a >= *b) {
			return false
		}
	}
	return true
}

func (*ManualExchangeRefundRule) compareDepartureDateFrom(a *string, b *string) bool {
	if a != nil {
		if b == nil || !(*a <= *b) {
			return false
		}
	}
	return true
}

func (*ManualExchangeRefundRule) compareDepartureDateTo(a *string, b *string) bool {
	if a != nil {
		if b == nil || !(*a > *b) {
			return false
		}
	}
	return true
}

func (*ManualExchangeRefundRule) compareFare(a *string, b *string) bool {
	if a != nil {
		if b == nil {
			return false
		}
		r, err := regexp.Compile(*a)
		if err != nil {
			return false
		}
		if !r.Match([]byte(*b)) {
			return false
		}
	}
	return true
}
