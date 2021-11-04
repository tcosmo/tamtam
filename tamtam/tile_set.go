package tamtam

import (
	"errors"
	"strconv"

	primes "github.com/fxtlabs/primes"
)

const NULL_GLUE string = ""

// North, East, South, West
type SquareGlues [4]string

func (glues SquareGlues) IsEqualTo(otherGlues SquareGlues) bool {
	for i, val := range glues {
		if otherGlues[i] != val {
			return false
		}
	}
	return true
}

type TileSet map[string]SquareGlues

func (tileSet TileSet) IsEqualTo(otherTileSet TileSet) bool {

	if len(tileSet) != len(otherTileSet) {
		return false
	}

	for name, val := range tileSet {
		if otherTileSet[name] != val {
			return false
		}
	}
	return true
}

func (tileSet TileSet) MatchTiles(glueConstraints SquareGlues, threshold int) (matches []SquareGlues) {

	for _, tileType := range tileSet {
		var count = 0
		for i := 0; i < 4; i += 1 {
			if glueConstraints[i] == NULL_GLUE {
				continue
			}

			if glueConstraints[i] != tileType[i] {
				break
			}
			count += 1
		}
		if count >= threshold {
			matches = append(matches, tileType)
		}
	}

	return matches
}

// Creates a Chinese Remainder Tile Set
func NewCrtTileSet(p int, q int) (tileSet TileSet, err error) {

	tileSet = make(TileSet)

	if !primes.Coprime(p, q) {
		return tileSet, errors.New("p and q must be co-primes")
	}

	for i := 0; i < p*q; i += 1 {
		tileSet[strconv.Itoa(i)] = SquareGlues{strconv.FormatInt(int64(i)/3, 10), strconv.FormatInt(int64(i)%3, 10), strconv.FormatInt(int64(i)%2, 10), strconv.FormatInt(int64(i)/2, 10)}
	}

	return tileSet, nil
}

// Returns the name of the tile or error if tile not in tile set
func (tileSet TileSet) GetTileName(tile SquareGlues) (tileName string, err error) {

	for name, tileType := range tileSet {
		if tileType.IsEqualTo(tile) {
			return name, nil
		}
	}

	return "", errors.New("The tile type is not in the tile set")
}
