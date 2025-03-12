package main

import (
	"context"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"myapp/backend"
)

// App adalah struct utama Wails
type App struct {
	ctx            context.Context
	ConfigManager  *backend.ConfigManager
}

// NewApp membuat instance baru aplikasi
func NewApp() *App {
	return &App{
		ConfigManager: backend.NewConfigManager(),
	}
}

// Startup adalah fungsi yang dijalankan saat aplikasi Wails dimulai
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// GetConfigs mengirim daftar konfigurasi ke frontend
func (a *App) GetConfigs() []backend.StreamConfig {
	configs, err := a.ConfigManager.GetConfigs()
	if err != nil {
		runtime.LogError(a.ctx, "Gagal mengambil konfigurasi: "+err.Error())
		return nil
	}
	return configs
}

// AddConfig menambahkan konfigurasi streaming baru
func (a *App) AddConfig(config backend.StreamConfig) string {
	err := a.ConfigManager.AddConfig(config)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Config added successfully"
}

// EditConfig memperbarui konfigurasi yang ada
func (a *App) EditConfig(config backend.StreamConfig) string {
	err := a.ConfigManager.EditConfig(config)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Config updated successfully"
}

// DeleteConfig menghapus konfigurasi berdasarkan Stream Key
func (a *App) DeleteConfig(streamKey string) string {
	err := a.ConfigManager.DeleteConfig(streamKey)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Config deleted successfully"
}

// Fungsi utama Wails
func main() {
	// Buat instance aplikasi
	app := NewApp()

	// Konfigurasi aplikasi Wails
	err := wails.Run(&options.App{
		Title:  "YouTube Streamer",
		Width:  800,
		Height: 600,
		Bind: []interface{}{
			app,
		},
	})

	// Handle error jika Wails gagal dijalankan
	if err != nil {
		log.Fatal("Error saat menjalankan Wails:", err)
	}
}
