package roaring

import (
	"testing"
)

func TestRoaringBitmap(t *testing.T) {
  r := NewRoaringBitset()

  for i := 0; i < 10_000; i++ {
    error := r.InsertOne(uint64(i))
    if error != nil {
      t.Fatalf("Error inserting %d: [%v]", i, error)
    }
  }


  for i := 0; i < 10_000; i++ {
    if !r.Has(uint64(i)) {
      t.Fatalf("Expected %d to be in the bitmap", i)
    }
  }

  for i := 90_000; i < 100_000; i++ {
    error := r.InsertOne(uint64(i))
    if error != nil {
      t.Fatalf("Error inserting %d: [%v]", i, error)
    }
  }

  for i := 90_000; i < 100_000; i++ {
    if !r.Has(uint64(i)) {
      t.Fatalf("Expected %d to be in the bitmap", i)
    }
  }

  for i := 10_000; i < 90_000; i++ {
    if r.Has(uint64(i)) {
      t.Fatalf("Expected %d to not be in the bitmap", i)
    }
  }
}

func BenchmarkRoaringBitmapInsert(b *testing.B) {
  r := NewRoaringBitset()

  for i := 0; i < b.N; i++ {
    r.InsertOne(uint64(i))
  }
}
