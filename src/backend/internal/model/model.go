package model

type Player struct {
	ID             uint    `json:"id" gorm:"primaryKey"`
	Name           string  `json:"name"`
	Class          string  `json:"class"`
	ItemLevel      float64 `json:"item_level"`
	RaidID         uint    `json:"raid_id"`
	GearScore      int     `json:"gear_score"`
	Specialization string  `json:"specialization"`
}

type RaidInfo struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`       // Raid name (e.g., "Nerub-ar Palace")
	Difficulty   string `json:"difficulty"` // Normal, Heroic, Mythic
	BossCount    int    `json:"boss_count"`
	MinItemLevel int    `json:"min_item_level"`
	Mechanics    string `json:"mechanics"` // JSON string of important raid mechanics
}

type RaidAnalysis struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	RaidInfoID  uint   `json:"raid_info_id"`
	PlayerCount int    `json:"player_count"`
	Composition string `json:"composition"` // JSON string of role distribution
	Predictions string `json:"predictions"` // JSON string of ML predictions
	Tips        string `json:"tips"`        // JSON string of generated tips
	CreatedAt   int64  `json:"created_at"`
}