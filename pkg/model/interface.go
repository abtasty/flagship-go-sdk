package model

// HitInterface express the interface for the hits
type HitInterface interface {
	Validate() []error
	SetBaseInfos(envID string, visitorID string)
	getBaseHit() BaseHit
	resetBaseHit()
	ComputeQueueTime()
}
