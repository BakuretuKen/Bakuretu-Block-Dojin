package main

// StageConfig は 1 ステージ分の調整値をまとめた構造体。
// 将来ステージ別の設定を追加する場合はこの struct にフィールドを足す。
type StageConfig struct {
	BallSpeed int // ボール速度（ピクセル/フレーム）
}

// defaultStageConfig は stageConfigs にエントリが無い場合のフォールバック。
var defaultStageConfig = StageConfig{
	BallSpeed: 9,
}

// stageConfigs は gameCode ("01"〜"06") をキーにしたステージ別設定。
// 面ごとの調整はこのマップを編集する。
var stageConfigs = map[string]StageConfig{
	"01": {BallSpeed: 9},
	"02": {BallSpeed: 9},
	"03": {BallSpeed: 9},
	"04": {BallSpeed: 9},
	"05": {BallSpeed: 9},
	"06": {BallSpeed: 12}, // シークレット面は少し速い
}

// GetStageConfig は gameCode に対応する設定を返す。
// 未登録の gameCode の場合は defaultStageConfig を返す。
func GetStageConfig(gameCode string) StageConfig {
	if cfg, ok := stageConfigs[gameCode]; ok {
		return cfg
	}
	return defaultStageConfig
}
