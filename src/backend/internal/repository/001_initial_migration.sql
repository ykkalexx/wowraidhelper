CREATE DATABASE wowraidhelper;
USE wowraidhelper;

CREATE TABLE raid_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    boss_count INT NOT NULL,
    min_item_level INT NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE player_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    class VARCHAR(50) NOT NULL,
    ilvl INT NOT NULL,
    raidId INT NOT NULL,
    gear_score INT NOT NULL,
    speci VARCHAR(255) NOT NULL,
)

CREATE TABLE raid_analysis (
    id INT AUTO_INCREMENT PRIMARY KEY,
    raid_id INT NOT NULL,
    player_count INT NOT NULL,
    avg_item_level FLOAT NOT NULL,
    composition JSON NOT NULL,  -- Store role distribution
    predictions JSON NOT NULL,  -- Store ML predictions
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (raid_id) REFERENCES raid_info(id)
);