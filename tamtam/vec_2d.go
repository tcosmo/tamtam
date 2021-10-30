package tamtam

// x, y
type Vec2Di [2]int

var North Vec2Di = Vec2Di{0, 1}
var East Vec2Di = Vec2Di{1, 0}
var South Vec2Di = Vec2Di{0, -1}
var West Vec2Di = Vec2Di{-1, 0}

var CardinalPoints = []Vec2Di{North, East, South, West}

func (a Vec2Di) Add(b Vec2Di) Vec2Di {
	return Vec2Di{a[0] + b[0], a[1] + b[1]}
}

func (a Vec2Di) Neighbors() (neighbors [4]Vec2Di) {
	for i, cardinalPoint := range CardinalPoints {
		neighbors[i] = a.Add(cardinalPoint)
	}
	return neighbors
}
