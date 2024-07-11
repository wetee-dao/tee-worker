package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"

	"wetee.app/worker/cmd/dkg/dkg"
	"wetee.app/worker/cmd/dkg/p2p"
)

func main() {
	// 初始化加密套件。
	suite := suites.MustFind("Ed25519")

	// 初始化参与者列表。
	participants := []kyber.Point{
		// 生成随机公钥。
		suite.Point().Mul(suite.Scalar().Pick(suite.RandomStream()), nil),
		// ...
	}

	// 获取阈值参数。
	// TODO: 从外部获取阈值参数。
	threshold, err := strconv.Atoi(os.Getenv("THRESHOLD"))
	if err != nil {
		fmt.Println("获取阈值参数失败:", err)
		os.Exit(1)
	}

	// 创建 DKG 实例。
	dkg, err := dkg.NewRabinDKG(suite, participants, threshold)
	if err != nil {
		fmt.Println("创建 DKG 实例失败:", err)
		os.Exit(1)
	}

	// 启动 P2P 网络。
	host, err := p2p.NewP2PNetwork(context.Background())
	if err != nil {
		fmt.Println("启动 P2P 网络失败:", err)
		os.Exit(1)
	}
	dkg.Host = host
	dkg.NodeID = host.ID()

	// 运行 DKG 协议。
	if err := dkg.Run(); err != nil {
		fmt.Println("运行 DKG 协议失败:", err)
		os.Exit(1)
	}

	// 获取密钥文件路径。
	// TODO: 从外部获取密钥文件路径。
	// keyFilePath := os.Getenv("KEY_FILE_PATH")

	// 保存密钥到文件。
	// ...
}

// import (
// 	"fmt"

// 	"github.com/consensys/gnark/frontend"

// 	"go.dedis.ch/kyber/v3"
// 	"go.dedis.ch/kyber/v3/share"
// 	dkg "go.dedis.ch/kyber/v3/share/dkg/pedersen"
// 	"go.dedis.ch/kyber/v3/suites"
// )

// // 门限方案参数
// type ThresholdScheme struct {
// 	Threshold int          // 恢复密钥所需的最小碎片数量
// 	Total     int          // 碎片总数
// 	Suite     suites.Suite // 密码学套件
// }

// // 生成密钥碎片
// func GenerateSecretShares(scheme *ThresholdScheme, secret kyber.Scalar) []*share.PriShare {
// 	pri := share.NewPriPoly(scheme.Suite, scheme.Threshold, secret, scheme.Suite.RandomStream())

// 	return pri.Shares(scheme.Total)
// }

// // 恢复密钥
// func RecoverSecret(scheme *ThresholdScheme, shares []*share.PriShare) (kyber.Scalar, error) {
// 	return share.RecoverSecret(scheme.Suite, shares, scheme.Threshold, scheme.Total)
// }

// // Gnark 算术电路
// type ThresholdCircuit struct {
// 	SecretShare share.PriShare
// 	Secret      kyber.Scalar
// }

// func (circuit *ThresholdCircuit) Define(curve frontend.API) error {
// 	// 将密钥碎片和密钥恢复过程定义为算术电路
// 	// ...

// 	return nil
// }

// type DkgNode struct {
// 	dkg         *dkg.DistKeyGenerator
// 	pubKey      kyber.Point
// 	privKey     kyber.Scalar
// 	deals       []*dkg.Deal
// 	resps       []*dkg.Response
// 	secretShare *share.PriShare
// }

// func main() {
// 	// 初始化密码学套件
// 	suite := suites.MustFind("Ed25519")

// 	// 设置门限方案参数
// 	// scheme := &ThresholdScheme{
// 	// 	Threshold: 3, // 至少需要 3 个碎片才能恢复密钥
// 	// 	Total:     5, // 总共 5 个碎片
// 	// 	Suite:     suite,
// 	// }

// 	// 生成密钥
// 	privKey := suite.Scalar().Pick(suite.RandomStream())
// 	// 生成公钥
// 	pubKey := suite.Point().Mul(privKey, nil)
// 	fmt.Printf("原始密钥: %x\n", privKey.String())

// 	// 创建DKG结点
// 	node := &DkgNode{
// 		pubKey:  pubKey,
// 		privKey: privKey,
// 		deals:   make([]*dkg.Deal, 0),
// 		resps:   make([]*dkg.Response, 0),
// 	}
// 	d, err := dkg.NewDistKeyGenerator(suite, privKey, []kyber.Point{pubKey}, 3)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	node.dkg = d

// 	// // 生成密钥碎片
// 	// shares := GenerateSecretShares(scheme, privKey)

// 	// // 打印密钥碎片
// 	// fmt.Println("密钥碎片:")
// 	// for _, share := range shares {
// 	// 	fmt.Printf("碎片 %d: %x\n", share.I, share.V.String())
// 	// }

// 	// // 要加密的数据
// 	// data := []byte("这是一个秘密信息")

// 	// // 加密数据
// 	// ciphertext, err := ecies.Encrypt(suite, pubKey, data, suite.Hash)
// 	// if err != nil {
// 	// 	fmt.Println("加密失败:", err)
// 	// 	return
// 	// }

// 	// // 解密数据
// 	// decryptedData, err := ecies.Decrypt(suite, privKey, ciphertext, suite.Hash)
// 	// if err != nil {
// 	// 	fmt.Println("解密失败:", err)
// 	// 	return
// 	// }

// 	// if string(decryptedData) != string(data) {
// 	// 	fmt.Println("解密数据与原始数据不一致")
// 	// 	return
// 	// }

// 	// // 生成 zk-SNARK 证明
// 	// // ...

// 	// // 恢复密钥
// 	// recoveredSecret, err := RecoverSecret(scheme, shares[0:3])
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return
// 	// }
// 	// fmt.Printf("恢复的密钥: %x\n", recoveredSecret.String())

// 	// // 验证密钥是否恢复成功
// 	// if recoveredSecret.Equal(privKey) {
// 	// 	fmt.Println("密钥恢复成功")
// 	// } else {
// 	// 	fmt.Println("密钥恢复失败")
// 	// }

// 	// 在 Substrate 中验证证明
// }
