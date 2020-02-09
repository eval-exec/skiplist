package skiplist_test

import (
	"github.com/eval-exec/skiplist"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkInsertRandom(b *testing.B) {
	list := skiplist.New()
	_ = list
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		k := rand.Int() % 100
		v := rand.Int() % 100
		_, _ = k, v
		list.Insert(k, v)
	}
	b.ReportAllocs()
}

func BenchmarkInsertInOrder(b *testing.B) {
	list := skiplist.New()
	for i := 0; i < b.N; i++ {
		list.Insert(i, i)
	}
	b.ReportAllocs()
}

func BenchmarkParallelWrite(b *testing.B) {
	list := skiplist.New()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			list.Insert(rand.Int(), rand.Int()/2)
		}
	})
	b.ReportAllocs()
}
