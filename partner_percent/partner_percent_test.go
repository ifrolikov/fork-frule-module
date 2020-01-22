package partner_percent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
	"time"
)

func TestPartnerPercentStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	partnerPercentFRule, err := NewPartnerPercentFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/partner_percent.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), partnerPercentFRule)

	dataStorage := partnerPercentFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 1)
	assert.Len(t, (*dataStorage)[2], 2)

	assert.Equal(t, 7, dataStorage.GetMaxRank())
}

func TestPartnerPercentResult(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	partnerPercentFRule, err := NewPartnerPercentFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/partner_percent.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, partnerPercentFRule)
	assert.NotNil(t, frule)

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	partner := "new_tt"
	dateOfPurchaseFrom := currentTime
	dateOfPurchaseTo := currentTime
	connectionGroup := "sabre"
	carrierId := int64(1062)
	fareType := "subsidy"
	assert.Equal(t, 0.1, frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		ConnectionGroup:    &connectionGroup,
		CarrierId:          &carrierId,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseTo,
		FareType:           &fareType,
	}).(PartnerPercentResult).Percent)

	fareType = "test"
	assert.NotEqual(t, 0.1, frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		ConnectionGroup:    &connectionGroup,
		CarrierId:          &carrierId,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseTo,
		FareType:           &fareType,
	}).(PartnerPercentResult).Percent)

	assert.EqualValues(t, partnerPercentFRule.GetDefaultValue(), frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseTo,
	}))
	dateOfPurchaseToAlternative := time.Date(2016, 11, 3, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")
	assert.Equal(t, 1.0, frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseToAlternative,
	}).(PartnerPercentResult).Percent)

	assert.Equal(t, partnerPercentFRule.GetDefaultValue(), frule.GetResult(PartnerPercentRule{Partner: &partner}))

	connectionGroup = "galileo"
	carrierId = int64(1111)
	assert.Equal(t, 0.4, frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		ConnectionGroup:    &connectionGroup,
		CarrierId:          &carrierId,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseTo,
	}).(PartnerPercentResult).Percent)
}
