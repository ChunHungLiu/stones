package core

import (
	"fmt"
)

type room struct {
	Offset, Size Offset
}

func (r room) rand() Offset {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(r)
			panic(e)
		}
	}()
	return Offset{r.Offset.X + RandRange(1, r.Size.X-2), r.Offset.Y + RandRange(1, r.Size.Y-2)}
}

func (r room) inside(o Offset) bool {
	return r.Offset.X <= o.X && o.X < r.Offset.X+r.Size.X &&
		r.Offset.Y <= o.Y && o.Y < r.Offset.Y+r.Size.Y
}

func Dungeon(numRooms, minRoomSize, maxRoomSize int, f BoolTileFactory) map[*Tile]struct{} {
	// TODO Added in better maze gen customization
	maze := abstractBraid(numRooms, .25, 0, 1)
	rooms := make(map[*mazenode]room)
	tiles := make(map[*Tile]struct{})

	// create rooms
	for _, nodes := range maze.Nodes {
		for _, node := range nodes {
			rooms[node] = room{
				Offset{maxRoomSize * node.Pos.X, maxRoomSize * node.Pos.Y},
				Offset{minRoomSize, minRoomSize},
			}
		}
	}

	// create room tiles
	for _, room := range rooms {
		for x := room.Offset.X; x < room.Offset.X+room.Size.X; x++ {
			for y := room.Offset.Y; y < room.Offset.Y+room.Size.Y; y++ {
				tiles[f(Offset{x, y}, true)] = struct{}{}
				// TODO connect room tiles
			}
		}
	}

	// create corridors
	frontier := []*mazenode{maze.Origin}
	enqued := map[*mazenode]struct{}{maze.Origin: {}}
	closed := map[*mazenode]struct{}{}
	for len(frontier) != 0 {
		curr := frontier[0]
		frontier = frontier[1:]

		if _, done := closed[curr]; done {
			continue
		}
		closed[curr] = struct{}{}

		for _, adj := range curr.Edges {
			if _, done := closed[adj]; done {
				continue
			}
			if _, seen := enqued[adj]; !seen {
				frontier = append(frontier, adj)
			}

			srcroom, srcok := rooms[curr]
			if !srcok {
				panic(curr)
			}
			dstroom, dstok := rooms[adj]
			if !dstok {
				panic(curr)
			}

			src, dst := srcroom.rand(), dstroom.rand()

			for src.X != dst.X {
				if !rooms[curr].inside(src) && !rooms[adj].inside(src) {
					tiles[f(src, true)] = struct{}{}
				}
				if src.X < dst.X {
					src.X++
				} else {
					src.X--
				}
			}
			for src.Y != dst.Y {
				if !rooms[curr].inside(src) && !rooms[adj].inside(src) {
					tiles[f(src, true)] = struct{}{}
				}
				if src.Y < dst.Y {
					src.Y++
				} else {
					src.Y--
				}
			}
			// TODO connect corridor tiles
		}
	}

	// TODO reuse maze connection stuff to fully connect dungeon and walls

	return tiles
}
