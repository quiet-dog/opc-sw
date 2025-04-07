package main

import "fmt"

func main() {

	a := []map[string]interface{}{}
	a = append(a, map[string]interface{}{
		"区域": "区域1",
	})
	a = append(a, map[string]interface{}{
		"区域": "区域2",
	})
	fmt.Println(a)
}
