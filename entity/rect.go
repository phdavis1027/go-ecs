package entity

import "github.com/phdavis1027/goecs/util"

type Rect struct {
  corner  util.Vec2 
  width   float64
  height  float64
}

func NewRect(corner util.Vec2, width, height float64) Rect {
  return Rect{corner: corner, width: width, height: height}
}
