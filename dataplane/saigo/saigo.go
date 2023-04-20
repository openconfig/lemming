package main

/*
#cgo CFLAGS: -I/usr/include/sai
#include <sai.h>
*/
import "C"
import "fmt"

//export go_sai_api_initialize
func go_sai_api_initialize(flags C.uint64_t, services *C.sai_service_method_table_t) C.sai_status_t {
	fmt.Println("hello go")
	return 0
}

func main() {}
