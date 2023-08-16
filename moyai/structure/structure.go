package structure

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type Structure struct {
	dimension [3]int
	blocks    map[cube.Pos]world.Block
}

func (s Structure) Dimensions() [3]int {
	return s.dimension
}

func (s Structure) At(x int, y int, z int, _ func(x int, y int, z int) world.Block) (world.Block, world.Liquid) {
	pos := cube.Pos{x, y, z}
	b := s.blocks[pos]
	return b, nil
}

func (s Structure) Set(x int, y int, z int, b world.Block) {
	pos := cube.Pos{x, y, z}
	s.blocks[pos] = b
}

func GenerateBoxStructure(dim [3]int, floor ...world.Block) world.Structure {
	s := Structure{dimension: dim, blocks: map[cube.Pos]world.Block{}}

	for x := 0; x < dim[0]; x++ {
		for y := 0; y < dim[1]; y++ {
			for z := 0; z < dim[2]; z++ {
				s.Set(x, y, z, block.Barrier{})
			}
		}
	}

	for x := 1; x < dim[0]-1; x++ {
		for y := 1; y < dim[1]-1; y++ {
			for z := 1; z < dim[2]-1; z++ {
				s.Set(x, y, z, block.Air{})
			}
		}
	}

	for x := 1; x < dim[0]-1; x++ {
		for z := 1; z < dim[2]-1; z++ {
			s.Set(x, 0, z, floor[rand.Intn(len(floor))])
			s.Set(x, -1, z, block.Bedrock{})
		}
	}
	return s
}
