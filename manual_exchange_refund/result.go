package manual_exchange_refund

type ManualExchangeRefundResult struct {
	IsFound                 bool
	Id                      int64
	FeeDestination          *FeeDestination
	ApplyStrategy           *ApplyStrategy
	FarePercent             *float64
	FarePercentDestination  *FarePercentDestination
	Amount                  *int64
	Currency                *string
	CalculationUnit         *CalculationUnit
	Brand                   *string
	TariffCalculationSource *TariffCalculationSource
	IsAvailable             bool
}

func NewManualExchangeRefundResult(
	id int64,
	feeDestination *FeeDestination,
	applyStrategy *ApplyStrategy,
	farePercent *float64,
	farePercentDestination *FarePercentDestination,
	amount *int64,
	currency *string,
	calculationUnit *CalculationUnit,
	brand *string,
	tariffCalculationSource *TariffCalculationSource,
	isAvailable bool,
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
		IsAvailable:             isAvailable,
	}
}
