package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MAP_SIZE = 20
)

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

type Position struct {
	x, y int32
}

type SnakeSegment struct {
	position Position
	next     *SnakeSegment
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

type Game struct {
	snakeHead *SnakeSegment
	direction Direction
	extend    bool
}

func NewGame() *Game {
	game := Game{
		snakeHead: &SnakeSegment{
			Position{10, 10}, nil,
		},

		direction: Left,
		extend:    false,
	}

	return &game
}

func (g *Game) Input() {
	if rl.IsKeyDown(rl.KeyRight) {
		if g.direction != Left {
			g.direction = Right
		}
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		if g.direction != Right {
			g.direction = Left
		}
	}
	if rl.IsKeyDown(rl.KeyDown) {
		if g.direction != Up {
			g.direction = Down
		}
	}
	if rl.IsKeyDown(rl.KeyUp) {
		if g.direction != Down {
			g.direction = Up
		}
	}

	if rl.IsKeyDown(rl.KeyA) {
		g.extend = true
	}

}

func (g *Game) Update() {
	assert(g.snakeHead != nil, "Snake not intialized")

	next := g.snakeHead.position

	switch g.direction {
	case Right:
		if next.x < MAP_SIZE-1 {
			next.x += 1
		}
		break
	case Left:
		if next.x > 0 {
			next.x -= 1
		}
		break
	case Down:
		if next.y < MAP_SIZE-1 {
			next.y += 1
		}
		break
	case Up:
		if next.y > 0 {
			next.y -= 1
		}
		break
	}

	if g.snakeHead.position != next {
		g.snakeHead.Update(next, g.extend)
		g.extend = false
	}
}

func (g *Game) Draw() {
	assert(g.snakeHead != nil, "Snake not intialized")

	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)
	g.snakeHead.Draw()

	rl.EndDrawing()
}

func main() {
	rl.InitWindow(
		800,
		600,
		"sanke",
	)
	defer rl.CloseWindow()

	rl.SetTargetFPS(10)

	game := NewGame()

	for !rl.WindowShouldClose() {
		game.Input()
		game.Update()
		game.Draw()
	}
}
