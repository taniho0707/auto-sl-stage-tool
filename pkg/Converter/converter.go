package Converter

import (
	"github.com/taniho0707/auto-sl-stage-tool/pkg/CommandArm"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreDeleste"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreSingleHand"
)

// ConvertToCommands は ScoreSingleHand.Note の配列を CommandArm.Command の配列に変換します
// 1つ目が左、2つ目が右
// bpm: 曲のテンポ (BPM)
// offset: 曲の開始オフセット (ミリ秒)
func ConvertToCommands(leftHand []ScoreSingleHand.Note, rightHand []ScoreSingleHand.Note, bpm float64, offset int) ([]CommandArm.Command, []CommandArm.Command, error) {
	left := []CommandArm.Command{}
	right := []CommandArm.Command{}

	// 左手のコマンドを生成
	leftCommands := generateHandCommands(leftHand, CommandArm.Left, bpm, offset)
	left = append(left, leftCommands...)

	// 右手のコマンドを生成
	rightCommands := generateHandCommands(rightHand, CommandArm.Right, bpm, offset)
	right = append(right, rightCommands...)

	// 時間順にソート
	sortCommands(left)
	sortCommands(right)

	return left, right, nil
}

// generateHandCommands は片手分のコマンドを生成します
func generateHandCommands(notes []ScoreSingleHand.Note, hand CommandArm.Hand, bpm float64, offset int) []CommandArm.Command {
	commands := []CommandArm.Command{}

	// 1拍あたりの時間（ミリ秒）
	beatTimeMs := 60000.0 / bpm

	for i, note := range notes {
		isFirstNote := i == 0
		isLastNote := i == len(notes)-1

		// ノートの時間（ミリ秒）を計算
		// 小節数 * 拍子 + 拍数 で全体の拍数を計算し、それに1拍の時間をかける
		timeMs := int((float64(note.Measure*note.BeatSet+note.Beat) * beatTimeMs)) + offset

		// TargetPos から Lane を決定
		lane := convertTargetPosToLane(note.TargetPos)

		currentNote := note
		var nextNote *ScoreSingleHand.Note = nil
		var lastNote *ScoreSingleHand.Note = nil
		if !isLastNote {
			nextNote = &notes[i+1]
		}
		if !isFirstNote {
			lastNote = &notes[i-1]
		}

		// 実処理
		// 前段移動
		if lastNote == nil || lastNote.Note == ScoreDeleste.None {
			moveCmd := CommandArm.NewCommandMove(
				timeMs-300,
				hand,
				lane,
				timeMs,
			)
			commands = append(commands, moveCmd)
		}

		// 本番移動・プッシュ・リリース
		if lastNote != nil && lastNote.Note == ScoreDeleste.LongStart {
			if currentNote.Note == ScoreDeleste.Tap {
				releaseCmd := CommandArm.NewCommandSolenoid(timeMs, hand, false)
				commands = append(commands, releaseCmd)
			} else if currentNote.Note == ScoreDeleste.LeftFlick {
				moveCmd := CommandArm.NewCommandMove(timeMs, hand, lane.Left(), 0)
				commands = append(commands, moveCmd)
			} else if currentNote.Note == ScoreDeleste.RightFlick {
				moveCmd := CommandArm.NewCommandMove(timeMs, hand, lane.Right(), 0)
				commands = append(commands, moveCmd)
			}
		} else {
			if currentNote.Note == ScoreDeleste.Tap || currentNote.Note == ScoreDeleste.LongStart {
				pressCmd := CommandArm.NewCommandSolenoid(timeMs, hand, true)
				commands = append(commands, pressCmd)
			} else if currentNote.Note == ScoreDeleste.LeftFlick {
				moveCmd := CommandArm.NewCommandMove(timeMs, hand, lane.Left(), 0)
				commands = append(commands, moveCmd)
			} else if currentNote.Note == ScoreDeleste.RightFlick {
				moveCmd := CommandArm.NewCommandMove(timeMs, hand, lane.Right(), 0)
				commands = append(commands, moveCmd)
			}
		}

		// 後段
		if currentNote.Note == ScoreDeleste.Tap || currentNote.Note == ScoreDeleste.LeftFlick || currentNote.Note == ScoreDeleste.RightFlick {
			// 移動
			if nextNote == nil || nextNote.Note == ScoreDeleste.None {
				var side CommandArm.Lane
				if hand == CommandArm.Left {
					side = CommandArm.LeftEdge
				} else {
					side = CommandArm.RightEdge
				}
				moveCmd := CommandArm.NewCommandMove(timeMs+10, hand, side, 0)
				commands = append(commands, moveCmd)
			} else {
				moveCmd := CommandArm.NewCommandMove(timeMs+10, hand, convertTargetPosToLane(nextNote.TargetPos), 0)
				commands = append(commands, moveCmd)
			}

			// リリース
			if nextNote != nil && nextNote.Note != ScoreDeleste.None && (nextNote.Note == ScoreDeleste.LeftFlick || nextNote.Note == ScoreDeleste.RightFlick) {
				// なにもしない
			} else {
				releaseCmd := CommandArm.NewCommandSolenoid(timeMs+10, hand, false)
				commands = append(commands, releaseCmd)
			}
		}
	}

	return commands
}

func convertTargetPosToLane(targetPos int) CommandArm.Lane {
	switch {
	case targetPos <= 0:
		return CommandArm.LeftEdge
	case targetPos == 1:
		return CommandArm.Lane1
	case targetPos == 2:
		return CommandArm.Lane2
	case targetPos == 3:
		return CommandArm.Lane3
	case targetPos == 4:
		return CommandArm.Lane4
	case targetPos == 5:
		return CommandArm.Lane5
	default:
		return CommandArm.RightEdge
	}
}

// sortCommands はコマンドを時間順にソートします
func sortCommands(commands []CommandArm.Command) {
	// 時間でソートするロジックを実装
	// 例: sort.Slice(commands, func(i, j int) bool {
	//     return commands[i].TimeMs() < commands[j].TimeMs()
	// })
}
