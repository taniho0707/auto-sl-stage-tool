package Converter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/taniho0707/auto-sl-stage-tool/pkg/CommandArm"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreDeleste"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreSingleHand"
)

func TestGenerateHandCommands(t *testing.T) {
	t.Run("empty notes list", func(t *testing.T) {
		notes := []ScoreSingleHand.Note{}
		commands := generateHandCommands(notes, CommandArm.Left, 120.0, 0)
		assert.Empty(t, commands)
	})

	t.Run("only tap note", func(t *testing.T) {
		notes := []ScoreSingleHand.Note{
			{
				Measure:   0,
				Beat:      0,
				BeatSet:   4,
				Note:      ScoreDeleste.Tap,
				TargetPos: 2,
			}, {
				Measure:   0,
				Beat:      1,
				BeatSet:   4,
				Note:      ScoreDeleste.Tap,
				TargetPos: 2,
			}, {
				Measure:   0,
				Beat:      2,
				BeatSet:   4,
				Note:      ScoreDeleste.Tap,
				TargetPos: 5,
			},
		}
		expected := []string{
			"M -300 L 2C 0",
			"S 0 L ON",
			"M 10 L 2C 0",
			"S 10 L OF",
			"S 500 L ON",
			"M 510 L 5C 0",
			"S 510 L OF",
			"S 1000 L ON",
			"M 1010 L LL 0",
			"S 1010 L OF",
		}

		commands := generateHandCommands(notes, CommandArm.Left, 120.0, 0)
		assert.Len(t, commands, 10)
		for i, command := range commands {
			assert.Equal(t, expected[i], command.Message())
		}
	})

	t.Run("tan and long note", func(t *testing.T) {
		notes := []ScoreSingleHand.Note{
			{
				Measure:   0,
				Beat:      0,
				BeatSet:   4,
				Note:      ScoreDeleste.Tap,
				TargetPos: 2,
			}, {
				Measure:   0,
				Beat:      1,
				BeatSet:   4,
				Note:      ScoreDeleste.LongStart,
				TargetPos: 3,
			}, {
				Measure:   0,
				Beat:      2,
				BeatSet:   4,
				Note:      ScoreDeleste.Tap,
				TargetPos: 3,
			},
		}
		expected := []string{
			"M -300 R 2C 0",
			"S 0 R ON",
			"S 10 R OF",
			"M 10 R 3C 0",
			"S 500 R ON",
			"S 1000 R OF",
			"M 1010 R RR 0",
		}

		commands := generateHandCommands(notes, CommandArm.Right, 120.0, 0)
		assert.Len(t, commands, 7)
		for i, command := range commands {
			assert.Equal(t, expected[i], command.Message())
		}
	})

	// t.Run("consecutive flick notes", func(t *testing.T) {
	// 	notes := []ScoreSingleHand.Note{
	// 		{
	// 			Measure:   0,
	// 			Beat:      0,
	// 			BeatSet:   4,
	// 			Note:      ScoreDeleste.LeftFlick,
	// 			TargetPos: 2,
	// 		},
	// 		{
	// 			Measure:   0,
	// 			Beat:      1,
	// 			BeatSet:   4,
	// 			Note:      ScoreDeleste.RightFlick,
	// 			TargetPos: 3,
	// 		},
	// 	}
	// 	commands := generateHandCommands(notes, CommandArm.Right, 120.0, 0)
	// 	assert.NotEmpty(t, commands)
	// })

	// t.Run("long note sequence", func(t *testing.T) {
	// 	notes := []ScoreSingleHand.Note{
	// 		{
	// 			Measure:   0,
	// 			Beat:      0,
	// 			BeatSet:   4,
	// 			Note:      ScoreDeleste.LongStart,
	// 			TargetPos: 2,
	// 		},
	// 		{
	// 			Measure:   0,
	// 			Beat:      1,
	// 			BeatSet:   4,
	// 			Note:      ScoreDeleste.Tap,
	// 			TargetPos: 2,
	// 		},
	// 	}
	// 	commands := generateHandCommands(notes, CommandArm.Left, 120.0, 100)
	// 	assert.NotEmpty(t, commands)
	// })

	// t.Run("with offset", func(t *testing.T) {
	// 	notes := []ScoreSingleHand.Note{
	// 		{
	// 			Measure:   0,
	// 			Beat:      0,
	// 			BeatSet:   4,
	// 			Note:      ScoreDeleste.Tap,
	// 			TargetPos: 2,
	// 		},
	// 	}
	// 	offset := 1000
	// 	commands := generateHandCommands(notes, CommandArm.Left, 120.0, offset)
	// 	assert.Equal(t, offset-300, commands[0].GetTime())
	// })

	// t.Run("edge lane positions", func(t *testing.T) {
	// 	notes := []ScoreSingleHand.Note{
	// 		{
	// 			Measure:   0,
	// 			Beat:      0,
	// 			BeatSet:   4,
	// 			Note:      ScoreDeleste.Tap,
	// 			TargetPos: 0,
	// 		},
	// 	}
	// 	commands := generateHandCommands(notes, CommandArm.Left, 120.0, 0)
	// 	lastCommand := commands[len(commands)-2]
	// 	assert.Equal(t, CommandArm.LeftEdge, lastCommand.GetLane())
	// })
}
