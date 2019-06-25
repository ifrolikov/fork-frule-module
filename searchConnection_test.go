package frule_module

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"stash.tutu.ru/golang/resources/db"
	"stash.tutu.ru/golang/resources/db/mysql"
	"strconv"
	"testing"
	"time"
)

func TestSearchConnectionDb(t *testing.T) {
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
	testRule := SearchConnectionRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		DepartureDate:   time.Now(),
	}

	ctx := context.Background()
	rule := NewFRule(ctx, NewSearchConnectionFRule(database))

	result := rule.GetResult(testRule)
	fmt.Println(result)
	//assert.True(t, result)
}

func TestSearchConnectionSpecs(t *testing.T) {
	database := mysql.NewDb()
	rule := NewSearchConnectionFRule(database)

	specs := []cronStrucString{
		{"50-59 23 * * *", "+3w"},
		{"0-15 20 * * *", "+100d"},
		{"16-30 20 * * *", "+1y"},
		{"* 0-19 * * *", "+2m"},
		{"* * * * *", "0"},
	}

	now, _ := time.Parse("2006-01-02 15:04:05", "2019-06-21 11:21:00")

	compareList := []compareStructure{
		{getLine(), rule.getSpecInterval(specs, parseTime("2019-10-10 01:00:00")), "+2m"},
		{getLine(), rule.getSpecInterval(specs, parseTime("2019-10-10 19:01:00")), "+2m"},
		{getLine(), rule.getSpecInterval(specs, parseTime("2019-10-10 20:17:00")), "+1y"},
		{getLine(), rule.getSpecInterval(specs, parseTime("2019-10-10 20:05:00")), "+100d"},
		{getLine(), rule.getSpecInterval(specs, parseTime("2019-10-10 23:50:00")), "+3w"},

		{getLine(), rule.getDateToCompare(rule.getSpecInterval(specs, parseTime("2019-10-10 01:00:00")), now).Format("2006-01-02"), "2019-08-21"},
		{getLine(), rule.getDateToCompare(rule.getSpecInterval(specs, parseTime("2019-10-10 19:01:00")), now).Format("2006-01-02"), "2019-08-21"},
		{getLine(), rule.getDateToCompare(rule.getSpecInterval(specs, parseTime("2019-10-10 20:17:00")), now).Format("2006-01-02"), "2020-06-21"},
		{getLine(), rule.getDateToCompare(rule.getSpecInterval(specs, parseTime("2019-10-10 20:05:00")), now).Format("2006-01-02"), "2019-09-29"},
		{getLine(), rule.getDateToCompare(rule.getSpecInterval(specs, parseTime("2019-10-10 23:50:00")), now).Format("2006-01-02"), "2019-07-12"},
	}

	compare(t, compareList)
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}

func getLine() int {
	_, _, l, _ := runtime.Caller(1)
	return l
}

type compareStructure struct {
	line int
	got  interface{}
	want interface{}
}

func compare(t *testing.T, compareList []compareStructure) {
	for i := range compareList {
		t.Run("Line"+strconv.Itoa(compareList[i].line),
			func(t *testing.T) {
				if !reflect.DeepEqual(compareList[i].got, compareList[i].want) {
					t.Errorf("got '%v' want '%v' structure on the %d line", compareList[i].got, compareList[i].want, compareList[i].line)
				}
			})
	}
}
