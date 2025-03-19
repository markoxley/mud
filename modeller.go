package dtorm

type Modeller interface {
	// StandingData returns the standing data for the model
	StandingData() []Modeller

	// GetID returns the ID of the model
	GetID() *string

	// IsNew returns true if the model has yet to be saved
	IsNew() bool

	// IsDeleted returns true if the model has been marked as deleted
	IsDeleted() bool

	Disable()
}
