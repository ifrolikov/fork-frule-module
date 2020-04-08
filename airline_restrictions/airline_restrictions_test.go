package airline_restrictions

import (
"context"
"github.com/stretchr/testify/assert"
"stash.tutu.ru/avia-search-common/frule-module"
"stash.tutu.ru/avia-search-common/repository"
"stash.tutu.ru/avia-search-common/utils/system"
"testing"
	"time"
)

func TestAirlineRestrictionStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	airlineRestrictionFRule, err := NewAirlineRestrictionFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/airline_restriction.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), airlineRestrictionFRule)

	dataStorage := airlineRestrictionFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 3)
	assert.Len(t, (*dataStorage)[3], 1)

	assert.Equal(t, 3, dataStorage.GetMaxRank())
}

func TestAirlineRestrictionResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	airlineRestrictionFRule, err := NewAirlineRestrictionFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/airline_restriction.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, airlineRestrictionFRule)
	assert.NotNil(t, frule)

	currentTime := time.Now().Format("2006-01-02")
	purchaseDateFrom := currentTime
	purchaseDateTo := currentTime
	assert.True(t, frule.GetResult(AirlineRestrictionsRule{PurchaseDateFrom: &purchaseDateFrom, PurchaseDateTo: &purchaseDateTo}).(bool))

	partner := "new_tt"
	gds := "galileo"
	validatingCarrierId := int64(6)
	marketingCarrierId := int64(6)
	operatingCarrierId := int64(6)
	departureCountryId := uint64(7)
	departureCityId := uint64(491)
	arrivalCountryId := uint64(7)
	arrivalCityId := uint64(10105)
	purchasePeriodFrom := int64(5)
	purchasePeriodTo := int64(5)
	assert.True(t, frule.GetResult(AirlineRestrictionsRule{
		PurchaseDateFrom: &purchaseDateFrom,
		PurchaseDateTo: &purchaseDateTo,
		Partner: &partner,
		Gds: &gds,
		ValidatingCarrierId: &validatingCarrierId,
		MarketingCarrierId: &marketingCarrierId,
		OperatingCarrierId: &operatingCarrierId,
		DepartureCountryId: &departureCountryId,
		DepartureCityId: &departureCityId,
		ArrivalCountryId: &arrivalCountryId,
		ArrivalCityId: &arrivalCityId,
		PurchasePeriodFrom: &purchasePeriodFrom,
		PurchasePeriodTo: &purchasePeriodTo,
	}).(bool))

	partner = "new_tt"
	gds = "galileo"
	validatingCarrierId = int64(7)
	assert.False(t, frule.GetResult(AirlineRestrictionsRule{Partner: &partner, Gds: &gds, ValidatingCarrierId: &validatingCarrierId}).(bool))

	partner = "new_tt"
	gds = "sabre"
	validatingCarrierId = int64(1111)
	assert.True(t, frule.GetResult(AirlineRestrictionsRule{Partner: &partner, Gds: &gds, ValidatingCarrierId: &validatingCarrierId}).(bool))

	purchaseDateFrom = "1970-11-23"
	assert.Equal(t, airlineRestrictionFRule.GetDefaultValue(), frule.GetResult(AirlineRestrictionsRule{PurchaseDateFrom: &purchaseDateFrom}).(bool))
}

