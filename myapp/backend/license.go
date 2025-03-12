package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	activationURL   = "https://omahbot.com/member/softsale/api/activate?key=%s&request[hardware-id]=%s"
	deactivationURL = "https://omahbot.com/member/softsale/api/deactivate?key=%s&request[hardware-id]=%s"
	licenseFileName = "activation.json"
)

type LicenseInfo struct {
	LicenseKey     string `json:"license_key"`
	OmahbotKey     string `json:"omahbotkey"`
	ActivationCode string `json:"activation_code"`
}

type LicenseManager struct {
	FilePath string
}

func NewLicenseManager() *LicenseManager {
	appData := os.Getenv("APPDATA")
	activationDir := filepath.Join(appData, "WailsStreamer")
	if err := os.MkdirAll(activationDir, 0755); err != nil {
		fmt.Printf("Gagal membuat direktori aktivasi: %v\n", err)
	}
	return &LicenseManager{FilePath: filepath.Join(activationDir, licenseFileName)}
}

func (lm *LicenseManager) Save(info LicenseInfo) error {
	data, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(lm.FilePath, data, 0644)
}

func (lm *LicenseManager) Load() (*LicenseInfo, error) {
	if _, err := os.Stat(lm.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file lisensi tidak ditemukan")
	}
	data, err := ioutil.ReadFile(lm.FilePath)
	if err != nil {
		return nil, err
	}
	var info LicenseInfo
	err = json.Unmarshal(data, &info)
	return &info, err
}

func (lm *LicenseManager) ValidateActivation() bool {
	info, err := lm.Load()
	if err != nil {
		return false
	}
	hwid, err := getHardwareID()
	if err != nil {
		return false
	}
	return info.OmahbotKey == hwid && info.ActivationCode != ""
}

func (lm *LicenseManager) ActivateLicense(licenseKey string) error {
	hwid, _ := getHardwareID()
	url := fmt.Sprintf(activationURL, licenseKey, hwid)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if code, ok := data["code"].(string); !ok || code != "ok" {
		return fmt.Errorf("aktivasi gagal")
	}

	info := LicenseInfo{
		LicenseKey:     licenseKey,
		OmahbotKey:     hwid,
		ActivationCode: data["activation_code"].(string),
	}
	return lm.Save(info)
}

func (lm *LicenseManager) DeactivateLicense() error {
	info, err := lm.Load()
	if err != nil {
		return err
	}
	url := fmt.Sprintf(deactivationURL, info.LicenseKey, info.OmahbotKey)
	_, err = http.Get(url)
	return err
}

func getHardwareID() (string, error) {
	hostname, _ := os.Hostname()
	hash := md5.Sum([]byte(hostname))
	return hex.EncodeToString(hash[:]), nil
}
