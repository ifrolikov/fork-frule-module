package refund_types

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/contracts/v2/gateSearch"
	frule_module "github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"stash.tutu.ru/golang/log"
	"testing"
	"time"
)


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

func TestRefundTypesStorage(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		ctx.Done()
	}()

	logger := log.Logger
	logger = logger.With().Str("context.type", "refund_types_rule").Logger()

	comparisonOrderImporterMock := NewMockComparisonOrderImporterInterface(ctrl)
	comparisonOrderImporterMock.EXPECT().getComparisonOrder(logger).Return(
		frule_module.ComparisonOrder{[]string{"plating_carrier_id", "issue_date_from", "issue_date_to"}},
		nil,
	).AnyTimes()

	refundTypesFRule, err := NewRefundTypesFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/refund_types.json")},
		comparisonOrderImporterMock,
	)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), refundTypesFRule)

	dataStorage := refundTypesFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 1)

	assert.Equal(t, 0, dataStorage.GetMaxRank())
}

func TestRefundTypesResultWithMockedImporter(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		ctx.Done()
	}()

	logger := log.Logger
	logger = logger.With().Str("context.type", "refund_types_frule").Logger()

	comparisonOrderImporterMock := NewMockComparisonOrderImporterInterface(ctrl)
	comparisonOrderImporterMock.EXPECT().getComparisonOrder(logger).Return(
		frule_module.ComparisonOrder{[]string{"plating_carrier_id", "issue_date_from", "issue_date_to"}},
		nil,
	).AnyTimes()

	refundTypesFRule, err := NewRefundTypesFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/refund_types.json")},
		comparisonOrderImporterMock,
	)
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, refundTypesFRule)
	assert.NotNil(t, frule)

	platingCarrierId := int64(1062)

	assert.EqualValues(t, refundTypesFRule.GetDefaultValue(), frule.GetResult(RefundTypesRule{
		PlatingCarrierId: &platingCarrierId,
	}))

	issueDate := "2020-06-02"

	assert.NotEqual(t, refundTypesFRule.GetDefaultValue(), frule.GetResult(RefundTypesRule{
		PlatingCarrierId: &platingCarrierId,
		IssueDateFrom:    &issueDate,
		IssueDateTo:      &issueDate,
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
	logger = logger.With().Str("context.type", "refund_types_frule").Logger()

	comparisonOrderUpdaterMock := NewMockComparisonOrderUpdaterInterface(ctrl)
	comparisonOrderUpdaterMock.EXPECT().update(logger).
		Return(
			&comparisonOrderContainer{
				frule_module.ComparisonOrder{
					[]string{"plating_carrier_id", "issue_date_from", "issue_date_to"}}},
			nil,
		).
		AnyTimes()

	comparisonOrderImporter := NewComparisonOrderImporter(time.Duration(0), comparisonOrderUpdaterMock, nil)

	refundTypesFRule, err := NewRefundTypesFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/refund_types.json")},
		comparisonOrderImporter,
	)
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, refundTypesFRule)
	assert.NotNil(t, frule)

	platingCarrierId := int64(1062)
	issueDate := "2020-06-02"

	assert.Equal(t, gateSearch.RefundType_airline_voucher, frule.GetResult(RefundTypesRule{
		PlatingCarrierId: &platingCarrierId,
		IssueDateFrom:    &issueDate,
		IssueDateTo:      &issueDate,
	}))
}
