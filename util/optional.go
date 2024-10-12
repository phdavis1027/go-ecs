package util

type Optional[T any] struct {
  Inner       T
  IsSome      bool
}
