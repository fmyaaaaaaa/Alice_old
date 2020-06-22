package enum

type OrderTriggerCondition string

const (
	DefaultOrderTriggerCondition = OrderPositionFill("DEFAULT")
	Inverse                      = OrderTriggerCondition("INVERSE")
	Bid                          = OrderTriggerCondition("BID")
	ASK                          = OrderTriggerCondition("ASK")
	Mid                          = OrderTriggerCondition("MID")
)
