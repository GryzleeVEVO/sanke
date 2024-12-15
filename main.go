package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	MAP_SIZE = 20
)

type Direction int

const (
	Left int = iota
	Right
	Up
	Down
)

type Position struct {
	x, y int32
}

func main() {
	rl.InitWindow(
		800,
		600,
		"sanke",
	)
	defer rl.CloseWindow()

	rl.SetTargetFPS(10)

	direction := Left
	pos := Position{10, 10}

	data := [MAP_SIZE][MAP_SIZE]int{}
	data[10][10] = 1

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeyRight) {
			if direction != Left {
				direction = Right
			}
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			if direction != Right {
				direction = Left
			}
		}
		if rl.IsKeyDown(rl.KeyDown) {
			if direction != Up {
				direction = Down
			}
		}
		if rl.IsKeyDown(rl.KeyUp) {
			if direction != Down {
				direction = Up
			}
		}

		data[pos.x][pos.y] = 0

		switch direction {
		case Right:
			if pos.x < MAP_SIZE-1 {
				pos.x += 1
			}
			break
		case Left:
			if pos.x > 0 {
				pos.x -= 1
			}
			break
		case Down:
			if pos.y < MAP_SIZE-1 {
				pos.y += 1
			}
			break
		case Up:
			if pos.y > 0 {
				pos.y -= 1
			}
			break
		}

		data[pos.x][pos.y] = 1

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawFPS(0, 0)

		for i, inner := range data {
			for j := range inner {
				if data[i][j] == 1 {
					rl.DrawRectangle(
						int32(i*10),
						int32(j*10),
						10, 10,
						rl.LightGray,
					)
				}
			}
		}

		rl.EndDrawing()
	}
}
