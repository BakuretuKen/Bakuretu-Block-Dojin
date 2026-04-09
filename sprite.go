// Ebitengine Sprite

package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	frames       []*ebiten.Image
	imageWidth   int
	imageHeight  int
	x            int
	y            int
	vx           int
	vy           int
	angle        int
	vAngle       int
	currentFrame int
	hitArea      int  // パーセントで指定（0-100）
	_frameCount  int  // フレーム数
	animeSpeed   int  // アニメーション速度
	isAnimeLoop  bool // アニメーションループフラグ
}

// NewSprite はスプライトを生成する
func NewSprite(sheet *ebiten.Image, count int) *Sprite {
	frames, frameWidth, frameHeight := splitFrames(sheet, count)
	return &Sprite{
		frames:      frames,
		imageWidth:  frameWidth,
		imageHeight: frameHeight,
		hitArea:     100,
		_frameCount: 0,
		animeSpeed:  0,
		isAnimeLoop: true,
	}
}

// splitFrames はスプライトシートを指定数のフレームに分割する
func splitFrames(sheet *ebiten.Image, count int) ([]*ebiten.Image, int, int) {
	bounds := sheet.Bounds()
	frameWidth := bounds.Dx() / count
	frameHeight := bounds.Dy()

	frames := make([]*ebiten.Image, count)
	for i := 0; i < count; i++ {
		frames[i] = ebiten.NewImage(frameWidth, frameHeight)
		frames[i].DrawImage(sheet.SubImage(image.Rect(i*frameWidth, 0, (i+1)*frameWidth, frameHeight)).(*ebiten.Image), nil)
	}

	return frames, frameWidth, frameHeight
}

func (s *Sprite) Update() {
	s._frameCount++

	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x = -s.x
		s.vx = -s.vx
	} else if mx := screenWidth - s.imageWidth; mx <= s.x {
		s.x = 2*mx - s.x
		s.vx = -s.vx
	}
	if s.y < 0 {
		s.y = -s.y
		s.vy = -s.vy
	} else if my := screenHeight - s.imageHeight; my <= s.y {
		s.y = 2*my - s.y
		s.vy = -s.vy
	}
	s.angle += s.vAngle
	if s.angle >= 256 {
		s.angle -= 256
	} else if s.angle < 0 {
		s.angle += 256
	}
	// アニメーション処理
	if s.animeSpeed > 0 {
		frameIndex := s._frameCount / s.animeSpeed
		if frameIndex >= len(s.frames) {
			if !s.isAnimeLoop {
				s.animeSpeed = 0
				return
			}
			// ループする場合の処理
			frameIndex = frameIndex % len(s.frames)
		}
		s.SetFrame(frameIndex)
	}
}

func (s *Sprite) SetAnimeSpeed(speed int, loop bool) {
	s.animeSpeed = speed
	s.isAnimeLoop = loop
}

func (s *Sprite) CurrentFrame() int {
	return s.currentFrame
}

func (s *Sprite) SetFrame(frame int) {
	if frame < 0 || frame >= len(s.frames) {
		return
	}
	if s.currentFrame == frame {
		return
	}
	s.currentFrame = frame
}

func (s *Sprite) Contains(x, y int) bool {
	// hitArea に応じて判定領域を縮小（中心基準）
	w := s.imageWidth * s.hitArea / 100
	h := s.imageHeight * s.hitArea / 100
	offsetX := (s.imageWidth - w) / 2
	offsetY := (s.imageHeight - h) / 2
	return x >= s.x+offsetX && x < s.x+offsetX+w &&
		y >= s.y+offsetY && y < s.y+offsetY+h
}

// HitBounds は判定領域の矩形を返す（hitAreaに応じて中心基準で縮小）
func (s *Sprite) HitBounds() (x, y, w, h int) {
	w = s.imageWidth * s.hitArea / 100
	h = s.imageHeight * s.hitArea / 100
	offsetX := (s.imageWidth - w) / 2
	offsetY := (s.imageHeight - h) / 2
	return s.x + offsetX, s.y + offsetY, w, h
}

func (s *Sprite) IsHitWith(other *Sprite) bool {
	// 自キャラの判定領域と相手キャラの判定領域が重なっているか（AABB衝突判定）
	x1, y1, w1, h1 := s.HitBounds()
	x2, y2, w2, h2 := other.HitBounds()

	// 矩形同士が重なっているかチェック
	return x1 < x2+w2 && x1+w1 > x2 &&
		y1 < y2+h2 && y1+h1 > y2
}

func (s *Sprite) GetImage() *ebiten.Image {
	return s.frames[s.currentFrame]
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := float64(s.imageWidth), float64(s.imageHeight)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(2 * math.Pi * float64(s.angle) / 256)
	op.GeoM.Translate(w/2, h/2)
	op.GeoM.Translate(float64(s.x), float64(s.y))
	screen.DrawImage(s.GetImage(), op)
}
