package payment_engine

import (
	"context"
	"github.com/stretchr/testify/assert"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestPaymentEngineResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	fareFRule, err := NewPaymentEngineFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/payment_engine.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, fareFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierId := int64(73)
	paymentMethod := "cardonline"

	res1 := "work_avia"
	res2 := "ntt_avia"

	assert.Equal(
		t,
		[]EngineConfig{
			{
				Engine:     "boxplat",
				ConfigType: &res1,
			},
			{
				Engine:     "gateline",
				ConfigType: &res2,
			},
		},
		frule.GetResult(PaymentEngineRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup, PaymentMethod: &paymentMethod}),
	)

	partner = "fake"
	connectionGroup = "fake_api"
	carrierId = int64(73)
	paymentMethod = "cardonline"

	assert.Equal(
		t,
		[]EngineConfig{
			{
				Engine:     "fakeCard",
				ConfigType: nil,
			},
		},
		frule.GetResult(PaymentEngineRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup, PaymentMethod: &paymentMethod}),
	)
}
