package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
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

func TestMultiWriteSingleRead(t *testing.T) {

	nThreads := 10

	wg := sync.WaitGroup{}
	wg.Add(nThreads)
	skipList := New()
	for i := 0; i < nThreads; i++ {
		go func(i int) {
			for ins := i * 100; ins < (i+1)*100; ins++ {
				skipList.Insert(ins, ins*2)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := 0; i < nThreads*100; i++ {
		if v, err := skipList.Search(i); err != nil || v != i*2 {
			t.Fatalf("expect get %d on key:%d, but get %d , err: %s", i*2, i, v, err)
		}
	}

}
func TestMultiWriteMultiRead(t *testing.T) {
	nThreads := 10
	skipList := New()
	{
		wg := sync.WaitGroup{}
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			go func(i int) {
				for ins := i * 100; ins < (i+1)*100; ins++ {
					skipList.Insert(ins, ins*2)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()

	}
	{
		nThreads := 1000
		wg := sync.WaitGroup{}
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			go func() {
				defer wg.Done()
				key := rand.Int() % 2000
				v, err := skipList.Search(key)
				if key < 0 || key >= 1000 {
					if err != ErrKeyNotFound {
						t.Fatalf("expect get error %s, but no error occured", ErrKeyNotFound)
					}
				} else if v != key*2 || err != nil {
					t.Fatalf("expect get %d on key:%d, but get %d , err: %s", key*2, key, v, err)
				}
			}()
		}
		wg.Wait()
	}
}
