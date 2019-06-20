package frule_module

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/golang/resources/db"
	"stash.tutu.ru/golang/resources/db/mysql"
	"testing"
)

func TestSearchRequest(t *testing.T) {
	database := mysql.NewDb()
	database.WithConfig(db.Config{
		DSN:   "webuser:qazxswedc@tcp(devel-02.mysql.avia.devel.tutu.ru:3306)/devel",
		Debug: true,
	})
	err := database.Init()
	if err != nil {
		t.Fatal(err)
	}

	serviceClass := "Y"
	connectionGroup := "fake"
	departureCityId := uint64(495)
	arrivalCityId := uint64(75)
	countryId := uint64(7)
	testRule := SearchRequest{
		ConnectionGroup:    &connectionGroup,
		DepartureCityId:    &departureCityId,
		ArrivalCityId:      &arrivalCityId,
		DepartureCountryId: &countryId,
		ArrivalCountryId:   &countryId,
		ServiceClass:       &serviceClass,
	}
	ctx := context.Background()
	searchRequestFrule := NewFRule(ctx, NewSearchRequest(database))
	result := searchRequestFrule.GetResult(testRule).(bool)
	assert.True(t, result)
}
