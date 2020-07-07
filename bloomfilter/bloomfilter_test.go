package bloomfilter

import (
	"testing"
)

func TestBloomFilter_Exists(t *testing.T) {
	bf := New()
	bf.Push("a")
	bf.Push("b")
	bf.Push("c")
	bf.Push("d")
	bf.Push("e")
	bf.Push("f")
	bf.Push("g")
	bf.Push("h")

	t.Logf("a exists? %v\n", bf.Exists("a"))
	t.Logf("b exists? %v\n", bf.Exists("b"))
	t.Logf("c exists? %v\n", bf.Exists("c"))
	t.Logf("d exists? %v\n", bf.Exists("d"))
	t.Logf("e exists? %v\n", bf.Exists("e"))
	t.Logf("f exists? %v\n", bf.Exists("f"))
	t.Logf("g exists? %v\n", bf.Exists("g"))
	t.Logf("h exists? %v\n", bf.Exists("h"))
	t.Logf("i exists? %v\n", bf.Exists("i"))
	t.Logf("j exists? %v\n", bf.Exists("j"))
	t.Logf("k exists? %v\n", bf.Exists("k"))
	t.Logf("l exists? %v\n", bf.Exists("l"))
	t.Logf("m exists? %v\n", bf.Exists("m"))
	t.Logf("n exists? %v\n", bf.Exists("n"))
	t.Logf("o exists? %v\n", bf.Exists("o"))
	t.Logf("p exists? %v\n", bf.Exists("p"))
	t.Logf("q exists? %v\n", bf.Exists("q"))
	t.Logf("r exists? %v\n", bf.Exists("r"))
	t.Logf("s exists? %v\n", bf.Exists("s"))
	t.Logf("t exists? %v\n", bf.Exists("t"))
	t.Logf("u exists? %v\n", bf.Exists("u"))
	t.Logf("v exists? %v\n", bf.Exists("v"))
	t.Logf("w exists? %v\n", bf.Exists("w"))
	t.Logf("x exists? %v\n", bf.Exists("x"))
	t.Logf("y exists? %v\n", bf.Exists("y"))
	t.Logf("z exists? %v\n", bf.Exists("z"))
}

func TestBloomFilter_GetFalsePositivesRate(t *testing.T) {
	bf := New()
	t.Logf("FalsePositivesRate= %v\n", bf.GetFalsePositivesRate(100))
}
