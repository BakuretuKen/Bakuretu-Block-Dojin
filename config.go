package main

// 1ステージ分の調整値をまとめた構造体
type StageConfig struct {
	BallSpeed     int // ボール速度（ピクセル/フレーム）
	AnimeSheet    int // アニメーションスプライトシートのフレーム数
	AnimeInterval int // アニメーション開始間隔
	AnimeSpeed    int // アニメーション速度
	AnimeX        int // アニメーション表示 X 座標
	AnimeY        int // アニメーション表示 Y 座標
}

// stageConfigs にエントリが無い場合のフォールバック
var defaultStageConfig = StageConfig{
	BallSpeed: 9,
}

// gameCode ("01"〜"06") をキーにしたステージ別設定
// 面ごとの調整はこのマップを編集する
var stageConfigs = map[string]StageConfig{
	"01": {
		BallSpeed:     9,
		AnimeSheet:    0,
		AnimeInterval: 0,
		AnimeSpeed:    0,
		AnimeX:        0,
		AnimeY:        0,
	},
	"02": {
		BallSpeed:     9,
		AnimeSheet:    0,
		AnimeInterval: 0,
		AnimeSpeed:    0,
		AnimeX:        0,
		AnimeY:        0,
	},
	"03": {
		BallSpeed:     9,
		AnimeSheet:    0,
		AnimeInterval: 0,
		AnimeSpeed:    0,
		AnimeX:        0,
		AnimeY:        0,
	},
	"04": {
		BallSpeed:     9,
		AnimeSheet:    0,
		AnimeInterval: 0,
		AnimeSpeed:    0,
		AnimeX:        0,
		AnimeY:        0,
	},
	"05": {
		BallSpeed:     9,
		AnimeSheet:    0,
		AnimeInterval: 0,
		AnimeSpeed:    0,
		AnimeX:        0,
		AnimeY:        0,
	},
	"06": {
		BallSpeed:     12, // シークレット面は少し速い
		AnimeSheet:    0,
		AnimeInterval: 0,
		AnimeSpeed:    0,
		AnimeX:        0,
		AnimeY:        0,
	},
}

// gameCode に対応する設定を取得
// 未登録の gameCode の場合は defaultStageConfig を返す
func GetStageConfig(gameCode string) StageConfig {
	if cfg, ok := stageConfigs[gameCode]; ok {
		return cfg
	}
	return defaultStageConfig
}
