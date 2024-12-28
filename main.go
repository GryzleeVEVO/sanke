/*
	TODO:
	- Fix bugs
	- Fix draw
	- Add menu
	- Add win/lose state
	- Add scoring system
*/

package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var winWidth, winHeight int32
var stepTime, deltaTime float32
var rng *rand.Rand

const (
	MapSize          = 20
	StepTime float32 = 0.5 // Time between state updates
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

// FIXME: Rectangles drawn instead of squares.
func (s *SnakeSegment) Draw() {
	rl.DrawRectangle(
		int32(s.position.x*(winWidth/MapSize)),
		int32(s.position.y*(winHeight/MapSize)),
		winWidth/MapSize, winHeight/MapSize,
		rl.LightGray,
	)

	if s.next != nil {
		s.next.Draw()
	}
}

// FRUIT

type Fruit struct {
	position Position
}

// FIXME: Rectangles drawn instead of squares.
func (f *Fruit) Draw() {
	rl.DrawRectangle(
		int32(f.position.x*(winWidth/MapSize)),
		int32(f.position.y*(winHeight/MapSize)),
		winWidth/MapSize, winHeight/MapSize,
		rl.Red)
}

// GAME

// TODO: Add win/loss states
// TODO: Add scoring system
// TODO: Add main menu
// TODO: Fancier in-game graphics. Nothing too complicated tho
type Game struct {
	snake     *SnakeSegment
	fruit     Fruit
	direction Direction
	extend    bool
}

func NewGame() *Game {
	game := Game{
		snake: &SnakeSegment{
			Position{10, 10}, nil,
		},
		fruit:     Fruit{Position{int32(RandomInt(0, MapSize)), int32(RandomInt(0, MapSize))}},
		direction: Left,
		extend:    false,
	}

	for game.snake.CollidesWith(game.fruit.position) {
		game.fruit.position.x = int32(RandomInt(0, MapSize))
		game.fruit.position.y = int32(RandomInt(0, MapSize))
	}

	return &game
}

// FIXME: U-turns can happen when pressing first a valid direction, then an invalid one
// e.g. if the snake is facing RIGTH, and before the update UP and LEFT are pressed in
// that order, the snake will u-turn
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

	if g.snake.position == next || g.snake.CollidesWith(next) {
		return
	}

	// FIXME: Next fruit spawns a frame AFTER the previous is eaten
	if g.snake.CollidesWith(g.fruit.position) {
		g.extend = true

		for g.snake.CollidesWith(g.fruit.position) {
			g.fruit.position.x = int32(RandomInt(0, MapSize))
			g.fruit.position.y = int32(RandomInt(0, MapSize))
		}
	}

	g.snake.Update(next, g.extend)
	g.extend = false
}

func (g *Game) Draw() {
	Assert(g.snake != nil, "Snake not intialized")

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	// rl.DrawFPS(winWidth-120, 10)
	g.fruit.Draw()
	g.snake.Draw()
	rl.EndDrawing()
}

// MAIN

func main() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
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

// Panic if condition is not true
func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

// Clamp f to [lo, hi]
func Clamp(f, lo, hi int32) int32 {
	if f < lo {
		return lo
	}
	if f > hi {
		return hi
	}
	return f
}

// Get random integer between [min, max)
func RandomInt(min, max int) int {
	return rng.Intn(max-min) + min
}
