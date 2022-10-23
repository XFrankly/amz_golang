package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Tasks struct {
	ClientParam *src.RequestsHeaders
	ClientConn  *src.HttpClientHandler
}

type AuthError error

func AuthErrors(msg string) *AuthError {
	var defaultErr AuthError = errors.New(fmt.Sprintf("Auth failed: %v", msg))
	return &defaultErr
}

type RestError error

func RestErrors(msg string) *RestError {
	var defaultErr RestError = errors.New(fmt.Sprintf("RestApi failed: %v", msg))
	return &defaultErr
}

var (
	Logg = log.New(os.Stderr, "INFO -:", 18)
)

type RestConnection struct {
	host       string
	ssl        bool
	token      string
	token_path string
}

func MakeRestConnection(host string, token string, ssl bool, path string) *RestConnection {
	return &RestConnection{host: host, token: token, ssl: ssl, token_path: path}
}

func (rc *RestConnection) GetConnection() string {
	if rc.ssl == true {
		return fmt.Sprintf("https://%v", rc.host)
	} else {
		return fmt.Sprintf("http://%v", rc.host)
	}

}

func (rc *RestConnection) Auth(user string, passwd string) *error {
	hpath := rc.GetConnection()
	payload := map[string]string{
		"type":     "login",
		"username": user,
		"password": passwd,
	}
	bodys := strings.NewReader(payload)
	rst := http.Post(h1+rc.token_path, "application/json", bodys)
}

func main() {
	Logg.Println(*AuthErrors("test error."))
	Logg.Println(*RestErrors("test rest error."))

}
