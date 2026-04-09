package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TitleScene はタイトル画面を表す
type LoseScene struct {
	bgImage           *ebiten.Image // 背景画像
	returnButtonImage *Sprite       // 戻るボタン画像
}

// NewTitleScene はタイトル画面を初期化して返す
func NewLoseScene(gameCode string) Scene {
	scene := &LoseScene{}

	// 背景画像を読み込む
	var err error
	scene.bgImage, err = LoadEbitenImage("assets/" + gameCode + "/lose_back_image.jpg")
	if err != nil {
		return NewErrorScene("Failed to load lose background image: " + err.Error())
	}

	// 戻る画像を読み込む
	returnButtonImage, err := LoadEbitenImage("assets/return_menu_button.png")
	if err != nil {
		return NewErrorScene("Failed to load return button image: " + err.Error())
	}
	scene.returnButtonImage = NewSprite(returnButtonImage, 2)
	scene.returnButtonImage.x = 20
	scene.returnButtonImage.y = 1080 - 84 - 20

	return scene
}

// Update は敗北画面のロジックを更新する
func (s *LoseScene) Update() Scene {
	HandleWindowRestore()

	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if s.returnButtonImage.Contains(mx, my) {
			StopBGM()
			PlaySoundEffect("assets/button_select_se.mp3", 0.1)
			return NewTitleScene()
		}
	}
	return nil // 敗北画面に留まる
}

// Draw は敗北画面を描画する
func (s *LoseScene) Draw(screen *ebiten.Image) {
	// 背景画像を描画
	screen.DrawImage(s.bgImage, &ebiten.DrawImageOptions{})

	// 戻るボタンを描画
	s.returnButtonImage.Draw(screen)
	mx, my := ebiten.CursorPosition()
	if s.returnButtonImage.Contains(mx, my) {
		s.returnButtonImage.SetFrame(1)
	} else {
		s.returnButtonImage.SetFrame(0)
	}
}
