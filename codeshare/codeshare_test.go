package codeshare

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestCodeshareStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	codeshareFRule, err := NewCodeshareFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/codeshare.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), codeshareFRule)

	dataStorage := codeshareFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 2)
	assert.Len(t, (*dataStorage)[1], 1)

	assert.Equal(t, 7, dataStorage.GetMaxRank())
}

func TestCodeshareResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	codeshareFRule, err := NewCodeshareFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/codeshare.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, codeshareFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	assert.False(t, frule.GetResult(CodeshareRule{Partner: &partner}).(bool))

	connectionGroup := "galileo"
	carrierOperating := int64(1106)
	carrierMarketing := int64(38)
	serviceClass := "Y"
	assert.True(t, frule.GetResult(CodeshareRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierOperating: &carrierOperating, CarrierMarketing: &carrierMarketing, ServiceClass: &serviceClass}).(bool))

	partner = "tt"
	connectionGroup = "galileo"
	carrierOperating = int64(1118)
	carrierMarketing = int64(1106)
	serviceClass = "Y"
	assert.False(t, frule.GetResult(CodeshareRule{Partner: &partner, ConnectionGroup: &connectionGroup, CarrierOperating: &carrierOperating, CarrierMarketing: &carrierMarketing, ServiceClass: &serviceClass}).(bool))

	assert.True(t, frule.GetResult(CodeshareRule{Partner: &partner}).(bool))

	partner = "unknown"
	assert.Equal(t, codeshareFRule.GetDefaultValue(), frule.GetResult(CodeshareRule{Partner: &partner}).(bool))
}
