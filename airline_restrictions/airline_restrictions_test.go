package airline_restrictions

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
	"time"
)

func TestAirlineRestrictionsStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	airlineRestrictionsFRule, err := NewAirlineRestrictionsFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/airline_restrictions.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), airlineRestrictionsFRule)

	dataStorage := airlineRestrictionsFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[1], 2)
	assert.Len(t, (*dataStorage)[133], 1)
	assert.Len(t, (*dataStorage)[1151], 1)
	assert.Len(t, (*dataStorage)[1071], 1)

	assert.Equal(t, 1151, dataStorage.GetMaxRank())
}

func TestAirlineRestrictionsResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	airlineRestrictionsFRule, err := NewAirlineRestrictionsFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/airline_restrictions.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, airlineRestrictionsFRule)
	assert.NotNil(t, frule)

	//общее разрешающее правило
	currentTime := time.Now().Format("2006-01-02")
	departureDateFrom := currentTime
	departureDateTo := currentTime
	assert.True(t, frule.GetResult(AirlineRestrictionsRule{DepartureDateFrom: &departureDateFrom, DepartureDateTo: &departureDateTo}).(bool))

	partner := "new_tt"
	gds := "galileo"
	connectionGroup := "galileo"
	platingCarrierId := int64(1106)
	marketingCarrierId := int64(1062)
	operatingCarrierId := int64(1062)
	departureCountryId := uint64(8)
	departureCityId := uint64(467)
	arrivalCountryId := uint64(7)
	arrivalCityId := uint64(491)
	departurePeriodFrom := int64(5)
	departurePeriodTo := int64(5)
	assert.False(t, frule.GetResult(AirlineRestrictionsRule{
		DepartureDateFrom:   &departureDateFrom,
		DepartureDateTo:     &departureDateTo,
		Partner:             &partner,
		Gds:                 &gds,
		ConnectionGroup:     &connectionGroup,
		PlatingCarrierId:    &platingCarrierId,
		MarketingCarrierId:  &marketingCarrierId,
		OperatingCarrierId:  &operatingCarrierId,
		DepartureCountryId:  &departureCountryId,
		DepartureCityId:     &departureCityId,
		ArrivalCountryId:    &arrivalCountryId,
		ArrivalCityId:       &arrivalCityId,
		DeparturePeriodFrom: &departurePeriodFrom,
		DeparturePeriodTo:   &departurePeriodTo,
	}).(bool))

	partner = "new_tt"
	gds = "galileo"
	connectionGroup = "galileo"
	platingCarrierId = int64(1062)
	marketingCarrierId = int64(6)
	operatingCarrierId = int64(6)
	departureCountryId = uint64(8)
	departureCityId = uint64(467)
	arrivalCountryId = uint64(7)
	arrivalCityId = uint64(491)
	departurePeriodFrom = int64(5)
	departurePeriodTo = int64(5)
	assert.True(t, frule.GetResult(AirlineRestrictionsRule{
		DepartureDateFrom:   &departureDateFrom,
		DepartureDateTo:     &departureDateTo,
		Partner:             &partner,
		Gds:                 &gds,
		ConnectionGroup:     &connectionGroup,
		PlatingCarrierId:    &platingCarrierId,
		MarketingCarrierId:  &marketingCarrierId,
		OperatingCarrierId:  &operatingCarrierId,
		DepartureCountryId:  &departureCountryId,
		DepartureCityId:     &departureCityId,
		ArrivalCountryId:    &arrivalCountryId,
		ArrivalCityId:       &arrivalCityId,
		DeparturePeriodFrom: &departurePeriodFrom,
		DeparturePeriodTo:   &departurePeriodTo,
	}).(bool))

	// проверка что нельзя S7 из sabre раньше чем через 2 недели
	partner = "new_tt"
	gds = "sabre"
	connectionGroup = "sabre"
	platingCarrierId = int64(1106)
	marketingCarrierId = int64(1106)
	operatingCarrierId = int64(1106)
	departureCountryId = uint64(8)
	departureCityId = uint64(467)
	arrivalCountryId = uint64(7)
	arrivalCityId = uint64(491)
	departurePeriodFrom = int64(5)
	departurePeriodTo = int64(5)
	assert.False(t, frule.GetResult(AirlineRestrictionsRule{
		DepartureDateFrom:   &departureDateFrom,
		DepartureDateTo:     &departureDateTo,
		Partner:             &partner,
		Gds:                 &gds,
		ConnectionGroup:     &connectionGroup,
		PlatingCarrierId:    &platingCarrierId,
		MarketingCarrierId:  &marketingCarrierId,
		OperatingCarrierId:  &operatingCarrierId,
		DepartureCountryId:  &departureCountryId,
		DepartureCityId:     &departureCityId,
		ArrivalCountryId:    &arrivalCountryId,
		ArrivalCityId:       &arrivalCityId,
		DeparturePeriodFrom: &departurePeriodFrom,
		DeparturePeriodTo:   &departurePeriodTo,
	}).(bool))

	departureDateFrom = "1970-11-23"
	assert.Equal(t, airlineRestrictionsFRule.GetDefaultValue(), frule.GetResult(AirlineRestrictionsRule{DepartureDateFrom: &departureDateFrom}).(bool))
}
