package main

import "fmt"

func main() {
	n := hashBucket("Go", 12)
	fmt.Println(n)
}

func hashBucket(word string, buckets int) int {
	letter := int(word[0]) + int(word[1])
	fmt.Println("letter % buckets", letter, buckets)
	bucket := letter % buckets
	return bucket
}
