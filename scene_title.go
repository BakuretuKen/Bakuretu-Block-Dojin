package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TitleScene はタイトル画面を表す
type TitleScene struct {
	bgImage         *ebiten.Image // 背景画像
	titleImage      *Sprite
	menu            []*Sprite // メニュー
	isWin           []bool
	frameIndex      int
	volumeImage     *Sprite
	albumIconImage  *Sprite // アルバムアイコン画像
	windowSizeImage *Sprite // ウィンドウサイズアイコン画像
}

// NewTitleScene はタイトル画面を初期化して返す
func NewTitleScene() Scene {
	scene := &TitleScene{}
	var err error

	// ゲームクリア状態を初期化
	scene.isWin = make([]bool, 6)

	// ゲーム設定確認
	if appDirNameValue == "" || saveDataFile == "" {
		return NewErrorScene("Game management code is not set: appDirNameValue, saveDataFile")
	}

	// ゲームクリア状態を読み込み
	gameClearStatusText, err := LoadGameData(saveDataFile)
	if err != nil {
		// 初回起動初期化
		if err := SaveGameData(saveDataFile, "000000"); err != nil {
			return NewErrorScene("Failed to initialize game data: " + err.Error())
		}
		gameClearStatusText = "000000"
	}
	// fmt.Println(gameClearStatusText)
	maxRead := len(gameClearStatusText)
	if maxRead > len(scene.isWin) {
		maxRead = len(scene.isWin)
	}
	for i := 0; i < maxRead; i++ {
		scene.isWin[i] = gameClearStatusText[i] == '1'
	}

	// 背景画像を読み込む
	scene.bgImage, err = LoadEbitenImage("assets/title_back_image.jpg")
	if err != nil {
		return NewErrorScene("Failed to load background image: " + err.Error())
	}
	// タイトル画像
	titleImage, err := LoadEbitenImage("assets/title_image.png")
	if err != nil {
		return NewErrorScene("Failed to load title image: " + err.Error())
	}
	scene.titleImage = NewSprite(titleImage, 2)
	scene.titleImage.SetAnimeSpeed(60, false)
	scene.titleImage.x = 360
	scene.titleImage.y = 20

	// メニュー画像を読み込む
	scene.menu = make([]*Sprite, 6)
	for i := 0; i < 3; i++ {
		spriteSheet, err := LoadEbitenImage("assets/" + fmt.Sprintf("%02d", i+1) + "/menu.jpg")
		if err != nil {
			return NewErrorScene("Failed to load menu image: " + err.Error())
		}
		scene.menu[i] = NewSprite(spriteSheet, 4)
		scene.menu[i].x = (380 * i) + 400
		scene.menu[i].y = 350
	}
	for i := 0; i < 3; i++ {
		spriteSheet, err := LoadEbitenImage("assets/" + fmt.Sprintf("%02d", i+4) + "/menu.jpg")
		if err != nil {
			return NewErrorScene("Failed to load menu image: " + err.Error())
		}
		scene.menu[i+3] = NewSprite(spriteSheet, 4)
		scene.menu[i+3].x = (380 * i) + 400
		scene.menu[i+3].y = 710
	}
	// BGMを読み込んで再生
	if err := PlayBGM("assets/title_bgm.mp3", 0.1); err != nil {
		// BGMの読み込み失敗はエラーシーンにしない（BGMなしで続行）
		fmt.Println(err.Error())
	}

	// album_icon_image.png を読み込む
	albumIconImage, err := LoadEbitenImage("assets/title_album_icon.png")
	if err != nil {
		return NewErrorScene("Failed to load album icon image: " + err.Error())
	}
	scene.albumIconImage = NewSprite(albumIconImage, 2)
	scene.albumIconImage.x = 20
	scene.albumIconImage.y = screenHeight - 180

	// window_size_icon_image.png を読み込む
	windowSizeImage, err := LoadEbitenImage("assets/winsow_size_button.png")
	if err != nil {
		return NewErrorScene("Failed to load window size icon image: " + err.Error())
	}
	scene.windowSizeImage = NewSprite(windowSizeImage, 4)
	scene.windowSizeImage.x = screenWidth - 200
	scene.windowSizeImage.y = screenHeight - 100

	return scene
}

// Update はタイトル画面のロジックを更新する
func (s *TitleScene) Update() Scene {
	// 現在のフレーム数はebitenのTickCount()で取得可能
	// 最大化/フルスクリーン中のESC解除処理（共通）
	HandleWindowRestore()

	// BGMのループ再生を管理
	EnsureBGMLooping()

	s.titleImage.Update()

	// スタートボタンがクリックされたらゲーム画面に遷移
	mx, my := ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if clicked {
		for i := 0; i < len(s.menu); i++ {
			if s.menu[i].Contains(mx, my) {
				// ボーナスステージは条件が達成されていないとクリック不可
				if i == 5 && !s.isClearAll() {
					continue
				}
				// BGMを停止
				StopBGM()
				PlaySoundEffect("assets/button_select_se.mp3", 0.1)
				return NewGameScene(fmt.Sprintf("%02d", i+1)) // ゲーム画面に遷移
			}
		}
		if s.albumIconImage.Contains(mx, my) {
			// BGMを停止
			StopBGM()
			PlaySoundEffect("assets/button_select_se.mp3", 0.1)
			return NewAlbumScene()
		}
		if s.windowSizeImage.Contains(mx, my) {
			toggleFullscreen()
		}
	}

	// ctrl + Alt + Dキークリックでセーブ状態初期化（デバッグ）
	if inpututil.IsKeyJustPressed(ebiten.KeyD) && ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) {
		fmt.Println("＜デバッグ：セーブデータクリア処理＞")
		SaveGameData(saveDataFile, "000000")
	}
	// ctrl + Alt + Aキークリックでセーブ状態初期化（デバッグ）
	if inpututil.IsKeyJustPressed(ebiten.KeyA) && ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) {
		fmt.Println("＜デバッグ：全ゲームクイア＞")
		SaveGameData(saveDataFile, "111111")
	}

	for i := 0; i < len(s.menu); i++ {
		s.menu[i].Update()
	}

	return nil // 現在のシーンを継続
}

// Draw はタイトル画面を描画する
func (s *TitleScene) Draw(screen *ebiten.Image) {
	// 背景画像を描画
	screen.DrawImage(s.bgImage, &ebiten.DrawImageOptions{})
	// タイトル画像を描画
	s.titleImage.Draw(screen)

	mx, my := ebiten.CursorPosition()
	for i := 0; i < len(s.menu); i++ {
		s.menu[i].SetFrame(0)
		if s.menu[i].Contains(mx, my) {
			if i != 5 || s.isClearAll() {
				s.menu[i].SetFrame(1) // マウスオーバー
			}
		}
		if s.isWin[i] {
			s.menu[i].SetFrame(2) // 勝利
			if s.menu[i].Contains(mx, my) {
				s.menu[i].SetFrame(3) // マウスオーバー
			}
		}
		// 時間差で表示
		if (i*2)+30 > s.frameIndex {
			continue
		}
		s.menu[i].Draw(screen)
	}
	s.frameIndex++

	// 時間差で表示
	if (7*2)+30 <= s.frameIndex {
		// アルバムアイコン画像を描画
		s.albumIconImage.SetFrame(0)
		if s.albumIconImage.Contains(mx, my) {
			s.albumIconImage.SetFrame(1)
		}
		s.albumIconImage.Draw(screen)
		// ウィンドウサイズアイコン画像を描画
		if ebiten.IsFullscreen() {
			s.windowSizeImage.SetFrame(2) // 通常サイズ
			if s.windowSizeImage.Contains(mx, my) {
				s.windowSizeImage.SetFrame(3)
			}
		} else {
			s.windowSizeImage.SetFrame(0) // フルスクリーン
			if s.windowSizeImage.Contains(mx, my) {
				s.windowSizeImage.SetFrame(1)
			}
		}
		s.windowSizeImage.Draw(screen)
	}
}

func (s *TitleScene) isClearAll() bool {
	for i := 0; i < 5; i++ {
		if !s.isWin[i] {
			return false
		}
	}
	return true
}
