package interline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestInterlineStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	interlineFRule, err := NewInterlineFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/interline.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), interlineFRule)

	dataStorage := interlineFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 2)
	assert.Len(t, (*dataStorage)[1], 1)

	assert.Equal(t, 7, dataStorage.GetMaxRank())
}

func TestInterlineResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	interlineFRule, err := NewInterlineFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/interline.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, interlineFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierPlating := int64(1347)
	pureInterline := true

	result := frule.GetResult(InterlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierPlating: &carrierPlating, PureInterline: &pureInterline, Carriers: []int64{1453}})
	assert.True(t, result.(bool))

	result = frule.GetResult(InterlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierPlating: &carrierPlating, Carriers: []int64{1453}})
	assert.False(t, result.(bool))

	result = frule.GetResult(InterlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierPlating: &carrierPlating, Carriers: []int64{1212,1062}})
	assert.False(t, result.(bool))

	result = frule.GetResult(InterlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierPlating: &carrierPlating, Carriers: []int64{1062}})
	assert.True(t, result.(bool))

	partner = "s7"
	result = frule.GetResult(InterlineRule{Partner: &partner})
	assert.False(t, result.(bool))

	result = frule.GetResult(InterlineRule{Partner: &partner, PureInterline: &pureInterline})
	assert.True(t, result.(bool))

	partner = "unknown"
	assert.Equal(t, interlineFRule.GetDefaultValue(), frule.GetResult(InterlineRule{Partner: &partner}).(bool))
}
