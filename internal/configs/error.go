package configs

import "fmt"

func MakeError(err error, file string) {
	msg := "\n💀 Error : \n"
	msg += "file : ~ " + file
	msg += "\nreason : ~ " + err.Error()
	msg += "\n❌\n"
	fmt.Print(msg)

	//panic(err)

}
