package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
)

// 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)

	return len(content), err
}

// 获取文件扩展名
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// 文件不存在返回true
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return errors.Is(err, fs.ErrNotExist)
}

// 检查文件目录的访问权限，没有权限则返回true
func CheckErrPermission(src string) bool {
	_, err := os.Stat(src)

	return errors.Is(err, fs.ErrPermission)
}

// 文件不存在则创建
func IsNotExistMkDir(src string) error {
	var err error = nil
	if notExist := CheckNotExist(src); notExist {
		err = MkDir(src)
	}

	return err
}

// 创建目录，权限为0777
func MkDir(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

// 打开文件，返回文件指针
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// 打开文件，没有则创建再打开，返回文件指针
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	noPerm := CheckErrPermission(src)
	if noPerm {
		return nil, fmt.Errorf("file.CheckErrPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
