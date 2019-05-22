package main

import (
	"context"
	"fmt"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/golang/resources/db/mysql"
)

func main() {

	db := mysql.NewDb()
	err := db.Init()
	if err != nil {
		fmt.Println(err)
	}
	// airline := frule_module.NewFRule(frule_module.NewAirlineFRule(db))
	partner := "iata"
	carrierId := 1116
	connectionGroup := "galileo"
	// fmt.Println(airline.GetResult(frule_module.AirlineRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup}))

	ctx := context.Background()

	partnerPercent := frule_module.NewFRule(ctx, frule_module.NewPartnerPercentFRule(db))

	partner = "new_tt"
	connectionGroup = "sig23_direct"
	countryId := 7
	from := "2019-05-03"
	to := "2019-05-03"
	fareType := "subsidy"

	result := partnerPercent.GetResult(frule_module.PartnerPercentRule{
		Partner:            &partner,
		ConnectionGroup:    &connectionGroup,
		CarrierId:          &carrierId,
		CarrierCountryId:   &countryId,
		DateOfPurchaseFrom: &from,
		DateOfPurchaseTo:   &to,
		FareType:           &fareType,
	})
	fmt.Println(result.(float64))
	a := make(chan struct{})
	<-a
}
