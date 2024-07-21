package main

import (
	"crypto/sha256"
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

	hash := sha256.Sum256([]byte("hello world"))
	reportBytes, err := enclave.GetRemoteReport(hash[:])

	report, err := enclave.VerifyRemoteReport(reportBytes)
	fmt.Println(report.UniqueID)
	fmt.Println(report.ProductID)
	// k, _, err = enclave.
}
