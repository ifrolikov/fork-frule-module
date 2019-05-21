package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
)

func main() {
	db, err := gorm.Open("mysql", "webuser:qazxswedc@tcp(devel.mysql.devel.tutu.ru:3306)/devel")
	if err != nil {
		fmt.Println(err)
	}
	airline := frule_module.NewFRule(frule_module.NewAirlineFRule(db))
	partner := "iata"
	carrierId := 213
	connectionGroup := "galileo"
	fmt.Println(airline.GetResult(frule_module.AirlineRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup}))
}
