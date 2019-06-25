package frule_module

import (
	"context"
	"fmt"
	"stash.tutu.ru/golang/resources/db"
	"stash.tutu.ru/golang/resources/db/mysql"
	"testing"
)

func TestInterlineDb(t *testing.T) {
	database := mysql.NewDb()
	database.WithConfig(db.Config{
		DSN:   "webuser:qazxswedc@tcp(devel-02.mysql.avia.devel.tutu.ru:3306)/devel",
		Debug: true,
	})
	err := database.Init()
	if err != nil {
		t.Fatal(err)
	}

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierPlating := int64(1347)
	testRule := InterlineRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierPlating:  &carrierPlating,
		Carriers:        []int64{1453},
	}

	ctx := context.Background()
	rule := NewFRule(ctx, NewInterlineFRule(database))

	result := rule.GetResult(testRule).(bool)
	fmt.Println(result)
	//assert.True(t, result)
}
