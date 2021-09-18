package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type player struct {
	id    int
	poin  int
	dices int
}

func newPlayer(id, dice int) *player {
	return &player{
		id:    id,
		dices: dice,
	}
}

func (p *player) RollDices(randomizer *rand.Rand, diceFromLeft int) (loseDice int, isWin bool) {
	var winDice int

	for i := 0; i < p.dices; i++ {
		dice := randomizer.Intn(6)
		if dice == 5 {
			winDice++
		}
		if dice == 0 {
			loseDice++
		}
	}

	p.poin += winDice
	p.dices -= (loseDice + winDice)
	p.dices += diceFromLeft
	return loseDice, p.dices < 1
}

type game struct {
	players    []*player
	playerLeft int
	randomizer *rand.Rand
}

func newGame(player, dice int) *game {
	s1 := rand.NewSource(time.Now().UnixNano())
	game := game{
		playerLeft: player,
		randomizer: rand.New(s1),
	}

	for i := 1; i < player+1; i++ {
		game.players = append(game.players, newPlayer(i, dice))
	}
	return &game
}

func (g *game) Run() {
	diceFromLeft := 0
	isWin := false

	for turn := 0; g.playerLeft > 1; turn++ {
		player := g.players[turn%g.playerLeft]
		diceFromLeft, isWin = player.RollDices(g.randomizer, diceFromLeft)
		log.Printf("player: %d, loseDice: %d, poin: %d, diceLeft: %d", player.id, diceFromLeft, player.poin, player.dices)
		if isWin {
			turn = turn % g.playerLeft
			g.leftGame(turn % g.playerLeft)
			log.Println("player left game:", player)
			turn--
		}
	}

	player := g.players[0]
	diceFromLeft, isWin = player.RollDices(g.randomizer, diceFromLeft)
	log.Printf("player: %d, loseDice: %d, poin: %d, diceLeft: %d", player.id, diceFromLeft, player.poin, player.dices)
}

func (g *game) leftGame(playerIndex int) {
	currentPlayer := g.players[playerIndex]
	for i := playerIndex; i < g.playerLeft-1; i++ {
		g.players[i] = g.players[i+1]
	}
	g.players[g.playerLeft-1] = currentPlayer
	g.playerLeft--
}

func (g *game) printWinner() {
	winner := g.players[0]
	for _, player := range g.players[1:] {
		if player.poin > winner.poin {
			winner = player
		}
	}

	log.Printf("Winner: player %d, poin: %d", winner.id, winner.poin)
}

func main() {
	if len(os.Args) < 3 {
		log.Println("Usage: app <num_player> <num_dice>")
		os.Exit(1)
	}

	numPlayer, err := strconv.Atoi(os.Args[1])
	if err != nil || numPlayer == 0 {
		log.Println("num_player should be int greater than 1!")
		os.Exit(1)
	}

	numDice, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Println("num_dice should be int!")
		os.Exit(1)
	}

	game := newGame(numPlayer, numDice)

	game.Run()
	game.printWinner()
}
