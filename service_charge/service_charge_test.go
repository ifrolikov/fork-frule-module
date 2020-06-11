package service_charge

import (
	"context"
	"github.com/stretchr/testify/assert"
	"stash.tutu.ru/avia-search-common/contracts/v2/base"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/avia-search-common/utils/system"
	"testing"
)

func TestServiceChargeStorage(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	serviceChargeRule, err := NewServiceChargeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/service_charge.json")})
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), serviceChargeRule)

	dataStorage := serviceChargeRule.GetDataStorage()
	assert.NotNil(t, dataStorage)

	assert.Len(t, (*dataStorage)[0], 0)
	assert.Len(t, (*dataStorage)[52], 2)

	assert.Equal(t, 81, len(*dataStorage))
	assert.Equal(t, 80, dataStorage.GetMaxRank())
}

func TestServiceChargeSimpleFormat(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	serviceChargeRule, err := NewServiceChargeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/service_charge.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, serviceChargeRule)
	assert.NotNil(t, frule)

	//обычный кейс
	suId := int64(1062)
	russiaId := uint64(7)
	params := ServiceChargeRule{
		CarrierId:              &suId,
		DepartureCountryId:     &russiaId,
		ArrivalCountryId:       &russiaId,
		TestOfferPrice:         base.Money{Amount: 6000},
	}
	result := frule.GetResult(params)
	assert.EqualValues(t, 3, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 10000, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 8000, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	//субсидированный тариф перебивает предыдущий кейс, хотя остальные параметры те же
	subsidyFare := "subsidy"
	params = ServiceChargeRule{
		CarrierId:          &suId,
		FareType:           &subsidyFare,
		TestOfferPrice:     base.Money{Amount: 7000},
		DepartureCountryId: &russiaId,
		ArrivalCountryId:   &russiaId,
	}
	result = frule.GetResult(params)
	assert.EqualValues(t, 2, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	//субсидированный тариф с измененными дургими параметрами перебивает первый кейс
	notRussiaId := uint64(34)
	params = ServiceChargeRule{
		CarrierId:          &suId,
		FareType:           &subsidyFare,
		DepartureCountryId: &notRussiaId,
		ArrivalCountryId:   &notRussiaId,
	}
	result = frule.GetResult(params)
	assert.EqualValues(t, 2, result.(ServiceChargeRuleResult).Id)

	//дефолтный вариант, попадающий в последнюю строку пирамиды
	params = ServiceChargeRule{
		CarrierId:          &suId,
		DepartureCountryId: &notRussiaId,
		ArrivalCountryId:   &notRussiaId,
	}
	result = frule.GetResult(params)
	assert.EqualValues(t, 1, result.(ServiceChargeRuleResult).Id)
	assert.Equal(t, base.Money{Amount: 20000, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 20000, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	//еще один дефолтный вариант, попадающий в последнюю строку пирамиды
	notSuId := int64(10)
	params = ServiceChargeRule{
		CarrierId:          &notSuId,
		DepartureCountryId: &russiaId,
		ArrivalCountryId:   &russiaId,
	}
	result = frule.GetResult(params)
	assert.EqualValues(t, 1, result.(ServiceChargeRuleResult).Id)
}

func TestServiceChargeComplexFormat(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	serviceChargeRule, err := NewServiceChargeFRule(ctx, &repository.Config{DataURI: system.GetFilePath("../testdata/service_charge.json")})
	assert.Nil(t, err)

	frule := frule_module.NewFRule(ctx, serviceChargeRule)
	assert.NotNil(t, frule)

	carrierId := int64(1111)
	russiaId := uint64(7)
	params := ServiceChargeRule{
		CarrierId:          &carrierId,
		DepartureCountryId: &russiaId,
		ArrivalCountryId:   &russiaId,
		TestOfferPrice:     base.Money{Amount: 60000, Currency: &base.Currency{Code: "RUB", Fraction: 100}},
	}
	result := frule.GetResult(params)
	assert.EqualValues(t, 4, result.(ServiceChargeRuleResult).Id)
	/**
	 * Фиксированная наценка
	 */
	assert.Equal(t, base.Money{Amount: 18840, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	assert.Equal(t, base.Money{Amount: 18840, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	/*
	 * К фиксированной наценке добавляется процент от тарифа и такс.
	 * Полкопейки округляются вверх до целой по математическому округлению
	 */
	params.TestOfferPrice = base.Money{Amount: 167500, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
	result = frule.GetResult(params)
	// 387.7RUR+2.3% = 38770 + 167500/100*2,3 = 38770 + 3852,5 = 42622,5
	assert.Equal(t, base.Money{Amount: 42623, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	// 387.7RUR+2.1% = 38770 + 167500/100*2,1 = 38770 + 3517,5 = 42287,5
	assert.Equal(t, base.Money{Amount: 42288, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	/*
	 * К фиксированной наценке добавляется процентная с ограничением.
	 * Взрослый упирается в ограничение, для него применяется верхняя граница
	 * Ребенок не упирается в ограничение
	 */
	params.TestOfferPrice = base.Money{Amount: 234300, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
	result = frule.GetResult(params)
	// 461.13RUR+2.3%<50.1RUR = 46113 + 234300/100*2,3<5010 = 46113 + 5010 = 51123
	assert.Equal(t, base.Money{Amount: 51123, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	// 461.13RUR+2.1%<50.1RUR = 46113 + 234300/100*2,1<5010 = 46113 + 4920,3 = 51033,3
	assert.Equal(t, base.Money{Amount: 51033, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	/*
	 * Фиксированная наценка нулевая. К ней добавляется процентная с ограничением.
	 * Взрослый упирается в ограничение, для него применяется верхняя граница
	 * Ребенок не упирается в ограничение
	 * Доли копейки округляются вниз по математическому округлению
	 */
	params.TestOfferPrice = base.Money{Amount: 315721, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
	result = frule.GetResult(params)
	// 0RUR+2.3%<67 = 315721/100*2,3<6700 = 6700
	assert.Equal(t, base.Money{Amount: 6700, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	// 0RUR+2.1%<67 = 315721/100*2,1<6700 = 6630,141
	assert.Equal(t, base.Money{Amount: 6630, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)

	/*
	 * Фиксированная наценка нулевая. К ней добавляется процентная без ограничения
	 * Доли копейки округляются по математическому округлению, у взрослого вверх, а ребенка - вниз
	 */
	params.TestOfferPrice = base.Money{Amount: 415722, Currency: &base.Currency{Code: "RUB", Fraction: 100}}
	result = frule.GetResult(params)
	// 0RUR+2.3% = 415722/100*2,3 = 9561,606
	assert.Equal(t, base.Money{Amount: 9562, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Full)
	// 0RUR+2.1% = 415722/100*2,1 = 8730,162
	assert.Equal(t, base.Money{Amount: 8730, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Child)
	assert.Equal(t, base.Money{Amount: 0, Currency: &base.Currency{Code: "RUB", Fraction: 100}}, result.(ServiceChargeRuleResult).Margin.Infant)
}
