package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var winWidth, winHeight int32
var stepTime, deltaTime float32

const (
	MapSize          = 20
	StepTime float32 = 0.1 // Time between state updates
)

// DIRECTION

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

// POSITION

type Position struct {
	x, y int32
}

// SNAKESEGMENT

// Segment of Snake. Implemented as a linked list
type SnakeSegment struct {
	position Position
	next     *SnakeSegment
}

func (s *SnakeSegment) CollidesWith(position Position) bool {
	collides := s.position == position

	if !collides && s.next != nil {
		return s.next.CollidesWith(position)
	}

	return collides
}

func (s *SnakeSegment) Update(position Position, extend bool) {
	prev := s.position
	s.position = position

	if s.next != nil {
		s.next.Update(prev, extend)
	} else if s.next == nil && extend == true {
		s.next = &SnakeSegment{prev, nil}
	}
}

func (s *SnakeSegment) Draw() {
	rl.DrawRectangle(
		int32(s.position.x*10),
		int32(s.position.y*10),
		10, 10,
		rl.LightGray,
	)

	if s.next != nil {
		s.next.Draw()
	}
}

// GAME

type Game struct {
	snake     *SnakeSegment
	direction Direction
	extend    bool
}

func NewGame() *Game {
	game := Game{
		snake: &SnakeSegment{
			Position{10, 10}, nil,
		},

		direction: Left,
		extend:    false,
	}

	return &game
}

func (g *Game) Input() {
	if rl.IsKeyDown(rl.KeyRight) && g.direction != Left {
		g.direction = Right
	}
	if rl.IsKeyDown(rl.KeyLeft) && g.direction != Right {
		g.direction = Left
	}
	if rl.IsKeyDown(rl.KeyDown) && g.direction != Up {
		g.direction = Down
	}
	if rl.IsKeyDown(rl.KeyUp) && g.direction != Down {
		g.direction = Up
	}
	if rl.IsKeyDown(rl.KeyA) {
		g.extend = true
	}
}

func (g *Game) Update() {
	Assert(g.snake != nil, "Snake not intialized")

	next := g.snake.position

	switch g.direction {
	case Right:
		next.x = Clamp(next.x+1, 0, MapSize-1)
		break
	case Left:
		next.x = Clamp(next.x-1, 0, MapSize-1)
		break
	case Down:
		next.y = Clamp(next.y+1, 0, MapSize-1)
		break
	case Up:
		next.y = Clamp(next.y-1, 0, MapSize-1)
		break
	}

	if g.snake.position == next {
		return
	}

	if g.snake.CollidesWith(next) {
		return
	}

	g.snake.Update(next, g.extend)
	g.extend = false
}

func (g *Game) Draw() {
	Assert(g.snake != nil, "Snake not intialized")

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawFPS(winWidth-120, 10)
	g.snake.Draw()
	rl.EndDrawing()
}

// MAIN

func main() {
	winWidth, winHeight = 800, 600

	// rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(
		winWidth,
		winHeight,
		"sanke",
	)
	defer rl.CloseWindow()
	// rl.SetTargetFPS(10)

	game := NewGame()
	stepTime = StepTime

	for !rl.WindowShouldClose() {
		winWidth, winHeight = int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
		deltaTime = rl.GetFrameTime()
		game.Input()
		stepTime += deltaTime

		if stepTime >= StepTime {
			stepTime = 0
			game.Update()
		}

		game.Draw()
	}
}

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func Clamp(f, lo, hi int32) int32 {
	if f < lo {
		return lo
	}
	if f > hi {
		return hi
	}
	return f
}
