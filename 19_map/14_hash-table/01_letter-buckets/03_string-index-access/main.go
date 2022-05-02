package main

import "fmt"

func main() {
	word := "Hello"
	words_zh := "您好"
	letter := rune(word[0])
	letter2 := rune(words_zh[0])
	fmt.Println(letter, word[0], string(letter), string(word[1]), word, len(word))
	fmt.Println(letter2, words_zh[0], string(letter2), string(words_zh[1]), words_zh, len(words_zh))

}
