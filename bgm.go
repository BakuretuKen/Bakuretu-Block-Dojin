package main

import "fmt"

// StopBGM は現在再生中のBGMを停止して破棄する
func StopBGM() {
	if bgmPlayer != nil {
		bgmPlayer.Close()
		bgmPlayer = nil
	}
}

// PlayBGM は既存BGMを停止してから新しいBGMを再生する
func PlayBGM(path string, volume float64) error {
	StopBGM()

	stream, err := LoadAudioStream(path)
	if err != nil {
		return fmt.Errorf("BGMの読み込みに失敗: %w", err)
	}

	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		return fmt.Errorf("BGMプレイヤーの作成に失敗: %w", err)
	}

	player.SetVolume(volume)
	player.Play()
	bgmPlayer = player

	return nil
}

// EnsureBGMLooping はBGMが停止していたら先頭に戻して再生する
func EnsureBGMLooping() {
	if bgmPlayer != nil && !bgmPlayer.IsPlaying() {
		bgmPlayer.Rewind()
		bgmPlayer.Play()
	}
}
