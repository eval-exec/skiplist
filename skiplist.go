package skiplist

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type node struct {
	key      int
	value    int
	forwards []*node
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newNode(key, value int, level int) *node {
	return &node{
		key:      key,
		value:    value,
		forwards: make([]*node, level+1),
	}
}

type SkipList struct {
	mu sync.RWMutex

	head      *node
	max_level int
}

var (
	MaxLevel int     = 64
	P        float64 = 0.5
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

func New() *SkipList {
	return &SkipList{
		head:      newNode(0, 0, MaxLevel+1),
		max_level: 0,
	}
}

func (sk *SkipList) Insert(key, value int) {
	sk.mu.Lock()
	defer sk.mu.Unlock()

	lvl := randomLevel()
	if lvl >= sk.max_level {
		sk.max_level = lvl + 1
	}
	var updates = make([]*node, lvl+1)
	now := sk.head
	for l := sk.max_level; l >= 1; l-- {
		for now.forwards[l] != nil && now.forwards[l].key < key {
			now = now.forwards[l]
		}
		if l <= lvl {
			updates[l] = now
		}
	}
	// [now] -> [next]
	next := now.forwards[1]

	if next != nil && next.key == key {
		next.value = value
	} else {
		newNode := newNode(key, value, lvl)
		for u := lvl; u >= 1; u-- {
			newNode.forwards[u] = updates[u].forwards[u]
			updates[u].forwards[u] = newNode
		}
	}
}

func (sk *SkipList) Update(key, value int) error {
	sk.mu.Lock()
	defer sk.mu.Unlock()

	now := sk.head.forwards[sk.max_level]
	for l := sk.max_level; l >= 1; l-- {
		for now.forwards != nil && now.forwards[l].key < key {
			now = now.forwards[l]
		}
	}
	// [now] -> [next]
	next := now.forwards[1]

	if next.key == key {
		next.value = value
		return nil
	}
	return ErrKeyNotFound
}

func (sk *SkipList) Search(key int) (value int, err error) {
	sk.mu.RLock()
	defer sk.mu.RUnlock()

	now := sk.head
	for l := sk.max_level; l >= 1; l-- {
		for now.forwards[l] != nil && now.forwards[l].key < key {
			now = now.forwards[l]
		}
	}
	// [now] -> [next]
	next := now.forwards[1]

	if next != nil && next.key == key {
		return next.value, nil
	}
	return 0, ErrKeyNotFound
}

func randomLevel() int {
	lvl := 1
	for lvl < MaxLevel && rand.Float64() < P {
		lvl += 1
	}
	return lvl
}

func print(sk *SkipList) {
	head := sk.head
	for head != nil {
		fmt.Printf("level: %d [%d,%d]\n", len(head.forwards), head.key, head.value)
		head = head.forwards[1]
	}
}
func printSimulate(sk *SkipList) {
	fmt.Println("print skip list :")
	defer fmt.Println("done")

	var matrix [][]string = make([][]string, sk.max_level+1)
	var length int
	{
		head := sk.head.forwards[1]
		for head != nil {
			length++
			head = head.forwards[1]
		}
		for i := range matrix {
			matrix[i] = make([]string, length+1)
		}
	}

	length = 0
	{
		head := sk.head
		for head != nil {
			for h := 0; h <= sk.max_level; h++ {
				var format = "%2d "
				var v string
				if h == 0 {
					v = fmt.Sprintf(format, head.key)
				} else if h < len(head.forwards) {
					v = " | "
				} else if h == len(head.forwards)+1 {
					v = fmt.Sprintf(format, head.value)
				} else if h == len(head.forwards) {
					v = " = "
				} else {
					v = "   "
				}
				matrix[h][length] = v
			}
			head = head.forwards[1]
			length++
		}
	}
	{
		length := sk.max_level + 1
		for i := 0; i < length/2; i++ {
			matrix[i], matrix[length-1-i] = matrix[length-1-i], matrix[i]
		}
	}
	for _, row := range matrix {
		for _, b := range row {
			fmt.Printf("%s", b)
		}
		fmt.Printf("\n")
	}

}
