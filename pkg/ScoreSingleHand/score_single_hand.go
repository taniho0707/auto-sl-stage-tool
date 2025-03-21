package ScoreSingleHand

import "github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreDeleste"

type Note struct {
	Measure   int // この音符が何小節目かを示す
	BeatSet   int // この小節が何拍子で表現されているかを示す
	Beat      int // この音符が何拍目かを示す
	Note      ScoreDeleste.NoteType
	TargetPos int
}

// ConvertFromDeleste は ScoreDeleste のデータを ScoreSingleHand.Score に変換します
// 1つ目の戻り値は奇数チャンネル（左手）、2つ目の戻り値は偶数チャンネル（右手）を抽出します
func ConvertFromDeleste(deleste *ScoreDeleste.Score) ([]Note, []Note, error) {
	result := make([][]Note, 2)

	for _, note := range deleste.Notes {
		// チャンネル番号の奇数/偶数判定
		channelIsRight := note.Channel%2 == 1

		measureNumber := note.Measure
		beatSet := len(note.Note)

		count := 0
		for beatNumber, beat := range note.Note {
			if beat == ScoreDeleste.None {
				continue
			}
			singleNote := Note{
				Measure:   measureNumber,
				BeatSet:   beatSet,
				Beat:      beatNumber,
				Note:      beat,
				TargetPos: note.TargetPos[count],
			}
			if channelIsRight {
				result[1] = append(result[1], singleNote)
			} else {
				result[0] = append(result[0], singleNote)
			}
			count++
		}
	}

	return result[0], result[1], nil
}
