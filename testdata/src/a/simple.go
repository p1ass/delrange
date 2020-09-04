package a

import (
	"fmt"
)

func whenKeyIsBasicLit() {
	m := map[int]int{1: 1, 2: 2}
	for key, value := range m {
		delete(m, 1) // want "function is called with a value different from range key"
		delete(m, key)
		fmt.Println(key, value)
	}
}

func whenUsingKey() {
	m := map[int]int{1: 1, 2: 2}
	for key, value := range m {
		delete(m, key)
		fmt.Println(key, value)
	}
}

func whenCopyKey() {
	m := map[int]int{1: 1, 2: 2}
	for key, value := range m {
		newKey := key
		delete(m, newKey) // want "function is called with a value different from range key"
		fmt.Println(key, value)
	}
}

func whenAssignBasicLit() {
	m := map[int]int{1: 1, 2: 2}
	for key, value := range m {
		key = 1
		delete(m, key)
		fmt.Println(key, value)
	}
}
