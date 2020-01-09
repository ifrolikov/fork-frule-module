package connection

import (
	"context"
	"github.com/stretchr/testify/assert"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestConnectionResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	fareFRule, err := NewConnectionFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/connection.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, fareFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierId := int64(1106)
	operation := "issue"
	paymentEngine := "mps_avia"

	assert.Equal(
		t,
		"new_tt_galileo_issue_s7",
		frule.GetResult(ConnectionRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup, Operation: &operation, PaymentEngine: &paymentEngine}),
	)

	partner = "new_tt"
	connectionGroup = "sig23"
	carrierId = int64(1106)
	operation = "issue"
	paymentEngine = "mps_avia"

	assert.Equal(
		t,
		"new_tt_sig23_mps",
		frule.GetResult(ConnectionRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup, Operation: &operation, PaymentEngine: &paymentEngine}),
	)

}
