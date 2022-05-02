package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func main() {
	// 获取错误
	res, err := http.Get("http://www.gutenberg.org/cache/epub/1661/pg1661.txt")//"http://www-01.sil.org/linguistics/wordlists/english/wordlist/wordsEn.txt")
	if err != nil {
		log.Fatalln(err)
	}

	bs, _ := ioutil.ReadAll(res.Body)
	fmt.Println("cat type:", reflect.TypeOf(bs), bs)
	log.Writer()
	str := string(bs)
	fmt.Println(str)
	fmt.Println("cat type:", reflect.TypeOf(str))

}
