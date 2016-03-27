package game

import (
	"encoding/json"
	"log"

	uuid "github.com/nu7hatch/gouuid"
)

// Game struct
type Game struct {
	ID        string
	Name      string
	Players   map[string]Client
	gameTurns chan (GameTurn)
	Scenario  Scenario
}

// GameTurn represent an action perform  by some entity (client, ia, hero)
type GameTurn struct {
	Type string `json:"type"`
}

// NewGame ceate new game with clients
func NewGame(name string) (Game, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Game{}, err
	}
	return Game{
		ID:      u4.String(),
		Players: make(map[string]Client),
		Name:    name,
	}, nil
}

// Join a client to Game
func (g *Game) Join(client Client) {
	if _, exist := g.Players[client.ID]; !exist {
		g.Players[client.ID] = client
		log.Printf("add client [%s:%s] to game [%s:%s]", client.ID, client.Name, g.ID, g.Name)
	}
}

func (g *Game) Launch() {
	log.Printf("Launch Game %s with %d players", g.ID, len(g.Players))
	go g.Scenario.Process(g.gameTurns)
	for _, c := range g.Players {
		log.Printf("Player in da Game %s ", c.String())
		if err := sendStartMessage(c); err != nil {
			continue
		}
		go c.Process(g.gameTurns)
	}
	go gameTurnsHandler(g.gameTurns)
}

func sendStartMessage(c Client) error {
	e := json.NewEncoder(c.Conn)
	return e.Encode("{insert start mesage here}")
}

func gameTurnsHandler(gts chan GameTurn) {
	for {
		gt := <-gts
		log.Printf("Game turn received : %s", gt.Type)
	}
}