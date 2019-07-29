package partner_percent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
	"time"
)

func TestPartnerPercentStorage(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/partner_percent.json")}
	ctx := context.Background()
	defer ctx.Done()

	partnerPercentFRule, err := NewPartnerPercentFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), partnerPercentFRule)

	dataStorage := partnerPercentFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 1)
	assert.Len(t, (*dataStorage)[2], 2)

	maxKey := 0
	for key := range *dataStorage {
		if key > maxKey {
			maxKey = key
		}
	}
	assert.Equal(t, 7, maxKey)
}

func TestPartnerPercentResult(t *testing.T) {
	pwd, _ := filepath.Abs("../")
	testConfig := &repository.Config{DataURI: filepath.ToSlash("file://" + pwd + "/testdata/partner_percent.json")}
	ctx := context.Background()
	defer ctx.Done()

	partnerPercentFRule, err := NewPartnerPercentFRule(ctx, testConfig)
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
	}))

	fareType = "test"
	assert.NotEqual(t, 0.1, frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		ConnectionGroup:    &connectionGroup,
		CarrierId:          &carrierId,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseTo,
		FareType:           &fareType,
	}))

	assert.Equal(t, partnerPercentFRule.GetDefaultValue(), frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseTo,
	}))
	dateOfPurchaseToAlternative := time.Date(2016, 11, 3, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")
	assert.Equal(t, float64(1), frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseToAlternative,
	}))

	assert.Equal(t, partnerPercentFRule.GetDefaultValue(), frule.GetResult(PartnerPercentRule{Partner: &partner}))

	connectionGroup = "galileo"
	carrierId = int64(1111)
	assert.Equal(t, 0.4, frule.GetResult(PartnerPercentRule{
		Partner:            &partner,
		ConnectionGroup:    &connectionGroup,
		CarrierId:          &carrierId,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseTo,
	}))
}
