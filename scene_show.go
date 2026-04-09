package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ShowScene はゲームクリア画面を表す
type ShowScene struct {
	bgImage  *ebiten.Image // 背景画像
	gameCode string        // ゲームコード
	number   int           // パターン番号
}

// NewTitleScene はタイトル画面を初期化して返す
func NewShowScene(gameCode string, number int) Scene {
	scene := &ShowScene{}
	scene.gameCode = gameCode
	scene.number = number
	var err error

	// 背景画像を読み込む
	var imageName string
	if scene.number == 0 {
		imageName = "game_image.jpg"
	} else if scene.number == 1 {
		imageName = "game_back_image.jpg"
	} else if scene.number == 2 {
		imageName = "win_back_image.jpg"
	}
	scene.bgImage, err = LoadEbitenImage("assets/" + gameCode + "/" + imageName)
	if err != nil {
		return NewErrorScene("Failed to load win background image: " + err.Error())
	}

	return scene
}

func (s *ShowScene) Update() Scene {
	HandleWindowRestore()

	// BGMのループ再生を管理（グローバルプレイヤーを使用）
	EnsureBGMLooping()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		PlaySoundEffect("assets/button_select_se.mp3", 0.1)
		if s.number == 2 {
			return NewAlbumScene()
		}
		return NewShowScene(s.gameCode, s.number+1) // 次の画像
	}
	return nil
}

func (s *ShowScene) Draw(screen *ebiten.Image) {
	// 指定画像を描画
	screen.DrawImage(s.bgImage, &ebiten.DrawImageOptions{})
}
