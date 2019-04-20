package sparsebit

import (
"fmt"
"math"
"math/bits"
"unsafe"
)

var (
  w = unsafe.Sizeof(uint64(0)) * 8 // Total width

  // Block len
  b3len = uint64(6)
  b2len = uint64(5)
  b1len = uint64(5)

  shift2 = b3len
  shift1 = b3len + b2len
  shift0 = b3len + b2len + b1len

  // Not really used, computed in the constructor
  // depending on required size
  // b0len = uint64(w) - shift0

  mask3 = (uint64(1) << b3len) - 1
  mask2 = (uint64(1) << b2len) - 1
  mask1 = (uint64(1) << b1len) - 1
)

type SparseBit struct {
  capacity uint64
  block0   [][][]uint64
}

func NewSparseBit(capacity uint64) (*SparseBit, error) {
  // Soft cap
  if math.MaxUint32 < capacity {
    return nil, fmt.Errorf("capacity (soft cap) < %d (actual: %d)", math.MaxUint32, capacity)
  }

  if capacity <= 0 {
    return nil, fmt.Errorf("capacity must be > 0 (actual %d)", capacity)
  }

  var block0Len int

  bitLen := int(w) - bits.LeadingZeros64(uint64(capacity))
  if (bitLen - int(shift0)) > 0 {
    block0Len = bitLen - int(shift0)
  }

  return &SparseBit{
    capacity: uint64(capacity),
    block0:   make([][][]uint64, 1<<uint(block0Len)),
  }, nil
}

func (sb *SparseBit) Set(value uint64) error {
  if value > sb.capacity {
    return fmt.Errorf("out of capacity (%d)", value)
  }

  if value < 0 {
    return fmt.Errorf("cannot be lesser than 0 (actual %d)", value)
  }

  b3value := uint64(1) << (value & mask3)
  sb.setb3(value, b3value)

  return nil
}

func (sb *SparseBit) setb3(value uint64, b3value uint64) {
  b0, b1, b2, _ := iblock(value)

  block1 := sb.block0[b0]
  if block1 == nil {
    block1 = make([][]uint64, 1<<b1len)
    sb.block0[b0] = block1
  }

  block2 := block1[b1]
  if block2 == nil {
    block2 = make([]uint64, 1<<b2len)
    block1[b1] = block2
  }

  block2[b2] |= b3value
}

func (sb *SparseBit) Get(value uint64) (bool, error) {
  if value > sb.capacity {
    return false, fmt.Errorf("out of capacity (%d)", value)
  }

  if value < 0 {
    return false, fmt.Errorf("cannot be lesser than 0 (actual %d)", value)
  }

  b3value := uint64(1) << (value & mask3)
  return (sb.getb3(value) & b3value) > 0, nil
}

func (sb *SparseBit) getb3(value uint64) uint64 {
  b0, b1, b2, _ := iblock(value)

  block1 := sb.block0[b0]
  if block1 == nil {
    return 0
  }

  block2 := block1[b1]
  if block2 == nil {
    return 0
  }

  return block2[b2]
}

// Block index
func iblock(value uint64) (uint64, uint64, uint64, uint64) {
  b1 := value >> shift0
  b2 := value >> shift1 & mask1
  b3 := value >> shift2 & mask2
  b4 := value & mask3

  return b1, b2, b3, b4
}
