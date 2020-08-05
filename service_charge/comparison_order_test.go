package service_charge

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestComparisonOrder(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	genetatedControlOrder := frule_module.ComparisonOrder{
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
	serviceChargeRule, _ := NewServiceChargeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/service_charge.json")})

	assert.EqualValues(t, genetatedControlOrder, serviceChargeRule.GetComparisonOrder())
}