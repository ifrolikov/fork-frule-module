package revenue

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/contracts/v2/base"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
	"time"
)

func TestRevenueStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	revenueFRule, err := NewRevenueFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/revenue.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), revenueFRule)

	dataStorage := revenueFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 0)
	assert.Len(t, (*dataStorage)[13], 2)

	assert.Equal(t, 242, dataStorage.GetMaxRank())
}

func TestRevenueData(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	revenueFRule, err := NewRevenueFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/revenue.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, revenueFRule)
	assert.NotNil(t, frule)

	partner := "fake"
	connectionGroup := "fake"
	carrierId := int64(1109)
	params := RevenueRule{
		Partner:                &partner,
		ConnectionGroup:        &connectionGroup,
		CarrierId:              &carrierId,
		TestOfferPrice:         base.Money{Amount: 6000},
		TestOfferPurchaseDate:  time.Now(),
		TestOfferDepartureDate: time.Now(),
	}
	result := frule.GetResult(params)
	assert.EqualValues(t, 190, result.(RevenueRuleResult).Id)
	assert.EqualValues(t, 50000, result.(RevenueRuleResult).Revenue.Full.Ticket.Amount)
	assert.EqualValues(t, 40000, result.(RevenueRuleResult).Revenue.Child.Ticket.Amount)
	assert.EqualValues(t, 20000, result.(RevenueRuleResult).Revenue.Infant.Ticket.Amount)

	partner = "new_tt"
	connectionGroup = "galileo"
	carrierId = int64(1062)
	fareType := "subsidy"
	departureCityId := uint64(21)
	arrivalCityId := uint64(100)
	params = RevenueRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		FareType:        &fareType,
		TestOfferPrice:  base.Money{Amount: 700000, Currency: &base.Currency{Code: "RUR", Fraction: 100}},
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
	}
	result = frule.GetResult(params)
	assert.EqualValues(t, 86012, result.(RevenueRuleResult).Id)
	assert.EqualValues(t, 21000, result.(RevenueRuleResult).Revenue.Full.Ticket.Amount)
	assert.EqualValues(t, 20000, result.(RevenueRuleResult).Revenue.Full.Segment.Amount)
	assert.EqualValues(t, -14000, result.(RevenueRuleResult).Revenue.Child.Ticket.Amount)
	assert.EqualValues(t, 10000, result.(RevenueRuleResult).Revenue.Child.Segment.Amount)
	assert.EqualValues(t, 5000, result.(RevenueRuleResult).Revenue.Infant.Ticket.Amount)
	assert.EqualValues(t, 1000, result.(RevenueRuleResult).Revenue.Infant.Segment.Amount)

	departureCityId = uint64(34)
	params = RevenueRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		FareType:        &fareType,
		DepartureCityId: &departureCityId,
		ArrivalCityId:   &arrivalCityId,
	}
	result = frule.GetResult(params)
	assert.EqualValues(t, 86013, result.(RevenueRuleResult).Id)

	departureCityId = uint64(34)
	params = RevenueRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierId:       &carrierId,
		FareType:        &fareType,
	}
	result = frule.GetResult(params)
	//fmt.Printf("%+v", result)
	assert.EqualValues(t, 86682, result.(RevenueRuleResult).Id)
	assert.EqualValues(t, 30000, result.(RevenueRuleResult).Revenue.Full.Ticket.Amount)
	assert.EqualValues(t, 30000, result.(RevenueRuleResult).Revenue.Child.Ticket.Amount)
	assert.EqualValues(t, 0, result.(RevenueRuleResult).Revenue.Infant.Ticket.Amount)
}
