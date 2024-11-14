package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"ykkalexx.com/wowraidhelper/internal/repository"
	"ykkalexx.com/wowraidhelper/pkg/wowapi"
)

type RaidAnalysisRequest struct {
    RaidName    string   `json:"raid_name"`
    PlayerNames []string `json:"player_names"` // Format: "name-realm"
}

type RaidService struct {
    wowClient  *wowapi.Client
    repository *repository.Repository
}

func NewRaidService(wowClient *wowapi.Client, repo *repository.Repository) *RaidService {
    return &RaidService{
        wowClient:  wowClient,
        repository: repo,
    }
}

func (s *RaidService) AnalyzeRaid(req RaidAnalysisRequest) (*repository.RaidAnalysis, error) {
    // 1. Fetch raid info from WoW API
    raidInfo, err := s.wowClient.GetRaidInfo(req.RaidName)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch raid info: %w", err)
    }

    // 2. Save or update raid info in database
    dbRaid := &repository.RaidInfo{
        Name:         raidInfo.Name,
        BossCount:    raidInfo.Bosses,
        MinItemLevel: raidInfo.MinItemLevel,
        Difficulty:   raidInfo.Difficulty,
    }
    
    if err := s.repository.SaveRaidInfo(dbRaid); err != nil {
        return nil, fmt.Errorf("failed to save raid info: %w", err)
    }

    // 3. Fetch all player details
    var totalItemLevel float64
    roleCount := map[string]int{
        "tank":   0,
        "healer": 0,
        "dps":    0,
    }

    for _, playerName := range req.PlayerNames {
        name, realm := parsePlayerString(playerName)
        player, err := s.wowClient.GetPlayerInfo(name, realm)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch player %s details: %w", playerName, err)
        }

        totalItemLevel += player.ItemLevel
        roleCount[strings.ToLower(player.Role)]++
    }

    // 4. Create basic analysis (without ML for now)
    composition, err := json.Marshal(roleCount)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal composition: %w", err)
    }

    // Create simple predictions for now
    basicPredictions := map[string]string{
        "status": "pending_ml_implementation",
        "note":   "Basic composition analysis only",
    }
    predictions, err := json.Marshal(basicPredictions)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal predictions: %w", err)
    }

    analysis := &repository.RaidAnalysis{
        RaidID:       dbRaid.ID,
        PlayerCount:  len(req.PlayerNames),
        AvgItemLevel: totalItemLevel / float64(len(req.PlayerNames)),
        Composition:  composition,
        Predictions:  predictions,
    }

    // 5. Save analysis
    if err := s.repository.SaveRaidAnalysis(analysis); err != nil {
        return nil, fmt.Errorf("failed to save analysis: %w", err)
    }

    return analysis, nil
}

func parsePlayerString(playerString string) (name, realm string) {
    parts := strings.Split(playerString, "-")
    if len(parts) != 2 {
        return playerString, "defaultRealm"
    }
    return parts[0], parts[1]
}