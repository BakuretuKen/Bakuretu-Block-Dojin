package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ErrorScene はエラーメッセージを表示するシーン
type ErrorScene struct {
	message    string
	cachedText *ebiten.Image
	cachedW    int
	cachedH    int
}

// NewErrorScene はエラーシーンを初期化して返す
func NewErrorScene(msg string) *ErrorScene {
	return &ErrorScene{message: msg}
}

// Update はエラーシーンのロジックを更新する
func (s *ErrorScene) Update() Scene {
	HandleWindowRestore()
	return nil // エラー画面に留まる
}

// Draw はエラーメッセージを画面に描画する
func (s *ErrorScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0x00, 0x00, 0xff}) // 暗い赤背景

	// テキスト描画用画像は再利用して毎フレーム生成を避ける
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if s.cachedText == nil || s.cachedW != w || s.cachedH != h {
		s.cachedText = ebiten.NewImage(w, h)
		s.cachedW = w
		s.cachedH = h
		ebitenutil.DebugPrint(s.cachedText, "\n Error: "+s.message)
	}

	// 一時画像をスケーリングして表示
	scale := 3.0
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(10, 20)
	screen.DrawImage(s.cachedText, op)
}

