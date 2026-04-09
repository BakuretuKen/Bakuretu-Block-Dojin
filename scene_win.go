package main

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TitleScene はタイトル画面を表す
type WinScene struct {
	bgImage           *ebiten.Image // 背景画像
	returnButtonImage *Sprite       // 戻るボタン画像
}

// NewTitleScene はタイトル画面を初期化して返す
func NewWinScene(gameCode string) Scene {
	scene := &WinScene{}
	var err error

	// ゲームデータを更新
	err = scene.updateAppData(gameCode)
	if err != nil {
		return NewErrorScene("Failed to update game data: " + err.Error())
	}
	// 背景画像を読み込む
	scene.bgImage, err = LoadEbitenImage("assets/" + gameCode + "/win_back_image.jpg")
	if err != nil {
		return NewErrorScene("Failed to load win background image: " + err.Error())
	}

	// 戻る画像を読み込む
	returnButtonImage, err := LoadEbitenImage("assets/return_menu_button.png")
	if err != nil {
		return NewErrorScene("Failed to load return button image: " + err.Error())
	}
	scene.returnButtonImage = NewSprite(returnButtonImage, 2)
	scene.returnButtonImage.x = 20
	scene.returnButtonImage.y = 1080 - 84 - 20

	// BGMを読み込んで再生
	if err := PlayBGM("assets/game_win_bgm.mp3", 0.1); err != nil {
		// BGMの読み込み失敗はエラーシーンにしない（BGMなしで続行）
		fmt.Println(err.Error())
	}

	return scene
}

// Update は勝利画面のロジックを更新する
func (s *WinScene) Update() Scene {
	HandleWindowRestore()

	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if s.returnButtonImage.Contains(mx, my) {
			StopBGM()
			PlaySoundEffect("assets/button_select_se.mp3", 0.1)
			return NewTitleScene()
		}
	}
	return nil // 勝利画面に留まる
}

// Draw は勝利画面を描画する
func (s *WinScene) Draw(screen *ebiten.Image) {
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

func (s *WinScene) updateAppData(gameCode string) error {
	gameClearStatusText, err := LoadGameData(saveDataFile)
	if err != nil {
		return err
	}
	const stageCount = 6
	if len(gameClearStatusText) < stageCount {
		for len(gameClearStatusText) < stageCount {
			gameClearStatusText += "0"
		}
	}
	gameTextPosition, err := strconv.Atoi(gameCode)
	if err != nil {
		return err
	}
	if gameTextPosition < 1 || gameTextPosition > stageCount {
		return fmt.Errorf("invalid game code: %s", gameCode)
	}
	data := []byte(gameClearStatusText)
	data[gameTextPosition-1] = '1'
	gameClearStatusText = string(data)

	fmt.Println("updateAppData: ", gameClearStatusText)

	err = SaveGameData(saveDataFile, gameClearStatusText)
	if err != nil {
		return err
	}
	return nil
}
