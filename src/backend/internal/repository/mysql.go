package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
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

func NewRepository(user, password, host, port, dbName string) (*Repository, error) {
    // Format: username:password@tcp(hostname:port)/dbname
    dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
        user, password, host, port, dbName)
        
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %v", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Verify connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to database: %v", err)
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

func (r *Repository) Close() error {
    return r.db.Close()
}
