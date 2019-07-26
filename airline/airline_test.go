package airline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestAirlineStorage(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/airline.json")}
	ctx := context.Background()
	defer ctx.Done()

	frule, err := NewAirlineFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), frule)

	dataStorage := frule.GetDataStorage()
	assert.NotNil(t, dataStorage)
	assert.Len(t, (*dataStorage)[0], 3)
	assert.Len(t, (*dataStorage)[3], 1)

	maxKey := 0
	for key := range *dataStorage {
		if key > maxKey {
			maxKey = key
		}
	}
	assert.Equal(t, 3, maxKey)
}

func TestAirlineResult(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/airline.json")}
	ctx := context.Background()
	defer ctx.Done()

	airlineFRule, err := NewAirlineFRule(ctx, testConfig)
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, airlineFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	assert.False(t, frule.GetResult(AirlineRule{Partner: &partner}).(bool))

	connectionGroup := "galileo"
	carrierId := int64(8)
	assert.True(t, frule.GetResult(AirlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierId: &carrierId}).(bool))

	partner = "new_tt"
	connectionGroup = "galileo"
	carrierId = int64(7)
	assert.False(t, frule.GetResult(AirlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierId: &carrierId}).(bool))

	partner = "new_tt"
	connectionGroup = "sabre"
	carrierId = int64(1111)
	assert.True(t, frule.GetResult(AirlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierId: &carrierId}).(bool))

	partner = "unknown"
	assert.Equal(t, airlineFRule.GetDefaultValue(), frule.GetResult(AirlineRule{Partner: &partner}).(bool))
}
