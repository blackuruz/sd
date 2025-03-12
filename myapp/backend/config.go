package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Lokasi penyimpanan data
const configFileName = "config.json"

// StreamConfig menyimpan informasi stream
type StreamConfig struct {
	StreamKey   string `json:"stream_key"`
	ChannelName string `json:"channel_name"`
	VideoFile   string `json:"video_file"`
}

// ConfigManager mengelola konfigurasi
type ConfigManager struct {
	FilePath string
}

// NewConfigManager membuat instance baru
func NewConfigManager() *ConfigManager {
	dir, _ := os.Getwd()
	return &ConfigManager{
		FilePath: filepath.Join(dir, configFileName),
	}
}

// LoadConfigs membaca konfigurasi dari file
func (cm *ConfigManager) LoadConfigs() ([]StreamConfig, error) {
	file, err := os.ReadFile(cm.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []StreamConfig{}, nil
		}
		return nil, err
	}

	var configs []StreamConfig
	err = json.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

// SaveConfigs menyimpan konfigurasi ke file
func (cm *ConfigManager) SaveConfigs(configs []StreamConfig) error {
	data, err := json.MarshalIndent(configs, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.FilePath, data, 0644)
}

// AddConfig menambahkan konfigurasi baru
func (cm *ConfigManager) AddConfig(config StreamConfig) error {
	configs, err := cm.LoadConfigs()
	if err != nil {
		return err
	}

	// Periksa apakah StreamKey sudah ada
	for _, existing := range configs {
		if existing.StreamKey == config.StreamKey {
			return errors.New("stream key already exists")
		}
	}

	configs = append(configs, config)
	return cm.SaveConfigs(configs)
}

// GetConfigs mengambil semua konfigurasi
func (cm *ConfigManager) GetConfigs() ([]StreamConfig, error) {
	return cm.LoadConfigs()
}

// EditConfig mengedit konfigurasi berdasarkan StreamKey
func (cm *ConfigManager) EditConfig(updatedConfig StreamConfig) error {
	configs, err := cm.LoadConfigs()
	if err != nil {
		return err
	}

	found := false
	for i, config := range configs {
		if config.StreamKey == updatedConfig.StreamKey {
			configs[i] = updatedConfig
			found = true
			break
		}
	}

	if !found {
		return errors.New("stream config not found")
	}

	return cm.SaveConfigs(configs)
}

// DeleteConfig menghapus konfigurasi berdasarkan StreamKey
func (cm *ConfigManager) DeleteConfig(streamKey string) error {
	configs, err := cm.LoadConfigs()
	if err != nil {
		return err
	}

	newConfigs := []StreamConfig{}
	for _, config := range configs {
		if config.StreamKey != streamKey {
			newConfigs = append(newConfigs, config)
		}
	}

	if len(newConfigs) == len(configs) {
		return errors.New("stream key not found")
	}

	return cm.SaveConfigs(newConfigs)
}
