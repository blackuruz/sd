package backend

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

// Konfigurasi file
const configFileName = "config.json"

// StreamConfig menyimpan informasi stream
type StreamConfig struct {
	StreamKey   string `json:"stream_key"`
	ChannelName string `json:"channel_name"`
	VideoFile   string `json:"video_file"`
	Quality     string `json:"quality"`
	// Jika diperlukan, Anda bisa menambahkan field lain, misal: durasi, tanggal, dsb.
	StartDate string `json:"start_date"`
	StartTime string `json:"start_time"`
	EndDate   string `json:"end_date"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}

// ConfigManager mengelola konfigurasi
type ConfigManager struct {
	FilePath string
}

// NewConfigManager membuat instance baru dan menetapkan path file config
func NewConfigManager() *ConfigManager {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error mendapatkan working directory:", err)
		dir = "."
	}
	log.Printf("Working directory: %s", dir)
	return &ConfigManager{
		FilePath: filepath.Join(dir, configFileName),
	}
}

// LoadConfigs membaca konfigurasi dari file
func (cm *ConfigManager) LoadConfigs() ([]StreamConfig, error) {
	if _, err := os.Stat(cm.FilePath); os.IsNotExist(err) {
		log.Printf("File config %s tidak ditemukan", cm.FilePath)
		return []StreamConfig{}, nil
	}
	data, err := os.ReadFile(cm.FilePath)
	if err != nil {
		return nil, err
	}
	log.Printf("Isi file config: %s", string(data))
	var configs []StreamConfig
	err = json.Unmarshal(data, &configs)
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
	log.Printf("Menyimpan config: %s", string(data))
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

// GetConfigs mengembalikan semua konfigurasi
func (cm *ConfigManager) GetConfigs() ([]StreamConfig, error) {
	return cm.LoadConfigs()
}

// EditConfig memperbarui konfigurasi berdasarkan StreamKey
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
