package warning

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
	"time"
)

func TestWarningStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	warningRule, err := NewWarningFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/warning.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), warningRule)

	dataStorage := warningRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[21], 3)

	assert.Equal(t, 24, dataStorage.GetMaxRank())
}

func TestWarningData(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	warningFRule, err := NewWarningFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/warning.json")})
	assert.Nil(t, err)
	frule := frule_module.NewFRule(ctx, warningFRule)
	assert.NotNil(t, frule)

	var departureCountryId int64 = 7
	var arrivalCountryId int64 = 72
	departureDate, _ := time.Parse("2006-01-02", "2016-01-01")
	lang := "rus"
	ruleResult := frule.GetResult(WarningRule{
		DepartureCountryId: &departureCountryId,
		ArrivalCountryId:   &arrivalCountryId,
		ParsedStartDate:    &departureDate,
		Lang:               &lang,
	})
	assert.NotNil(t, ruleResult)
	assert.Equal(t, "egypt_2015", *ruleResult.(*RuleResult).Result[0].Name)

	notActualDepartureDate, _ := time.Parse("2006-01-02", "2019-01-01")
	ruleResult = frule.GetResult(WarningRule{
		DepartureCountryId: &departureCountryId,
		ArrivalCountryId:   &arrivalCountryId,
		ParsedStartDate:    &notActualDepartureDate,
		Lang:               &lang,
	})
	assert.Nil(t, ruleResult)
}