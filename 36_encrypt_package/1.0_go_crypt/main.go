package main

import (
	"fmt"

	"github.com/xxtea/xxtea-go/xxtea"
)

func main() {
	str := "Hello word. Ph EU ch"
	key := "1234567890"
	encrypt_data := xxtea.Encrypt([]byte(str), []byte(key))
	decrypt_data := string(xxtea.Decrypt(encrypt_data, []byte(key)))
	fmt.Println("encrypt_data_string:", encrypt_data)

	if str == decrypt_data {
		fmt.Println("success", decrypt_data)
	} else {
		fmt.Println("faild.", decrypt_data)
	}

	encrypt_data_string := xxtea.EncryptString(str, key)
	decrypt_string, err := xxtea.DecryptString(encrypt_data_string, key)
	fmt.Println("encrypt_data_string:", encrypt_data_string)
	if err != nil {
		fmt.Println("failed,", err)
	} else {
		fmt.Println("success:", decrypt_string)
	}
}
