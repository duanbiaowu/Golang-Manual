// See the document
// https://zhuanlan.zhihu.com/p/398886752

package testify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Name string
	Age  int
}

type Crawler interface {
	GetUserList() ([]*User, error)
}

type MyCrawler struct {
	url string
}

type IExample interface {
	Hello(n int) int
}

type Example struct {
}

func (e *Example) Hello(n int) int {
	fmt.Printf("Hello with %d\n", n)
	return n
}

func ExampleFunc(e IExample) {
	for n := 1; n <= 3; n++ {
		for i := 0; i < n; i++ {
			e.Hello(n)
		}
	}
}

func (c *MyCrawler) GetUserList() ([]*User, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userList []*User
	err = json.Unmarshal(data, &userList)
	if err != nil {
		return nil, err
	}

	return userList, nil
}

func GetAndPrintUsers(crawler Crawler) {
	users, err := crawler.GetUserList()
	if err != nil {
		return
	}

	for _, u := range users {
		fmt.Println(u)
	}
}
