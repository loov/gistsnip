package testdata

import "fmt"

//gistsnip:start:example
func Example() {
	for i := 0; i < 10; i++ {
		//gistsnip:start:for
		fmt.Println(i)
		//gistsnip:end:for
	}
}
//gistsnip:end:example
