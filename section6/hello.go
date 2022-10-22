package section6

import "fmt"

func Hello() {
	fmt.Println("Hello")
}

func Unordered() {
	for _, v := range []int{1, 2, 3} {
		fmt.Println(v)
	}
}

func ShuffleWillBeFailed() {
	x := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range x {
		fmt.Printf("k=%s v=%d\n", k, v)
	}
}
