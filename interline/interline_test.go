package interline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestInterlineStorage(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/interline.json")}
	ctx := context.Background()
	defer ctx.Done()

	interlineFRule, err := NewInterlineFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), interlineFRule)

	dataStorage := interlineFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 13)
	assert.Len(t, (*dataStorage)[1], 147)

	maxKey := 0
	for key := range *dataStorage {
		if key > maxKey {
			maxKey = key
		}
	}
	assert.Equal(t, 7, maxKey)
}

func TestInterlineResult(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/interline.json")}
	ctx := context.Background()
	defer ctx.Done()

	interlineFRule, err := NewInterlineFRule(ctx, testConfig)
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, interlineFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierPlating := int64(1347)

	result := frule.GetResult(InterlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierPlating: &carrierPlating, Carriers: []int64{1453}})
	assert.False(t, result.(bool))

	result = frule.GetResult(InterlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierPlating: &carrierPlating, Carriers: []int64{1212,1062}})
	assert.False(t, result.(bool))

	result = frule.GetResult(InterlineRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierPlating: &carrierPlating, Carriers: []int64{1062}})
	assert.True(t, result.(bool))

	partner = "s7"
	result = frule.GetResult(InterlineRule{Partner: &partner})
	assert.False(t, result.(bool))

	pureInterline := true
	result = frule.GetResult(InterlineRule{Partner: &partner, PureInterline: &pureInterline})
	assert.True(t, result.(bool))

	partner = "unknown"
	assert.Equal(t, interlineFRule.GetDefaultValue(), frule.GetResult(InterlineRule{Partner: &partner}).(bool))
}
