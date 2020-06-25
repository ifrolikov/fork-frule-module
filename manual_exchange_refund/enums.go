package manual_exchange_refund

const (
	UsedTypeNotUsed                  UsedType                = "not_used"
	UsedTypePartiallyUsed            UsedType                = "partially_used"
	TariffStartTypeSinceBuy          TariffStartType         = "since_buy"
	TariffStartTypeSinceFly          TariffStartType         = "since_fly"
	TariffCalculationSourceCurrent   TariffCalculationSource = "current"
	TariffCalculationSourceBase      TariffCalculationSource = "base"
	PenaltyStrategySelfSegment       PenaltyStrategy         = "self_segment"
	PenaltyStrategyNearestSegment    PenaltyStrategy         = "nearest_segment"
	FlightTypeInternal               FlightType              = "internal"
	FlightTypeInternational          FlightType              = "international"
	CalculationUnitDefault           CalculationUnit         = "default"
	CalculationUnitPU                CalculationUnit         = "pu"
	ApplyStrategySummary             ApplyStrategy           = "summary"
	ApplyStrategyExpensive           ApplyStrategy           = "expensive"
	FeeDestinationTicket             FeeDestination          = "ticket"
	FeeDestinationRoute              FeeDestination          = "route"
	FeeDestinationSegment            FeeDestination          = "segment"
	FarePercentDestinationSegment    FarePercentDestination  = "segment"
	FarePercentDestinationRoute      FarePercentDestination  = "route"
	FarePercentDestinationFlight     FarePercentDestination  = "flight"
	FarePercentDestinationHalfFlight FarePercentDestination  = "half_flight"
	ContextExchange                  Context                 = "exchange"
	ContextRefund                    Context                 = "refund"
)

type FarePercentDestination string
type FeeDestination string
type ApplyStrategy string
type CalculationUnit string
type FlightType string
type PenaltyStrategy string
type TariffCalculationSource string
type TariffStartType string
type UsedType string
type Context string
