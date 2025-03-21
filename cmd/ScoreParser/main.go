package main

import (
	"fmt"

	"github.com/taniho0707/auto-sl-stage-tool/pkg/Converter"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreDeleste"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreSingleHand"
)

func main() {
	score, err := ScoreDeleste.ParseScore("star.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(score)

	scoreLeft, scoreRight, err := ScoreSingleHand.ConvertFromDeleste(score)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(scoreLeft)
	fmt.Println(scoreRight)

	cmdLeft, cmdRight, err := Converter.ConvertToCommands(scoreLeft, scoreRight, 120.0, 0)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(cmdLeft)
	fmt.Println(cmdRight)
}
