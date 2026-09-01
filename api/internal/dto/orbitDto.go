package dto

type EnterOrbitRequest struct {
	OrbiterID string
	OrbitedID string
}

type EnterOrbitResponse struct {
	IsOrbiting bool `json:"is_orbiting"`
}

type LeaveOrbitRequest struct {
	OrbiterID string
	OrbitedID string
}

type LeaveOrbitResponse struct {
	IsOrbiting bool `json:"is_orbiting"`
}

type IsOrbitingRequest struct {
	OrbiterID string
	OrbitedID string
}

type IsOrbitingResponse struct {
	IsOrbiting bool `json:"is_orbiting"`
}
