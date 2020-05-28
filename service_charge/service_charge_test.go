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

func TestServiceChargeSimpleFormat(t *testing.T) {
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
	assert.EqualValues(t, 74, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 21137, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 21137, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
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
	assert.EqualValues(t, 2000, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 49428, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 49428, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
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
	assert.EqualValues(t, 2000, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 58344, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 58344, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
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
	assert.EqualValues(t, 4179, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)
}

func TestServiceChargeComplexFormat(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	serviceChargeRule, err := NewServiceChargeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/service_charge.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, serviceChargeRule)
	assert.NotNil(t, frule)

	partner := "nba"
	connectionGroup := "nba"
	carrierId := int64(1011)
	departureCityId := uint64(39)
	arrivalCityId := uint64(491)
	params := ServiceChargeRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
		TestOfferPrice:  base.Money{Amount: 60000, Currency: &base.Currency{Code: "RUB", Fraction: 100}},
	}
	result := frule.GetResult(params)
/*	fmt.Printf("1 %+v", result)
	fmt.Println()*/
	assert.EqualValues(t, 87, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 18840, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 18840, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	params.TestOfferPrice = base.Money{Amount: 167500, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
	result = frule.GetResult(params)
	// 387.7RUR+2.3% = 38770 + 167500/100*2,3 = 38770 + 3852,5 = 42622,5
	assert.Equal(t, base.Money{Amount: 42623, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	// 387.7RUR+2.1% = 38770 + 167500/100*2,1 = 38770 + 3517,5 = 42287,5
	assert.Equal(t, base.Money{Amount: 42288, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	params.TestOfferPrice = base.Money{Amount: 234300, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
	result = frule.GetResult(params)
	// 461.13RUR+2.3%<50.1RUR = 46113 + 234300/100*2,3<5010 = 46113 + 5010 = 51123
	assert.Equal(t, base.Money{Amount: 51123, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	// 461.13RUR+2.1%<50.1RUR = 46113 + 234300/100*2,1<5010 = 46113 + 4920,3 = 51033,3
	assert.Equal(t, base.Money{Amount: 51033, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	params.TestOfferPrice = base.Money{Amount: 315721, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
	result = frule.GetResult(params)
	// 0RUR+2.3%<67 = 315721/100*2,3<6700 = 6700
	assert.Equal(t, base.Money{Amount: 6700, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	// 0RUR+2.1%<67 = 315721/100*2,1<6700 = 6630,141
	assert.Equal(t, base.Money{Amount: 6630, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	params.TestOfferPrice = base.Money{Amount: 415722, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
	result = frule.GetResult(params)
	// 0RUR+2.3% = 415722/100*2,3 = 9561,606
	assert.Equal(t, base.Money{Amount: 9562, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	// 0RUR+2.1% = 415722/100*2,1 = 8730,162
	assert.Equal(t, base.Money{Amount: 8730, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)
}
