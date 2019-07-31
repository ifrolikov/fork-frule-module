package airline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestAirlineStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	airlineFRule, err := NewAirlineFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/airline.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), airlineFRule)

	dataStorage := airlineFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 3)
	assert.Len(t, (*dataStorage)[3], 1)

	assert.Equal(t, 3, dataStorage.GetMaxRank())
}

func TestAirlineResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	airlineFRule, err := NewAirlineFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/airline.json")})
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
