package manual_exchange_refund

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

type DummyComparisonOrderImporter struct {
}

func (importer *DummyComparisonOrderImporter) getComparisonOrder(logger zerolog.Logger) (frule_module.ComparisonOrder, error) {
	return frule_module.ComparisonOrder{
		[]string{"carrier_id", "context", "fare", "passenger_type"},
	}, nil
}

func TestManualExchangeRefundStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	manualExchangeRefundFRule, err := NewManualExchangeRefundFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/manual_exchange_refund.json")},
		&DummyComparisonOrderImporter{},
	)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), manualExchangeRefundFRule)

	dataStorage := manualExchangeRefundFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 2)

	assert.Equal(t, 0, dataStorage.GetMaxRank())
}

func TestManualExchangeRefundResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	manualExchangeRefundFRule, err := NewManualExchangeRefundFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/manual_exchange_refund.json")},
		&DummyComparisonOrderImporter{},
	)
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, manualExchangeRefundFRule)
	assert.NotNil(t, frule)

	fare := "SOMEFARE"
	passengerType := "full"
	carrierId := int32(1062)
	fruleContext := ContextExchange

	assert.Equal(t, int32(1), frule.GetResult(ManualExchangeRefundRule{
		CarrierId:     &carrierId,
		Context:       &fruleContext,
		Fare:          &fare,
		PassengerType: &passengerType,
	}).(ManualExchangeRefundResult).Id)

	passengerType = "child"
	assert.Equal(t, int32(2), frule.GetResult(ManualExchangeRefundRule{
		CarrierId:     &carrierId,
		Context:       &fruleContext,
		Fare:          &fare,
		PassengerType: &passengerType,
	}).(ManualExchangeRefundResult).Id)

	assert.EqualValues(t, manualExchangeRefundFRule.GetDefaultValue(), frule.GetResult(ManualExchangeRefundRule{
		CarrierId: &carrierId,
		Context:   &fruleContext,
		Fare:      &fare,
	}))
}
