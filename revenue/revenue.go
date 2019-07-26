package revenue

import (
	"context"
	"encoding/json"
	"reflect"
	"regexp"
	"stash.tutu.ru/avia-search-common/contracts/base"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
	"time"
)

type RevenueRule struct {
	Id                     int     `json:"id"`
	CarrierId              *int64  `json:"carrier_id"`
	Partner                *string `json:"partner"`
	ConnectionGroup        *string `json:"connection_group"`
	TicketingConnection    *string `json:"ticketing_connection"`
	DaysToDepartureMin     *int64  `json:"days_to_departure_min"`
	DaysToDepartureMax     *int64  `json:"days_to_departure_max"`
	FareType               *string `json:"tariff"`
	ABVariant              *string `json:"ab_variant"`
	DepartureCountryId     *int64  `json:"departure_country_id"`
	ArrivalCountryId       *int64  `json:"arrival_country_id"`
	DepartureCityId        *int64  `json:"departure_city_id"`
	ArrivalCityId          *int64  `json:"arrival_city_id"`
	Revenue                *string `json:"revenue"`
	Margin                 *string `json:"margin"`
	TestOfferPrice         base.Money
	TestOfferPurchaseDate  time.Time
	TestOfferDepartureDate time.Time
	repo                   *frule_module.Repository
}

type Revenue struct {
	Full   []ConditionRevenueResult `json:"full"`
	Child  []ConditionRevenueResult `json:"child"`
	Infant []ConditionRevenueResult `json:"infant"`
}

type Conditions struct {
	PriceRange        *string `json:"price_range"`
	DepartureCronSpec *string `json:"departure_cron_spec"`
	PurchaseCronSpec  *string `json:"purchase_cron_spec"`
}

type Result struct {
	Ticket  *string `json:"ticket"`
	Segment *string `json:"segment"`
}

type ConditionRevenueResult struct {
	Conditions Conditions `json:"conditions"`
	Result     Result     `json:"result"`
}

type ConditionMarginResult struct {
	Conditions Conditions `json:"conditions"`
	Result     string     `json:"result"`
}

type Margin struct {
	Full   []ConditionMarginResult `json:"full"`
	Child  []ConditionMarginResult `json:"child"`
	Infant []ConditionMarginResult `json:"infant"`
}

type RevenueRuleResult struct {
	Id      int
	Revenue struct {
		Full struct {
			Ticket  base.Money
			Segment base.Money
		}
		Child struct {
			Ticket  base.Money
			Segment base.Money
		}
		Infant struct {
			Ticket  base.Money
			Segment base.Money
		}
	}
	Margin struct {
		Full   base.Money
		Child  base.Money
		Infant base.Money
	}
}

func NewRevenueFRule(ctx context.Context, config *repository.Config) (*RevenueRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &RevenueRule{repo: repo}, nil
}

var moneySpec = regexp.MustCompile("([0-9]+)([A-Z]+)")

func parseMoneySpec(spec *string) base.Money {

	if spec == nil {
		return base.Money{
			Currency: &base.Currency{
				Code:     "RUB",
				Fraction: 100,
			},
		}
	}
	parsedData := moneySpec.FindStringSubmatch(*spec)
	if len(parsedData) == 3 {
		amount, err := strconv.ParseInt(parsedData[1], 10, 64)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Parsing money")
		}
		return base.Money{
			Amount: amount,
			Currency: &base.Currency{ // TODO: load from DB by code
				Code:     parsedData[2],
				Fraction: 100,
			},
		}
	} else {
		return base.Money{
			Currency: &base.Currency{
				Code:     "RUB",
				Fraction: 100,
			},
		}
	}
}

func selectRevenueRow(choices []ConditionRevenueResult, testRule RevenueRule) Result {
	for _, choice := range choices {
		if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			frule_module.CronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) {
			return choice.Result
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			choice.Conditions.PurchaseCronSpec == nil {
			return choice.Result
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) &&
			choice.Conditions.DepartureCronSpec == nil {
			return choice.Result
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			choice.Conditions.DepartureCronSpec == nil && choice.Conditions.PurchaseCronSpec == nil {
			return choice.Result
		}
	}
	return Result{}
}

func selectMarginRow(choices []ConditionMarginResult, testRule RevenueRule) string {
	for _, choice := range choices {
		if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			frule_module.CronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) {
			return choice.Result
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			choice.Conditions.PurchaseCronSpec == nil {
			return choice.Result
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) &&
			choice.Conditions.DepartureCronSpec == nil {
			return choice.Result
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			choice.Conditions.DepartureCronSpec == nil && choice.Conditions.PurchaseCronSpec == nil {
			return choice.Result
		}
	}
	return ""
}

func (rule *RevenueRule) GetResultValue(testRule interface{}) interface{} {
	var revenueResult Revenue
	var marginResult Margin
	result := RevenueRuleResult{
		Id: rule.Id,
	}
	if rule.Revenue != nil && *rule.Revenue != "[]" {
		if err := json.Unmarshal([]byte(*rule.Revenue), &revenueResult); err != nil {
			log.Logger.Error().Err(err).Msg("Unmarshal revenue")
		}
		fullResult := selectRevenueRow(revenueResult.Full, testRule.(RevenueRule))
		result.Revenue.Full.Ticket = parseMoneySpec(fullResult.Ticket)
		result.Revenue.Full.Segment = parseMoneySpec(fullResult.Segment)
		childResult := selectRevenueRow(revenueResult.Child, testRule.(RevenueRule))
		result.Revenue.Child.Ticket = parseMoneySpec(childResult.Ticket)
		result.Revenue.Child.Segment = parseMoneySpec(childResult.Segment)
		infantResult := selectRevenueRow(revenueResult.Infant, testRule.(RevenueRule))
		result.Revenue.Infant.Ticket = parseMoneySpec(infantResult.Ticket)
		result.Revenue.Infant.Segment = parseMoneySpec(infantResult.Segment)
	}
	if rule.Margin != nil && *rule.Margin != "[]" {
		if err := json.Unmarshal([]byte(*rule.Margin), &marginResult); err != nil {
			log.Logger.Error().Err(err).Msg("Unmarshal margin")
		}
		fullResult := selectMarginRow(marginResult.Full, testRule.(RevenueRule))
		result.Margin.Full = parseMoneySpec(&fullResult)
		childResult := selectMarginRow(marginResult.Child, testRule.(RevenueRule))
		result.Margin.Child = parseMoneySpec(&childResult)
		infantResult := selectMarginRow(marginResult.Infant, testRule.(RevenueRule))
		result.Margin.Infant = parseMoneySpec(&infantResult)
	}
	return result
}

func (rule *RevenueRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return frule_module.ComparisonOrder{
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id",
			"carrier_id", "tariff",
		},

		// []string{город=>город} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_city_id", "carrier_id", "tariff"},

		// []string{город=>город} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_city_id", "carrier_id", "tariff"},

		// []string{город=>страна} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id",
			"carrier_id", "tariff",
		},

		// []string{город=>страна} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_country_id", "carrier_id", "tariff"},

		// []string{город=>страна} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_country_id", "carrier_id", "tariff"},

		// []string{страна=>город} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id",
			"carrier_id", "tariff",
		},

		// []string{страна=>город} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_city_id", "carrier_id", "tariff"},

		// []string{страна=>город} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_city_id", "carrier_id", "tariff"},

		// []string{страна=>страна} и конкретная а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id",
			"carrier_id", "tariff"},

		// []string{страна=>страна} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "carrier_id", "tariff"},

		// []string{страна=>страна} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_country_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к из конкретного города
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к из конкретного города без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к из конкретного города без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к из конкретной страны
		[]string{"partner", "connection_group", "ticketing_connection", "departure_country_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к из конкретной страны без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к из конкретной страны без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к в конкретный город
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_city_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к в конкретный город без ticketing_connection
		[]string{"partner", "connection_group", "arrival_city_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к в конкретный город без ticketing_connection, connection_group
		[]string{"partner", "arrival_city_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к в конкретную страну
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_country_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к в конкретную страну без ticketing_connection
		[]string{"partner", "connection_group", "arrival_country_id", "carrier_id", "tariff"},

		// все рейсы конкретной а/к в конкретную страну без ticketing_connection, connection_group
		[]string{"partner", "arrival_country_id", "carrier_id", "tariff"},

		// все рейсы
		[]string{"partner", "connection_group", "ticketing_connection", "carrier_id", "tariff"},

		// все рейсы без ticketing_connection
		[]string{"partner", "connection_group", "carrier_id", "tariff"},

		// все рейсы без ticketing_connection, connection_group
		[]string{"partner", "carrier_id", "tariff"},

		// основной набор правил []string{направления: города и страны} +аб кампания
		// []string{город=>город} и конкретная а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id",
			"carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{город=>город} и конкретная а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "arrival_city_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{город=>город} и конкретная а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_city_id", "arrival_city_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// []string{город=>страна} и конкретная а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id",
			"carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{город=>страна} и конкретная а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "arrival_country_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{город=>страна} и конкретная а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_city_id", "arrival_country_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// []string{страна=>город} и конкретная а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id",
			"carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{страна=>город} и конкретная а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "arrival_city_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{страна=>город} и конкретная а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_country_id", "arrival_city_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// []string{страна=>страна} и конкретная а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id",
			"carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{страна=>страна} и конкретная а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "arrival_country_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{страна=>страна} и конкретная а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_country_id", "arrival_country_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// []string{город=>город} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id",
			"carrier_id", "ab_variant",
		},

		// []string{город=>город} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_city_id", "carrier_id", "ab_variant"},

		// []string{город=>город} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_city_id", "carrier_id", "ab_variant"},

		// []string{город=>страна} и конкретная а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id",
			"carrier_id", "ab_variant"},

		// []string{город=>страна} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_country_id", "carrier_id", "ab_variant"},

		// []string{город=>страна} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_country_id", "carrier_id", "ab_variant"},

		// []string{страна=>город} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id",
			"carrier_id", "ab_variant",
		},

		// []string{страна=>город} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_city_id", "carrier_id", "ab_variant"},

		// []string{страна=>город} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_city_id", "carrier_id", "ab_variant"},

		// []string{страна=>страна} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id",
			"carrier_id", "ab_variant",
		},

		// []string{страна=>страна} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "carrier_id", "ab_variant"},

		// []string{страна=>страна} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_country_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к из конкретного города с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к из конкретного города с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к из конкретного города с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы конкретной а/к из конкретной страны с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к из конкретной страны с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к из конкретной страны с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы конкретной а/к из конкретного города
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к из конкретного города без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к из конкретного города без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к из конкретной страны
		[]string{"partner", "connection_group", "ticketing_connection", "departure_country_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к из конкретной страны без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к из конкретной страны без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к в конкретный город с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "arrival_city_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к в конкретный город с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "arrival_city_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к в конкретный город с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "arrival_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы конкретной а/к в конкретную страну с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "arrival_country_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к в конкретную страну с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "arrival_country_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к в конкретную страну с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "arrival_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы конкретной а/к в конкретный город
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_city_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к в конкретный город без ticketing_connection
		[]string{"partner", "connection_group", "arrival_city_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к в конкретный город без ticketing_connection, connection_group
		[]string{"partner", "arrival_city_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к в конкретную страну
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_country_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к в конкретную страну без ticketing_connection
		[]string{"partner", "connection_group", "arrival_country_id", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к в конкретную страну без ticketing_connection, connection_group
		[]string{"partner", "arrival_country_id", "carrier_id", "ab_variant"},

		// []string{город=>город} любых а/к с днями до вылета
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant"},

		// []string{город=>город} любых а/к с днями до вылета без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_city_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant"},

		// []string{город=>город} любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_city_id", "arrival_city_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// []string{город=>страна} любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{город=>страна} любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "arrival_country_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// []string{город=>страна} любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_city_id", "arrival_country_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// []string{страна=>город} любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{страна=>город} любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "arrival_city_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// []string{страна=>город} любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_country_id", "arrival_city_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// []string{страна=>страна} любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id",
			"days_to_departure_min", "days_to_departure_max", "ab_variant",
		},

		// []string{страна=>страна} любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "arrival_country_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// []string{страна=>страна} любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_country_id", "arrival_country_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// []string{город=>город} любых а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id",
			"ab_variant",
		},

		// []string{город=>город} любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_city_id", "ab_variant"},

		// []string{город=>город} любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_city_id", "ab_variant"},

		// []string{город=>страна} любых а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id",
			"ab_variant",
		},

		// []string{город=>страна} любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_country_id", "ab_variant"},

		// []string{город=>страна} любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_country_id", "ab_variant"},

		// []string{страна=>город} любых а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id",
			"ab_variant",
		},

		// []string{страна=>город} любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_city_id", "ab_variant"},

		// []string{страна=>город} любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_city_id", "ab_variant"},

		// []string{страна=>страна} любых а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id",
			"ab_variant",
		},

		// []string{страна=>страна} любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "ab_variant"},

		// []string{страна=>страна} любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_country_id", "ab_variant"},

		// все рейсы из конкретного города любых а/к с днями до вылета
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant"},

		// все рейсы из конкретного города любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы из конкретного города любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},

		// все рейсы из конкретной страны любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// все рейсы из конкретной страны любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы из конкретной страны любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},

		// все рейсы из конкретного города любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "ab_variant"},

		// все рейсы из конкретного города любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "ab_variant"},

		// все рейсы из конкретного города любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "ab_variant"},

		// все рейсы из конкретной страны любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_country_id", "ab_variant"},

		// все рейсы из конкретной страны любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "ab_variant"},

		// все рейсы из конкретной страны любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "ab_variant"},

		// все рейсы в конкретный город любых а/к с днями до вылета
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_city_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant"},

		// все рейсы в конкретный город любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "arrival_city_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы в конкретный город любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "arrival_city_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},

		// все рейсы в конкретную страну любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "arrival_country_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// все рейсы в конкретную страну любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "arrival_country_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы в конкретную страну любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "arrival_country_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},

		// все рейсы в конкретный город любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_city_id", "ab_variant"},

		// все рейсы в конкретный город любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "arrival_city_id", "ab_variant"},

		// все рейсы в конкретный город любых а/к без ticketing_connection, connection_group
		[]string{"partner", "arrival_city_id", "ab_variant"},

		// все рейсы в конкретную страну любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_country_id", "ab_variant"},

		// все рейсы в конкретную страну любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "arrival_country_id", "ab_variant"},

		// все рейсы в конкретную страну любых а/к без ticketing_connection, connection_group
		[]string{"partner", "arrival_country_id", "ab_variant"},

		// все рейсы конкретной а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "carrier_id", "days_to_departure_min",
			"days_to_departure_max", "ab_variant",
		},

		// все рейсы конкретной а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "carrier_id", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы конкретной а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "carrier_id", "days_to_departure_min", "days_to_departure_max", "ab_variant"},

		// все рейсы конкретной а/к
		[]string{"partner", "connection_group", "ticketing_connection", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к без ticketing_connection
		[]string{"partner", "connection_group", "carrier_id", "ab_variant"},

		// все рейсы конкретной а/к без ticketing_connection, connection_group
		[]string{"partner", "carrier_id", "ab_variant"},

		// все рейсы с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "days_to_departure_min", "days_to_departure_max",
			"ab_variant",
		},

		// все рейсы с днями до вылета без ticketing_connection
		[]string{"partner", "connection_group", "days_to_departure_min", "days_to_departure_max", "ab_variant"},

		// все рейсы с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "days_to_departure_min", "days_to_departure_max", "ab_variant"},

		// все рейсы
		[]string{"partner", "connection_group", "ticketing_connection", "ab_variant"},

		// все рейсы без ticketing_connection
		[]string{"partner", "connection_group", "ab_variant"},

		// все рейсы без ticketing_connection, connection_group
		[]string{"partner", "ab_variant"},

		// основной набор правил []string{направления: города и страны}
		// []string{город=>город} и конкретная а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id",
			"carrier_id", "days_to_departure_min", "days_to_departure_max",
		},

		// []string{город=>город} и конкретная а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "arrival_city_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// []string{город=>город} и конкретная а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_city_id", "arrival_city_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// []string{город=>страна} и конкретная а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id",
			"carrier_id", "days_to_departure_min", "days_to_departure_max",
		},

		// []string{город=>страна} и конкретная а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "arrival_country_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// []string{город=>страна} и конкретная а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_city_id", "arrival_country_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// []string{страна=>город} и конкретная а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id",
			"carrier_id", "days_to_departure_min", "days_to_departure_max",
		},

		// []string{страна=>город} и конкретная а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "arrival_city_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// []string{страна=>город} и конкретная а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_country_id", "arrival_city_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// []string{страна=>страна} и конкретная а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id",
			"carrier_id", "days_to_departure_min", "days_to_departure_max",
		},

		// []string{страна=>страна} и конкретная а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "arrival_country_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// []string{страна=>страна} и конкретная а/к с днями до вылета без ticketing_connection, connection_group
		[]string{
			"partner", "departure_country_id", "arrival_country_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// []string{город=>город} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id",
			"carrier_id",
		},

		// []string{город=>город} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_city_id", "carrier_id"},

		// []string{город=>город} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_city_id", "carrier_id"},

		// []string{город=>страна} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id",
			"carrier_id",
		},

		// []string{город=>страна} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_country_id", "carrier_id"},

		// []string{город=>страна} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_country_id", "carrier_id"},

		// []string{страна=>город} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id",
			"carrier_id",
		},

		// []string{страна=>город} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_city_id", "carrier_id"},

		// []string{страна=>город} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_city_id", "carrier_id"},

		// []string{страна=>страна} и конкретная а/к
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id",
			"carrier_id",
		},

		// []string{страна=>страна} и конкретная а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id", "carrier_id"},

		// []string{страна=>страна} и конкретная а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_country_id", "carrier_id"},

		// все рейсы конкретной а/к из конкретного города с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// все рейсы конкретной а/к из конкретного города с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// все рейсы конкретной а/к из конкретного города с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы конкретной а/к из конкретной страны с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// все рейсы конкретной а/к из конкретной страны с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max"},

		// все рейсы конкретной а/к из конкретной страны с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы конкретной а/к из конкретного города
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "carrier_id"},

		// все рейсы конкретной а/к из конкретного города без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "carrier_id"},

		// все рейсы конкретной а/к из конкретного города без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "carrier_id"},

		// все рейсы конкретной а/к из конкретной страны
		[]string{"partner", "connection_group", "ticketing_connection", "departure_country_id", "carrier_id"},

		// все рейсы конкретной а/к из конкретной страны без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "carrier_id"},

		// все рейсы конкретной а/к из конкретной страны без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "carrier_id"},

		// все рейсы конкретной а/к в конкретный город с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "arrival_city_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// все рейсы конкретной а/к в конкретный город с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "arrival_city_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// все рейсы конкретной а/к в конкретный город с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "arrival_city_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы конкретной а/к в конкретную страну с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "arrival_country_id", "carrier_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// все рейсы конкретной а/к в конкретную страну с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "arrival_country_id", "carrier_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// все рейсы конкретной а/к в конкретную страну с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "arrival_country_id", "carrier_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы конкретной а/к в конкретный город
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_city_id", "carrier_id"},

		// все рейсы конкретной а/к в конкретный город без ticketing_connection
		[]string{"partner", "connection_group", "arrival_city_id", "carrier_id"},

		// все рейсы конкретной а/к в конкретный город без ticketing_connection, connection_group
		[]string{"partner", "arrival_city_id", "carrier_id"},

		// все рейсы конкретной а/к в конкретную страну
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_country_id", "carrier_id"},

		// все рейсы конкретной а/к в конкретную страну без ticketing_connection
		[]string{"partner", "connection_group", "arrival_country_id", "carrier_id"},

		// все рейсы конкретной а/к в конкретную страну без ticketing_connection, connection_group
		[]string{"partner", "arrival_country_id", "carrier_id"},

		// []string{город=>город} любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// []string{город=>город} любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "arrival_city_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// []string{город=>город} любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_city_id", "days_to_departure_min", "days_to_departure_max"},

		// []string{город=>страна} любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// []string{город=>страна} любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_city_id", "arrival_country_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// []string{город=>страна} любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_country_id", "days_to_departure_min", "days_to_departure_max"},

		// []string{страна=>город} любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// []string{страна=>город} любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "arrival_city_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// []string{страна=>город} любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_city_id", "days_to_departure_min", "days_to_departure_max"},

		// []string{страна=>страна} любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id",
			"days_to_departure_min", "days_to_departure_max",
		},

		// []string{страна=>страна} любых а/к с днями до вылета без ticketing_connection
		[]string{
			"partner", "connection_group", "departure_country_id", "arrival_country_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// []string{страна=>страна} любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_country_id", "days_to_departure_min", "days_to_departure_max"},

		// []string{город=>город} любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_city_id"},

		// []string{город=>город} любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_city_id"},

		// []string{город=>город} любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_city_id"},

		// []string{город=>страна} любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id", "arrival_country_id"},

		// []string{город=>страна} любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "arrival_country_id"},

		// []string{город=>страна} любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "arrival_country_id"},

		// []string{страна=>город} любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_city_id"},

		// []string{страна=>город} любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_city_id"},

		// []string{страна=>город} любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_city_id"},

		// []string{страна=>страна} любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_country_id", "arrival_country_id"},

		// []string{страна=>страна} любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "arrival_country_id"},

		// []string{страна=>страна} любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "arrival_country_id"},

		// все рейсы из конкретного города любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_city_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// все рейсы из конкретного города любых а/к с днями до вылета без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы из конкретного города любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы из конкретной страны любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "departure_country_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// все рейсы из конкретной страны любых а/к с днями до вылета без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы из конкретной страны любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы из конкретного города любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_city_id"},

		// все рейсы из конкретного города любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_city_id"},

		// все рейсы из конкретного города любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_city_id"},

		// все рейсы из конкретной страны любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "departure_country_id"},

		// все рейсы из конкретной страны любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "departure_country_id"},

		// все рейсы из конкретной страны любых а/к без ticketing_connection, connection_group
		[]string{"partner", "departure_country_id"},

		// все рейсы в конкретный город любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "arrival_city_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// все рейсы в конкретный город любых а/к с днями до вылета без ticketing_connection
		[]string{"partner", "connection_group", "arrival_city_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы в конкретный город любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "arrival_city_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы в конкретную страну любых а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "arrival_country_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// все рейсы в конкретную страну любых а/к с днями до вылета без ticketing_connection
		[]string{"partner", "connection_group", "arrival_country_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы в конкретную страну любых а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "arrival_country_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы в конкретный город любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_city_id"},

		// все рейсы в конкретный город любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "arrival_city_id"},

		// все рейсы в конкретный город любых а/к без ticketing_connection, connection_group
		[]string{"partner", "arrival_city_id"},

		// все рейсы в конкретную страну любых а/к
		[]string{"partner", "connection_group", "ticketing_connection", "arrival_country_id"},

		// все рейсы в конкретную страну любых а/к без ticketing_connection
		[]string{"partner", "connection_group", "arrival_country_id"},

		// все рейсы в конкретную страну любых а/к без ticketing_connection, connection_group
		[]string{"partner", "arrival_country_id"},

		// все рейсы конкретной а/к с днями до вылета
		[]string{
			"partner", "connection_group", "ticketing_connection", "carrier_id", "days_to_departure_min",
			"days_to_departure_max",
		},

		// все рейсы конкретной а/к с днями до вылета без ticketing_connection
		[]string{"partner", "connection_group", "carrier_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы конкретной а/к с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "carrier_id", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы конкретной а/к
		[]string{"partner", "connection_group", "ticketing_connection", "carrier_id"},

		// все рейсы конкретной а/к без ticketing_connection
		[]string{"partner", "connection_group", "carrier_id"},

		// все рейсы конкретной а/к без ticketing_connection, connection_group
		[]string{"partner", "carrier_id"},

		// все рейсы с днями до вылета
		[]string{"partner", "connection_group", "ticketing_connection", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы с днями до вылета без ticketing_connection
		[]string{"partner", "connection_group", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы с днями до вылета без ticketing_connection, connection_group
		[]string{"partner", "days_to_departure_min", "days_to_departure_max"},

		// все рейсы
		[]string{"partner", "connection_group", "ticketing_connection"},

		// все рейсы без ticketing_connection
		[]string{"partner", "connection_group"},

		// все рейсы без ticketing_connection, connection_group
		[]string{"partner"},
	}
}

func (rule *RevenueRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return frule_module.ComparisonOperators{
		"ab_variant": func(a, b reflect.Value) bool {
			if a.IsNil() {
				return true
			}
			offerABCampaigns, ok := b.Elem().Interface().([]string)
			if !ok {
				return false
			}
			return frule_module.InSlice(a.Elem().Interface().(string), offerABCampaigns)
		},
		"days_to_departure_min": func(a, b reflect.Value) bool {
			if a.IsNil() {
				return true
			}
			return a.Elem().Interface().(string) <= b.Elem().Interface().(string)
		},
		"days_to_departure_max": func(a, b reflect.Value) bool {
			if a.IsNil() {
				return true
			}
			return a.Elem().Interface().(string) > b.Elem().Interface().(string)
		},
	}
}

func (rule *RevenueRule) GetStrategyKeys() []string {
	return []string{
		"partner",
		"connection_group",
		"ticketing_connection",
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
}

func (rule *RevenueRule) GetDefaultValue() interface{} {
	zeroRub := base.Money{
		Currency: &base.Currency{
			Fraction: 100,
			Code:     "RUB",
		},
	}
	return RevenueRuleResult{
		Id: -1,
		Revenue: struct {
			Full struct {
				Ticket  base.Money
				Segment base.Money
			}
			Child struct {
				Ticket  base.Money
				Segment base.Money
			}
			Infant struct {
				Ticket  base.Money
				Segment base.Money
			}
		}{
			Full: struct {
				Ticket  base.Money
				Segment base.Money
			}{
				Ticket:  zeroRub,
				Segment: zeroRub,
			},
			Child: struct {
				Ticket  base.Money
				Segment base.Money
			}{
				Ticket:  zeroRub,
				Segment: zeroRub,
			},
			Infant: struct {
				Ticket  base.Money
				Segment base.Money
			}{
				Ticket:  zeroRub,
				Segment: zeroRub,
			},
		},
		Margin: struct {
			Full   base.Money
			Child  base.Money
			Infant base.Money
		}{
			Full:   zeroRub,
			Child:  zeroRub,
			Infant: zeroRub,
		},
	}
}

func (rule *RevenueRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *RevenueRule) GetNotificationChannel() chan error {
	return rule.repo.NotificationChannel
}
