package search_request

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestSearchRequest(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/search_request.json")}
	ctx := context.Background()
	defer ctx.Done()

	searchRequestFRule, err := NewSearchRequestFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), searchRequestFRule)

	dataStorage := searchRequestFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	frule := frule_module.NewFRule(ctx, searchRequestFRule)
	assert.NotNil(t, frule)

	serviceClass := "Y"
	connectionGroup := "fake"
	departureCityId := uint64(495)
	arrivalCityId := uint64(75)
	countryId := uint64(7)
	testRule := SearchRequestRule{
		ConnectionGroup:    &connectionGroup,
		DepartureCityId:    &departureCityId,
		ArrivalCityId:      &arrivalCityId,
		DepartureCountryId: &countryId,
		ArrivalCountryId:   &countryId,
		ServiceClass:       &serviceClass,
	}

	result := frule.GetResult(testRule).(bool)
	assert.False(t, result)
}
