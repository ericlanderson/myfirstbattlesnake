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

type Boardmeta struct {
	Height int
	Width  int
}

// Enums the golang way

type Element int

const (
	Empty Element = iota
	Head
	Body
	Food
	OOB // Out Of Bounds
)

// type GridElement interface {
// 	Set(a Coord) Element
// 	Get(a Coord) Element
// 	GetRight(a Coord) Element
// 	GetLeft(a Coord) Element
// 	GetAbove(a Coord) Element
// 	GetBelow(a Coord) Element
// }

type Grid struct {
	Board  [][]Element
	Height int
	Width  int
}

func (g Grid) Get(a Coord) Element {
	return g.Board[a.X][a.Y]
}

func (g Grid) Set(a Coord, e Element) {
	g.Board[a.X][a.Y] = e
}

func (g Grid) GetRight(a Coord) Element {
	if a.X+1 >= g.Width {
		return OOB
	}
	return g.Board[a.X+1][a.Y]
}
func (g Grid) GetLeft(a Coord) Element {
	if a.X-1 < 0 {
		return OOB
	}
	return g.Board[a.X-1][a.Y]
}

func (g Grid) GetAbove(a Coord) Element {
	if a.Y+1 >= g.Height {
		return OOB
	}
	return g.Board[a.X][a.Y+1]
}

func (g Grid) GetBelow(a Coord) Element {
	if a.Y-1 < 0 {
		return OOB
	}
	return g.Board[a.X][a.Y-1]
}

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

	var myBoard = Boardmeta{state.Board.Height, state.Board.Width}

	// This is our board state where we record where each object (body or food)
	// is located. Makes lookups much easier/faster/efficient.
	var grid Grid
	grid.Board = make([][]Element, myBoard.Height)
	for i := range grid.Board {
		grid.Board[i] = make([]Element, myBoard.Width)
	}
	grid.Height = state.Board.Height
	grid.Width = state.Board.Width

	// record where each snake body element is on the grid.
	for _, snake := range state.Board.Snakes {
		// Can use this later to find "head" that are adjacent to safe spaces.
		// Likely want to avoid those spaces.
		grid.Set(Coord{snake.Head.X, snake.Head.Y}, Head)
		// Now record the remainder of the body coordinates
		for _, body := range snake.Body[1:] {
			grid.Set(Coord{body.X, body.Y}, Body)
		}
	}

	// record where each food element is on the grid.
	for _, food := range state.Board.Food {
		grid.Set(Coord{food.X, food.Y}, Food)
	}

	// Avoid hitting the walls
	log.Printf("Avoiding Walls ...")
	for move, _ := range possibleMoves {
		switch move {
		case "up":
			if myHead.Y == myBoard.Height-1 {
				possibleMoves["up"] = false
			}
		case "down":
			if myHead.Y == 0 {
				possibleMoves["down"] = false
			}
		case "right":
			if myHead.X == myBoard.Width-1 {
				possibleMoves["right"] = false
			}
		case "left":
			if myHead.X == 0 {
				possibleMoves["left"] = false
			}
		}

	}
	log.Printf(" Complete\n")

	// Avoid hitting myself or other snakes
	log.Printf("Avoiding Snakes ...")
	for move, isSafe := range possibleMoves {
		if isSafe {
			switch move {
			case "up":
				if grid.GetAbove(myHead) == Head || grid.GetAbove(myHead) == Body {
					possibleMoves["up"] = false
				}
			case "down":
				if grid.GetBelow(myHead) == Head || grid.GetBelow(myHead) == Body {
					possibleMoves["down"] = false
				}
			case "right":
				if grid.GetRight(myHead) == Head || grid.GetRight(myHead) == Body {
					possibleMoves["right"] = false
				}
			case "left":
				if grid.GetLeft(myHead) == Head || grid.GetLeft(myHead) == Body {
					possibleMoves["left"] = false
				}
			}
		}
	}
	log.Printf(" Complete\n")

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

	// The list of safe moves
	safeMoves := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}
	log.Printf("Safe moves: %s", safeMoves)

	// The list of moves which are both safe and desired
	desiredMoves := []string{}
	for move, isDesired := range foodMoves {
		if isDesired && possibleMoves[move] {
			desiredMoves = append(desiredMoves, move)
		}
	}
	log.Printf("Desired moves: %s", desiredMoves)

	finalMoves := []string{}
	// Avoid moving adjacent to another snake head
	for _, move := range desiredMoves {
		var adjacent []Element
		switch move {
		case "up":
			c := Coord{myHead.X, myHead.Y + 1}
			adjacent = append(adjacent, grid.GetAbove(c))
			adjacent = append(adjacent, grid.GetRight(c))
			adjacent = append(adjacent, grid.GetLeft(c))
			if !find(adjacent, Head) {
				finalMoves = append(finalMoves, "up")
			}
		case "down":
			c := Coord{myHead.X, myHead.Y - 1}
			adjacent = append(adjacent, grid.GetBelow(c))
			adjacent = append(adjacent, grid.GetRight(c))
			adjacent = append(adjacent, grid.GetLeft(c))
			if !find(adjacent, Head) {
				finalMoves = append(finalMoves, "down")
			}
		case "right":
			c := Coord{myHead.X + 1, myHead.Y}
			adjacent = append(adjacent, grid.GetRight(c))
			adjacent = append(adjacent, grid.GetAbove(c))
			adjacent = append(adjacent, grid.GetBelow(c))
			if !find(adjacent, Head) {
				finalMoves = append(finalMoves, "right")
			}
		case "left":
			c := Coord{myHead.X - 1, myHead.Y}
			adjacent = append(adjacent, grid.GetLeft(c))
			adjacent = append(adjacent, grid.GetAbove(c))
			adjacent = append(adjacent, grid.GetBelow(c))
			if !find(adjacent, Head) {
				finalMoves = append(finalMoves, "left")
			}
		}
	}

	// Finally, choose a move from the available safe moves.
	var nextMove string

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

func distanceBetween(a Coord, b Coord) float64 {
	dX := float64(b.X - a.X)
	dY := float64(b.Y - a.Y)
	return math.Sqrt(math.Pow(dX, 2) + math.Pow(dY, 2))
}

func find(slice []Element, e Element) bool {
	for _, item := range slice {
		if item == e {
			return true
		}
	}
	return false
}
