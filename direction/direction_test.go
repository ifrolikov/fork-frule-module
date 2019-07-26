package direction

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestDirectionStorage(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/direction.json")}
	ctx := context.Background()
	defer ctx.Done()

	frule, err := NewDirectionFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), frule)

	dataStorage := frule.GetDataStorage()
	assert.NotNil(t, dataStorage)
	assert.Len(t, (*dataStorage)[0], 2)
	assert.Len(t, (*dataStorage)[1], 1)

	maxKey := 0
	for key := range *dataStorage {
		if key > maxKey {
			maxKey = key
		}
	}
	assert.Equal(t, maxKey, 7)
}

func TestDirectionData(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/direction.json")}
	ctx := context.Background()
	defer ctx.Done()

	directionFRule, err := NewDirectionFRule(ctx, testConfig)
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