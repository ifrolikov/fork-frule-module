package revenue

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/contracts/base"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
	"time"
)

func TestRevenue(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/revenue.json")}
	ctx := context.Background()
	defer ctx.Done()

	revenueFRule, err := NewRevenueFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), revenueFRule)

	dataStorage := revenueFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

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
	fmt.Printf("%+v", result)
	assert.Equal(t, 190, result.(RevenueRuleResult).Id)
	assert.Equal(t, 500, result.(RevenueRuleResult).Revenue.Full.Ticket.Amount)
}
