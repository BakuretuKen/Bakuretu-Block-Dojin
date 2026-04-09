package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

var audioContext *audio.Context
var bgmPlayer *audio.Player // ゲームBGMプレイヤー（scene_game、scene_win、scene_loseで共有）
var sePlayers = map[string]*audio.Player{}

// 画面サイズの定数定義
const (
	screenWidth  = 1920
	screenHeight = 1080
)

// セーブデータ保存フォルダ名・ファイル名
const (
	appDirNameValue = "BakuretuProject"
	saveDataFile    = "block_game_01.txt"
)

// Game はEbitengineのゲームインターフェースを実装するメイン構造体
type Game struct {
	scene Scene // 現在のシーン
}

// NewGame はGame構造体を初期化して返す
func NewGame() *Game {
	return &Game{
		scene: NewTitleScene(), // 初期シーンはタイトル画面
	}
}

// Update はゲームロジックを更新する
func (g *Game) Update() error {
	HandleQuit()
	if ShouldQuit {
		return ebiten.Termination
	}

	// シーンのUpdateを呼び出し、次のシーンがあれば切り替え
	if next := g.scene.Update(); next != nil {
		g.scene = next
	}
	return nil
}

// Draw は画面描画を行う
func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

// Layout は論理的な画面サイズを返す
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// appDirNameValue、saveDataFile が空文字の場合はエラー終了
	if appDirNameValue == "" || saveDataFile == "" {
		fmt.Println("appDirNameValue or saveDataFile is empty")
		os.Exit(1)
	}
	SetAppDirName(appDirNameValue)

	// オーディオコンテキストを初期化
	audioContext = audio.NewContext(44100)
	defer StopBGM()
	defer StopAllSE()

	// ウィンドウサイズは最初に設定する必要がある
	ebiten.SetWindowSize(screenWidth/2, screenHeight/2)
	ebiten.SetWindowTitle("爆裂ブロック崩し")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowFloating(false)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
