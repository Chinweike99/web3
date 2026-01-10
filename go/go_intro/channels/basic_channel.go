package main
import "fmt"

func main_() {
	ch := make(chan string)

	go func(){
		ch <- "hello from go"
	}()

	msg := <-ch
	fmt.Println(msg)

}