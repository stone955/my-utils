package bloomfilter

import (
	"hash"
	"hash/fnv"
	"math"

	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	Bytes []byte
	Hash  []hash.Hash64
}

func New() *BloomFilter {
	bf := &BloomFilter{
		Bytes: make([]byte, math.MaxInt16),
		Hash: []hash.Hash64{
			fnv.New64(),
			murmur3.New64(),
		},
	}
	return bf
}

// Push 向 BloomFilter 添加元素
func (bf *BloomFilter) Push(str string) {
	for _, h := range bf.Hash {
		// 重置种子
		h.Reset()
		// 存入值
		_, _ = h.Write([]byte(str))
		// 返回 hash
		sum := h.Sum64()
		// 行/列
		x, y := sum%8, sum%uint64(len(bf.Bytes))
		// 这么写只适用于 Intel 小端模式，大端模式会有问题
		bf.Bytes[y] |= 1 << x
	}
}

// Exists 判定元素是否可能存在, 返回 false 一定不存在, 返回 true 可能存在
func (bf *BloomFilter) Exists(str string) bool {
	for _, h := range bf.Hash {
		// 重置种子
		h.Reset()
		// 存入值
		_, _ = h.Write([]byte(str))
		// 返回 hash
		sum := h.Sum64()
		// 行/列
		x, y := sum%8, sum%uint64(len(bf.Bytes))
		// 这么写只适用于 Intel 小端模式，大端模式会有问题
		if bf.Bytes[y]|1<<x != bf.Bytes[y] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) GetFalsePositivesRate(n int) float64 {
	a := 1 / float64(len(bf.Bytes))
	b := float64(len(bf.Hash)) * float64(n)
	c := float64(len(bf.Hash))
	d := 1 - math.Pow(1-a, b)
	return math.Pow(d, c)
}
