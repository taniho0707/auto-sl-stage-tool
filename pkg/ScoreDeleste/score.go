package ScoreDeleste

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Score struct {
	Header Header
	Notes  []Note
}

type Difficulty int

const (
	Debug Difficulty = iota + 1
	Regular
	Pro
	Master
	MasterPlus
)

type Attribute int

const (
	Cu Attribute = iota + 1
	Co
	Pa
	All
)

type Header struct {
	Title       string     // 曲のタイトル
	Lyricist    string     // 作詞者
	Composer    string     // 作曲者
	Background  string     // 背景ファイルパス
	Song        string     // 音楽ファイルパス
	Lyrics      string     // 歌詞ファイルパス
	BPM         float64    // テンポ
	Offset      int        // 譜面オフセット(ms)
	SongOffset  int        // 曲オフセット(ms)
	MovieOffset int        // 動画オフセット(ms)
	Difficulty  Difficulty // 難易度
	Level       int        // 楽曲レベル (1-30)
	BGMVolume   int        // 音楽音量 (0-100)
	SEVolume    int        // 効果音音量 (0-100)
	Attribute   Attribute  // 楽曲属性
	Brightness  int        // 背景明るさ (0-255)
}

type NoteType int

const (
	None NoteType = iota
	LeftFlick
	Tap
	RightFlick
	LongStart
	Slide
)

func (e *NoteType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "None":
		*e = None
	case "LeftFlick":
		*e = LeftFlick
	case "Tap":
		*e = Tap
	case "RightFlick":
		*e = RightFlick
	case "LongStart":
		*e = LongStart
	case "Slide":
		*e = Slide
	default:
		return fmt.Errorf("invalid note type: %s", string(text))
	}
	return nil
}

type Note struct {
	Channel   int        // チャンネル番号
	Measure   int        // 小節数
	Note      []NoteType // ノートタイプ
	StartPos  []int      // 出現位置 (1-5)
	TargetPos []int      // 目標位置 (1-5)
}

func ParseScore(filepath string) (*Score, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// BOMチェック用のバッファ
	bom := make([]byte, 3)
	if _, err := file.Read(bom); err != nil {
		return nil, err
	}
	// ファイルポインタを先頭に戻す
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	var reader *bufio.Reader
	// UTF-8 BOMチェック (0xEF,0xBB,0xBF)
	if bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
		// BOMがある場合は3バイトスキップ
		if _, err := file.Seek(3, 0); err != nil {
			return nil, err
		}
		reader = bufio.NewReader(file)
	} else {
		// BOMがない場合はそのまま読み込み
		reader = bufio.NewReader(file)
	}

	score := &Score{}
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			// #の次が数字の場合はノート情報
			if len(line) > 1 && unicode.IsDigit(rune(line[1])) {
				note, err := parseNote(line)
				if err != nil {
					return nil, err
				}
				score.Notes = append(score.Notes, *note)
			} else {
				// #の次が文字の場合はヘッダー情報
				parseHeader(line, score)
			}
			continue
		}
	}

	return score, scanner.Err()
}

func parseHeader(line string, score *Score) {
	line = strings.TrimPrefix(line, "#")
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return
	}

	value := strings.TrimSpace(parts[1])

	switch parts[0] {
	case "Title":
		score.Header.Title = value
	case "Lyricist":
		score.Header.Lyricist = value
	case "Composer":
		score.Header.Composer = value
	case "Background":
		score.Header.Background = value
	case "Song":
		score.Header.Song = value
	case "Lyrics":
		score.Header.Lyrics = value
	case "BPM":
		if bpm, err := strconv.ParseFloat(value, 64); err == nil {
			score.Header.BPM = bpm
		}
	case "Offset":
		if offset, err := strconv.Atoi(value); err == nil {
			score.Header.Offset = offset
		}
	case "SongOffset":
		if offset, err := strconv.Atoi(value); err == nil {
			score.Header.SongOffset = offset
		}
	case "MovieOffset":
		if offset, err := strconv.Atoi(value); err == nil {
			score.Header.MovieOffset = offset
		}
	case "Level":
		if level, err := strconv.Atoi(value); err == nil {
			score.Header.Level = level
		}
	case "BGMVolume":
		if volume, err := strconv.Atoi(value); err == nil {
			score.Header.BGMVolume = volume
		}
	case "SEVolume":
		if volume, err := strconv.Atoi(value); err == nil {
			score.Header.SEVolume = volume
		}
	case "Attribute":
		if attribute, err := strconv.Atoi(value); err == nil {
			score.Header.Attribute = Attribute(attribute)
		}
	case "Brightness":
		if brightness, err := strconv.Atoi(value); err == nil {
			score.Header.Brightness = brightness
		}
	}
}

func parseNote(line string) (*Note, error) {
	// #<チャンネル>,<小節数>:<タイミング>:<出現位置>:<目標位置>
	line = strings.TrimPrefix(line, "#")
	parts := strings.Split(line, ":")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid note format")
	}

	headerParts := strings.Split(parts[0], ",")
	if len(headerParts) != 2 {
		return nil, fmt.Errorf("invalid note header format")
	}

	channel, err := strconv.Atoi(headerParts[0])
	if err != nil {
		return nil, err
	}

	measure, err := strconv.Atoi(headerParts[1])
	if err != nil {
		return nil, err
	}

	note := &Note{
		Channel: channel,
		Measure: measure,
	}

	if len(parts) > 1 {
		note.Note = parseTimingString(parts[1])
	}

	if len(parts) > 2 {
		note.StartPos = parsePositionString(parts[2])
	}

	if len(parts) > 3 {
		note.TargetPos = parsePositionString(parts[3])
	}

	return note, nil
}

func parseTimingString(timing string) []NoteType {
	result := make([]NoteType, len(timing))
	for i, c := range timing {
		n, _ := strconv.Atoi(string(c))
		result[i] = NoteType(n)
	}
	return result
}

func parsePositionString(pos string) []int {
	result := make([]int, len(pos))
	for i, c := range pos {
		result[i], _ = strconv.Atoi(string(c))
	}
	return result
}
