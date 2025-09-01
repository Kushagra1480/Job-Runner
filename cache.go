package main

import (
	"os"
	"path/filepath"
	"time"
	"encoding/json"
)

type CacheData struct {
	Date 				string `json:"date"`
	CategoryCompanyMap 	map[string][]Company `json:"category_company_map"`
	CategoriesLinkMap 	map[string]string `json:"categories_link_map"`
}

func GetCacheFilePath() string {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "startup-gallery")
	os.MkdirAll(cacheDir, 0755)
	return filepath.Join(cacheDir, "data.json")
}

func IsCacheValid(cacheFilePath string) bool {
	info, err := os.Stat(cacheFilePath)
	if err != nil {
		return false
	}
	today := time.Now().Format("2006-01-02")
	fileDate := info.ModTime().Format("2006-01-02")

	return today == fileDate
}

func LoadCache(cacheFilePath string) (*CacheData, error) {
	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return nil, err
	}
	var cache CacheData
	err = json.Unmarshal(data, &cache)
	return &cache, err
}

func SaveData(cacheFilePath string, cache *CacheData) error {
	data, err := json.MarshalIndent(cache, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(cacheFilePath, data, 0644)
}
