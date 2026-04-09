package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TitleScene はタイトル画面を表す
type AlbumScene struct {
	bgImage           *ebiten.Image // 背景画像
	titleImage        *Sprite
	menu              []*Sprite // メニュー
	isWin             []bool
	frameIndex        int
	volumeImage       *Sprite
	returnButtonImage *Sprite // 戻るボタン画像
}

// NewTitleScene はタイトル画面を初期化して返す
func NewAlbumScene() Scene {
	scene := &AlbumScene{}
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
	scene.bgImage, err = LoadEbitenImage("assets/album_back_image.jpg")
	if err != nil {
		return NewErrorScene("Failed to load background image: " + err.Error())
	}
	// タイトル画像
	titleImage, err := LoadEbitenImage("assets/album_title_image.png")
	if err != nil {
		return NewErrorScene("Failed to load title image: " + err.Error())
	}
	scene.titleImage = NewSprite(titleImage, 1)
	scene.titleImage.x = 60
	scene.titleImage.y = 20

	// メニュー画像を読み込む
	scene.menu = make([]*Sprite, 6)
	for i := 0; i < 3; i++ {
		spriteSheet, err := LoadEbitenImage("assets/" + fmt.Sprintf("%02d", i+1) + "/album.jpg")
		if err != nil {
			return NewErrorScene("Failed to load menu image: " + err.Error())
		}
		scene.menu[i] = NewSprite(spriteSheet, 4)
		scene.menu[i].x = (380 * i) + 400
		scene.menu[i].y = 250
	}
	for i := 0; i < 3; i++ {
		spriteSheet, err := LoadEbitenImage("assets/" + fmt.Sprintf("%02d", i+4) + "/album.jpg")
		if err != nil {
			return NewErrorScene("Failed to load menu image: " + err.Error())
		}
		scene.menu[i+3] = NewSprite(spriteSheet, 4)
		scene.menu[i+3].x = (380 * i) + 400
		scene.menu[i+3].y = 610
	}

	// 戻るボタン画像を読み込む
	returnButtonImage, err := LoadEbitenImage("assets/return_menu_button.png")
	if err != nil {
		return NewErrorScene("Failed to load return button image: " + err.Error())
	}
	scene.returnButtonImage = NewSprite(returnButtonImage, 2)
	scene.returnButtonImage.x = 20
	scene.returnButtonImage.y = 1080 - 84 - 20

	// BGMを読み込んで再生（グローバルプレイヤーを使用）
	if err := PlayBGM("assets/album_bgm.mp3", 0.1); err != nil {
		// BGMの読み込み失敗はエラーシーンにしない（BGMなしで続行）
		fmt.Println(err.Error())
	}

	return scene
}

// Update はタイトル画面のロジックを更新する
func (s *AlbumScene) Update() Scene {
	// 現在のフレーム数はebitenのTickCount()で取得可能
	// 最大化/フルスクリーン中のESC解除処理（共通）
	HandleWindowRestore()

	// BGMのループ再生を管理（グローバルプレイヤーを使用）
	EnsureBGMLooping()

	// スタートボタンがクリックされたら画像表示画面に遷移
	mx, my := ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if clicked {
		for i := 0; i < len(s.menu); i++ {
			if s.menu[i].Contains(mx, my) {
				// クリアしていないとクリック不可
				if !s.isWin[i] {
					continue
				}
				// BGMは停止せずに継続（scene_showでも同じBGMを使用）
				PlaySoundEffect("assets/button_select_se.mp3", 0.1)
				return NewShowScene(fmt.Sprintf("%02d", i+1), 0) // 画像表示画面に遷移
			}
		}
		if s.returnButtonImage.Contains(mx, my) {
			// BGMを停止（タイトル画面に戻るため）
			StopBGM()
			PlaySoundEffect("assets/button_select_se.mp3", 0.1)
			return NewTitleScene()
		}
	}

	for i := 0; i < len(s.menu); i++ {
		s.menu[i].Update()
	}

	return nil // 現在のシーンを継続
}

// Draw はタイトル画面を描画する
func (s *AlbumScene) Draw(screen *ebiten.Image) {
	// 背景画像を描画
	screen.DrawImage(s.bgImage, &ebiten.DrawImageOptions{})

	mx, my := ebiten.CursorPosition()
	for i := 0; i < len(s.menu); i++ {
		s.menu[i].SetFrame(0)
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

	// タイトル表示
	if (6*2)+30 < s.frameIndex {
		s.titleImage.Draw(screen)
	}
	s.frameIndex++

	// 戻るボタンを描画
	s.returnButtonImage.Draw(screen)
	if s.returnButtonImage.Contains(mx, my) {
		s.returnButtonImage.SetFrame(1)
	} else {
		s.returnButtonImage.SetFrame(0)
	}
}
