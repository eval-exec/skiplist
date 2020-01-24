package skiplist

import (
	"math/rand"
	"testing"
	"time"
)

var (
	testSkipList *SkipList
)

func initTest() {
	testSkipList = New()
	rand.Seed(time.Now().UnixNano())
	var mm = make(map[int]int)
	for i := 1; i < 100; i++ {
		n := int(rand.Uint32() % 100)
		k := n
		v := int(rand.Uint32() % 100)
		mm[k] = v
		testSkipList.Insert(k, v)
	}
}

func validateOrder() bool {
	head := testSkipList.head.forwards[1]
	for head != nil {
		if head.forwards[1] != nil {
			if head.key > head.forwards[1].key {
				return false
			}
		}
		head = head.forwards[1]
	}
	return true
}

func TestOrder(t *testing.T) {
	initTest()
	if !validateOrder() {
		t.Fatal("validate order")
	}

}

