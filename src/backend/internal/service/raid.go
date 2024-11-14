package service

import (
	"context"
	"fmt"

	"ykkalexx.com/wowraidhelper/internal/model"
	"ykkalexx.com/wowraidhelper/pkg/wowapi"
)

type RaidAnalysisRequest struct {
    RaidName    string   `json:"raid_name"`
    PlayerNames []string `json:"player_names"` // Format: "name-realm"
}

type RaidService struct {
    wowClient  *wowapi.Client
    repository *repository.MySQLRepository
    mlClient   *MLClient  // ill implement this later for Python service communication
}

func (s *RaidService) AnalyzeRaid(ctx context.Context, req RaidAnalysisRequest) (*model.RaidAnalysis, error) {
    // 1. Fetch raid info
    raidInfo, err := s.wowClient.GetRaidInfo(ctx, req.RaidName)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch raid info: %w", err)
    }

    // 2. Fetch all player details
    players := make([]*wowapi.PlayerDetails, 0, len(req.PlayerNames))
    for _, playerName := range req.PlayerNames {
        name, realm := parsePlayerString(playerName) // "name-realm" format
        player, err := s.wowClient.GetPlayerDetails(ctx, name, realm)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch player %s details: %w", playerName, err)
        }
        players = append(players, player)
    }

    // 3. Prepare data for ML analysis
    mlData := prepareMlData(raidInfo, players)

    // 4. Send to ML service and get predictions (ill implement this later)
    analysis, err := s.mlClient.AnalyzeRaid(ctx, mlData)
    if err != nil {
        return nil, fmt.Errorf("failed to perform ML analysis: %w", err)
    }

    // 5. Save and return analysis
    return s.repository.SaveAnalysis(ctx, analysis)
}

func parsePlayerString(playerString string) (name, realm string) {
    // Implementation to split "name-realm" string
    return
}

func prepareMlData(raid *model.RaidInfo, players []*wowapi.PlayerDetails) map[string]interface{} {
    // Implementation to prepare data for ML service
    return nil
}