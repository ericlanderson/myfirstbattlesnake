package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
	"math"
	"math/rand"
)

// Enums the golang way

type GridContent int

const (
	Empty GridContent = iota
	Head
	Body
	Food
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",        // TODO: Your Battlesnake username
		Color:      "#880088", // TODO: Personalize
		Head:       "default", // TODO: Personalize
		Tail:       "default", // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {
	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	foodMoves := map[string]bool{
		"up":    false,
		"down":  false,
		"left":  false,
		"right": false,
	}

	myHead := state.You.Body[0] // Coordinates of your head; same as state.You.Head
	// myNeck := state.You.Body[1] // Coordinates of body piece directly behind your head (your "neck")

	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	// This is our board state where we record where each object (body or food)
	// is located. Makes lookups much easier/faster/efficient.
	grid := make([][]GridContent, boardHeight)
	for i := range grid {
		grid[i] = make([]GridContent, boardWidth)
	}

	// record where each snake body element is on the grid.
	for _, snake := range state.Board.Snakes {
		// Can use this later to find "head" that are adjacent to safe spaces.
		// Likely want to avoid those spaces.
		grid[snake.Head.X][snake.Head.Y] = Head
		// Now record the remainder of the body coordinates
		for _, body := range snake.Body[1:] {
			grid[body.X][body.Y] = Body
		}
	}

	// record where each food element is on the grid.
	for _, food := range state.Board.Food {
		grid[food.X][food.Y] = Food
	}

	// Avoid hitting the walls
	for move, _ := range possibleMoves {
		switch move {
		case "up":
			if myHead.Y == boardHeight-1 {
				possibleMoves["up"] = false
			}
		case "down":
			if myHead.Y == 0 {
				possibleMoves["down"] = false
			}
		case "right":
			if myHead.X == boardWidth-1 {
				possibleMoves["right"] = false
			}
		case "left":
			if myHead.X == 0 {
				possibleMoves["left"] = false
			}
		}

	}

	// Avoid hitting myself or other snakes
	for move, isSafe := range possibleMoves {
		if isSafe {
			switch move {
			case "up":
				up := Coord{myHead.X, myHead.Y + 1}
				if doesMoveContain(Head, up, grid) || doesMoveContain(Body, up, grid) {
					possibleMoves["up"] = false
				}
			case "down":
				down := Coord{myHead.X, myHead.Y - 1}
				if doesMoveContain(Head, down, grid) || doesMoveContain(Body, down, grid) {
					possibleMoves["down"] = false
				}
			case "right":
				right := Coord{myHead.X + 1, myHead.Y}
				if doesMoveContain(Head, right, grid) || doesMoveContain(Body, right, grid) {
					possibleMoves["right"] = false
				}
			case "left":
				left := Coord{myHead.X - 1, myHead.Y}
				if doesMoveContain(Head, left, grid) || doesMoveContain(Body, left, grid) {
					possibleMoves["left"] = false
				}
			}
		}
	}

	// Use information in GameState to seek out and find food.
	// Assume Food[0] is the closest.
	closestFood := state.Board.Food[0]
	closestFoodDistance := distanceBetween(myHead, closestFood)
	// Now check the rest of the food and find which is closest
	for _, food := range state.Board.Food[1:] {
		nextFoodDistance := distanceBetween(myHead, food)
		if nextFoodDistance < closestFoodDistance {
			closestFood = food
			closestFoodDistance = nextFoodDistance
		}
	}

	// Set our desired moves based on the current closest food
	if myHead.X < closestFood.X {
		foodMoves["right"] = true
	}
	if myHead.X > closestFood.X {
		foodMoves["left"] = true
	}
	if myHead.Y < closestFood.Y {
		foodMoves["up"] = true
	}
	if myHead.Y > closestFood.Y {
		foodMoves["down"] = true
	}

	// Finally, choose a move from the available safe moves.
	var nextMove string

	// The list of safe moves
	safeMoves := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	// The list of desired and safe moves
	desiredMoves := []string{}
	for move, isDesired := range foodMoves {
		if isDesired && possibleMoves[move] {
			desiredMoves = append(desiredMoves, move)
		}
	}

	finalMoves := []string{}
	// Avoid moving adjacent to another snake head
	for _, move := range desiredMoves {
		if isMoveNotAdjacentTo(Head, myHead, move, grid) {
			finalMoves = append(finalMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		if len(finalMoves) == 0 {
			nextMove = safeMoves[rand.Intn(len(safeMoves))]
			log.Printf("%s MOVE %d: No desired moves detected! Making random safe move: %s\n", state.Game.ID, state.Turn, nextMove)
		} else {
			nextMove = finalMoves[rand.Intn(len(finalMoves))]
			log.Printf("%s MOVE %d: Making random desired move: %s\n", state.Game.ID, state.Turn, nextMove)
		}
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}

func doesMoveContain(item GridContent, a Coord, grid [][]GridContent) bool {
	return grid[a.X][a.Y] == item
}

func distanceBetween(a Coord, b Coord) float64 {
	dX := float64(b.X - a.X)
	dY := float64(b.Y - a.Y)
	return math.Sqrt(math.Pow(dX, 2) + math.Pow(dY, 2))
}

func isMoveNotAdjacentTo(item GridContent, myHead Coord, move string, grid [][]GridContent) bool {
	switch move {
	case "up":
		if grid[myHead.X][myHead.Y+2] == item {
			return false
		}
		if grid[myHead.X-1][myHead.Y+1] == item {
			return false
		}
		if grid[myHead.X+1][myHead.Y+1] == item {
			return false
		}
	case "down":
		if grid[myHead.X][myHead.Y-2] == item {
			return false
		}
		if grid[myHead.X-1][myHead.Y-1] == item {
			return false
		}
		if grid[myHead.X+1][myHead.Y-1] == item {
			return false
		}
	case "right":
		if grid[myHead.X+2][myHead.Y] == item {
			return false
		}
		if grid[myHead.X+1][myHead.Y-1] == item {
			return false
		}
		if grid[myHead.X+1][myHead.Y+1] == item {
			return false
		}
	case "left":
		if grid[myHead.X-2][myHead.Y] == item {
			return false
		}
		if grid[myHead.X-1][myHead.Y+1] == item {
			return false
		}
		if grid[myHead.X-1][myHead.Y-1] == item {
			return false
		}
	}
	return true
}
