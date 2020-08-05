package revenue

import (
	"context"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"regexp"
	"stash.tutu.ru/avia-search-common/contracts/v2/base"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
	"time"
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

type RevenueRule struct {
	Id                     int32       `json:"id"`
	CarrierId              *int64      `json:"carrier_id"`
	Partner                *string     `json:"partner"`
	ConnectionGroup        *string     `json:"connection_group"`
	TicketingConnection    *string     `json:"ticketing_connection"`
	DaysToDepartureMin     *int64      `json:"days_to_departure_min"`
	DaysToDepartureMax     *int64      `json:"days_to_departure_max"`
	FareType               *string     `json:"tariff"`
	ABVariant              interface{} `json:"ab_variant"`
	DepartureCountryId     *uint64     `json:"departure_country_id"`
	ArrivalCountryId       *uint64     `json:"arrival_country_id"`
	DepartureCityId        *uint64     `json:"departure_city_id"`
	ArrivalCityId          *uint64     `json:"arrival_city_id"`
	Revenue                *string     `json:"result_revenue"`
	RevenueParsed          *Revenue
	Margin                 *string `json:"result_margin"`
	MarginParsed           *Margin
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

type ConditionRevenueResult struct {
	Conditions Conditions `json:"conditions"`
	Result     Result     `json:"result"`
}

type Conditions struct {
	PriceRange        *string `json:"price_range"`
	DepartureCronSpec *string `json:"departure_cron_spec"`
	PurchaseCronSpec  *string `json:"purchase_cron_spec"`
}

type Result struct {
	Ticket        *string `json:"ticket"`
	TicketParsed  MoneyParsed
	Segment       *string `json:"segment"`
	SegmentParsed MoneyParsed
}

type MoneyParsed struct {
	Money   base.Money
	Percent float64
}

type Margin struct {
	Full   []ConditionMarginResult `json:"full"`
	Child  []ConditionMarginResult `json:"child"`
	Infant []ConditionMarginResult `json:"infant"`
}

type ConditionMarginResult struct {
	Conditions   Conditions `json:"conditions"`
	Result       string     `json:"result"`
	ResultParsed base.Money
}

type RevenueRuleResult struct {
	Id      int32
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

var moneyAmountSpec = regexp.MustCompile(`([0-9\.]+)([A-Z]+)`)
var moneyPercentSpec = regexp.MustCompile(`([\-0-9\.]+)\%`)

func parseMoneySpec(spec *string) MoneyParsed {
	if spec != nil {
		if parsedData := moneyAmountSpec.FindStringSubmatch(*spec); len(parsedData) == 3 {
			amount, err := strconv.ParseFloat(parsedData[1], 64)
			if err != nil {
				log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot parse money %s", *spec)).Msg("parsing revenue amount")
			}
			return MoneyParsed{
				Money: base.Money{
					Amount:   int64(math.Round(amount * 100)),
					Currency: &base.Currency{Code: parsedData[2], Fraction: 100},
				},
			}
		} else if parsedData := moneyPercentSpec.FindStringSubmatch(*spec); len(parsedData) == 2 {
			percent, err := strconv.ParseFloat(parsedData[1], 64)
			if err != nil {
				log.Logger.Error().Stack().Err(errors.Wrapf(err, "cannot parse money %s", *spec)).Msg("parsing revenue percent")
			}
			return MoneyParsed{
				Money:   *base.CreateZeroRubMoney(),
				Percent: percent,
			}
		}
	}
	return MoneyParsed{Money: *base.CreateZeroRubMoney()}
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

func selectMarginRow(choices []ConditionMarginResult, testRule RevenueRule) base.Money {
	for _, choice := range choices {
		if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			frule_module.CronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) {
			return choice.ResultParsed
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.DepartureCronSpec, testRule.TestOfferDepartureDate) &&
			choice.Conditions.PurchaseCronSpec == nil {
			return choice.ResultParsed
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			frule_module.CronSpec(choice.Conditions.PurchaseCronSpec, testRule.TestOfferPurchaseDate) &&
			choice.Conditions.DepartureCronSpec == nil {
			return choice.ResultParsed
		} else if frule_module.PriceRange(choice.Conditions.PriceRange, testRule.TestOfferPrice) &&
			choice.Conditions.DepartureCronSpec == nil && choice.Conditions.PurchaseCronSpec == nil {
			return choice.ResultParsed
		}
	}
	return base.Money{Amount: 0, Currency: &base.Currency{Fraction: 100, Code: "RUB"}}
}

func (rule *RevenueRule) GetResultValue(testRule interface{}) interface{} {
	result := RevenueRuleResult{
		Id: rule.Id,
	}
	if rule.RevenueParsed != nil {
		revenueParams := testRule.(RevenueRule)
		fullResult := selectRevenueRow(rule.RevenueParsed.Full, revenueParams)
		result.Revenue.Full.Ticket = rule.CalculateRevenue(fullResult.TicketParsed, revenueParams.TestOfferPrice)
		result.Revenue.Full.Segment = rule.CalculateRevenue(fullResult.SegmentParsed, revenueParams.TestOfferPrice)

		childResult := selectRevenueRow(rule.RevenueParsed.Child, revenueParams)
		result.Revenue.Child.Ticket = rule.CalculateRevenue(childResult.TicketParsed, revenueParams.TestOfferPrice)
		result.Revenue.Child.Segment = rule.CalculateRevenue(childResult.SegmentParsed, revenueParams.TestOfferPrice)

		infantResult := selectRevenueRow(rule.RevenueParsed.Infant, revenueParams)
		result.Revenue.Infant.Ticket = rule.CalculateRevenue(infantResult.TicketParsed, revenueParams.TestOfferPrice)
		result.Revenue.Infant.Segment = rule.CalculateRevenue(infantResult.SegmentParsed, revenueParams.TestOfferPrice)
	}
	if rule.MarginParsed != nil {
		result.Margin.Full = selectMarginRow(rule.MarginParsed.Full, testRule.(RevenueRule))
		result.Margin.Child = selectMarginRow(rule.MarginParsed.Child, testRule.(RevenueRule))
		result.Margin.Infant = selectMarginRow(rule.MarginParsed.Infant, testRule.(RevenueRule))
	}
	return result
}

func (rule *RevenueRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

func (rule *RevenueRule) CalculateRevenue(moneyParsed MoneyParsed, price base.Money) base.Money {
	if moneyParsed.Percent != 0 && price.Validate() {
		money := base.CloneMoney(&price)
		money.MultiplyFloat64(moneyParsed.Percent / 100)
		return *money
	} else {
		return moneyParsed.Money
	}
}

func (rule *RevenueRule) GetComparisonOrder() frule_module.ComparisonOrder {
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

func (rule *RevenueRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

func (rule *RevenueRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *RevenueRule) GetDefaultValue() interface{} {
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
				Ticket:  *base.CreateZeroRubMoney(),
				Segment: *base.CreateZeroRubMoney(),
			},
			Child: struct {
				Ticket  base.Money
				Segment base.Money
			}{
				Ticket:  *base.CreateZeroRubMoney(),
				Segment: *base.CreateZeroRubMoney(),
			},
			Infant: struct {
				Ticket  base.Money
				Segment base.Money
			}{
				Ticket:  *base.CreateZeroRubMoney(),
				Segment: *base.CreateZeroRubMoney(),
			},
		},
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

func (rule *RevenueRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *RevenueRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *RevenueRule) GetRuleName() string {
	return "Revenue"
}
