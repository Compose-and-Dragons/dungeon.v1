package types

type NPCType string

// Define different types of NPCs
const (
	Merchant NPCType = "merchant"
	Guard    NPCType = "guard"
	Sorcerer NPCType = "sorcerer"
	Healer   NPCType = "healer"
)

// Non-player character
type NonPlayerCharacter struct {
	Type     NPCType     `json:"type"`
	Name     string      `json:"name"`
	Race     string      `json:"race"`
	Position Coordinates `json:"position"`
	RoomID   string      `json:"room_id"`
}

// var AllNPCTypes = []NPCType{
// 	Merchant,
// 	Guard,
// 	Sorcerer,
// 	Healer,
// }

// func NPCTypeCount() int {
// 	return len(AllNPCTypes)
// }