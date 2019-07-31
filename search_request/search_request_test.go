package search_request

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestSearchConnectionStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	searchRequestFRule, err := NewSearchRequestFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/search_request.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), searchRequestFRule)

	dataStorage := searchRequestFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 4)

	assert.Equal(t, 17, dataStorage.GetMaxRank())
}

func TestSearchRequestData(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	searchRequestFRule, err := NewSearchRequestFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/search_request.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, searchRequestFRule)
	assert.NotNil(t, frule)

	serviceClass := "Y"
	connectionGroup := "fake"
	departureCityId := uint64(495)
	arrivalCityId := uint64(75)
	countryId := uint64(7)
	assert.True(t, frule.GetResult(SearchRequestRule{
		ConnectionGroup:    &connectionGroup,
		DepartureCityId:    &departureCityId,
		ArrivalCityId:      &arrivalCityId,
		DepartureCountryId: &countryId,
		ArrivalCountryId:   &countryId,
		ServiceClass:       &serviceClass,
	}).(bool))

	countryId = uint64(1)
	assert.False(t, frule.GetResult(SearchRequestRule{
		ConnectionGroup:    &connectionGroup,
		DepartureCityId:    &departureCityId,
		ArrivalCityId:      &arrivalCityId,
		DepartureCountryId: &countryId,
		ArrivalCountryId:   &countryId,
		ServiceClass:       &serviceClass,
	}).(bool))

	connectionGroup = "galileo"
	assert.True(t, frule.GetResult(SearchRequestRule{ConnectionGroup:    &connectionGroup}).(bool))

	connectionGroup = "unknown"
	assert.Equal(t, searchRequestFRule.GetDefaultValue(), frule.GetResult(SearchRequestRule{ConnectionGroup:    &connectionGroup}).(bool))
}
