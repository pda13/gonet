package slices

import "fmt"

func PrettyPrint[T any](sl []T) {
	fmt.Println("[")
	for _, v := range sl {
		fmt.Printf("\t%v\n", v)
	}
	fmt.Println("]")
}
