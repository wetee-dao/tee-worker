package main

import (
	"fmt"

	"github.com/edgelesssys/ego/enclave"
)

func main() {
	k, _, err := enclave.GetProductSealKey()
	if err != nil {
		fmt.Println("GetKey error", err)
	}
	fmt.Println("GetKey", k)
	k, _, err = enclave.GetUniqueSealKey()
	if err != nil {
		fmt.Println("GetKey error", err)
	}
	fmt.Println("GetKey", k)
	fmt.Println("xxxx")
	// k, _, err = enclave.
}
