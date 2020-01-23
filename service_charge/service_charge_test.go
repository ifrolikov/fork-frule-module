package service_charge

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/contracts/v2/base"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestServiceChargeStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	serviceChargeRule, err := NewServiceChargeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/service_charge.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), serviceChargeRule)

	dataStorage := serviceChargeRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 0)
	assert.Len(t, (*dataStorage)[26], 3)

	assert.Equal(t, 242, dataStorage.GetMaxRank())
}

func TestServiceChargeData(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	serviceChargeRule, err := NewServiceChargeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/service_charge.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, serviceChargeRule)
	assert.NotNil(t, frule)

	partner := "nba"
	connectionGroup := "nba"
	carrierId := int64(1062)
	departureCityId := uint64(29)
	arrivalCityId := uint64(491)
	params := ServiceChargeRule{
		Partner:                &partner,
		ConnectionGroup:        &connectionGroup,
		CarrierId:              &carrierId,
		DepartureCityId:        &departureCityId,
		ArrivalCityId:          &arrivalCityId,
		TestOfferPrice:         base.Money{Amount: 6000},
	}
	result := frule.GetResult(params)
	assert.EqualValues(t, 74, result.(ServiceChargeRuleResult).Id)

	partner = "new_tt"
	connectionGroup = "galileo"
	carrierId = int64(1062)
	fareType := "subsidy"
	departureCityId = uint64(21)
	arrivalCityId = uint64(100)
	params = ServiceChargeRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		FareType:        &fareType,
		TestOfferPrice:  base.Money{Amount: 7000},
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
	}
	result = frule.GetResult(params)
	assert.EqualValues(t, 4179, result.(ServiceChargeRuleResult).Id)

	departureCityId = uint64(34)
	params = ServiceChargeRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		FareType:        &fareType,
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
	}
	result = frule.GetResult(params)
	assert.EqualValues(t, 4179, result.(ServiceChargeRuleResult).Id)

	departureCityId = uint64(34)
	params = ServiceChargeRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		FareType:        &fareType,
	}
	result = frule.GetResult(params)
	//fmt.Printf("%+v", result)
	assert.EqualValues(t, 4179, result.(ServiceChargeRuleResult).Id)
}