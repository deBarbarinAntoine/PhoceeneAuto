package data

import "time"

type CarCatalog struct {
	CatID            uint
	CatCreatedAt     time.Time
	CatUpdatedAt     time.Time
	Make             string
	Model            string
	Cylinders        uint
	Drive            string
	EngineDescriptor string
	Fuel1            string
	Fuel2            string
	LuggageVolume    float32
	PassengerVolume  float32
	Transmission     string
	SizeClass        string
	Year             uint
	ElectricMotor    float32
	BaseModel        string
}
