package main

import "fmt"

func main() {
	//// 字符串转换为列表处理
	src_Words := []rune("Music should strike fire from the heart of man, and bring tears from the eyes of woman. * Music is the mediator between the spiritual and the sensual life. * Music is the one incorporeal entrance into the higher world of knowledge which comprehends mankind but which mankind cannot comprehend. *  , you men who think or say that I am malevolent, stubborn or misanthropic, how greatly do you wrong me. You do not know the secret cause which makes me seem that way to you, and I would have ended my life - it was only my art that held me back. Ah, it seemed impossible to leave the world until I had brought forth all that I felt was within me. * Music from my fourth year began to be the first of my youthful occupations. Thus early acquainted with the gracious muse who tuned my soul to pure harmonies, I became fond of her, and, as it often seemed to me, she of me. * A great poet is the most precious jewel of a nation. * Music comes to me more readily than words. * The barriers are not erected which can say to aspiring talents and industry, 'Thus far and no farther.' * I only live in my music, and I have scarcely begun one thing when I start on another. As I am now working, I am often engaged on three or four things at the same time. * There ought to be an artistic depot where the artist need only hand in his artwork in order to receive what he asks for. As things are, one must be half a business man, and how can one understand - good heavens! - that's what I really call troublesome. *")

	a_words := ""
	for i := 0; i < len(src_Words); i++ {
		if string(src_Words[i]) == "*" {
			fmt.Println(a_words)
			fmt.Println("++++++++++++++++++++++++++++++++++++++")
			fmt.Println()
			a_words = ""
		} else {
			a_words = a_words + string(src_Words[i])
		}

	}
}
