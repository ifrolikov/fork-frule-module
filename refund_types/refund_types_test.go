package refund_types

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
		[]string{"plating_carrier_id", "issue_date_from", "issue_date_to"},
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

	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	time.Sleep(110*time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, defaultComparisonOrder, result)
	assert.Nil(t, err)

	time.Sleep(110*time.Millisecond)
	result, err = comparisonOrderImporter.getComparisonOrder(log.Logger)
	assert.Equal(t, comparisonOrderFromUpdater, result)
	assert.Nil(t, err)
}

func TestRefundTypesStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	refundTypesFRule, err := NewRefundTypesFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/refund_types.json")},
		&DummyComparisonOrderImporter{},
	)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), refundTypesFRule)

	dataStorage := refundTypesFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 1)

	assert.Equal(t, 0, dataStorage.GetMaxRank())
}

func TestRefundTypesResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	refundTypesFRule, err := NewRefundTypesFRule(
		ctx,
		&repository.Config{DataURI: system.GetFilePath("../testdata/refund_types.json")},
		&DummyComparisonOrderImporter{},
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
