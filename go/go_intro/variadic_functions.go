package main
import "fmt"

// func main() {
// 		fullname := func(s1, s2 string)string{
// 			return fmt.Sprintf("%s %s", s1, s2)
// 		}
// 		WelcomeString := sayHello("code", "Chinweike", fullname)
// 		fmt.Println(WelcomeString)
// }

// func sayHello(first, last string, fn func(string, string)string) string{
// 	fullname := fn(first, last)
// 	return fmt.Sprintf("Welcome %s", fullname)
// }


/**
 * Variadic functions
*/
func main (){
	fmt.Println(sum(1, 2))        
	fmt.Println(sum(1, 2, 3, 4)) 
	fmt.Println(sum()) 

	greet("Hello", "Code", "Chinweike", "Akwolu")
}

func sum(numbers ...int) int {
	total := 0
	for _, n := range numbers{
		total += n
	}
	return total
}


func greet(prefix string, names ...string) {
	for _, name := range names {
		fmt.Println(prefix, name)
	}
}


