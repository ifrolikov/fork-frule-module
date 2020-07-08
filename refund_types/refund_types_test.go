package refund_types

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
		[]string{"plating_carrier_id", "issue_date_from", "issue_date_to"},
	}, nil
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

	r1 := refundTypesFRule.GetDefaultValue()
	r2 := frule.GetResult(RefundTypesRule{
		PlatingCarrierId: &platingCarrierId,
		IssueDateFrom: &issueDate,
		IssueDateTo: &issueDate,
	})
	assert.NotEqual(t, r1, r2)
}
