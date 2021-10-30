package tamtam

import "encoding/json"

type TileMap map[Vec2Di]SquareGlues

func (tiles TileMap) MarshalJSON() ([]byte, error) {
	var toMarshal map[string]SquareGlues = make(map[string]SquareGlues)
	for pos, tile := range tiles {
		marshaledPos, err := json.Marshal(pos)
		if err != nil {
			return nil, err
		}
		toMarshal[string(marshaledPos)] = tile
	}
	return json.Marshal(toMarshal)
}

func (tiles *TileMap) UnmarshalJSON(b []byte) error {
	var marshaled map[string]SquareGlues
	err := json.Unmarshal(b, &marshaled)

	if err != nil {
		return err
	}

	*tiles = make(TileMap)

	for encodedPos, tile := range marshaled {
		var pos Vec2Di

		err := json.Unmarshal([]byte(encodedPos), &pos)

		if err != nil {
			return err
		}

		(*tiles)[pos] = tile
	}

	return nil
}

func (tiles TileMap) IsEqualTo(otherTiles TileMap) bool {
	if len(tiles) != len(otherTiles) {
		return false
	}

	for key, value := range tiles {
		if otherValue, ok := otherTiles[key]; !ok {
			return false
		} else {
			if !value.IsEqualTo(otherValue) {
				return false
			}
		}
	}

	return true
}
