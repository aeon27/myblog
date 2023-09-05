package article_service

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/aeon27/myblog/pkg/file"
	"github.com/aeon27/myblog/pkg/qrcode"
)

type ArticlePoster struct {
	*Article
	Name   string
	QRCode *qrcode.QRCode
}

func NewArticlePoster(posterName string, article *Article, q *qrcode.QRCode) *ArticlePoster {
	return &ArticlePoster{
		Article: article,
		Name:    posterName,
		QRCode:  q,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (ap *ArticlePoster) CheckNotExistMergedImg(path string) bool {
	if file.CheckNotExist(path + ap.Name) {
		return true
	}

	return false
}

func (ap *ArticlePoster) OpenMergedImg(path string) (*os.File, error) {
	f, err := file.MustOpen(ap.Name, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type ArticlePosterBG struct {
	Name string
	*ArticlePoster
	*Rect
	*Point
}

type Rect struct {
	Name           string
	X0, Y0, X1, Y1 int
}

type Point struct {
	X, Y int
}

func NewArticlePosterBG(name string, ap *ArticlePoster, rect *Rect, p *Point) *ArticlePosterBG {
	return &ArticlePosterBG{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Point:         p,
	}
}

func (apb *ArticlePosterBG) OpenBackGroundImg(path string) (*os.File, error) {
	f, err := file.MustOpen(apb.Name, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (apb *ArticlePosterBG) Generate() (string, string, error) {
	fullPath := qrcode.GetQRCodeFullPath()
	// 生成二维码
	fileName, path, err := apb.QRCode.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	if apb.CheckNotExistMergedImg(path) {
		// 打开背景图文件
		bgFile, err := apb.OpenBackGroundImg(path)
		if err != nil {
			return "", "", err
		}
		defer bgFile.Close()

		// 打开二维码文件
		qrFile, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrFile.Close()

		// 获取背景图
		bgImg, err := jpeg.Decode(bgFile)
		if err != nil {
			return "", "", err
		}

		// 获取二维码
		qrImg, err := jpeg.Decode(qrFile)
		if err != nil {
			return "", "", err
		}

		// 合并
		mergedImg := image.NewRGBA(image.Rect(apb.Rect.X0, apb.Rect.Y0, apb.Rect.X1, apb.Rect.Y1))
		draw.Draw(mergedImg, mergedImg.Bounds(), bgImg, bgImg.Bounds().Min, draw.Over)
		draw.Draw(mergedImg, mergedImg.Bounds(), qrImg, qrImg.Bounds().Min.Sub(image.Pt(apb.Point.X, apb.Point.Y)), draw.Over)

		// 打开合并文件
		mergedFile, err := apb.ArticlePoster.OpenMergedImg(path)
		if err != nil {
			return "", "", err
		}
		defer mergedFile.Close()

		err = jpeg.Encode(mergedFile, mergedImg, nil)
		if err != nil {
			return "", "", err
		}

		return apb.ArticlePoster.Name, path, nil
	}

	return "", "", fmt.Errorf("*ArticlePosterBG.Generate: The poster named %s have already existed.", apb.ArticlePoster.Name)
}
