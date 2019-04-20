package sparsebit;

import (
  "math"
  "math/rand"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestSparseBitShouldExplode1(t *testing.T) { testShouldExplode(t, 0) }

// func TestShouldExplode2(t *testing.T) { testShouldExplode(t, -1) }
func TestSparseBitShouldExplode3(t *testing.T) { testShouldExplode(t, 1<<32) }

func testShouldExplode(t *testing.T, cap int64) {
  sb, err := NewSparseBit(uint64(cap))
  assert.Error(t, err, "bad cap")
  assert.Nil(t, sb)
}

func testSparseBit(t *testing.T, size uint64) {
  sb, _ := NewSparseBit(size)
  for i := uint64(0); i < size; i++ {
    if rand.Intn(1000) < 100 {
      _ = sb.Set(i)
      isset, err := sb.Get(i)
      assert.Nil(t, err)
      assert.Equal(t, true, isset, i)
    } else {
      isset, err := sb.Get(i)
      assert.Nil(t, err)
      assert.Equal(t, false, isset, i)
    }
  }
}

func TestSparseBit_1K(t *testing.T) { testSparseBit(t, 1000) }
func TestSparseBit_1M(t *testing.T) { testSparseBit(t, 1000000) }

// func TestSparseBit_1B(t *testing.T)    { testSparseBit(t, 1000000000) }
func TestSparseBitUint16(t *testing.T) { testSparseBit(t, math.MaxUint16) }

func BenchmarkSparseSet(b *testing.B) {
  sb, _ := NewSparseBit(uint64(b.N))
  for i := 0; i < b.N; i++ {
    _ = sb.Set(uint64(i))
  }
}

func BenchmarkSparseGet(b *testing.B) {
  sb, _ := NewSparseBit(uint64(b.N))
  for i := 0; i < b.N; i++ {
    sb.Get(uint64(i))
  }
}

func BenchmarkGoMapSet(b *testing.B) {
  a := make([]int, b.N)
  for i := 0; i < b.N; i++ {
    a[i] = 0
  }
}

func BenchmarkGoMapGet(b *testing.B) {
  a := make([]int, b.N)
  for i := 0; i < b.N; i++ {
    _ = a[i]
  }
}
