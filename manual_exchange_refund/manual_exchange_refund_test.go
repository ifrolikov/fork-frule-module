package manual_exchange_refund

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"stash.tutu.ru/golang/log"
	"testing"
	"time"
)

// Тут адовый замес
func TestComparisonOrderImporter(t *testing.T) {
	updaterErr := errors.New("test")
	defaultComparisonOrder := frule_module.ComparisonOrder{[]string{"test"}}
	comparisonOrderFromUpdater := frule_module.ComparisonOrder{[]string{"from updater"}}
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	// если нужен изначальный импорт и он не сработал
	updaterMock := NewMockComparisonOrderUpdaterInterface(ctrl)
	updaterMock.EXPECT().update(log.Logger).Return(nil, updaterErr).AnyTimes()

	comparisonOrderImporter := NewComparisonOrderImporter(
		time.Duration(0),
		updaterMock,
		nil,
	)
	result, err := comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, updaterErr, err)

	// Если нужен изначальный импорт и он сработал
	updaterMock = NewMockComparisonOrderUpdaterInterface(ctrl)
	updaterMock.EXPECT().update(log.Logger).Return(
		&comparisonOrderContainer{comparisonOrderFromUpdater},
		nil,
	).AnyTimes()
	comparisonOrderImporter = NewComparisonOrderImporter(
		time.Duration(0),
		updaterMock,
		nil,
	)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, comparisonOrderFromUpdater, result)
	assert.Nil(t, err)

	// Если не нужен импорт, проверяем ошибочный импорт
	updaterMock = NewMockComparisonOrderUpdaterInterface(ctrl)
	updaterMock.EXPECT().update(log.Logger).Return(
		nil,
		updaterErr,
	).AnyTimes()
	comparisonOrderImporter = NewComparisonOrderImporter(
		time.Duration(1*time.Nanosecond),
		updaterMock,
		&comparisonOrderContainer{defaultComparisonOrder},
	)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	// Если не нужен импорт, проверяем поведение пока импорт идет, и пирамида уже есть + то что она поменяется в итоге
	updaterMock = NewMockComparisonOrderUpdaterInterface(ctrl)
	updaterMock.EXPECT().update(log.Logger).DoAndReturn(func(logger zerolog.Logger) (*comparisonOrderContainer, error) {
		time.Sleep(100 * time.Millisecond)
		return &comparisonOrderContainer{comparisonOrderFromUpdater}, nil
	}).AnyTimes()
	comparisonOrderImporter = NewComparisonOrderImporter(
		time.Duration(0),
		updaterMock,
		&comparisonOrderContainer{defaultComparisonOrder},
	)
	// Тут отдаст дефолтную и запустит апдейтер
	time.Sleep(110 * time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	// Тут отдаст дефолтную т к апдейт еще идет
	time.Sleep(50 * time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	// Тут отдаст новую - апдейт уже закончился
	time.Sleep(60 * time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, comparisonOrderFromUpdater, result)
	assert.Nil(t, err)
}

func TestManualExchangeRefundStorage(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		ctx.Done()
	}()

	logger := log.Logger
	logger = logger.With().Str("context.type", "manual_exchange_refund_frule").Logger()

	comparisonOrderImporterMock := NewMockComparisonOrderImporterInterface(ctrl)
	comparisonOrderImporterMock.EXPECT().getComparisonOrder(logger).Return(
		frule_module.ComparisonOrder{[]string{"carrier_id", "context", "fare", "passenger_type"}},
		nil,
	).AnyTimes()

	manualExchangeRefundFRule, err := NewManualExchangeRefundFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/manual_exchange_refund.json")},
		comparisonOrderImporterMock,
	)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), manualExchangeRefundFRule)

	dataStorage := manualExchangeRefundFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 2)

	assert.Equal(t, 0, dataStorage.GetMaxRank())
}

// Тут проверяется что именно импортер возвращает пирамиду фрулу
func TestManualExchangeRefundResultWithMockedImporter(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		ctx.Done()
	}()

	logger := log.Logger
	logger = logger.With().Str("context.type", "manual_exchange_refund_frule").Logger()

	comparisonOrderImporterMock := NewMockComparisonOrderImporterInterface(ctrl)
	comparisonOrderImporterMock.EXPECT().getComparisonOrder(logger).Return(
		frule_module.ComparisonOrder{[]string{"carrier_id", "context", "fare", "passenger_type"}},
		nil,
	).AnyTimes()

	manualExchangeRefundFRule, err := NewManualExchangeRefundFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/manual_exchange_refund.json")},
		comparisonOrderImporterMock,
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

// Тут проверяется что именно из апдейтера через импортер тянется пирамида во фрул
func TestManualExchangeRefundResultWithMockedUpdater(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		ctx.Done()
	}()

	logger := log.Logger
	logger = logger.With().Str("context.type", "manual_exchange_refund_frule").Logger()

	comparisonOrderUpdaterMock := NewMockComparisonOrderUpdaterInterface(ctrl)
	comparisonOrderUpdaterMock.EXPECT().update(logger).
		Return(
			&comparisonOrderContainer{frule_module.ComparisonOrder{[]string{"carrier_id", "context", "fare", "passenger_type"}}},
			nil,
		).
		AnyTimes()

	comparisonOrderImporter := NewComparisonOrderImporter(time.Duration(0), comparisonOrderUpdaterMock, nil)

	manualExchangeRefundFRule, err := NewManualExchangeRefundFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/manual_exchange_refund.json")},
		comparisonOrderImporter,
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
}
