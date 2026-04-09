package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Scene はゲーム内の各画面（タイトル、ゲームなど）を表すインターフェース
type Scene interface {
	// Update はシーンのロジックを更新する
	// 戻り値: 次のシーン（nilなら現在のシーンを継続）
	Update() Scene

	// Draw はシーンを描画する
	Draw(screen *ebiten.Image)
}

// isExpanded はウィンドウが最大化またはフルスクリーン状態かを返す
func isExpanded() bool {
	return ebiten.IsWindowMaximized() || ebiten.IsFullscreen()
}

// restoreWindow は最大化/フルスクリーンを解除する
func restoreWindow() {
	if ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
	} else {
		ebiten.RestoreWindow()
	}
}

// toggleFullscreen はフルスクリーン表示を切り替える
func toggleFullscreen() {
	ebiten.SetFullscreen(!ebiten.IsFullscreen())
}

// HandleWindowRestore は最大化/フルスクリーン中にESCキーで元に戻す共通処理
// 各シーンのUpdate()の先頭で呼び出す
func HandleWindowRestore() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		toggleFullscreen()
	}

	if isExpanded() && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		restoreWindow()
	}
}

// ShouldQuit はゲーム終了フラグ
var ShouldQuit bool

// HandleQuit は Ctrl+W / Command+W でゲーム終了フラグをセットする
func HandleQuit() {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta) {
			ShouldQuit = true
		}
	}
}
