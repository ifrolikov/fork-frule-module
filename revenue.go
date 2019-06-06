package frule_module

import (
	"encoding/json"
	"reflect"
	"regexp"
	"stash.tutu.ru/avia-search-common/contracts/base"
	"stash.tutu.ru/golang/log"
	"stash.tutu.ru/golang/resources/db"
	"strconv"
	"time"
)

type RevenueRule struct {
	Id                     int     `gorm:"column:id"`
	CarrierId              *int64  `gorm:"column:carrier_id"`
	Partner                *string `gorm:"column:partner"`
	ConnectionGroup        *string `gorm:"column:connection_group"`
	TicketingConnection    *string `gorm:"column:ticketing_connection"`
	DaysToDepartureMin     *int64  `gorm:"column:days_to_departure_min"`
	DaysToDepartureMax     *int64  `gorm:"column:days_to_departure_max"`
	FareType               *string `gorm:"column:tariff"`
	ABVariant              *string `gorm:"column:ab_variant"`
	DepartureCountryId     *int64  `gorm:"column:departure_country_id"`
	ArrivalCountryId       *int64  `gorm:"column:arrival_country_id"`
	DepartureCityId        *int64  `gorm:"column:departure_city_id"`
	ArrivalCityId          *int64  `gorm:"column:arrival_city_id"`
	Revenue                *string `gorm:"column:result_revenue"`
	Margin                 *string `gorm:"column:result_margin"`
	TestOfferPrice         base.Money
	TestOfferPurchaseDate  time.Time
	TestOfferDepartureDate time.Time
	db                     *db.Database
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
	Conditions Conditions `gorm:"column:conditions"`
	Result     Result     `gorm:"column:result"`
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

func NewRevenueFRule(db *db.Database) RevenueRule {
	return RevenueRule{
		db: db,
	}
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
		if priceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			cronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			cronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) {
			return choice.Result
		} else if priceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			cronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			choice.Conditions.PurchaseCronSpec == nil {
			return choice.Result
		} else if priceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			cronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) &&
			choice.Conditions.DepartureCronSpec == nil {
			return choice.Result
		} else if priceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			choice.Conditions.DepartureCronSpec == nil && choice.Conditions.PurchaseCronSpec == nil {
			return choice.Result
		}
	}
	return Result{}
}

func selectMarginRow(choices []ConditionMarginResult, testRule RevenueRule) string {
	for _, choice := range choices {
		if priceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			cronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			cronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) {
			return choice.Result
		} else if priceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			cronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			choice.Conditions.PurchaseCronSpec == nil {
			return choice.Result
		} else if priceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			cronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) &&
			choice.Conditions.DepartureCronSpec == nil {
			return choice.Result
		} else if priceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			choice.Conditions.DepartureCronSpec == nil && choice.Conditions.PurchaseCronSpec == nil {
			return choice.Result
		}
	}
	return ""
}

func (a RevenueRule) GetResultValue(testRule interface{}) interface{} {
	var revenueResult Revenue
	var marginResult Margin
	result := RevenueRuleResult{
		Id: a.Id,
	}
	if a.Revenue != nil && *a.Revenue != "[]" {
		if err := json.Unmarshal([]byte(*a.Revenue), &revenueResult); err != nil {
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
	if a.Margin != nil && *a.Margin != "[]" {
		if err := json.Unmarshal([]byte(*a.Margin), &marginResult); err != nil {
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

func (a RevenueRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
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

func (a RevenueRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{
		"ab_variant": func(a, b reflect.Value) bool {
			if b.IsNil() {
				return false
			}
			offerABCampaigns, ok := b.Elem().Interface().([]string)
			if !ok {
				return false
			}
			return inSlice(a.Elem().Interface().(string), offerABCampaigns)
		},
		"days_to_departure_min": func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) <= b.Elem().Interface().(string)
		},
		"days_to_departure_max": func(a, b reflect.Value) bool {
			return a.Elem().Interface().(string) > b.Elem().Interface().(string)
		},
	}
}

func (a RevenueRule) getStrategyKeys() []string {
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

func (a RevenueRule) getTableName() string {
	return "rm_frule_revenue"
}

func (a RevenueRule) GetDefaultValue() interface{} {
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

func (a RevenueRule) GetLastUpdateTime() time.Time {
	return getLastUpdateTime("revenue", a.db)
}

func (a RevenueRule) GetDataStorage() (map[int][]FRuler, error) {
	result := make(map[int][]FRuler)
	for rank, fieldList := range a.GetComparisonOrder() {
		query := a.db.Table(a.getTableName())
		for _, field := range a.getStrategyKeys() {
			if inSlice(field, fieldList) {
				query = query.Where(field + " IS NOT NULL")
			} else {
				query = query.Where(field + " IS NULL")
			}
		}
		rows, err := query.Rows()
		if err != nil {
			return result, err
		}

		for rows.Next() {
			var rowData RevenueRule

			if err := a.db.ScanRows(rows, &rowData); err != nil {
				return result, err
			}
			result[rank] = append(result[rank], rowData)

		}
		err = rows.Close()
		if err != nil {
			return result, err
		}
	}
	return result, nil
}
