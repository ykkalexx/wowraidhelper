package repository

import (
	"database/sql"
	"encoding/json"
)

type Repository struct {
	db *sql.DB
}

type RaidInfo struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    BossCount    int    `json:"boss_count"`
    MinItemLevel int    `json:"min_item_level"`
    Difficulty   string `json:"difficulty"`
}

type RaidAnalysis struct {
    ID           int             `json:"id"`
    RaidID       int             `json:"raid_id"`
    PlayerCount  int             `json:"player_count"`
    AvgItemLevel float64         `json:"avg_item_level"`
    Composition  json.RawMessage `json:"composition"`
    Predictions  json.RawMessage `json:"predictions"`
}

func NewRepository(dsn string) (*Repository, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return &Repository{db: db}, nil
}

func (r *Repository) SaveRaidInfo(info *RaidInfo) error {
    query := `
        INSERT INTO raid_info (name, boss_count, min_item_level, difficulty)
        VALUES (?, ?, ?, ?)
    `
    
    result, err := r.db.Exec(query, info.Name, info.BossCount, info.MinItemLevel, info.Difficulty)
    if err != nil {
        return err
    }

    id, _ := result.LastInsertId()
    info.ID = int(id)
    return nil
}

func (r *Repository) SaveRaidAnalysis(analysis *RaidAnalysis) error {
    query := `
        INSERT INTO raid_analysis 
        (raid_id, player_count, avg_item_level, composition, predictions)
        VALUES (?, ?, ?, ?, ?)
    `
    
    result, err := r.db.Exec(query,
        analysis.RaidID,
        analysis.PlayerCount,
        analysis.AvgItemLevel,
        analysis.Composition,
        analysis.Predictions,
    )
    if err != nil {
        return err
    }

    id, _ := result.LastInsertId()
    analysis.ID = int(id)
    return nil
}

func (r *Repository) GetRaidInfo(name string) (*RaidInfo, error) {
    var info RaidInfo
    query := `
        SELECT id, name, boss_count, min_item_level, difficulty
        FROM raid_info
        WHERE name = ?
    `
    
    err := r.db.QueryRow(query, name).Scan(
        &info.ID,
        &info.Name,
        &info.BossCount,
        &info.MinItemLevel,
        &info.Difficulty,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &info, nil
}