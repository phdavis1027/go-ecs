package entity

type Tile struct {
	ID       Entity 
	Walkable bool
}

type TileMap struct {
	Width  int
	Height int
	Tiles  []Tile
}

func NewTileMap(width int, height int) *TileMap {
	tiles := make([]Tile, width * height)
	return &TileMap{Width: width, Height: height, Tiles: tiles}
}

func (tm *TileMap) GetTile(x int, y int) *Tile {
	r := y * tm.Width + x
	return &tm.Tiles[r]
}

func (tm *TileMap) SetTile(x int, y int, tile Tile) {
	r := y * tm.Width + x
	tm.Tiles[r] = tile
}

func (tm *TileMap) FillQuadBuf(quadBuf []float32) {
}
