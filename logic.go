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
		Color:      "#888888", // TODO: Personalize
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

	// Step 0: Don't let your Battlesnake move back in on it's own neck
	myHead := state.You.Body[0] // Coordinates of your head; same as state.You.Head
	// myNeck := state.You.Body[1] // Coordinates of body piece directly behind your head (your "neck")

	// We do this in step 2.

	// TODO: Step 1 - Don't hit walls.
	// Use information in GameState to prevent your Battlesnake from moving beyond the boundaries of the board.
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

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

	// TODO: Step 2 - Don't hit yourself.
	// Use information in GameState to prevent your Battlesnake from colliding with itself.
	// mybody := state.You.Body

	torso := state.You.Body[1:]
	mysnake := Battlesnake{Body: torso}
	snakes := state.Board.Snakes
	snakes = append(snakes, mysnake)

	for move, isSafe := range possibleMoves {
		if isSafe {
			switch move {
			case "up":
				if checkNextAgainstSnakes(Coord{myHead.X, myHead.Y + 1}, snakes) {
					possibleMoves["up"] = false
				}
			case "down":
				if checkNextAgainstSnakes(Coord{myHead.X, myHead.Y - 1}, snakes) {
					possibleMoves["down"] = false
				}
			case "right":
				if checkNextAgainstSnakes(Coord{myHead.X + 1, myHead.Y}, snakes) {
					possibleMoves["right"] = false
				}
			case "left":
				if checkNextAgainstSnakes(Coord{myHead.X - 1, myHead.Y}, snakes) {
					possibleMoves["left"] = false
				}
			}
		}
	}

	// TODO: Step 3 - Don't collide with others.
	// Use information in GameState to prevent your Battlesnake from colliding with others.

	// See number two above. We check both our torso and all the other snakes.

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.
	closestFood := state.Board.Food[0]
	closestFoodDistance := distance(myHead, closestFood)
	for _, food := range state.Board.Food[1:] {
		nextFoodDistance := distance(myHead, food)
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
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
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

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		if len(desiredMoves) == 0 {
			nextMove = safeMoves[rand.Intn(len(safeMoves))]

		} else {
			nextMove = desiredMoves[rand.Intn(len(desiredMoves))]
		}
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}

func checkNextAgainstSnakes(nextCoord Coord, snakes []Battlesnake) bool {
	for _, snake := range snakes {
		for _, bodyCoord := range snake.Body {
			if nextCoord == bodyCoord {
				return true
			}
		}
	}
	return false
}

func distance(a Coord, b Coord) float64 {
	return math.Sqrt(math.Pow(float64(b.X-a.X), 2) + math.Pow(float64(b.Y-a.Y), 2))
}
