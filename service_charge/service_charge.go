package service_charge

import (
	"context"
	"math"
	"reflect"
	"regexp"
	"stash.tutu.ru/avia-search-common/contracts/v2/base"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
)

var comparisonOrder = frule_module.ComparisonOrder{
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

var strategyKeys = []string{
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

type ServiceChargeRule struct {
	Id                  int32       `json:"id"`
	CarrierId           *int64      `json:"carrier_id"`
	Partner             *string     `json:"partner"`
	ConnectionGroup     *string     `json:"connection_group"`
	TicketingConnection *string     `json:"ticketing_connection"`
	DaysToDepartureMin  *int64      `json:"days_to_departure_min"`
	DaysToDepartureMax  *int64      `json:"days_to_departure_max"`
	FareType            *string     `json:"tariff"`
	ABVariant           interface{} `json:"ab_variant"`
	DepartureCountryId  *uint64     `json:"departure_country_id"`
	ArrivalCountryId    *uint64     `json:"arrival_country_id"`
	DepartureCityId     *uint64     `json:"departure_city_id"`
	ArrivalCityId       *uint64     `json:"arrival_city_id"`
	Margin              *string     `json:"result_margin"`
	MarginParsed        *Margin
	TestOfferPrice      base.Money
	repo                *frule_module.Repository
}

type Conditions struct {
	PriceRange *string `json:"price_range"`
}

type Result struct {
	Ticket  *string `json:"ticket"`
	Segment *string `json:"segment"`
}

type ConditionMarginResult struct {
	Conditions   Conditions `json:"conditions"`
	Result       string     `json:"result"`
	ResultParsed base.Money
}

type Margin struct {
	Full   []ConditionMarginResult `json:"full"`
	Child  []ConditionMarginResult `json:"child"`
	Infant []ConditionMarginResult `json:"infant"`
}

type ServiceChargeRuleResult struct {
	Id     int32
	Margin struct {
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

var moneySpec = regexp.MustCompile(`([0-9\.]+)([A-Z]+)`)

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
		amount, err := strconv.ParseFloat(parsedData[1], 64)
		if err != nil {
			log.Logger.Error().Stack().Err(err).Msg("parsing money")
		}
		return base.Money{
			Amount: int64(math.Round(amount * 100)), // TODO вынести defaultFraction
			Currency: &base.Currency{ // TODO: load from DB by code
				Code:     parsedData[2],
				Fraction: 100,
			},
		}
	} else {
		return base.Money{
			Currency: &base.Currency{ // TODO нужна factory для Currency??
				Code:     "RUB",
				Fraction: 100,
			},
		}
	}
}

func selectMarginRow(choices []ConditionMarginResult, testRule ServiceChargeRule) base.Money {
	for _, choice := range choices {
		if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) {
			return choice.ResultParsed
		}
	}
	return base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
}

func (rule *ServiceChargeRule) GetResultValue(testRule interface{}) interface{} {
	result := ServiceChargeRuleResult{
		Id: rule.Id,
	}
	if rule.MarginParsed != nil {
		result.Margin.Full = selectMarginRow(rule.MarginParsed.Full, testRule.(ServiceChargeRule))
		result.Margin.Child = selectMarginRow(rule.MarginParsed.Child, testRule.(ServiceChargeRule))
		result.Margin.Infant = selectMarginRow(rule.MarginParsed.Infant, testRule.(ServiceChargeRule))
	}
	return result
}

func (rule *ServiceChargeRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{
	"ab_variant": func(a, b reflect.Value) bool {
		offerABCampaigns, ok := b.Elem().Interface().([]string)
		if !ok {
			return false
		}
		return frule_module.InSlice(a.Elem().Interface().(string), offerABCampaigns)
	},
	"days_to_departure_min": func(a, b reflect.Value) bool {
		return a.Elem().Interface().(int64) <= b.Elem().Interface().(int64)
	},
	"days_to_departure_max": func(a, b reflect.Value) bool {
		return a.Elem().Interface().(int64) > b.Elem().Interface().(int64)
	},
}

func (rule *ServiceChargeRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

func (rule *ServiceChargeRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *ServiceChargeRule) GetDefaultValue() interface{} {
	zeroRub := base.Money{
		Currency: &base.Currency{
			Fraction: 100,
			Code:     "RUB",
		},
	}
	return ServiceChargeRuleResult{
		Id: -1,
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

func (rule *ServiceChargeRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *ServiceChargeRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *ServiceChargeRule) GetRuleName() string {
	return "ServiceCharge"
}
