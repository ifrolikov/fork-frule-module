package card_in_pnr

import (
	"context"
	"github.com/stretchr/testify/assert"
	frule_module "github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestCardInPnrResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	fareFRule, err := NewCardInPnrFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/card_in_pnr.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, fareFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierId := int64(73)

	assert.Equal(
		t,
		false,
		frule.GetResult(CardInPnrRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup}),
	)

	partner = "new_tt"
	connectionGroup = "sig23_devel"
	carrierId = int64(1116)

	assert.Equal(
		t,
		true,
		frule.GetResult(CardInPnrRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup}),
	)

}
