package airline_restrictions

import (
"context"
"github.com/stretchr/testify/assert"
"stash.tutu.ru/avia-search-common/frule-module"
"stash.tutu.ru/avia-search-common/repository"
"stash.tutu.ru/avia-search-common/utils/system"
"testing"
)

func TestAirlineRestrictionStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	airlineRestrictionFRule, err := NewAirlineRestrictionFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/airline_restriction.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), airlineRestrictionFRule)

	dataStorage := airlineRestrictionFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 3)
	assert.Len(t, (*dataStorage)[3], 1)

	assert.Equal(t, 3, dataStorage.GetMaxRank())
}

func TestAirlineRestrictionResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	airlineRestrictionFRule, err := NewAirlineRestrictionFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/airline_restriction.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, airlineRestrictionFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	assert.False(t, frule.GetResult(AirlineRestrictionsRule{Partner: &partner}).(bool))

	gds := "galileo"
	validatingCarrierId := int64(8)
	assert.True(t, frule.GetResult(AirlineRestrictionsRule{Partner: &partner, Gds: &gds, ValidatingCarrierId: &validatingCarrierId}).(bool))

	partner = "new_tt"
	gds = "galileo"
	validatingCarrierId = int64(7)
	assert.False(t, frule.GetResult(AirlineRestrictionsRule{Partner: &partner, Gds: &gds, ValidatingCarrierId: &validatingCarrierId}).(bool))

	partner = "new_tt"
	gds = "sabre"
	validatingCarrierId = int64(1111)
	assert.True(t, frule.GetResult(AirlineRestrictionsRule{Partner: &partner, Gds: &gds, ValidatingCarrierId: &validatingCarrierId}).(bool))

	partner = "unknown"
	assert.Equal(t, airlineRestrictionFRule.GetDefaultValue(), frule.GetResult(AirlineRestrictionsRule{Partner: &partner}).(bool))
}

