package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/Converter"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreDeleste"
	"github.com/taniho0707/auto-sl-stage-tool/pkg/ScoreSingleHand"
	"github.com/zserge/lorca"
)

func prev(n int, format *beep.Format, pos int) int {
	threeSecSamples := format.SampleRate.N(time.Second * time.Duration(n))
	newPos := max(pos-threeSecSamples, 0)
	fmt.Printf("currentPos: %d, newPos: %d\n", pos, newPos)
	return newPos
}

func proc(n int, format *beep.Format, buffer *beep.Buffer, pos int) int {
	threeSecSamples := format.SampleRate.N(time.Second * time.Duration(n))
	newPos := min(pos+threeSecSamples, buffer.Len())
	fmt.Printf("currentPos: %d, newPos: %d\n", pos, newPos)
	return newPos
}

func main() {
	// Set up the audio
	f, err := os.Open("S:\\git\\auto-sl-stage-tool\\star.wav")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/50))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	var offsetStreamShot int = 0
	offsetScore := 50
	shot := buffer.Streamer(0, buffer.Len())
	ctrl := &beep.Ctrl{Streamer: shot, Paused: true}
	speaker.Play(ctrl)

	// Set up the Score
	score, err := ScoreDeleste.ParseScore("S:\\git\\auto-sl-stage-tool\\star.txt")
	if err != nil {
		panic(err)
	}
	scoreLeft, scoreRight, err := ScoreSingleHand.ConvertFromDeleste(score)
	if err != nil {
		panic(err)
	}
	cmdLeft, cmdRight, err := Converter.ConvertToCommands(scoreLeft, scoreRight, 178.0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(cmdLeft)
	fmt.Println(cmdRight)

	var indexLeft int = 0
	var indexRight int = 0
	var positionLeft float64 = 0
	var positionRight float64 = 0
	var pushLeft bool = false
	var pushRight bool = false

	// Set up the UI
	ui, err := lorca.New("", "", 850, 535, "--remote-allow-origins=*")
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	startPause := make(chan bool)

	ui.Bind("startPause", func() { startPause <- true })
	// ui.Bind("reset", func() {
	// 	atomic.StoreUint32(&ticks, 0)
	// 	ui.Eval(`document.querySelector('.timer').innerText = '0'`)
	// })
	ui.Bind("drawCircle", func() {
		ui.Eval(`
			const ctx = document.getElementById("canvas").getContext("2d");
			ctx.clearRect(0, 0, canvas.width, canvas.height);
			for (let i = 0; i < 5; i++) {
				ctx.beginPath();
				ctx.arc(100 + 150 * i, 300, 50, 0, 2 * Math.PI);
				ctx.strokeStyle = 'gray';
				ctx.lineWidth = 2;
				ctx.stroke();
			}
		`)
	})
	ui.Bind("prev3s", func() {
		speaker.Lock()
		newPos := prev(3, &format, shot.Position()+offsetStreamShot)
		offsetStreamShot = newPos
		shot = buffer.Streamer(newPos, buffer.Len())
		ctrl = &beep.Ctrl{Streamer: shot, Paused: ctrl.Paused}
		speaker.Unlock()
		speaker.Clear()
		speaker.Play(ctrl)
	})
	ui.Bind("prev15s", func() {
		speaker.Lock()
		newPos := prev(15, &format, shot.Position()+offsetStreamShot)
		offsetStreamShot = newPos
		shot = buffer.Streamer(newPos, buffer.Len())
		ctrl = &beep.Ctrl{Streamer: shot, Paused: ctrl.Paused}
		speaker.Unlock()
		speaker.Clear()
		speaker.Play(ctrl)
	})
	ui.Bind("proc3s", func() {
		speaker.Lock()
		newPos := proc(3, &format, buffer, shot.Position()+offsetStreamShot)
		offsetStreamShot = newPos
		shot = buffer.Streamer(newPos, buffer.Len())
		ctrl = &beep.Ctrl{Streamer: shot, Paused: ctrl.Paused}
		speaker.Unlock()
		speaker.Clear()
		speaker.Play(ctrl)
	})
	ui.Bind("proc15s", func() {
		speaker.Lock()
		newPos := proc(15, &format, buffer, shot.Position()+offsetStreamShot)
		offsetStreamShot = newPos
		shot = buffer.Streamer(newPos, buffer.Len())
		ctrl = &beep.Ctrl{Streamer: shot, Paused: ctrl.Paused}
		speaker.Unlock()
		speaker.Clear()
		speaker.Play(ctrl)
	})

	ui.Bind("drawLefthand", func() {
		color := "gray"
		if pushLeft {
			color = "red"
		}
		pos := int(100 + 150*positionLeft)

		ui.Eval(fmt.Sprintf(`
			ctx.clearRect(0, 100-25, canvas.width, 50);
			ctx.fillStyle = '%s';
			ctx.fillRect(%d-25, 100-25, 50, 50);
		`, color, pos))
	})
	ui.Bind("drawRighthand", func() {
		color := "gray"
		if pushRight {
			color = "red"
		}
		pos := int(100 + 150*positionRight)

		ui.Eval(fmt.Sprintf(`
			ctx.clearRect(0, 200-25, canvas.width, 50);
			ctx.fillStyle = '%s';
			ctx.fillRect(%d-25, 200-25, 50, 50);
		`, color, pos))
	})

	ui.Load("data:text/html," + url.PathEscape(`
	<html>
		<head>
		    <title>Command Simulator</title>
		    <style>
		        canvas {
		            border: 0px;
		            margin: 0px;
		        }
		    </style>
		</head>
		<body>
			<div class="timer" onclick="toggle()"></div>
			<canvas id="canvas" width="800" height="400"></canvas>

			<button onclick="prev15s()">Prev 15s</button>
			<button onclick="prev3s()">Prev 3s</button>
			<button onclick="startPause()">Start/Pause</button>
			<button onclick="proc3s()">Proc 3s</button>
			<button onclick="proc15s()">Proc 15s</button>

			<script>
				const draw = () => {
					drawLefthand();
					drawRighthand();
					window.requestAnimationFrame(draw);
				};

				window.onload = function() {
					drawCircle();
					window.requestAnimationFrame(draw);
				};

			</script>
		</body>
	</html>
	`))

	go func() {
		statusStartPause := true
		t := time.NewTicker(5 * time.Millisecond)
		for {
			select {
			case <-t.C:
				currentPosition := (shot.Position()+offsetStreamShot)*1000/int(format.SampleRate) - offsetScore
				if indexLeft < len(cmdLeft) && cmdLeft[indexLeft].TimeMs() <= currentPosition {
					cmd := strings.Split(cmdLeft[indexLeft].Message(), " ")
					switch {
					case cmd[0] == "S" && len(cmd) == 4:
						if cmd[3] == "ON" {
							pushLeft = true
						} else if cmd[3] == "OF" {
							pushLeft = false
						}
					case cmd[0] == "M" && len(cmd) == 5:
						switch cmd[3] {
						case "LL":
							positionLeft = -0.5
						case "1L":
							positionLeft = -0.2
						case "1C":
							positionLeft = 0
						case "1R":
							positionLeft = 0.2
						case "2L":
							positionLeft = 0.8
						case "2C":
							positionLeft = 1
						case "2R":
							positionLeft = 1.2
						case "3L":
							positionLeft = 1.8
						case "3C":
							positionLeft = 2
						case "3R":
							positionLeft = 2.2
						case "4L":
							positionLeft = 2.8
						case "4C":
							positionLeft = 3
						case "4R":
							positionLeft = 3
						case "5L":
							positionLeft = 3.8
						case "5C":
							positionLeft = 4
						case "5R":
							positionLeft = 4.2
						case "RR":
							positionLeft = 4.5
						default:
							panic("Unknown command: " + cmd[3])
						}
					default:
						panic("Unknown command: " + cmd[0])
					}

					fmt.Println(cmdLeft[indexLeft])
					indexLeft++
				}
				if indexRight < len(cmdRight) && cmdRight[indexRight].TimeMs() <= currentPosition {
					cmd := strings.Split(cmdRight[indexRight].Message(), " ")
					switch {
					case cmd[0] == "S" && len(cmd) == 4:
						if cmd[3] == "ON" {
							pushRight = true
						} else if cmd[3] == "OF" {
							pushRight = false
						}
					case cmd[0] == "M" && len(cmd) == 5:
						switch cmd[3] {
						case "LL":
							positionRight = -1
						case "1L", "1C", "1R":
							positionRight = 0
						case "2L", "2C", "2R":
							positionRight = 1
						case "3L", "3C", "3R":
							positionRight = 2
						case "4L", "4C", "4R":
							positionRight = 3
						case "5L", "5C", "5R":
							positionRight = 4
						case "RR":
							positionRight = 5
						default:
							panic("Unknown command: " + cmd[3])
						}
					default:
						panic("Unknown command: " + cmd[0])
					}

					fmt.Println(cmdRight[indexRight])
					indexRight++
				}
			case <-startPause:
				statusStartPause = !statusStartPause
				speaker.Lock()
				ctrl.Paused = statusStartPause
				speaker.Unlock()
			}
		}
	}()
	<-ui.Done()
}
