package types

type MeasurementDataType string

const (
	Temperature    MeasurementDataType = "temperature"
	Humidity                           = "cap"
	BatteryVoltage                     = "batteryVolt"
	Signal                             = "signal"
)
