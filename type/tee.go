package types

type TeeParam struct {
	// sign address
	Address string
	// report time
	Time int64
	// 0: sgx, 1: sev 2: tdx 3: sev-snp
	TeeType uint8
	// report data
	Data []byte
	// report
	Report []byte
}
