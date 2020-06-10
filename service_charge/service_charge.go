package service_charge

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
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
	Version             int32       `json:"version"`
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
	CurrencyConvertor   base.CurrencyConvertor
	repo                *frule_module.Repository
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

func calculatePassengerServiceCharge(moneyParsed MoneyParsed, price base.Money, currencyConvertor base.CurrencyConvertor) base.Money {
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
		if currencyConvertor != nil {
			if err := passengerServiceCharge.ConvertAndAdd(currencyConvertor, moneyParsed.Money); err != nil {
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
			if currencyConvertor != nil {
				if moreThanLimit, _ := tariffAddition.ConvertAndCompare(currencyConvertor, moneyParsed.Limit); moreThanLimit == 1 {
					tariffAddition.Copy(moneyParsed.Limit)
				}
			} else {
				if moreThanLimit, _ := tariffAddition.More(moneyParsed.Limit); moreThanLimit {
					tariffAddition.Copy(moneyParsed.Limit)
				}
			}
		}

		if currencyConvertor != nil {
			if err := passengerServiceCharge.ConvertAndAdd(currencyConvertor, tariffAddition); err != nil {
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
			serviceChargeParams.CurrencyConvertor,
		)

		result.Margin.Child = calculatePassengerServiceCharge(
			findPricingRangeValue(rule.MarginParsed.Child, serviceChargeParams),
			serviceChargeParams.TestOfferPrice,
			serviceChargeParams.CurrencyConvertor,
		)

		result.Margin.Infant = calculatePassengerServiceCharge(
			findPricingRangeValue(rule.MarginParsed.Infant, serviceChargeParams),
			serviceChargeParams.TestOfferPrice,
			serviceChargeParams.CurrencyConvertor,
		)
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
