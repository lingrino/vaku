package main

import (
	"fmt"
	"strings"
)

func search(d map[string]interface{}, s string) bool {
	for i, j := range d {
		if strings.Contains(i, s) {
			fmt.Printf("found in %s\n", i)
			return true
		}
		t, ok := j.(string)
		if ok {
			if strings.Contains(t, s) {
				fmt.Printf("found in %s\n", t)
				return true
			}
		}
		// } else {
		// 	for k, l := range j {
		// 		_, ok := j.(string)
		// 	}
		// }
		// fmt.Println(i)
		// fmt.Println(j)
	}
	return false
}

func main() {
	data := map[string]interface{}{
		"tls": "hello",
		"inner": []string{
			"ewtiuooox",
			"eptuwre",
			"ernqpokc",
		},
		"im": map[string]interface{}{
			"welkjfw": "weriouwe",
			"riewurw": []string{
				"weoijce",
				"ewroijcccx",
			},
			"weioasjc": "eriwouf",
		},
	}

	search(data, "ls")
}
