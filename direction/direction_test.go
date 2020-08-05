package direction

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestDirectionStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	directionFRule, err := NewDirectionFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/direction.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), directionFRule)

	dataStorage := directionFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 2)
	assert.Len(t, (*dataStorage)[1], 1)

	assert.Equal(t, 35, dataStorage.GetMaxRank())
}

func TestDirectionData(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	directionFRule, err := NewDirectionFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/direction.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, directionFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	assert.True(t, frule.GetResult(DirectionRule{Partner: &partner}).(bool))

	connectionGroup := "galileo"
	carrierId := int64(1106)
	departureCountryId := uint64(7)
	departureCityId := uint64(491)
	arrivalCountryId := uint64(7)
	arrivalCityId := uint64(10105)

	assert.True(t,
		frule.GetResult(
			DirectionRule{
				Partner:            &partner,
				ConnectionGroup:    &connectionGroup,
				CarrierId:          &carrierId,
				DepartureCountryId: &departureCountryId,
				DepartureCityId:    &departureCityId,
				ArrivalCountryId:   &arrivalCountryId,
				ArrivalCityId:      &arrivalCityId,
			}).(bool))

	partner = "tt"
	assert.False(t, frule.GetResult(DirectionRule{Partner: &partner}).(bool))

	assert.False(t,
		frule.GetResult(
			DirectionRule{
				Partner:            &partner,
				ConnectionGroup:    &connectionGroup,
				CarrierId:          &carrierId,
				DepartureCountryId: &departureCountryId,
				DepartureCityId:    &departureCityId,
				ArrivalCountryId:   &arrivalCountryId,
				ArrivalCityId:      &arrivalCityId,
			}).(bool))

	partner = "unknown"
	assert.Equal(t, directionFRule.GetDefaultValue(), frule.GetResult(DirectionRule{Partner: &partner}).(bool))
}