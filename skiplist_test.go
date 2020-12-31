package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

var (
	testSkipList *SkipList
)

func validateProb() {
	head := testSkipList.head.forwards[1]
	mm := make(map[int]int)
	for head != nil {
		if head.forwards[1] != nil {
			mm[len(head.forwards)] += 1
		}
		head = head.forwards[1]
	}
	lvl := MaxLevel
	for lvl >= 1 {
		n := mm[lvl]
		for n > 0 {
			fmt.Printf("=")
			n--
		}
		fmt.Println()
		lvl--
	}
}

func validateOrder() bool {
	head := testSkipList.head.forwards[1]
	for head != nil {
		if head.forwards[1] != nil && head.key >= head.forwards[1].key {
			return false
		}
		head = head.forwards[1]
	}
	return true
}

func TestOrder(t *testing.T) {
	testSkipList = New()
	rand.Seed(time.Now().UnixNano())
	var mm = make(map[int]int)
	for i := 1; i < 1e6; i++ {
		n := int(rand.Uint32() % 100)
		k := n
		v := int(rand.Uint32() % 100)
		mm[k] = v
		testSkipList.Insert(k, v)
	}
	if !validateOrder() {
		t.Fatal("validate order")
	}
	for k, v := range mm {
		if got, err := testSkipList.Search(k); err != nil {
			t.FailNow()
		} else if got != v {
			t.FailNow()
		}
	}
	validateProb()
}

func TestLevel(t *testing.T) {
	N := 1e7
	mm := make(map[int]int)
	for N > 0 {
		mm[randomLevel()]++
		N--
	}
	var GET100 bool
	for i := 2; i <= MaxLevel; i++ {
		if mm[i] < 100 {
			break
		}
		GET100 = true
		if math.Abs(float64(mm[i-1])/float64(mm[i])-2) > 0.5 {
			t.Fatalf("[%d: %d], [%d: %d] \n", i, mm[i], i-1, mm[i-1])
		}
	}
	if !GET100 {
		t.Fatalf("sample count less than 100")
	}
}
