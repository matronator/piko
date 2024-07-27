package interpreter

import (
	"fmt"
	"math"
	"strings"
)

type Direction string

const (
	Up    Direction = "up"
	Down  Direction = "down"
	Left  Direction = "left"
	Right Direction = "right"
)

type OutputRegister interface{}

type Pointer struct {
	X                  int
	Y                  int
	Direction          Direction
	CurrentCell        string
	StringMode         bool
	ConditionMode      bool
	Comparator         *float64
	ComparisonOperator *string
	Operator           *string
	Stack              float64
	StringStack        []string
	OutputStack        OutputRegister
	StringModeInit     *string
}

func NewPointer(x int, y int, direction Direction) *Pointer {
	return &Pointer{
		X:             x,
		Y:             y,
		Direction:     direction,
		CurrentCell:   " ",
		StringMode:    false,
		ConditionMode: false,
		Stack:         0,
		StringStack:   []string{},
		OutputStack:   nil,
	}
}

type Point struct {
	X int
	Y int
}

type TokenPoint struct {
	Point
	Token string
}

type Line []TokenPoint
type Grid []Line

type Operation struct {
	Done   bool
	Output OutputRegister
}

type Parser struct {
	Program        string
	Width          int
	Height         int
	Grid           Grid
	VisualGrid     Grid
	Pointer        *Pointer
	CurrentPoint   Point
	OutputRegister OutputRegister
}

func NewParser(input string) *Parser {
	parser := &Parser{
		Program:        input,
		Grid:           Grid{},
		VisualGrid:     Grid{},
		OutputRegister: nil,
	}
	lines := strings.Split(parser.Program, "\n")
	parser.Height = len(lines)
	parser.Width = 0
	for _, line := range lines {
		if len(line) > parser.Width {
			parser.Width = len(line)
		}
	}

	for i, line := range lines {
		if len(line) < parser.Width {
			lines[i] = line + strings.Repeat(" ", parser.Width-len(line))
		}
	}

	tempGrid := lines
	x, y := 0, 0
	pointerX, pointerY := 0, 0
	pointerSet := false

	for _, row := range tempGrid {
		cols := []rune(row)
		line := Line{}
		for _, col := range cols {
			if isPointer(string(col)) {
				if pointerSet == true {
					panic("Multiple pointers found")
				}
				pointerSet = true
				pointerX = x
				pointerY = y
				point := TokenPoint{Point{x, y}, " "}
				line = append(line, point)
				x++
				continue
			}
			point := TokenPoint{Point{x, y}, string(col)}
			line = append(line, point)
			x++
		}
		parser.Grid = append(parser.Grid, line)
		y++
		x = 0
	}

	parser.VisualGrid = make(Grid, len(parser.Grid))
	for i, line := range parser.Grid {
		parser.VisualGrid[i] = make(Line, len(line))
		for j, col := range line {
			parser.VisualGrid[i][j] = col
		}
	}

	parser.Pointer = NewPointer(pointerX, pointerY, Right)
	parser.CurrentPoint = Point{pointerX, pointerY}

	return parser
}

func (p *Parser) Run() {
	operation := Operation{Done: false, Output: nil}
	for !operation.Done {
		operation = p.Step(operation)
	}
	fmt.Println(operation.Output)
}

func (p *Parser) Step(op Operation) Operation {
	if op.Done {
		return op
	}
	p.move()
	cell := p.Grid[p.Pointer.Y][p.Pointer.X]

	shouldContinue := p.stringModeCheck(cell)
	if shouldContinue {
		return op
	}

	shouldContinue = p.conditionModeCheck(cell)
	if shouldContinue {
		return op
	}

	switch cell.Token {
	case ";":
		op.Done = true
	case "&":
		if p.OutputRegister == nil {
			p.OutputRegister = ""
		}
		for _, char := range p.Pointer.StringStack {
			p.OutputRegister = fmt.Sprintf("%s%s", p.OutputRegister, char)
		}
		p.Pointer.StringStack = []string{}
	case ":":
		fmt.Println(p.OutputRegister)
		p.OutputRegister = nil
	case "~":
		if len(p.Pointer.StringStack) > 0 {
			if p.OutputRegister == nil {
				p.OutputRegister = ""
			}
			p.OutputRegister = fmt.Sprintf("%s%s", p.OutputRegister, p.Pointer.StringStack[len(p.Pointer.StringStack)-1])
			p.Pointer.StringStack = p.Pointer.StringStack[:len(p.Pointer.StringStack)-1]
		}
	case "!":
		p.Pointer.Stack = 0
		p.Pointer.StringStack = []string{}
		p.OutputRegister = nil
	case "'":
		p.stringModeCheck(cell)
	case `"`:
		p.stringModeCheck(cell)
	case "v":
		p.Pointer.Direction = Down
	case "^":
		p.Pointer.Direction = Up
	case "<":
		p.Pointer.Direction = Left
	case ">":
		p.Pointer.Direction = Right
	case "=":
		fmt.Println(p.Pointer.Stack)
	case "?":
		p.Pointer.ConditionMode = true
	case "+":
		if p.Pointer.Operator == nil {
			p.Pointer.Operator = &cell.Token
		} else {
			p.operatorCheck(cell)
		}
	case "-":
		if p.Pointer.Operator == nil {
			p.Pointer.Operator = &cell.Token
		} else {
			p.operatorCheck(cell)
		}
	case "*":
		if p.Pointer.Operator == nil {
			p.Pointer.Operator = &cell.Token
		} else {
			p.operatorCheck(cell)
		}
	case "/":
		if p.Pointer.Operator == nil {
			p.Pointer.Operator = &cell.Token
		} else {
			p.operatorCheck(cell)
		}
	case "_":
		p.Pointer.Stack = math.Floor(p.Pointer.Stack)
	case " ":
		// no-op
	default:
		p.operatorCheck(cell)
	}

	op.Output = p.OutputRegister
	if op.Output == nil {
		op.Output = p.Pointer.Stack
	}
	return op
}

func (p *Parser) move() {
	oldX := p.Pointer.X
	oldY := p.Pointer.Y
	currentCell := p.Pointer.CurrentCell
	p.VisualGrid[oldY][oldX].Token = currentCell

	switch p.Pointer.Direction {
	case Up:
		p.Pointer.Y--
		if p.Pointer.Y < 0 {
			p.Pointer.Y = p.Height - 1
		}
	case Down:
		p.Pointer.Y++
		if p.Pointer.Y >= p.Height {
			p.Pointer.Y = 0
		}
	case Left:
		p.Pointer.X--
		if p.Pointer.X < 0 {
			p.Pointer.X = p.Width - 1
		}
	case Right:
		p.Pointer.X++
		if p.Pointer.X >= p.Width {
			p.Pointer.X = 0
		}
	}

	p.Pointer.CurrentCell = p.Grid[p.Pointer.Y][p.Pointer.X].Token
	p.VisualGrid[p.Pointer.Y][p.Pointer.X].Token = "@"
	p.CurrentPoint = Point{p.Pointer.X, p.Pointer.Y}
}

func (p *Parser) stringModeCheck(cell TokenPoint) bool {
	if p.Pointer.StringMode {
		if cell.Token != *p.Pointer.StringModeInit {
			p.Pointer.StringStack = append(p.Pointer.StringStack, cell.Token)
			return true
		}
		p.Pointer.StringMode = false
		p.Pointer.StringModeInit = nil
		return true
	} else if cell.Token == "'" || cell.Token == `"` {
		p.Pointer.StringMode = true
		p.Pointer.StringModeInit = &cell.Token
		return true
	}

	return false
}

func (p *Parser) conditionModeCheck(cell TokenPoint) bool {
	if p.Pointer.ConditionMode {
		if isComparisonOperator(cell.Token) {
			if p.Pointer.ComparisonOperator == nil {
				p.Pointer.ComparisonOperator = &cell.Token
				return true
			}
		}

		if cell.Token == "?" && p.Pointer.ComparisonOperator != nil {
			if p.Pointer.Comparator == nil {
				comparator := float64(0)
				p.Pointer.Comparator = &comparator
			}
			switch *p.Pointer.ComparisonOperator {
			case "=":
				if p.Pointer.Stack != *p.Pointer.Comparator {
					p.Pointer.Direction = p.rotatePointer(false)
				}
			case "<":
				if p.Pointer.Stack >= *p.Pointer.Comparator {
					p.Pointer.Direction = p.rotatePointer(true)
				}
			case ">":
				if p.Pointer.Stack <= *p.Pointer.Comparator {
					p.Pointer.Direction = p.rotatePointer(false)
				}
			case "!":
				if p.Pointer.Stack == *p.Pointer.Comparator {
					p.Pointer.Direction = p.rotatePointer(true)
				}
			}
			p.Pointer.ConditionMode = false
			p.Pointer.ComparisonOperator = nil
			p.Pointer.Comparator = nil
			return true
		}

		if isDigit(cell.Token) {
			comparator := float64(cell.Token[0] - '0')
			p.Pointer.Comparator = &comparator
		} else {
			if cell.Token != " " {
				comparator := float64(cell.Token[0])
				p.Pointer.Comparator = &comparator
			}
		}
		return true
	}

	return false
}

func (p *Parser) rotatePointer(anticlockwise bool) Direction {
	directions := []Direction{Up, Right, Down, Left}
	currentDirection := p.Pointer.Direction
	currentIndex := -1
	for i, dir := range directions {
		if dir == currentDirection {
			currentIndex = i
			break
		}
	}
	newIndex := 0
	if anticlockwise {
		newIndex = currentIndex - 1
	} else {
		newIndex = currentIndex + 1
	}
	if newIndex < 0 {
		return directions[3]
	} else if newIndex > 3 {
		return directions[0]
	}
	return directions[newIndex]
}

func (p *Parser) operatorCheck(cell TokenPoint) bool {
	if p.Pointer.Operator != nil {
		switch *p.Pointer.Operator {
		case "+":
			p.Pointer.Stack = float64(p.Pointer.Stack) + float64(p.getCellAsValue(cell))
		case "-":
			p.Pointer.Stack = float64(p.Pointer.Stack) - float64(p.getCellAsValue(cell))
		case "*":
			p.Pointer.Stack = float64(p.Pointer.Stack) * float64(p.getCellAsValue(cell))
		case "/":
			p.Pointer.Stack = float64(p.Pointer.Stack) / float64(p.getCellAsValue(cell))
		}
		p.Pointer.Operator = nil
		return true
	}
	p.Pointer.Stack = p.getCellAsValue(cell)
	return true
}

func (p *Parser) getCellAsValue(cell TokenPoint) float64 {
	if isDigit(cell.Token) {
		return float64(cell.Token[0] - '0')
	}
	return float64(cell.Token[0])
}

func isPointer(token string) bool {
	return token == "#"
}

func isComparisonOperator(token string) bool {
	return token == "=" || token == "<" || token == ">" || token == "!"
}

func isDigit(token string) bool {
	return token >= "0" && token <= "9"
}
