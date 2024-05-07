package main

import (
	"fmt"
)

var p1 = fmt.Println // just an alias

/*
--> execution begins at main
--> in go , var x datatype = value is the syntax , if datatype is not specified then it is inferred from initialized value
--> variable with ':=' are accessible at function level, while variable with '=' are accessible at package level

*/

func main() {
	p1("Enter something in console!!")
	//reader := bufio.NewReader(os.Stdin) // to read keyboard input
	// var x int = 25 // new declaration so :=
	// var y int = 50
	// x = 10
	// fmt.Printf("x:%d  y:%d", x, y) // use printf for adding %d,%s and println for usual print statement

	// swap(&x, &y)
	// fmt.Println("After swapping X", x, " Y ", y)
	// fmt.Println(x)

	res := add() // Made reference of the function , since it is an

	fmt.Println("Sum of anonumous func", res(100, 1))
	// name, err := reader.ReadString('\n')
	// str := "Mayank Raj"

	// // Display the length of the string
	// fmt.Printf("Length of the string is:%d",
	// 	len(str))
	// if err == nil {
	// 	p1("Hello ", name)
	// } else {
	// 	log.Fatal(err)
	// }

}

// func swap(a *int, b *int) {
// 	var temp int
// 	temp = *a
// 	*a = *b
// 	*b = temp
// }

func add() func(x int, y int) int { // example of function taking anonymous function as an argument
	ans := func(x int, y int) int {
		return x + y
	}
	return ans
}

func swap(a, b *int) {
	var o int
	o = *a
	*a = *b
	*b = o

	//return o
}
