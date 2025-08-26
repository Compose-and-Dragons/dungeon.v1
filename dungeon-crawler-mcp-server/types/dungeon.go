package types

type Dungeon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Rooms       []Room `json:"rooms"`
	EntranceCoords Coordinates `json:"entrance_coords"`
	ExitCoords Coordinates `json:"exit_coords"`
}

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Room struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	IsEntrance  bool        `json:"is_entrance"`
	IsExit      bool        `json:"is_exit"`
	Coordinates Coordinates `json:"coordinates"`
	Visited     bool        `json:"visited"`
	HasMonster  bool        `json:"has_monster"`
	HasNPC      bool        `json:"has_npc"`
	HasTreasure bool        `json:"has_treasure"`
}
