package store

import (
	"github.com/nutsdb/nutsdb"
)

// AddToList 将键值对添加到指定的 bucket 中
func AddToList(bucket string, key []byte, val []byte) error {
	// 使用 SealWithProductKey 函数对值进行处理
	val, err := SealWithProductKey(val, nil)
	if err != nil {
		// 如果处理过程中出现错误，则返回相应错误
		return err
	}
	// 检查 bucket 是否有效
	err = checkBucket(bucket, nutsdb.DataStructureList)
	if err != nil {
		// 如果 bucket 无效，则返回相应错误
		return err
	}

	// 更新数据库，将 key 和处理后的值 val 添加到指定 bucket 中
	return DB.Update(
		func(tx *nutsdb.Tx) error {
			// 使用 LPush 命令将 key 和 val 添加到 bucket 中
			err := tx.LPush(bucket, key, val)
			return err
		},
	)
}

// GetList 根据指定的 bucket、key、page 和 size，从 NutsDB 中获取列表数据，进行解密封装后返回
func GetList(bucket string, key []byte, page int, size int) ([][]byte, error) {
	// 检查 bucket 是否有效，内容是否为列表
	err := checkBucket(bucket, nutsdb.DataStructureList)
	if err != nil {
		// 如果 bucket 无效或类型错误，返回错误
		return nil, err
	}

	// 初始化一个大小为 0，容量为 size 的 [][]byte 切片
	list := make([][]byte, 0, size)

	// 使用事务来读取数据库的内容
	err = DB.View(
		func(tx *nutsdb.Tx) error {
			// 当 page 大于 1，根据 page 和 size 自动计算列表数据的读取开始和结束索引
			var start = 0
			var end = size
			if page > 1 {
				start = (page - 1) * size
				end = start + size
			}

			// 使用 LRange 命令从指定的 bucket 中获取 key 对应的列表数据
			clist, err2 := tx.LRange(bucket, key, start, end)

			// 遍历读取到的列表数据，对每个元素进行解密封装后，添加到新的列表中
			for _, v := range clist {
				// 调用 Unseal 函数对值进行解密封装并添加到新的列表中
				item, err := Unseal(v, nil)
				if err != nil {
					//如果解密封装失败，返回错误
					return err
				}
				// 将解密封装后的数据添加到 list 切片中
				list = append(list, item)
			}
			// 返回读取过程中可能发生的错误
			return err2
		},
	)
	// 返回读取到的列表数据以及读取过程中可能发生的错误
	return list, err
}
