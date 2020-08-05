package service_charge

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"regexp"
	"stash.tutu.ru/avia-search-common/contracts/v2/base"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
)

var comparisonOrder = frule_module.ComparisonOrder{
	//поле тарифа используется для субсидированных билетов, идет в связке вместе с маршрутом и а/к и перебивает настройки АБ и дней до вылета [направления: города и страны]
	//[город=>город] и конкретная а/к
	[]string{"departure_city_id", "arrival_city_id", "carrier_id", "tariff"},
	//[город=>страна] и конкретная а/к
	[]string{"departure_city_id", "arrival_country_id", "carrier_id", "tariff"},
	//[страна=>город] и конкретная а/к
	[]string{"departure_country_id", "arrival_city_id", "carrier_id", "tariff"},
	//[страна=>страна] и конкретная а/к
	[]string{"departure_country_id", "arrival_country_id", "carrier_id", "tariff"},
	//все рейсы конкретной а/к из конкретного города
	[]string{"departure_city_id", "carrier_id", "tariff"},
	//все рейсы конкретной а/к из конкретной страны
	[]string{"departure_country_id", "carrier_id", "tariff"},
	//все рейсы конкретной а/к в конкретный город
	[]string{"arrival_city_id", "carrier_id", "tariff"},
	//все рейсы конкретной а/к в конкретную страну
	[]string{"arrival_country_id", "carrier_id", "tariff"},
	//все рейсы
	[]string{"carrier_id", "tariff"},

	//основной набор правил [направления: города и страны] +аб кампания
	//[город=>город] и конкретная а/к с днями до вылета
	[]string{"departure_city_id", "arrival_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//[город=>страна] и конкретная а/к с днями до вылета
	[]string{"departure_city_id", "arrival_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//[страна=>город] и конкретная а/к с днями до вылета
	[]string{"departure_country_id", "arrival_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//[страна=>страна] и конкретная а/к с днями до вылета
	[]string{"departure_country_id", "arrival_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//[город=>город] и конкретная а/к
	[]string{"departure_city_id", "arrival_city_id", "carrier_id", "ab_variant"},
	//[город=>страна] и конкретная а/к
	[]string{"departure_city_id", "arrival_country_id", "carrier_id", "ab_variant"},
	//[страна=>город] и конкретная а/к
	[]string{"departure_country_id", "arrival_city_id", "carrier_id", "ab_variant"},
	//[страна=>страна] и конкретная а/к
	[]string{"departure_country_id", "arrival_country_id", "carrier_id", "ab_variant"},
	//все рейсы конкретной а/к из конкретного города с днями до вылета
	[]string{"departure_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы конкретной а/к из конкретной страны с днями до вылета
	[]string{"departure_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы конкретной а/к из конкретного города
	[]string{"departure_city_id", "carrier_id", "ab_variant"},
	//все рейсы конкретной а/к из конкретной страны
	[]string{"departure_country_id", "carrier_id", "ab_variant"},
	//все рейсы конкретной а/к в конкретный город с днями до вылета
	[]string{"arrival_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы конкретной а/к в конкретную страну с днями до вылета
	[]string{"arrival_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы конкретной а/к в конкретный город
	[]string{"arrival_city_id", "carrier_id", "ab_variant"},
	//все рейсы конкретной а/к в конкретную страну
	[]string{"arrival_country_id", "carrier_id", "ab_variant"},
	//[город=>город] любых а/к с днями до вылета
	[]string{"departure_city_id", "arrival_city_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//[город=>страна] любых а/к с днями до вылета
	[]string{"departure_city_id", "arrival_country_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//[страна=>город] любых а/к с днями до вылета
	[]string{"departure_country_id", "arrival_city_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//[страна=>страна] любых а/к с днями до вылета
	[]string{"departure_country_id", "arrival_country_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//[город=>город] любых а/к
	[]string{"departure_city_id", "arrival_city_id", "ab_variant"},
	//[город=>страна] любых а/к
	[]string{"departure_city_id", "arrival_country_id", "ab_variant"},
	//[страна=>город] любых а/к
	[]string{"departure_country_id", "arrival_city_id", "ab_variant"},
	//[страна=>страна] любых а/к
	[]string{"departure_country_id", "arrival_country_id", "ab_variant"},
	//все рейсы из конкретного города любых а/к с днями до вылета
	[]string{"departure_city_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы из конкретной страны любых а/к с днями до вылета
	[]string{"departure_country_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы из конкретного города любых а/к
	[]string{"departure_city_id", "ab_variant"},
	//все рейсы из конкретной страны любых а/к
	[]string{"departure_country_id", "ab_variant"},
	//все рейсы в конкретный город любых а/к с днями до вылета
	[]string{"arrival_city_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы в конкретную страну любых а/к с днями до вылета
	[]string{"arrival_country_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы в конкретный город любых а/к
	[]string{"arrival_city_id", "ab_variant"},
	//все рейсы в конкретную страну любых а/к
	[]string{"arrival_country_id", "ab_variant"},
	//все рейсы конкретной а/к с днями до вылета
	[]string{"carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы конкретной а/к
	[]string{"carrier_id", "ab_variant"},
	//все рейсы с днями до вылета
	[]string{"days_to_departure_min", "days_to_departure_max", "ab_variant"},
	//все рейсы
	[]string{"ab_variant"},

	//основной набор правил [направления: города и страны]
	//[город=>город] и конкретная а/к с днями до вылета
	[]string{"departure_city_id", "arrival_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//[город=>страна] и конкретная а/к с днями до вылета
	[]string{"departure_city_id", "arrival_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//[страна=>город] и конкретная а/к с днями до вылета
	[]string{"departure_country_id", "arrival_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//[страна=>страна] и конкретная а/к с днями до вылета
	[]string{"departure_country_id", "arrival_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//[город=>город] и конкретная а/к
	[]string{"departure_city_id", "arrival_city_id", "carrier_id"},
	//[город=>страна] и конкретная а/к
	[]string{"departure_city_id", "arrival_country_id", "carrier_id"},
	//[страна=>город] и конкретная а/к
	[]string{"departure_country_id", "arrival_city_id", "carrier_id"},
	//[страна=>страна] и конкретная а/к
	[]string{"departure_country_id", "arrival_country_id", "carrier_id"},
	//все рейсы конкретной а/к из конкретного города с днями до вылета
	[]string{"departure_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы конкретной а/к из конкретной страны с днями до вылета
	[]string{"departure_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы конкретной а/к из конкретного города
	[]string{"departure_city_id", "carrier_id"},
	//все рейсы конкретной а/к из конкретной страны
	[]string{"departure_country_id", "carrier_id"},
	//все рейсы конкретной а/к в конкретный город с днями до вылета
	[]string{"arrival_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы конкретной а/к в конкретную страну с днями до вылета
	[]string{"arrival_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы конкретной а/к в конкретный город
	[]string{"arrival_city_id", "carrier_id"},
	//все рейсы конкретной а/к в конкретную страну
	[]string{"arrival_country_id", "carrier_id"},
	//[город=>город] любых а/к с днями до вылета
	[]string{"departure_city_id", "arrival_city_id", "days_to_departure_min", "days_to_departure_max"},
	//[город=>страна] любых а/к с днями до вылета
	[]string{"departure_city_id", "arrival_country_id", "days_to_departure_min", "days_to_departure_max"},
	//[страна=>город] любых а/к с днями до вылета
	[]string{"departure_country_id", "arrival_city_id", "days_to_departure_min", "days_to_departure_max"},
	//[страна=>страна] любых а/к с днями до вылета
	[]string{"departure_country_id", "arrival_country_id", "days_to_departure_min", "days_to_departure_max"},
	//[город=>город] любых а/к
	[]string{"departure_city_id", "arrival_city_id"},
	//[город=>страна] любых а/к
	[]string{"departure_city_id", "arrival_country_id"},
	//[страна=>город] любых а/к
	[]string{"departure_country_id", "arrival_city_id"},
	//[страна=>страна] любых а/к
	[]string{"departure_country_id", "arrival_country_id"},
	//все рейсы из конкретного города любых а/к с днями до вылета
	[]string{"departure_city_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы из конкретной страны любых а/к с днями до вылета
	[]string{"departure_country_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы из конкретного города любых а/к
	[]string{"departure_city_id"},
	//все рейсы из конкретной страны любых а/к
	[]string{"departure_country_id"},
	//все рейсы в конкретный город любых а/к с днями до вылета
	[]string{"arrival_city_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы в конкретную страну любых а/к с днями до вылета
	[]string{"arrival_country_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы в конкретный город любых а/к
	[]string{"arrival_city_id"},
	//все рейсы в конкретную страну любых а/к
	[]string{"arrival_country_id"},
	//все рейсы конкретной а/к с днями до вылета
	[]string{"carrier_id", "days_to_departure_min", "days_to_departure_max"},
	//все рейсы конкретной а/к
	[]string{"carrier_id"},
	//все рейсы с днями до вылета
	[]string{"days_to_departure_min", "days_to_departure_max"},
	//все рейсы
	[]string{},
}

var strategyKeys = []string{
	"carrier_id",
	"tariff",
	"ab_variant",
	"departure_city_id",
	"arrival_city_id",
	"departure_country_id",
	"arrival_country_id",
	"days_to_departure_min",
	"days_to_departure_max",
}

type ServiceChargeRule struct {
	Id                 int32       `json:"id"`
	Version            int32       `json:"version"`
	CarrierId          *int64      `json:"carrier_id"`
	DaysToDepartureMin *int64      `json:"days_to_departure_min"`
	DaysToDepartureMax *int64      `json:"days_to_departure_max"`
	FareType           *string     `json:"tariff"`
	ABVariant          interface{} `json:"ab_variant"`
	DepartureCountryId *uint64     `json:"departure_country_id"`
	ArrivalCountryId   *uint64     `json:"arrival_country_id"`
	DepartureCityId    *uint64     `json:"departure_city_id"`
	ArrivalCityId      *uint64     `json:"arrival_city_id"`
	Margin             *string     `json:"result_margin"`
	MarginParsed       *Margin
	TestOfferPrice     base.Money
	CurrencyConverter  base.CurrencyConverter
	repo               *frule_module.Repository
}

type Conditions struct {
	PriceRange *string `json:"price_range"`
}

type MoneyParsed struct {
	Percent float64
	Limit   *base.Money
	Money   *base.Money
}

type ConditionMarginResult struct {
	Conditions   Conditions `json:"conditions"`
	Result       *string    `json:"result"`
	ResultParsed MoneyParsed
}

type Margin struct {
	Full   []ConditionMarginResult `json:"full"`
	Child  []ConditionMarginResult `json:"child"`
	Infant []ConditionMarginResult `json:"infant"`
}

type ServiceChargeRuleResult struct {
	Id      int32
	Version int32
	Margin  struct {
		Full   base.Money
		Child  base.Money
		Infant base.Money
	}
}

func NewServiceChargeFRule(ctx context.Context, config *repository.Config) (*ServiceChargeRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &ServiceChargeRule{repo: repo}, nil
}

/* определение строки, которая начинается на число% */
var startFromPercentSpec = regexp.MustCompile(`^([0-9\.]+)%`)

/*
парсинг строки вида:
583.44RUB+2.1%<1000.12RUB, где 583.44RUB - абсолютное значение, 2.1%<1000.12RUB - процент от тарифа, но не более 1000.12RUB

примеры:
0RUB
583RUB
583.44RUB
2.1%
2.1%<1000.12RUB
583.44RUB+2.1%
583.44RUB+2.1%<1000.12
583.44RUB+2.1%<1000.12RUB
*/
var moneySpec = regexp.MustCompile(`^([0-9\.]+)([A-Z]{3})\+?([0-9\.]*)%?<?([0-9\.]*)([A-Z]{0,3})$`)

func parseMoneySpec(spec *string) MoneyParsed {
	moneyParsed := MoneyParsed{}
	if spec != nil {
		specString := *spec
		if specString != "" {
			// костылек - если в начале строки указывается % от тарифа, то для совпадения с основным регулярным выражением добавим 0RUB
			if percentParsedData := startFromPercentSpec.FindStringSubmatch(specString); len(percentParsedData) > 0 {
				specString = "0RUB+" + specString
			}

			/**
			Парсинг строки вида 583.44RUB+2.1%<1000.12RUB:
			Full match	583.44RUB+2.1%<1000.12RUB
			Group 1.	583.44
			Group 2.	RUB
			Group 3.	2.1
			Group 4.	1000.12
			Group 5.	RUB
			*/
			parsedData := moneySpec.FindStringSubmatch(specString)
			if len(parsedData) == 6 {
				if parsedData[1] != "" {
					amount, err := strconv.ParseFloat(parsedData[1], 64)
					if err != nil {
						log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot parse string %s", specString)).Msg("parsing service charge amount")
					}
					moneyParsed.Money = &base.Money{
						Amount:   int64(math.Round(amount * 100)),
						Currency: &base.Currency{Code: fixCurrencyCode(parsedData[2]), Fraction: 100},
					}
				}
				if parsedData[3] != "" {
					percent, err := strconv.ParseFloat(parsedData[3], 64)
					if err != nil {
						log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot parse string %s", specString)).Msg("parsing service charge percent")
					}
					if percent > 0 {
						moneyParsed.Percent = percent
					}
				}
				if parsedData[4] != "" {
					amount, err := strconv.ParseFloat(parsedData[4], 64)
					if err != nil {
						log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot parse string %s", specString)).Msg("parsing service charge limit")
					}
					moneyParsed.Limit = &base.Money{
						Amount:   int64(math.Round(amount * 100)),
						Currency: &base.Currency{Code: fixCurrencyCode(parsedData[5]), Fraction: 100},
					}
				}
			} else {
				log.Logger.Error().Stack().Err(fmt.Errorf("cannot parse string %s", specString)).Msg("parsing service charge")
			}
		} else {
			log.Logger.Error().Stack().Err(errors.New("spec string is empty")).Msg("parsing service charge")
		}
	}
	return moneyParsed
}

func fixCurrencyCode(code string) string {
	if code == "" || code == "RUR" {
		return "RUB"
	}
	return code
}

func findPricingRangeValue(choices []ConditionMarginResult, testRule ServiceChargeRule) MoneyParsed {
	for _, choice := range choices {
		if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) {
			return choice.ResultParsed
		}
	}
	return MoneyParsed{}
}

func calculatePassengerServiceCharge(moneyParsed MoneyParsed, price base.Money, currencyConverter base.CurrencyConverter) base.Money {
	var passengerServiceCharge *base.Money

	if price.Validate() {
		passengerServiceCharge = base.CloneMoney(&price)
		passengerServiceCharge.Amount = 0
	} else {
		if moneyParsed.Money != nil {
			passengerServiceCharge = base.CloneMoney(moneyParsed.Money)
			passengerServiceCharge.Amount = 0
		} else {
			passengerServiceCharge = base.CreateZeroRubMoney()
		}
	}

	if moneyParsed.Money != nil {
		if currencyConverter != nil {
			if err := passengerServiceCharge.ConvertAndAdd(currencyConverter, moneyParsed.Money); err != nil {
				log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot add amount %v", *moneyParsed.Money)).Msg("calculate service charge")
			}
		} else {
			if err := passengerServiceCharge.Add(moneyParsed.Money); err != nil {
				log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot add amount %v", *moneyParsed.Money)).Msg("calculate service charge")
			}
		}
	}

	if moneyParsed.Percent != 0 && price.Validate() {
		tariffAddition := base.CloneMoney(&price)
		tariffAddition.MultiplyFloat64(moneyParsed.Percent / 100)

		if moneyParsed.Limit != nil {
			if currencyConverter != nil {
				if moreThanLimit, _ := tariffAddition.ConvertAndCompare(currencyConverter, moneyParsed.Limit); moreThanLimit == 1 {
					tariffAddition.Copy(moneyParsed.Limit)
				}
			} else {
				if moreThanLimit, _ := tariffAddition.More(moneyParsed.Limit); moreThanLimit {
					tariffAddition.Copy(moneyParsed.Limit)
				}
			}
		}

		if currencyConverter != nil {
			if err := passengerServiceCharge.ConvertAndAdd(currencyConverter, tariffAddition); err != nil {
				log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot add amount %v", *tariffAddition)).Msg("calculate service charge")
			}
		} else {
			if err := passengerServiceCharge.Add(tariffAddition); err != nil {
				log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot add amount %v", *tariffAddition)).Msg("calculate service charge")
			}
		}
	}

	return *passengerServiceCharge
}

func (rule *ServiceChargeRule) GetResultValue(testRule interface{}) interface{} {
	result := ServiceChargeRuleResult{
		Id:      rule.Id,
		Version: rule.Version,
	}
	if rule.MarginParsed != nil {
		serviceChargeParams := testRule.(ServiceChargeRule)

		result.Margin.Full = calculatePassengerServiceCharge(
			findPricingRangeValue(rule.MarginParsed.Full, serviceChargeParams),
			serviceChargeParams.TestOfferPrice,
			serviceChargeParams.CurrencyConverter,
		)

		result.Margin.Child = calculatePassengerServiceCharge(
			findPricingRangeValue(rule.MarginParsed.Child, serviceChargeParams),
			serviceChargeParams.TestOfferPrice,
			serviceChargeParams.CurrencyConverter,
		)

		result.Margin.Infant = calculatePassengerServiceCharge(
			findPricingRangeValue(rule.MarginParsed.Infant, serviceChargeParams),
			serviceChargeParams.TestOfferPrice,
			serviceChargeParams.CurrencyConverter,
		)
	}
	return result
}

func (rule *ServiceChargeRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

func (rule *ServiceChargeRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{
	{
		Field: "days_to_departure_min",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(int64) <= b.Elem().Interface().(int64)
		},
	},
	{
		Field: "days_to_departure_max",
		Function: func(a, b reflect.Value) bool {
			return a.Elem().Interface().(int64) > b.Elem().Interface().(int64)
		},
	},
	{
		Field: "ab_variant",
		Function: func(a, b reflect.Value) bool {
			offerABCampaigns, ok := b.Elem().Interface().([]string)
			if !ok {
				return false
			}
			return frule_module.InSlice(a.Elem().Interface().(string), offerABCampaigns)
		},
	},
}

func (rule *ServiceChargeRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

func (rule *ServiceChargeRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *ServiceChargeRule) GetDefaultValue() interface{} {
	return ServiceChargeRuleResult{
		Id:      -1,
		Version: -1,
		Margin: struct {
			Full   base.Money
			Child  base.Money
			Infant base.Money
		}{
			Full:   *base.CreateZeroRubMoney(),
			Child:  *base.CreateZeroRubMoney(),
			Infant: *base.CreateZeroRubMoney(),
		},
	}
}

func (rule *ServiceChargeRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *ServiceChargeRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *ServiceChargeRule) GetRuleName() string {
	return "ServiceCharge"
}
