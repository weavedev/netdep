//nolint
package main

func wrappedGetCall(url string) {
	location := url + "endpoint"
	//http.Get(location)
	print(location)
}

// target: GET example.com
func main() {
	wrappedGetCall("http://example.com/")
}
