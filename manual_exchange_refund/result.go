package manual_exchange_refund

type ManualExchangeRefundResult struct {
	IsFound                 bool
	Id                      int32
	FeeDestination          FeeDestination
	ApplyStrategy           ApplyStrategy
	FarePercent             *float64
	FarePercentDestination  *FarePercentDestination
	Amount                  *int32
	Currency                *string
	CalculationUnit         *CalculationUnit
	Brand                   *string
	TariffCalculationSource *TariffCalculationSource
}

func NewManualExchangeRefundResult(
	id int32,
	feeDestination FeeDestination,
	applyStrategy ApplyStrategy,
	farePercent *float64,
	farePercentDestination *FarePercentDestination,
	amount *int32,
	currency *string,
	calculationUnit *CalculationUnit,
	brand *string,
	tariffCalculationSource *TariffCalculationSource,
) ManualExchangeRefundResult {
	return ManualExchangeRefundResult{
		IsFound:                 true,
		Id:                      id,
		FeeDestination:          feeDestination,
		ApplyStrategy:           applyStrategy,
		FarePercent:             farePercent,
		FarePercentDestination:  farePercentDestination,
		Amount:                  amount,
		Currency:                currency,
		CalculationUnit:         calculationUnit,
		Brand:                   brand,
		TariffCalculationSource: tariffCalculationSource,
	}
}
