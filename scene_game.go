package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	blockSize            = 64   // ブロックサイズ（ピクセル）
	blockThreshold       = 0.10 // 不透明率しきい値（n%）
	ballSpeed            = 9    // デフォルトボール速度（ピクセル/フレーム）
	ballSpeedAdditional  = 3    // デフォルトボール速度追加（ピクセル/フレーム）
	panelWidth           = 240  // 反射板の幅（ピクセル）
	panelY               = 42   // 反射板のY位置（下端からの距離）
	panelHitFrames       = 10   // 反射板ヒット演出フレーム数
	panelCenterRange     = 3    // 中央判定の半径（ピクセル）
	panelEdgeRange       = 35   // 端判定の範囲（ピクセル）
	panelCollisionHeight = 16   // 反射板の衝突判定高さ（ピクセル）
	panelImageHeight     = 32   // 反射板画像の高さ（ピクセル）
	panelBallGap         = 10   // 反射板とボールの隙間（ピクセル）

	gameStateStart   = 0 // スタート待ち
	gameStatePlaying = 1 // ゲーム中
)

// GameScene はゲーム画面を表す
type GameScene struct {
	bgImage        *ebiten.Image // 背景画像
	fgImage        *ebiten.Image // 前景画像
	ballImage      *ebiten.Image // ボール画像
	superBallImage *ebiten.Image // スーパーボール画像
	isSuperBall    bool          // スーパーボールモードフラグ

	ballX  int // ボールの現在X座標
	ballY  int // ボールの現在Y座標
	ballVX int // ボールのX方向速度
	ballVY int // ボールのY方向速度

	panelImage       *ebiten.Image // 反射板画像（通常）
	panelHitImage    *ebiten.Image // 反射板画像（ヒット時）
	panelX           int           // 反射板のX座標
	panelHit         int           // ヒット演出カウンタ

	blockMap        [][]int // ブロック判定マップ [行][列]
	blockCols       int     // ブロック列数
	blockRows       int     // ブロック行数
	remainingBlocks int     // 残りブロック数

	gameState        int           // ゲーム状態
	startButtonImage *ebiten.Image // スタートボタン画像

	lives     int           // 残機
	lifeImage *ebiten.Image // 残機画像

	gameCode string // ゲームコード

	ballSpeed           int // ボール速度（ピクセル/フレーム）
	ballSpeedAdditional int // ボール速度追加（ピクセル/フレーム）
}

// NewGameScene はゲーム画面を初期化して返す
func NewGameScene(gameCode string) Scene {
	scene := &GameScene{
		ballX:               0,
		ballY:               0,
		ballVX:              ballSpeed,
		ballVY:              ballSpeed,
		gameState:           gameStateStart,
		lives:               3,
		gameCode:            gameCode,
		ballSpeed:           ballSpeed,
		ballSpeedAdditional: ballSpeedAdditional,
	}
	// ボーナスステージのみボール速度アップ
	if gameCode == "06" {
		scene.ballSpeed = ballSpeed + 3
	}

	// 背景画像を読み込む
	var err error
	scene.bgImage, err = LoadEbitenImage("assets/" + gameCode + "/game_back_image.jpg")
	if err != nil {
		return NewErrorScene("Failed to load background image: " + err.Error())
	}

	// 前景画像を読み込む（ブロックマップ生成用にraw imageも取得）
	fgImg, err := LoadImage("assets/" + gameCode + "/game_front_image.png")
	if err != nil {
		return NewErrorScene("Failed to load foreground image: " + err.Error())
	}
	scene.fgImage = ebiten.NewImageFromImage(fgImg)

	// ブロックマップを生成
	bounds := fgImg.Bounds()
	scene.blockCols = (bounds.Dx() + blockSize - 1) / blockSize
	scene.blockRows = (bounds.Dy() + blockSize - 1) / blockSize
	scene.blockMap = make([][]int, scene.blockRows)

	for row := 0; row < scene.blockRows; row++ {
		scene.blockMap[row] = make([]int, scene.blockCols)
		for col := 0; col < scene.blockCols; col++ {
			// このブロック領域の不透明ピクセルをカウント
			opaqueCount := 0
			totalCount := 0

			for y := row * blockSize; y < min((row+1)*blockSize, bounds.Dy()); y++ {
				for x := col * blockSize; x < min((col+1)*blockSize, bounds.Dx()); x++ {
					_, _, _, a := fgImg.At(x, y).RGBA()
					if a > 0 {
						opaqueCount++
					}
					totalCount++
				}
			}

			if float64(opaqueCount)/float64(totalCount) > blockThreshold {
				scene.blockMap[row][col] = 1 // ブロック化する画像
				continue
			}
			// 指定しきい値を超えたらブロック
			if opaqueCount > 0 {
				scene.blockMap[row][col] = 2 // ブロック化しない画像
				continue
			}
			scene.blockMap[row][col] = 0 // ブロックなし
		}
	}
	// blockMap を標準出力 DEBUG & ブロック数をカウント
	scene.countBlockMap()

	// ボール画像を読み込み（左側22x22を使用）
	ballImg, err := LoadEbitenImage("assets/boll_image.png")
	if err != nil {
		return NewErrorScene("Failed to load ball image: " + err.Error())
	}
	scene.ballImage = ballImg.SubImage(image.Rect(0, 0, 44, 44)).(*ebiten.Image)
	scene.superBallImage = ballImg.SubImage(image.Rect(44, 0, 88, 44)).(*ebiten.Image)

	// 反射板画像を読み込み（上半分: 通常、下半分: ヒット時）
	panelImg, err := LoadEbitenImage("assets/panel_image.png")
	if err != nil {
		return NewErrorScene("反射板画像の読み込みに失敗: " + err.Error())
	}
	scene.panelImage = panelImg.SubImage(image.Rect(0, 0, panelWidth, panelImageHeight)).(*ebiten.Image)
	scene.panelHitImage = panelImg.SubImage(image.Rect(0, panelImageHeight, panelWidth, panelImageHeight*2)).(*ebiten.Image)
	scene.panelX = screenWidth / 2

	// スタートボタン画像を読み込み
	scene.startButtonImage, err = LoadEbitenImage("assets/start_button.png")
	if err != nil {
		return NewErrorScene("Failed to load start button image: " + err.Error())
	}

	// 残機画像を読み込み
	scene.lifeImage, err = LoadEbitenImage("assets/life_image.png")
	if err != nil {
		return NewErrorScene("Failed to load lives image: " + err.Error())
	}

	// BGMを読み込んで再生（グローバルプレイヤーを使用）
	if err := PlayBGM("assets/game_bgm.mp3", 0.08); err != nil {
		// BGMの読み込み失敗はエラーシーンにしない（BGMなしで続行）
		fmt.Println(err.Error())
	}

	return scene
}

// Update はゲーム画面のロジックを更新する
func (s *GameScene) Update() Scene {
	// 最大化/フルスクリーン中のESC解除処理（共通）
	HandleWindowRestore()

	// BGMのループ再生を管理（グローバルプレイヤーを使用）
	EnsureBGMLooping()

	s.updatePanelPosition()
	if s.updateStartState() {
		return nil
	}

	// ゲーム中の処理
	return s.updatePlayingState()
}

func (s *GameScene) updatePanelPosition() {
	// 反射板をマウスX座標に追従（中央を合わせる）
	mouseX, _ := ebiten.CursorPosition()
	s.panelX = mouseX - panelWidth/2
	// 画面内に制限
	if s.panelX < -panelWidth/2 {
		s.panelX = -panelWidth / 2
	}
	if s.panelX > screenWidth-panelWidth/2 {
		s.panelX = screenWidth - panelWidth/2
	}
}

func (s *GameScene) updateStartState() bool {
	// スタート待ち状態
	if s.gameState != gameStateStart {
		return false
	}

	size := s.ballImage.Bounds().Size()
	panelYPos := screenHeight - panelY
	// ボールを反射板の上に配置
	s.ballX = s.panelX + panelWidth/2 - size.X/2
	s.ballY = panelYPos - size.Y - panelBallGap

	// マウスクリックでゲーム開始
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.gameState = gameStatePlaying
		s.ballVX = s.ballSpeed  // 右方向
		s.ballVY = -s.ballSpeed // 上方向
	}
	return true
}

func (s *GameScene) updatePlayingState() Scene {
	size := s.ballImage.Bounds().Size()
	panelYPos := screenHeight - panelY

	s.updateBallVelocityFallback()
	s.updatePanelHitCounter()
	s.updateBallMovementAndBlockCollision(size)
	s.updatePaddleCollision(size, panelYPos)
	s.updateWallCollision(size)
	return s.updateLoseOrWinTransition()
}

func (s *GameScene) updateBallVelocityFallback() {
	// もしボール速度が0の場合はballSpeedに設定
	if s.ballVX == 0 {
		s.ballVX = s.ballSpeed
	}
	if s.ballVY == 0 {
		s.ballVY = s.ballSpeed
	}
}

func (s *GameScene) updatePanelHitCounter() {
	// ヒット演出カウンタを減算
	if s.panelHit > 0 {
		s.panelHit--
	}
}

func (s *GameScene) updateBallMovementAndBlockCollision(size image.Point) {
	// 現在のボール位置を保存
	prevBallX := s.ballX
	prevBallY := s.ballY

	// X方向のみ移動した位置でブロック判定
	s.ballX += s.ballVX
	centerX := s.ballX + size.X/2
	centerY := prevBallY + size.Y/2
	col, row := s.GetBlockIndex(centerX, centerY)

	if s.IsBlockAt(centerX, centerY) {
		// ブロックを破壊
		s.destroyBlock(col, row)

		if !s.isSuperBall {
			// 通常モード: X方向の反射処理
			s.ballVX = -s.ballVX
			// X方向の移動を取り消し
			s.ballX = prevBallX
		}
		// スーパーボールモード: 反射なし（ブロック破壊のみ、移動はそのまま）
	}

	// Y方向のみ移動した位置でブロック判定
	s.ballY += s.ballVY
	centerX = s.ballX + size.X/2
	centerY = s.ballY + size.Y/2
	col, row = s.GetBlockIndex(centerX, centerY)

	if s.IsBlockImagePartsAt(centerX, centerY) {
		// ブロック化されていない画像を削除
		s.destroyBlockImageParts(col, row)
	}

	if s.IsBlockAt(centerX, centerY) {
		// ブロックを破壊
		s.destroyBlock(col, row)

		if !s.isSuperBall {
			// 通常モード: Y方向の反射処理
			s.ballVY = -s.ballVY
			// Y方向の移動を取り消し
			s.ballY = prevBallY
		}
		// スーパーボールモード: 反射なし（ブロック破壊のみ、移動はそのまま）
	}
}

func (s *GameScene) updatePaddleCollision(size image.Point, panelYPos int) {
	// 反射板との衝突判定
	ballBottom := s.ballY + size.Y
	if s.ballVY > 0 && ballBottom >= panelYPos && s.ballY < panelYPos+panelCollisionHeight {
		// ボールのX範囲と反射板のX範囲が重なっているか
		if s.ballX+size.X >= s.panelX && s.ballX <= s.panelX+panelWidth {
			s.ballY = panelYPos - size.Y
			s.panelHit = panelHitFrames
			s.playPanelHitSE()

			// ボールの中心X座標と反射板上での相対位置
			ballCenterX := s.ballX + size.X/2
			relativeX := ballCenterX - s.panelX

			// 中央判定
			centerStart := panelWidth/2 - panelCenterRange
			centerEnd := panelWidth/2 + panelCenterRange

			if relativeX >= centerStart && relativeX <= centerEnd {
				// 中央: スーパーボールモード ON
				s.isSuperBall = true
				// 通常速度に戻す
				if s.ballVX > 0 {
					s.ballVX = s.ballSpeed
				} else {
					s.ballVX = -s.ballSpeed
				}
				s.ballVY = -s.ballSpeed
			} else if relativeX < panelEdgeRange || relativeX > panelWidth-panelEdgeRange {
				// 端: 横に速く、縦に遅く
				s.isSuperBall = false
				if s.ballVX > 0 {
					s.ballVX = s.ballSpeed + s.ballSpeedAdditional
				} else {
					s.ballVX = -(s.ballSpeed + s.ballSpeedAdditional)
				}
				s.ballVY = -(s.ballSpeed - s.ballSpeedAdditional)
			} else {
				// その他: 通常速度、スーパーモード OFF
				s.isSuperBall = false
				if s.ballVX > 0 {
					s.ballVX = s.ballSpeed
				} else {
					s.ballVX = -s.ballSpeed
				}
				s.ballVY = -s.ballSpeed
			}
		}
	}
}

func (s *GameScene) updateWallCollision(size image.Point) {
	// 左端・上端での跳ね返り
	if s.ballX < 0 {
		s.ballX = -s.ballX
		s.ballVX = -s.ballVX
		s.playPanelHitSE()
	}
	if s.ballY < 0 {
		s.ballY = -s.ballY
		s.ballVY = -s.ballVY
		s.playPanelHitSE()
	}

	// 右端での跳ね返り
	maxX := screenWidth - size.X
	if maxX < s.ballX {
		s.ballX = -s.ballX + 2*maxX
		s.ballVX = -s.ballVX
		s.playPanelHitSE()
	}
}

func (s *GameScene) updateLoseOrWinTransition() Scene {
	// 下端（床）に落ちた場合
	if s.ballY > screenHeight {
		s.lives--
		s.isSuperBall = false
		s.gameState = gameStateStart

		// 残機が0になったらLoseシーンに遷移
		if s.lives <= 0 {
			// BGMは停止せずに継続（scene_loseでも同じBGMを使用）
			return NewLoseScene(s.gameCode)
		}
	}

	// 全てのブロックを削除したらWinシーンに遷移
	if s.remainingBlocks <= 0 {
		// BGMは停止せずに継続（scene_winでも同じBGMを使用）
		return NewWinScene(s.gameCode)
	}

	return nil // 現在のシーンを継続
}

// destroyBlock はブロックを破壊し、背景画像で前景を上書きする
func (s *GameScene) destroyBlock(col, row int) {
	if row < 0 || row >= s.blockRows || col < 0 || col >= s.blockCols {
		return
	}
	if s.blockMap[row][col] == 0 {
		return // 既に破壊済み
	}
	s.blockMap[row][col] = 0
	s.remainingBlocks--
	// fmt.Printf("残りブロック数: %d\n", s.remainingBlocks)
	s.playBlockBreakSE()
	s.eraseBlockCell(col, row)
}

// destroyBlockImageParts はブロック化されていない画像を削除
func (s *GameScene) destroyBlockImageParts(col, row int) {
	if row < 0 || row >= s.blockRows || col < 0 || col >= s.blockCols {
		return
	}
	s.blockMap[row][col] = 0

	s.eraseBlockCell(col, row)
}

func (s *GameScene) eraseBlockCell(col, row int) {
	// 背景画像からブロック範囲を前景画像に描画
	x, y := col*blockSize, row*blockSize
	blockRect := image.Rect(x, y, x+blockSize, y+blockSize)
	bgBlock := s.bgImage.SubImage(blockRect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	s.fgImage.DrawImage(bgBlock, op)
}

func (s *GameScene) playPanelHitSE() {
	_ = PlaySoundEffect("assets/game_panel_se.mp3", 0.1)
}

func (s *GameScene) playBlockBreakSE() {
	_ = PlaySoundEffect("assets/game_block_se.mp3", 0.1)
}

// Draw はゲーム画面を描画する
func (s *GameScene) Draw(screen *ebiten.Image) {
	// 背景を描画
	screen.DrawImage(s.bgImage, &ebiten.DrawImageOptions{})
	// 前景を描画
	screen.DrawImage(s.fgImage, &ebiten.DrawImageOptions{})

	// 残機画像を右上に表示（lives - 1 個）
	lifeSize := s.lifeImage.Bounds().Size()
	for i := 0; i < s.lives-1; i++ {
		lifeOp := &ebiten.DrawImageOptions{}
		lifeOp.GeoM.Translate(float64(screenWidth-lifeSize.X*(i+1)-4*(i+1)), 4)
		screen.DrawImage(s.lifeImage, lifeOp)
	}

	// 反射板を描画
	panelOp := &ebiten.DrawImageOptions{}
	panelOp.GeoM.Translate(float64(s.panelX), float64(screenHeight-panelY))
	if s.panelHit > 0 {
		screen.DrawImage(s.panelHitImage, panelOp)
	} else {
		screen.DrawImage(s.panelImage, panelOp)
	}

	// ボールを描画
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.ballX), float64(s.ballY))
	if s.isSuperBall {
		screen.DrawImage(s.superBallImage, op)
	} else {
		screen.DrawImage(s.ballImage, op)
	}

	// スタート待ち状態ならスタートボタンを描画
	if s.gameState == gameStateStart {
		btnSize := s.startButtonImage.Bounds().Size()
		btnOp := &ebiten.DrawImageOptions{}
		btnOp.GeoM.Translate(float64(screenWidth-btnSize.X)/2, float64(screenHeight-btnSize.Y)/2)
		screen.DrawImage(s.startButtonImage, btnOp)
	}
}

// IsBlockAt はピクセル座標がブロック内かどうかを返す
func (s *GameScene) IsBlockAt(px, py int) bool {
	col, row := px/blockSize, py/blockSize
	if row < 0 || row >= s.blockRows || col < 0 || col >= s.blockCols {
		return false
	}
	return s.blockMap[row][col] == 1
}

// IsBlockAt はピクセル座標がブロック化されていない画像かどうかを返す
func (s *GameScene) IsBlockImagePartsAt(px, py int) bool {
	col, row := px/blockSize, py/blockSize
	if row < 0 || row >= s.blockRows || col < 0 || col >= s.blockCols {
		return false
	}
	return s.blockMap[row][col] == 2
}

// GetBlockIndex はピクセル座標からブロックインデックスを返す
func (s *GameScene) GetBlockIndex(px, py int) (col, row int) {
	return px / blockSize, py / blockSize
}

func (s *GameScene) countBlockMap() {
	// ブロック数をカウント
	for row := range s.blockMap {
		for col := range s.blockMap[row] {
			if s.blockMap[row][col] == 1 {
				// fmt.Print("■")
				s.remainingBlocks++
			} else if s.blockMap[row][col] == 2 {
				// fmt.Print("□")
			} else {
				// fmt.Print("・")
			}
		}
		// fmt.Println()
	}
	// fmt.Printf("ブロック数: %d\n", scene.remainingBlocks)
}
