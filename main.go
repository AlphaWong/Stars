package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var m map[string]int
	var t []map[string]interface{}
	response, err := http.Get("https://api.github.com/users/" + "alphawong" + "/starred?page=" + strconv.Itoa(0))
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		json.NewDecoder(bytes.NewReader(contents)).Decode(&t)
		for _, v := range t {
			lang := v["language"].(string)
			if _, ok := m[lang]; !ok {
				fmt.Print(ok)
				m[lang] = 0
			}
			m[lang] = m[lang] + 1
		}
		fmt.Print(m)
	}
}
