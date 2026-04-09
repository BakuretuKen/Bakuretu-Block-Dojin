// Windows用 AppData 読み書きライブラリ
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var appDirName string

// SetAppDirName アプリケーションのディレクトリ名を設定します
// main関数の最初で一度だけ呼び出してください
func SetAppDirName(name string) {
	appDirName = name
}

func getUserSaveDir(appName string) (string, error) {
	if appName == "" {
		return "", fmt.Errorf("appDirName is not initialized: call SetAppDirName at startup")
	}

	// Windows の場合
	appData := os.Getenv("LOCALAPPDATA")
	if appData != "" {
		saveDir := filepath.Join(appData, appName)
		return saveDir, nil
	}
	// Macの場合
	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		saveDir := filepath.Join(homeDir, "Library", "Application Support", appName)
		return saveDir, nil
	}
	return "", fmt.Errorf("USER_HOME not defined")
}

func SaveGameData(filename string, data string) error {
	dataBytes := []byte(data)
	dir, err := getUserSaveDir(appDirName)
	if err != nil {
		return err
	}
	// ディレクトリが存在しない場合は作る
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, filename)
	return os.WriteFile(path, dataBytes, 0644)
}

func LoadGameData(filename string) (string, error) {
	dir, err := getUserSaveDir(appDirName)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, filename)
	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(dataBytes), nil
}
