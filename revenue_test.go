package frule_module

import (
	"context"
	"fmt"
	"stash.tutu.ru/avia-search-common/contracts/base"
	"stash.tutu.ru/golang/resources/db"
	"stash.tutu.ru/golang/resources/db/mysql"
	"testing"
	"time"
)

func TestRevenue(t *testing.T) {
	database := mysql.NewDb()
	database.WithConfig(db.Config{
		DSN:   "webuser:qazxswedc@tcp(devel-02.mysql.avia.devel.tutu.ru:3306)/devel",
		Debug: true,
	})
	err := database.Init()
	if err != nil {
		t.Fatal(err)
	}

	partner := "fake"
	connectionGroup := "fake"
	carrierId := 1109
	testRule := RevenueRule{
		Partner:                &partner,
		ConnectionGroup:        &connectionGroup,
		CarrierId:              &carrierId,
		TestOfferPrice:         base.Money{Amount: 6000},
		TestOfferPurchaseDate:  time.Now(),
		TestOfferDepartureDate: time.Now(),
	}
	ctx := context.Background()
	revenueFrule := NewFRule(ctx, NewRevenueFRule(database))
	result := revenueFrule.GetResult(testRule)
	fmt.Printf("%+v", result)
	if result.(RevenueRuleResult).Id != 190 {
		t.Error("Wrong rule selected")
	}
	if result.(RevenueRuleResult).Revenue.Full.Ticket.Amount != 500 {
		t.Error("Wrong revenue for full passenger")
	}

}
