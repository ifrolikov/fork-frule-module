package search_scheme

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestSearchSchemeStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	searchRequestFRule, err := NewSearchSchemeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/search_scheme.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), searchRequestFRule)

	dataStorage := searchRequestFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 21)
	assert.Len(t, (*dataStorage)[2], 0)

	assert.Equal(t, 13, dataStorage.GetMaxRank())
}

func TestSearchSchemeData(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	searchRequestFRule, err := NewSearchSchemeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/search_scheme.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, searchRequestFRule)
	assert.NotNil(t, frule)

	connectionGroup := "fake"
	departureCityId := uint64(495)
	arrivalCityId := uint64(75)
	countryId := uint64(7)
	testRule := SearchSchemeRule{
		ConnectionGroup:    &connectionGroup,
		DepartureCityId:    &departureCityId,
		ArrivalCityId:      &arrivalCityId,
		DepartureCountryId: &countryId,
		ArrivalCountryId:   &countryId,
	}

	result := frule.GetResult(testRule).([]string)
	assert.Equal(t, []string{"default"}, result)

	connectionGroup = "galileo"
	departureCityId = uint64(491)
	arrivalCityId = uint64(34)
	countryId = uint64(7)
	testRule = SearchSchemeRule{
		ConnectionGroup:    &connectionGroup,
		DepartureCityId:    &departureCityId,
		ArrivalCityId:      &arrivalCityId,
		DepartureCountryId: &countryId,
		ArrivalCountryId:   &countryId,
	}

	result = frule.GetResult(testRule).([]string)
	assert.Contains(t, result, "onlySU")
	assert.Contains(t, result, "notSU")
}
