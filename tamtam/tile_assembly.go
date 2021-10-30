package tamtam

import (
	"encoding/json"
	"errors"
)

type PosAndTile struct {
	Pos  Vec2Di
	Tile SquareGlues
}

type TileAssembly struct {
	tileSet                      TileSet
	tileMap                      TileMap
	threshold                    int
	emptyPositionsAboveThreshold map[Vec2Di]bool
	newlyAddedTiles              []PosAndTile
}

func (assembly TileAssembly) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TileSet   TileSet `json:"tile_set"`
		TileMap   TileMap `json:"tile_map"`
		Threshold int     `json:"threshold"`
	}{
		TileSet:   assembly.tileSet,
		TileMap:   assembly.tileMap,
		Threshold: assembly.threshold,
	})
}

func (assembly *TileAssembly) UnmarshalJSON(b []byte) error {

	var rawAssembly struct {
		TileSet   TileSet `json:"tile_set"`
		TileMap   TileMap `json:"tile_map"`
		Threshold int     `json:"threshold"`
	}
	err := json.Unmarshal(b, &rawAssembly)

	if err != nil {
		return err
	}

	*assembly = NewAssembly(rawAssembly.TileSet, rawAssembly.TileMap, rawAssembly.Threshold)

	return nil
}

func NewAssembly(tileSet TileSet, initialTiles map[Vec2Di]SquareGlues, threshold int) (assembly TileAssembly) {
	assembly.tileSet = tileSet
	assembly.threshold = threshold

	assembly.tileMap = make(map[Vec2Di]SquareGlues)
	assembly.emptyPositionsAboveThreshold = make(map[Vec2Di]bool)

	for pos, tile := range initialTiles {
		assembly.AddTile(pos, tile)
	}

	return assembly
}

func (assembly TileAssembly) Size() int {
	return len(assembly.tileMap)
}

func (assembly TileAssembly) neighboringGlues(pos Vec2Di) (glues SquareGlues) {
	for i, nei := range pos.Neighbors() {
		if val, ok := assembly.tileMap[nei]; ok {
			glues[i] = val[(i+2)%4]
		} else {
			glues[i] = NULL_GLUE
		}
	}
	return glues
}

func (assembly TileAssembly) isPosAboveThreshold(pos Vec2Di) bool {
	var count = 0
	for _, glue := range assembly.neighboringGlues(pos) {
		if glue != NULL_GLUE {
			count += 1
		}
	}
	return count >= assembly.threshold
}

func (assembly *TileAssembly) AddTile(pos Vec2Di, tile SquareGlues) {
	assembly.tileMap[pos] = tile

	assembly.newlyAddedTiles = append(assembly.newlyAddedTiles, PosAndTile{Pos: pos, Tile: tile})
	delete(assembly.emptyPositionsAboveThreshold, pos)

	for _, nei := range pos.Neighbors() {
		if _, ok := assembly.tileMap[nei]; assembly.isPosAboveThreshold(nei) && !ok {
			assembly.emptyPositionsAboveThreshold[nei] = true
		}
	}

}

// Performs a synchronous growth step
func (assembly *TileAssembly) GrowSync(directed bool) (bool, error) {

	var toAdd []PosAndTile

	for pos := range assembly.emptyPositionsAboveThreshold {
		var matches = assembly.tileSet.MatchTiles(assembly.neighboringGlues(pos), assembly.threshold)

		if len(matches) > 1 && directed {
			return false, errors.New("two different tiles fit the same position, this is not allowed in directed setting")
		}

		for _, tile := range matches {
			toAdd = append(toAdd, PosAndTile{Pos: pos, Tile: tile})
		}
	}

	for _, posAndTile := range toAdd {
		assembly.AddTile(posAndTile.Pos, posAndTile.Tile)
	}

	var anyGrowth = len(toAdd) >= 1
	return anyGrowth, nil
}

func (assembly TileAssembly) IsEqualTo(otherAssembly TileAssembly) bool {
	return assembly.threshold == otherAssembly.threshold && assembly.tileSet.IsEqualTo(otherAssembly.tileSet) && assembly.tileMap.IsEqualTo(otherAssembly.tileMap)
}

// Returns the tiles (with their position) that were added at the last round of growth
func (assembly TileAssembly) GetNewlyAddedTiles() []PosAndTile {

	return assembly.newlyAddedTiles
}

func (assembly *TileAssembly) FlushNewlyAddedTiles() {
	assembly.newlyAddedTiles = []PosAndTile{}
}
