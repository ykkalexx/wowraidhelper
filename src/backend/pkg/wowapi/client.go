package wowapi

import (
	"context"

	"ykkalexx.com/wowraidhelper/internal/model"
)

type Client struct {
	clientID     string
	clientSecret string
	accessToken  string
	baseURL      string
}

// PlayerDetails represents temporary player data fetched from WoW API
type PlayerDetails struct {
	Name           string   `json:"name"`
	Class          string   `json:"class"`
	Spec           string   `json:"spec"`
	ItemLevel      float64  `json:"item_level"`
	Role           string   `json:"role"` // Tank, Healer, DPS
	RaidExperience []string `json:"raid_experience"`
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      "https://api.blizzard.com",
	}
}

func (c *Client) GetRaidInfo(ctx context.Context, raidName string) (*model.RaidInfo, error) {
	// Implementation to fetch raid information from WoW API
	return nil, nil
}

func (c *Client) GetPlayerDetails(ctx context.Context, playerName string, realm string) (*PlayerDetails, error) {
	// Implementation to fetch player details from WoW API
	return nil, nil
}