package interline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestInterline(t *testing.T) {
	pwd, _ := filepath.Abs("./")
	testConfig := &repository.Config{
		DataURI: "file://" + pwd + "/../testdata/interline.json",
	}
	ctx := context.Background()
	defer ctx.Done()

	interlineFRule, err := NewInterlineFRule(ctx, testConfig)
	assert.Nil(t, err)

	dataStorage := interlineFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	frule := frule_module.NewFRule(ctx, interlineFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierPlating := int64(1347)
	testRule := InterlineRule{
		Partner:         &partner,
		ConnectionGroup: &connectionGroup,
		CarrierPlating:  &carrierPlating,
		Carriers:        []int64{1453},
	}
	result := frule.GetResult(testRule)
	assert.False(t, result.(bool))
}
