package model

type Orbit struct {
	BaseModel
	OrbiterID string `gorm:"type:uuid;uniqueIndex:idx_orbiter_orbited"`
	OrbitedID string `gorm:"type:uuid;uniqueIndex:idx_orbiter_orbited"`
}
