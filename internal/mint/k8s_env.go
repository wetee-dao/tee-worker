package mint

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	contracts "github.com/wetee-dao/tee-dsecret/pkg/chains/revive"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type PodData struct {
	InitDatas map[string]string
	Files     map[int]map[string]string
	Encrypts  map[int]map[string]uint64
	Disks     map[int]map[string]uint64
}

// Build Envs
// 获取配置文件
func (m *Minter) BuildEnvsFromSettings(podId uint64, nameSpace string, index int, settingEnvs []model.Env, disks []model.ContainerDisk, init *PodData) ([]corev1.EnvVar, error) {
	envs := []corev1.EnvVar{}

	if init.Files[index] == nil {
		init.Files[index] = make(map[string]string)
	}
	if init.Encrypts[index] == nil {
		init.Encrypts[index] = make(map[string]uint64)
	}
	if init.Disks[index] == nil {
		init.Disks[index] = make(map[string]uint64)
	}

	for _, setting := range settingEnvs {
		if setting.Env != nil {
			envs = append(envs, corev1.EnvVar{
				Name:  string(setting.Env.F0),
				Value: string(setting.Env.F1),
			})
		} else if setting.File != nil {
			init.Files[index][string(setting.File.F0)] = hex.EncodeToString(setting.File.F1)
		} else if setting.Encrypt != nil {
			init.Encrypts[index][string(setting.Encrypt.F0)] = setting.Encrypt.F1
		}
	}

	for _, disk := range disks {
		init.Disks[index][string(disk.Path)] = uint64(disk.Id)
	}

	return envs, nil
}

func (m *Minter) WrapDeploymentInitData(deployment *appsv1.Deployment, version model.TEEType, initData *PodData) error {
	encrypts, _ := json.Marshal(initData.Encrypts)
	files, _ := json.Marshal(initData.Files)
	disks, _ := json.Marshal(initData.Disks)

	params := map[string]string{}
	params["polkadot_cloud_addr"] = contracts.GetCloudAddress()
	paramStr, _ := json.Marshal(params)

	if version.SGX != nil {
		for k, v := range initData.InitDatas {
			deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
				Name:  k,
				Value: v,
			})
		}

		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "__ENCRYPTS__",
			Value: string(encrypts),
		})

		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "__FILES__",
			Value: string(files),
		})

		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "__DISKS__",
			Value: string(disks),
		})

		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "__PARAMS__",
			Value: string(paramStr),
		})
	} else if version.CVM != nil {
		initToml := CVMInitData{
			Version:   "0.1.0",
			Algorithm: "sha256",
			Data:      map[string]string{},
		}

		fmt.Println("11111111111111111")

		maps.Copy(initToml.Data, initData.InitDatas)
		initToml.Data["__ENCRYPTS__"] = string(encrypts)
		initToml.Data["__FILES__"] = string(files)
		initToml.Data["__DISKS__"] = string(disks)

		// TODO
		initToml.Data["__PARAMS__"] = string(paramStr)

		deployment.Spec.Template.ObjectMeta.Annotations["io.katacontainers.config.runtime.cc_init_data"] = initToml.base64()
	}

	return nil
}

func (m *Minter) WrapPodInitData(pod *corev1.Pod, version model.TEEType, initData *PodData) {
	for k, v := range initData.InitDatas {
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	chain_addr := strings.Join(chains.MainChain.GetChainUrls(), ",")
	pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
		Name:  "CHAIN_ADDR",
		Value: chain_addr,
	})

	encrypts, _ := json.Marshal(initData.Encrypts)
	pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
		Name:  "__ENCRYPTS__",
		Value: string(encrypts),
	})

	files, _ := json.Marshal(initData.Files)
	pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
		Name:  "__FILES__",
		Value: string(files),
	})
}

// WrapNodeService
// 包装环境变量
func (m *Minter) WrapEnvs(envs []corev1.EnvVar, nameSpace, name string, nodeSers *corev1.Service) error {
	mdata := make(map[string]string)
	mdata["cluster_domain"] = m.HostDomain
	mdata["project_domain"] = nameSpace + ".svc.cluster.local"
	// mdata["server_ssl"] = strings.Join(util.GetSslRoot(), "|")
	for i, port := range nodeSers.Spec.Ports {
		if port.NodePort != 0 {
			mdata["ser_"+fmt.Sprint(i)+"_nodeport"] = fmt.Sprint(port.NodePort)
		}
	}

	for i, env := range envs {
		if strings.Contains(env.Value, "{{.") {
			v, err := renderTemplate(env.Value, mdata)
			if err != nil {
				return err
			}
			envs[i].Value = v
		}
	}

	return nil
}

// cvm initdata
type CVMInitData struct {
	Version   string            `toml:"version"`   // 对应 version = "0.1.0"
	Algorithm string            `toml:"algorithm"` // 对应 algorithm = "sha256"
	Data      map[string]string `toml:"data"`      // 对应 [data] 区块下的键值对
}

func (d CVMInitData) base64() string {
	yamlContent, err := toml.Marshal(d)
	if err != nil {
		panic(err)
	}

	// 2. Gzip 压缩（对应 gzip 命令）
	var gzipBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzipBuffer)
	// 将文件内容写入 gzip 压缩流
	if _, err := gzipWriter.Write(yamlContent); err != nil {
		fmt.Printf("Gzip 压缩失败: %v\n", err)
		panic(err)
	}
	// 关闭 gzip 写入器（必须调用，否则压缩数据不完整）
	if err := gzipWriter.Close(); err != nil {
		fmt.Printf("关闭 Gzip 写入器失败: %v\n", err)
		panic(err)
	}

	// 3. Base64 编码（无换行，对应 base64 -w0）
	// StdEncoding.EncodeToString 默认不换行，与 -w0 效果一致
	return base64.StdEncoding.EncodeToString(gzipBuffer.Bytes())
}
