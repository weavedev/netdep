//nolint
package main

import "fmt"

// target: GET example.com and POST example2.com
func main() {
	url := make([]string, 5)
	url[0] = "foo"
	url[1] = " "
	url[2] = "bar"
	print(url[0] + url[1] + url[2])
}

func fail2() {
	fmt.Sprintf("")
}
