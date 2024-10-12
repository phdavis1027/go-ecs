package roaring

type container interface {
  IntoSortedArray() []uint16
  IntoBitmap()      [1024]uint64
}
