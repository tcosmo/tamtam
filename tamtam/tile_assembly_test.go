package tamtam

import (
	"encoding/json"
	"testing"
)

// Testing successful assembly of simple scenario with a CRT tile set
// Only high level property of the final assembly is tested (its size)
// We are mainly testing against errors that can occur at each assembly steps
// Also testing that JSON serializing/deserializing is working well
func TestAssemblyAndSerialization(t *testing.T) {

	SIZE := 20
	FINAL_ASSEMBLY_SIZE := SIZE*SIZE + 2*SIZE
	tileSet, err := NewCrtTileSet(2, 11)

	if err != nil {
		t.Fatalf(`%v`, err)
		return
	}

	var initialAssembly = make(map[Vec2Di]SquareGlues)

	for i := 0; i < SIZE; i += 1 {
		initialAssembly[Vec2Di{-1, i}] = SquareGlues{NULL_GLUE, "0", NULL_GLUE, NULL_GLUE}
	}

	initialAssembly[Vec2Di{0, -1}] = SquareGlues{"1", NULL_GLUE, NULL_GLUE, NULL_GLUE}

	for i := 0; i < SIZE-1; i += 1 {
		initialAssembly[Vec2Di{1 + i, -1}] = SquareGlues{"0", NULL_GLUE, NULL_GLUE, NULL_GLUE}
	}

	var assembly = NewAssembly(tileSet, initialAssembly, 2)

	didGrow, err := assembly.GrowSync(true)

	for didGrow && err == nil {
		didGrow, err = assembly.GrowSync(true)
	}

	if err != nil {
		t.Fatalf(`%v`, err)
		return
	}

	if assembly.Size() != FINAL_ASSEMBLY_SIZE {
		t.Fatalf(`Assembly size %d != %d`, assembly.Size(), FINAL_ASSEMBLY_SIZE)
	}

	b, err := json.Marshal(assembly)

	if err != nil {
		t.Fatalf(`%v`, err)
	}

	var newAssembly TileAssembly

	err = json.Unmarshal(b, &newAssembly)

	if err != nil {
		t.Fatalf(`%v`, err)
	}

	if !newAssembly.IsEqualTo(assembly) {
		t.Fatalf(`%v`, err)
	}
}
