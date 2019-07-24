package partner_percent

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
	"time"
)

func TestPartnerPercent(t *testing.T) {
	pwd, _ := filepath.Abs("./")
	testConfig := &repository.Config{
		DataURI: "file://" + pwd + "/../testdata/partner_percent.json",
	}
	ctx := context.Background()
	defer ctx.Done()

	partnerPercentFRule, err := NewPartnerPercentFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), partnerPercentFRule)

	dataStorage := partnerPercentFRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	frule := frule_module.NewFRule(ctx, partnerPercentFRule)
	assert.NotNil(t, frule)

	partner := "new_tt"
	connectionGroup := "galileo"
	carrierId := int64(1347)
	dateOfPurchaseFrom := "2019-01-01"
	dateOfPurchaseTo := time.Now().Format("2006-01-02 15:04:05")

	testRule := PartnerPercentRule{
		Partner:            &partner,
		ConnectionGroup:    &connectionGroup,
		CarrierId:          &carrierId,
		DateOfPurchaseFrom: &dateOfPurchaseFrom,
		DateOfPurchaseTo:   &dateOfPurchaseTo,
	}

	result := frule.GetResult(testRule)
	fmt.Printf("%+v", result)
	assert.Equal(t, 0.4, result)
}
