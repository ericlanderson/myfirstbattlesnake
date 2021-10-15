package main

import (
	"testing"
)

func TestNeckAvoidance(t *testing.T) {
	// Arrange
	me := Battlesnake{
		// Length 3, facing right
		Head: Coord{X: 2, Y: 0},
		Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	above := Battlesnake{
		// Length 2, facing right
		// Above "me"
		Head: Coord{X: 1, Y: 1},
		Body: []Coord{{X: 1, Y: 1}, {X: 0, Y: 1}},
	}
	right := Battlesnake{
		// Length 3, facing left
		// To the right of me
		Head: Coord{X: 3, Y: 0},
		Body: []Coord{{X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Height: 11,
			Width:  11,
			Snakes: []Battlesnake{me, above, right},
			Food:   []Coord{{0, 0}, {10, 10}},
		},
		You: me,
	}

	// Act 1,000x (this isn't a great way to test, but it's okay for starting out)
	for i := 0; i < 10; i++ {
		_ = move(state)
		// nextMove := move(state)
		// Assert never move left
		// if nextMove.Move == "left" {
		// 	t.Errorf("snake moved onto its own neck, %s", nextMove.Move)
		// }
	}
}

// TODO: More GameState test cases!
