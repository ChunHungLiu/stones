package main

import (
	"github.com/jlund3/stones/core"
)

func main() {
	core.MustTermInit()
	defer core.TermDone()

	cols, rows := 20, 10
	tiles := core.GenStub(cols, rows)

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			core.TermDraw(x, y, tiles[x][y].Face)
		}
	}
	core.TermRefresh()

	core.GetKey()
}
