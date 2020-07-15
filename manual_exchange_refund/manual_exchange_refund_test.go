package manual_exchange_refund

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"stash.tutu.ru/golang/log"
	"testing"
	"time"
)

type DummyComparisonOrderImporter struct {
}

func (importer *DummyComparisonOrderImporter) getComparisonOrder(logger zerolog.Logger) (frule_module.ComparisonOrder, error) {
	return frule_module.ComparisonOrder{
		[]string{"carrier_id", "context", "fare", "passenger_type"},
	}, nil
}

type DummyComparisonOrderUpdater struct {
	sleepingTime time.Duration
	result       *comparisonOrderContainer
	err          error
}

func (updater *DummyComparisonOrderUpdater) update(logger zerolog.Logger) (*comparisonOrderContainer, error) {
	time.Sleep(updater.sleepingTime)
	return updater.result, updater.err
}

func TestComparisonOrderImporter(t *testing.T) {
	updaterErr := errors.New("test")
	defaultComparisonOrder := frule_module.ComparisonOrder{[]string{"test"}}
	comparisonOrderFromUpdater := frule_module.ComparisonOrder{[]string{"from updater"}}

	// если нужен изначальный импорт и он не сработал
	comparisonOrderImporter := NewComparisonOrderImporter(
		time.Duration(1*time.Second),
		&DummyComparisonOrderUpdater{
			time.Duration(1*time.Microsecond),
			nil,
			updaterErr,
		},
		nil,
	)
	result, err := comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, updaterErr, err)

	// Если нужен изначальный импорт и он сработал
	comparisonOrderImporter = NewComparisonOrderImporter(
		time.Duration(1*time.Second),
		&DummyComparisonOrderUpdater{
			time.Duration(1*time.Microsecond),
			&comparisonOrderContainer{comparisonOrderFromUpdater},
			nil,
		},
		nil,
	)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, comparisonOrderFromUpdater, result)
	assert.Nil(t, err)

	// Если не нужен импорт, проверяем ошибочный импорт
	comparisonOrderImporter = NewComparisonOrderImporter(
		time.Duration(100*time.Millisecond),
		&DummyComparisonOrderUpdater{
			time.Duration(100*time.Millisecond),
			nil,
			updaterErr,
		},
		&comparisonOrderContainer{defaultComparisonOrder},
	)
	time.Sleep(110*time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	// Если не нужен импорт, проверяем поведение пока импорт идет, и пирамида уже есть + то что она поменяется в итоге
	comparisonOrderImporter = NewComparisonOrderImporter(
		time.Duration(100*time.Millisecond),
		&DummyComparisonOrderUpdater{
			time.Duration(100*time.Millisecond),
			&comparisonOrderContainer{comparisonOrderFromUpdater},
			nil,
		},
		&comparisonOrderContainer{defaultComparisonOrder},
	)
	// Тут отдаст дефолтную и не запустит апдейтер
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	// Тут отдаст дефолтную и запустит апдейтер
	time.Sleep(110*time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	// Тут отдаст дефолтную т к апдейт еще идет
	time.Sleep(50*time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	// Тут отдаст новую - апдейт уже закончился
	time.Sleep(60*time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, comparisonOrderFromUpdater, result)
	assert.Nil(t, err)
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
	carrierId := int64(1062)
	fruleContext := ContextExchange

	assert.Equal(t, int64(1), frule.GetResult(ManualExchangeRefundRule{
		CarrierId:     &carrierId,
		Context:       &fruleContext,
		Fare:          &fare,
		PassengerType: &passengerType,
	}).(ManualExchangeRefundResult).Id)

	passengerType = "child"
	assert.Equal(t, int64(2), frule.GetResult(ManualExchangeRefundRule{
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
