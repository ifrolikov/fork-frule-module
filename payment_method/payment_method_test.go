package payment_method

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/contracts/pricing"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestPaymentMethodStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	fareFRule, err := NewPaymentMethodFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/payment_method.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), fareFRule)

	dataStorage := fareFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)
	assert.Len(t, (*dataStorage)[1], 13)
	assert.Len(t, (*dataStorage)[3], 15)

	assert.Equal(t, 7, dataStorage.GetMaxRank())
}

func TestPaymentMethodResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	fareFRule, err := NewPaymentMethodFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/payment_method.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, fareFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierId := int64(73)
	autoticketing := false

	assert.Equal(t, []string{pricing.PAYMENT_METHOD_CARDONLINE, pricing.PAYMENT_METHOD_APPLE_PAY}, frule.GetResult(PaymentMethodRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup, Autoticketing: &autoticketing, TestDaysTillDeparture: 4}))

	partner = "fake"
	connectionGroup = "fake"
	carrierId = int64(73)
	autoticketing = false

	assert.Equal(t, []string{pricing.PAYMENT_METHOD_CARDONLINE, pricing.PAYMENT_METHOD_APPLE_PAY}, frule.GetResult(PaymentMethodRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup, Autoticketing: &autoticketing, TestDaysTillDeparture: 4}))

	partner = "new_tt"
	connectionGroup = "sirena_direct_ut"
	carrierId = int64(73)
	autoticketing = false

	assert.Equal(t, []string{""}, frule.GetResult(PaymentMethodRule{Partner: &partner, CarrierId: &carrierId, ConnectionGroup: &connectionGroup, Autoticketing: &autoticketing, TestDaysTillDeparture: 4}))

}
