package service_charge

import (
	"context"
	"fmt"
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
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
		TestOfferPrice:  base.Money{Amount: 60000},
	}
	result := frule.GetResult(params)
	fmt.Printf("1 %+v", result)
	fmt.Println()
	assert.EqualValues(t, 74, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 21137, Currency: &base.Currency{Code: "RUR", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 21137, Currency: &base.Currency{Code: "RUR", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	partner = "new_tt"
	connectionGroup = "galileo"
	carrierId = int64(1062)
	departureCityId = uint64(21)
	arrivalCityId = uint64(100)
	params = ServiceChargeRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
		TestOfferPrice:  base.Money{Amount: 700000},
	}
	result = frule.GetResult(params)
	fmt.Printf("2 %+v", result)
	fmt.Println()
	assert.EqualValues(t, 2000, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 49428, Currency: &base.Currency{Code: "RUR", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 49428, Currency: &base.Currency{Code: "RUR", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	departureCityId = uint64(34)
	params = ServiceChargeRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
		TestOfferPrice:  base.Money{Amount: 1000000},
	}
	result = frule.GetResult(params)
	fmt.Printf("3 %+v", result)
	fmt.Println()
	assert.EqualValues(t, 2000, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 58344, Currency: &base.Currency{Code: "RUR", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 58344, Currency: &base.Currency{Code: "RUR", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	departureCityId = uint64(34)
	fareType := "subsidy"
	params = ServiceChargeRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		FareType:        &fareType,
	}
	result = frule.GetResult(params)
	fmt.Printf("4 %+v", result)
	fmt.Println()
	assert.EqualValues(t, 4179, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUR", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUR", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)
}
