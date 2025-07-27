package greeting

import ("fmt" 
	"errors")


func Say_hello(name string)(string, error){
	if name == ""{
		return "", errors.New("empty name")
	}

	msg := fmt.Sprintf("Hello %v, wellcome!", name)
	return msg, nil
}
