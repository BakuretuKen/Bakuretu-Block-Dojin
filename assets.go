package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

//go:embed assets/*
var assetsFS embed.FS

// LoadImage は埋め込みアセットから生の画像を読み込む
func LoadImage(path string) (image.Image, error) {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// LoadEbitenImage は埋め込みアセットから画像を読み込む
func LoadEbitenImage(path string) (*ebiten.Image, error) {
	img, err := LoadImage(path)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

// LoadAudioStream は埋め込みアセットからオーディオストリームを読み込む
func LoadAudioStream(path string) (io.ReadSeeker, error) {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp3":
		return mp3.DecodeWithoutResampling(bytes.NewReader(data))
	case ".wav":
		return wav.DecodeWithoutResampling(bytes.NewReader(data))
	default:
		return nil, fmt.Errorf("未対応の音声形式です: %s", ext)
	}
}

// LoadAudioBytes は埋め込みアセットからデコード済みPCMを読み込む
func LoadAudioBytes(path string) ([]byte, error) {
	stream, err := LoadAudioStream(path)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(stream)
}

// PlaySoundEffect は効果音を一度だけ再生する
// filePath: 効果音ファイルのパス（例: "assets/button_click.mp3"）
// volume: 音量（0.0〜1.0）
func PlaySoundEffect(filePath string, volume float64) error {
	player, err := getOrCreateSEPlayer(filePath)
	if err != nil {
		return err
	}
	player.SetVolume(volume)
	player.Rewind()
	player.Play()
	return nil
}

func getOrCreateSEPlayer(filePath string) (*audio.Player, error) {
	if player, ok := sePlayers[filePath]; ok {
		return player, nil
	}

	stream, err := LoadAudioStream(filePath)
	if err != nil {
		return nil, err
	}

	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		return nil, err
	}
	sePlayers[filePath] = player
	return player, nil
}

// PlaySoundEffectFromBytes はデコード済みPCMを再生する
func PlaySoundEffectFromBytes(data []byte, volume float64) error {
	player, err := audioContext.NewPlayer(bytes.NewReader(data))
	if err != nil {
		return err
	}

	player.SetVolume(volume)
	player.Play()

	go func() {
		for player.IsPlaying() {
			time.Sleep(10 * time.Millisecond)
		}
		player.Close()
	}()

	return nil
}

func StopAllSE() {
	for key, player := range sePlayers {
		if player != nil {
			player.Close()
		}
		delete(sePlayers, key)
	}
}
