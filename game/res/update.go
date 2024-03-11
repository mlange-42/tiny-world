package res

// UpdateInterval resource.
type UpdateInterval struct {
	// The interval between updates of entities, in game ticks. Usually equal to TPS.
	Interval int64
	// Number of intervals used in (production) countdowns. Usually 60, resulting in 1 minute.
	Countdown int
}
