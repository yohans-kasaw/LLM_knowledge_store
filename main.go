package main
import(
	"fmt"
	"log"
	"starter/greeting"
)

func main(){
	log.SetPrefix("greeting => ")
	log.SetFlags(0)


	// name := "jo"
	name := ""
	msg, err := greeting.Say_hello(name)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(msg)
}
