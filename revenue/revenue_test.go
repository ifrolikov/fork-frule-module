package revenue

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/contracts/base"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
	"time"
)

func TestRevenueStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	revenueFRule, err := NewRevenueFRule(ctx, &repository.Config{DataURI: frule_module.GetFilePath("../testdata/revenue.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), revenueFRule)

	dataStorage := revenueFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 0)
	assert.Len(t, (*dataStorage)[13], 2)

	assert.Equal(t, 236, dataStorage.GetMaxRank())
}

func TestRevenueData(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	revenueFRule, err := NewRevenueFRule(ctx, &repository.Config{DataURI: frule_module.GetFilePath("../testdata/revenue.json")})
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
	//fmt.Printf("%+v", result)
	assert.Equal(t, 190, result.(RevenueRuleResult).Id)
	assert.Equal(t, int64(500), result.(RevenueRuleResult).Revenue.Full.Ticket.Amount)
}