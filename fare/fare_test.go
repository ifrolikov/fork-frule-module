package fare

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestFareStorage(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/fare.json")}
	ctx := context.Background()
	defer ctx.Done()

	frule, err := NewFareFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), frule)

	dataStorage := frule.GetDataStorage()
	assert.NotNil(t, dataStorage)
	assert.Len(t, (*dataStorage)[0], 1)
	assert.Len(t, (*dataStorage)[16], 3)

	maxKey := 0
	for key := range *dataStorage {
		if key > maxKey {
			maxKey = key
		}
	}
	assert.Equal(t, 17, maxKey)
}

func TestFareResult(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/fare.json")}
	ctx := context.Background()
	defer ctx.Done()

	fareFRule, err := NewFareFRule(ctx, testConfig)
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, fareFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierId := int64(73)
	departureCityId := uint64(491)
	arrivalCityId := uint64(75)
	fareSpec := "TEST"

	assert.Equal(t, "closed", frule.GetResult(FareRule{Partner: &partner, CarrierId: &carrierId, FareSpec: &fareSpec}))

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
