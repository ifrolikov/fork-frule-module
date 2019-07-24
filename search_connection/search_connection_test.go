package search_connection

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
	"time"
)

func TestSearchConnection(t *testing.T) {
	pwd, _ := filepath.Abs("./")
	testConfig := &repository.Config{
		DataURI: "file://" + pwd + "/../testdata/search_connection.json",
	}
	ctx := context.Background()
	defer ctx.Done()

	searchConnectionFRule, err := NewSearchConnectionFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), searchConnectionFRule)

	dataStorage := searchConnectionFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	frule := frule_module.NewFRule(ctx, searchConnectionFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	testRule := SearchConnectionRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		DepartureDate:   time.Now(),
	}

	result := frule.GetResult(testRule)
	fmt.Println(result)
	//assert.True(t, result)
}

func TestSearchConnectionSpecs(t *testing.T) {
	pwd, _ := filepath.Abs("./")
	testConfig := &repository.Config{
		DataURI: "file://" + pwd + "/../testdata/search_request.json",
	}
	ctx := context.Background()
	defer ctx.Done()

	searchConnectionFRule, err := NewSearchConnectionFRule(ctx, testConfig)
	assert.Nil(t, err)

	specs := []frule_module.CronStructString{
		{"50-59 23 * * *", "+3w"},
		{"0-15 20 * * *", "+100d"},
		{"16-30 20 * * *", "+1y"},
		{"* 0-19 * * *", "+2m"},
		{"* * * * *", "0"},
	}

	now, _ := time.Parse("2006-01-02 15:04:05", "2019-06-21 11:21:00")

	assert.Equal(t, "+2m", searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 01:00:00")))
	assert.Equal(t, "+2m", searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 19:01:00")))
	assert.Equal(t, "+1y", searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 20:17:00")))
	assert.Equal(t, "+100d", searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 20:05:00")))
	assert.Equal(t, "+3w", searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 23:50:00")))

	assert.Equal(t, "2019-08-21", searchConnectionFRule.getDateToCompare(searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 01:00:00")), now).Format("2006-01-02"))
	assert.Equal(t, "2019-08-21", searchConnectionFRule.getDateToCompare(searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 19:01:00")), now).Format("2006-01-02"))
	assert.Equal(t, "2020-06-21", searchConnectionFRule.getDateToCompare(searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 20:17:00")), now).Format("2006-01-02"))
	assert.Equal(t, "2019-09-29", searchConnectionFRule.getDateToCompare(searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 20:05:00")), now).Format("2006-01-02"))
	assert.Equal(t, "2019-07-12", searchConnectionFRule.getDateToCompare(searchConnectionFRule.getSpecInterval(specs, parseTime("2019-10-10 23:50:00")), now).Format("2006-01-02"))
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}
