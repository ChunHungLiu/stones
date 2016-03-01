package core

// MapGenInt generates Tiles for int values to form dungeon maps.
type MapGenInt func(o Offset, tiletype int) *Tile

const (
	TileTypeRoom = 1 << iota
	TileTypeWall
	TileTypeDoor
	TileTypeCorridor
)

type room struct {
	X, Y, W, H int
}

func (r room) ConnectX(o room, f MapGenInt) []*Tile {
	tiles := make([]*Tile, 0)

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
	if r.InBounds(midX, srcY) || o.InBounds(midX, dstY) {
		midX = (r.X + r.W + o.X) / 2
	}
	if r.InBounds(midX, srcY) || o.InBounds(midX, dstY) {
		midX = (r.X + o.X + o.W) / 2
	}

	for srcX != midX {
		srcX += Signum(midX - srcX)
		if !r.InBounds(srcX, srcY) && !o.InBounds(srcX, srcY) {
			tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeCorridor))
		}
	}
	for srcY != dstY {
		srcY += Signum(dstY - srcY)
		if !r.InBounds(srcX, srcY) && !o.InBounds(srcX, srcY) {
			tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeCorridor))
		}
	}
	for srcX != dstX {
		srcX += Signum(dstX - srcX)
		if !r.InBounds(srcX, srcY) && !o.InBounds(srcX, srcY) {
			tiles = append(tiles, f(Offset{srcX, srcY}, TileTypeCorridor))
		}
	}

	return tiles
}

func (r room) InBounds(x, y int) bool {
	return InBounds(x-r.X, y-r.Y, r.W, r.H)
}

// Dungeon stub - will eventually generate room and corridor maps.
func Dungeon(numRooms, minRoomSize, maxRoomSize int, f MapGenInt) []*Tile {
	// TODO Added in better maze gen customization
	maze := abstractBraid(numRooms, .25, 0, 1)
	rooms := make(map[*mazenode]room)
	tiles := make([]*Tile, 0)
	gridSize := maxRoomSize + minRoomSize

	// create rooms
	for _, nodes := range maze.Nodes {
		for _, node := range nodes {
			w := RandRange(minRoomSize, maxRoomSize)
			h := RandRange(minRoomSize, maxRoomSize)
			x := RandRange(gridSize*node.Pos.X, gridSize*(node.Pos.X+1)-w-1)
			y := RandRange(gridSize*node.Pos.Y, gridSize*(node.Pos.Y+1)-h-1)
			rooms[node] = room{x, y, w, h}
		}
	}

	// create room tiles
	for _, room := range rooms {
		for x := room.X; x < room.X+room.W; x++ {
			for y := room.Y; y < room.Y+room.H; y++ {
				tiles = append(tiles, f(Offset{x, y}, TileTypeRoom))
				// TODO connect room tiles
			}
		}
	}

	// create corridors
	origin := maze.GetArbitraryNode()
	if origin == nil {
		panic("WTF")
	}
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

			if step.X != 0 {
				tiles = append(tiles, rooms[curr].ConnectX(rooms[adj], f)...)
			} else {
				//tiles = append(tiles, rooms[curr].ConnectY(rooms[adj], f)...)
			}

			// TODO connect corridor tiles
		}
	}

	// TODO reuse maze connection stuff to fully connect dungeon and walls

	return tiles
}
