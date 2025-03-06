package os

import "os"

// exist0 判断路径是否存在. 存在返回 true, 否则返回 false
func exist0(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, f != nil && (err == nil || os.IsExist(err))
}

// Exist 判断路径是否存在. 存在返回 true, 否则返回 false
func Exist(path string) bool {
	_, exist := exist0(path)
	return exist
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	f, exist := exist0(path)
	return exist && f.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	f, exist := exist0(path)
	return exist && !f.IsDir()
}

// ReadFile 从磁盘读取文件, 失败抛出异常
func ReadFile(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}

// WriteFile 把数据写入指定路径的文件. 文件不存在自动创建, 存在则覆盖原有文件
func WriteFile(path, data string) {
	// create or override file
	destFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()
	// write data to dest file
	if _, err = destFile.WriteString(data); err != nil {
		panic(err)
	}
}
