package wowapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
    clientID     string
    clientSecret string
    accessToken  string
}

type PlayerInfo struct {
    Name      string  `json:"name"`
    Class     string  `json:"class"`
    ItemLevel float64 `json:"average_item_level"`
    Spec      string  `json:"active_spec"`
    Role      string  `json:"role"`        // Tank, Healer, DPS
}

type RaidInfo struct {
    Name        string   `json:"name"`
    Bosses      int      `json:"boss_count"`
    Difficulty  string   `json:"difficulty"`
    MinItemLevel int     `json:"minimum_item_level"`
}

func NewClient(clientID, clientSecret string) *Client {
    return &Client{
        clientID:     clientID,
        clientSecret: clientSecret,
    }
}

// Gets token from Blizzard
func (c *Client) authenticate() error {
    data := strings.NewReader("grant_type=client_credentials")
    req, err := http.NewRequest("POST", "https://oauth.battle.net/token", data)
    if err != nil {
        return err
    }

    req.SetBasicAuth(c.clientID, c.clientSecret)
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var result struct {
        AccessToken string `json:"access_token"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return err
    }

    c.accessToken = result.AccessToken
    return nil
}

// GetPlayerInfo fetches basic player info from EU servers
func (c *Client) GetPlayerInfo(name, realm string) (*PlayerInfo, error) {
    if c.accessToken == "" {
        if err := c.authenticate(); err != nil {
            return nil, err
        }
    }

    url := fmt.Sprintf("https://eu.api.blizzard.com/profile/wow/character/%s/%s?namespace=profile-eu",
        strings.ToLower(realm),
        strings.ToLower(name))

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("Authorization", "Bearer "+c.accessToken)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == 401 { // Token expired
        if err := c.authenticate(); err != nil {
            return nil, err
        }
        // Retry request with new token
        req.Header.Set("Authorization", "Bearer "+c.accessToken)
        resp, err = http.DefaultClient.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
    }

    var player PlayerInfo
    if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
        return nil, err
    }

    return &player, nil
}

// GetRaidInfo fetches raid information
func (c *Client) GetRaidInfo(raidName string) (*RaidInfo, error) {
    if c.accessToken == "" {
        if err := c.authenticate(); err != nil {
            return nil, err
        }
    }

    url := fmt.Sprintf("https://eu.api.blizzard.com/data/wow/raid/%s?namespace=static-eu",
        strings.ToLower(raidName))

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("Authorization", "Bearer "+c.accessToken)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == 401 { // Token expired
        if err := c.authenticate(); err != nil {
            return nil, err
        }
        // Retry request with new token
        req.Header.Set("Authorization", "Bearer "+c.accessToken)
        resp, err = http.DefaultClient.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
    }

    var raid RaidInfo
    if err := json.NewDecoder(resp.Body).Decode(&raid); err != nil {
        return nil, err
    }

    return &raid, nil
}