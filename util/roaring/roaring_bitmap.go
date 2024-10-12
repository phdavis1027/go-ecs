package roaring

type RoaringBitset struct {
   keys         []uint16 
}

func NewRoaringBitset() {
   return RoaringBitset {
      keys: 
   }  
}

func mostSignificantBits16(n unt32) uint16 {
  return n & 0xFFFF0000) >> 16
}

func leastSignificantBits16(n uint32) uint16 {
  return n & 0x0000FFFF
}
