package data

import "github.com/openai/openai-go/v2"

func GetRoomSchema() openai.ResponseFormatJSONSchemaJSONSchemaParam {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type": "string",
			},
			"description": map[string]any{
				"type": "string",
			},
		},
		"required": []string{"name", "description"},
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "room_info",
		Description: openai.String("name and description of the room"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}
	return schemaParam
}

func GetMonsterSchema() openai.ResponseFormatJSONSchemaJSONSchemaParam {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"kind": map[string]any{
				"type": "string",
				"enum": []string{
					"skeleton", "zombie", "goblin", "orc", "troll", "dragon", "werewolf", "vampire",
				},
			},
			"name": map[string]any{
				"type": "string",
			},
			"description": map[string]any{
				"type": "string",
			},
			"health": map[string]any{
				"type":    "integer",
				"minimum": 0,
			},
			"strength": map[string]any{
				"type":    "integer",
				"minimum": 0,
			},
		},
		"required": []string{"kind", "name", "description", "health", "strength"},
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "monster_info",
		Description: openai.String("details about the monster in the room"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}
	return schemaParam
}