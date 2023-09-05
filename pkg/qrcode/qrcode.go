package qrcode

import (
	"image/jpeg"

	"github.com/aeon27/myblog/pkg/file"
	"github.com/aeon27/myblog/pkg/setting"
	"github.com/aeon27/myblog/pkg/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QRCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const EXT_JPEG = ".jpeg"

func NewQRCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QRCode {
	return &QRCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPEG,
	}
}

func GetQRCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQRCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetQRCodePath()
}

func GetQRCodeFullURL(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQRCodePath() + name
}

func GetQRCodeMD5Name(value string) string {
	return util.EncodeMD5(value)
}

func (q *QRCode) GetQRCodeExt() string {
	return q.Ext
}

func (q *QRCode) CheckEncode(path string) bool {
	src := path + GetQRCodeMD5Name(q.URL) + q.GetQRCodeExt()
	if file.CheckNotExist(src) {
		return false
	}

	return true
}

// Encode
// 1、创建二维码生成路径
// 2、生成二维码
// 3、缩放二维码到指定大小
// 4、新建存放二维码图片的文件
// 5、将二维码图像以JPEG 4:2:0 的基线格式写入文件
func (q *QRCode) Encode(path string) (string, string, error) {
	name := GetQRCodeMD5Name(q.URL) + q.GetQRCodeExt()
	src := path + name
	if file.CheckNotExist(src) {
		bc, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		bc, err = barcode.Scale(bc, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path)
		if err != nil {
			return "", "", err
		}
		defer f.Close()

		err = jpeg.Encode(f, bc, nil)
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}
