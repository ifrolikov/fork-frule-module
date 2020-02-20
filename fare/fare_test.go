package fare

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestFareStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	fareFRule, err := NewFareFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/fare.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), fareFRule)

	dataStorage := fareFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)
	assert.Len(t, (*dataStorage)[1], 1)
	assert.Len(t, (*dataStorage)[33], 3)

	assert.Equal(t, 35, dataStorage.GetMaxRank())
}

func TestFareResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	fareFRule, err := NewFareFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/fare.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, fareFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierId := int64(73)
	departureCityId := uint64(491)
	arrivalCityId := uint64(75)
	fareSpec := "TEST"
	fareAccessGroup := "test_access_group"

	assert.Equal(t, "closed", frule.GetResult(FareRule{Partner: &partner, CarrierId: &carrierId, FareSpec: &fareSpec, FareAccessGroup: nil}))
	assert.Equal(t, "", frule.GetResult(FareRule{Partner: &partner, CarrierId: &carrierId, FareSpec: &fareSpec, FareAccessGroup: &fareAccessGroup}))

	assert.Equal(t, fareFRule.GetDefaultValue(), frule.GetResult(FareRule{Partner: &partner, CarrierId: &carrierId}))

	assert.NotEqual(t, "subsidy", frule.GetResult(FareRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
		FareSpec:        &fareSpec,
	}))

	fareSpec = "TBTPACRAN"
	assert.Equal(t, "subsidy", frule.GetResult(FareRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
		FareSpec:        &fareSpec,
	}))

	carrierId = int64(19)
	fareSpec = "TTARUF"
	assert.Equal(t, "closed", frule.GetResult(FareRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		FareSpec:        &fareSpec,
	}))
}
