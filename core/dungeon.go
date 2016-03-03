package core

// MapGenInt generates Tiles for int values to form dungeon maps.
type MapGenInt func(o Offset, tiletype int) *Tile

// Constants used by Dungeon as tile type values passed to MapGenInt. It is
// likely that any user defined MapGenInt will switch on these constants.
const (
	TileTypeRoom = 1 << iota
	TileTypeCorridor
	TileTypeWall
	TileTypeDoor
)

type room struct {
	X, Y, W, H int
	Tiles      []*Tile
}

func (r *room) ConnectX(o *room, f MapGenInt) []*Tile {
	var tiles []*Tile

	minY := Max(r.Y, o.Y) + 1
	maxY := Min(r.Y+r.H, o.Y+o.H) - 2
	var srcY, dstY int
	if minY < maxY && RandBool() {
		srcY = RandRange(minY, maxY)
		dstY = srcY
	} else {
		srcY = RandRange(r.Y+1, r.Y+r.H-2)
		dstY = RandRange(o.Y+1, o.Y+o.H-2)
	}
	srcX, dstX := r.X+r.W/2, o.X+o.W/2

	midX := (srcX + dstX) / 2
	if InRange(midX, r.X, r.X+r.W) || InRange(midX, o.X, o.X+o.W) {
		midX = (r.X + r.W + o.X) / 2
	}
	if InRange(midX, r.X, r.X+r.W) || InRange(midX, o.X, o.X+o.W) {
		midX = (r.X + o.X + o.W) / 2
	}

	for srcX != midX {
		nextX := srcX + Signum(midX-srcX)
		if InRange(srcX, r.X, r.X+r.W) {
			if !InRange(nextX, r.X, r.X+r.W) {
				tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeDoor))
			}
		} else {
			tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeCorridor))
		}
		srcX = nextX
	}
	tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeCorridor))

	for srcY != dstY {
		srcY += Signum(dstY - srcY)
		tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeCorridor))
	}

	for srcX != dstX {
		srcX += Signum(dstX - srcX)
		if InRange(srcX, o.X, o.X+o.W) {
			tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeDoor))
			break
		} else {
			tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeCorridor))
		}
	}

	return tiles
}

func (r *room) Transpose() *room {
	return &room{r.Y, r.X, r.H, r.W, nil}
}

func (r *room) ConnectY(o *room, f MapGenInt) []*Tile {
	fTranspose := MapGenInt(func(o Offset, tiletype int) *Tile {
		return f(Offset{o.Y, o.X}, tiletype)
	})
	return r.Transpose().ConnectX(o.Transpose(), fTranspose)
}

func (r *room) CreateTiles(f MapGenInt) []*Tile {
	r.Tiles = createTileGrid(r.W-2, r.H-2, Offset{r.X + 1, r.Y + 1}, func(o Offset) *Tile {
		return f(o, TileTypeRoom)
	})
	return r.Tiles
}

func (r *room) ConnectDoor(door *Tile) {
	for _, tile := range r.Tiles {
		if step := door.Offset.Sub(tile.Offset); step.Chebyshev() == 1 {
			tile.Adjacent[step] = door
			door.Adjacent[step.Neg()] = tile
		}
	}
}

// Dungeon stub - will eventually generate room and corridor maps.
func Dungeon(numRooms, minRoomSize, maxRoomSize int, f MapGenInt) []*Tile {
	var tiles []*Tile

	maze := abstractBraid(numRooms, .25, 0, 1)
	rooms := make(map[*mazenode]*room)
	gridSize := maxRoomSize + minRoomSize

	// create rooms
	for _, nodes := range maze.Nodes {
		for _, node := range nodes {
			w := RandRange(minRoomSize, maxRoomSize)
			h := RandRange(minRoomSize, maxRoomSize)
			x := RandRange(gridSize*node.Pos.X, gridSize*(node.Pos.X+1)-w-1)
			y := RandRange(gridSize*node.Pos.Y, gridSize*(node.Pos.Y+1)-h-1)
			rooms[node] = &room{x, y, w, h, nil}
		}
	}

	// create room tiles
	for _, room := range rooms {
		tiles = append(tiles, room.CreateTiles(f)...)
	}

	// create corridors
	origin := maze.GetArbitraryNode()
	frontier := []*mazenode{origin}
	enqued := map[*mazenode]struct{}{origin: {}}
	closed := map[*mazenode]struct{}{}
	for len(frontier) != 0 {
		curr := frontier[0]
		frontier = frontier[1:]

		if _, done := closed[curr]; done {
			continue
		}
		closed[curr] = struct{}{}

		for step, adj := range curr.Edges {
			if _, done := closed[adj]; done {
				continue
			}
			if _, seen := enqued[adj]; !seen {
				frontier = append(frontier, adj)
			}

			var corridor []*Tile
			currRoom, adjRoom := rooms[curr], rooms[adj]
			if step.X != 0 {
				corridor = currRoom.ConnectX(adjRoom, f)
			} else {
				corridor = currRoom.ConnectY(adjRoom, f)
			}

			currRoom.ConnectDoor(corridor[0])
			adjRoom.ConnectDoor(corridor[len(corridor)-1])
			for i := 0; i < len(corridor)-1; i++ {
				j := i + 1
				step := corridor[j].Offset.Sub(corridor[i].Offset)
				corridor[i].Adjacent[step] = corridor[j]
				corridor[j].Adjacent[step.Neg()] = corridor[i]
			}

			tiles = append(tiles, corridor...)
		}
	}

	connectDiagonals(tiles[0])
	walls := addWalls(tiles[0], func(o Offset) *Tile {
		return f(o, TileTypeWall)
	})
	tiles = append(tiles, walls...)

	return tiles
}
