package CommandArm

import "fmt"

type Command interface {
	TimeMs() int
	Message() string
}

type Hand int

const (
	Left Hand = iota + 1
	Right
)

func (h Hand) String() string {
	switch h {
	case Left:
		return "L"
	case Right:
		return "R"
	default:
		return "X"
	}
}

type Lane int

const (
	LeftEdge Lane = iota + 1
	Lane1Left
	Lane1
	Lane1Right
	Lane2Left
	Lane2
	Lane2Right
	Lane3Left
	Lane3
	Lane3Right
	Lane4Left
	Lane4
	Lane4Right
	Lane5Left
	Lane5
	Lane5Right
	RightEdge
)

func (l Lane) Left() Lane {
	return l - 1
}

func (l Lane) Right() Lane {
	return l + 1
}

func (l Lane) String() string {
	switch l {
	case LeftEdge:
		return "LL"
	case Lane1Left:
		return "1L"
	case Lane1:
		return "1C"
	case Lane1Right:
		return "1R"
	case Lane2Left:
		return "2L"
	case Lane2:
		return "2C"
	case Lane2Right:
		return "2R"
	case Lane3Left:
		return "3L"
	case Lane3:
		return "3C"
	case Lane3Right:
		return "3R"
	case Lane4Left:
		return "4L"
	case Lane4:
		return "4C"
	case Lane4Right:
		return "4R"
	case Lane5Left:
		return "5L"
	case Lane5:
		return "5C"
	case Lane5Right:
		return "5R"
	case RightEdge:
		return "RR"
	default:
		return "XX"
	}
}

type CommandSolenoid struct {
	time  int
	hand  Hand
	state bool
}

func (c *CommandSolenoid) TimeMs() int {
	return c.time
}

func (c *CommandSolenoid) Message() string {
	var state string
	if c.state {
		state = "ON"
	} else {
		state = "OF"
	}
	return fmt.Sprintf("S %d %s %s", c.time, c.hand, state)
}

func (c *CommandSolenoid) String() string {
	return c.Message()
}

type CommandMove struct {
	time    int
	hand    Hand
	lane    Lane
	endTime int
}

func (c *CommandMove) TimeMs() int {
	return c.time
}

func (c *CommandMove) Message() string {
	return fmt.Sprintf("M %d %s %s %d", c.time, c.hand, c.lane, c.endTime)
}

func (c *CommandMove) String() string {
	return c.Message()
}

func NewCommandSolenoid(time int, hand Hand, state bool) Command {
	return &CommandSolenoid{
		time:  time,
		hand:  hand,
		state: state,
	}
}

func NewCommandMove(time int, hand Hand, lane Lane, endTime int) Command {
	return &CommandMove{
		time:    time,
		hand:    hand,
		lane:    lane,
		endTime: endTime,
	}
}
