package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"

	// 注册图片解码器
	_ "image/gif"

	"github.com/gin-gonic/gin"

	"life_journey/response"
)

// CompressImage 图片压缩工具
// 接收 multipart 文件上传，返回压缩后的图片文件
// 参数：
//   - file: 图片文件（必须）
//   - quality: JPEG 压缩质量 1-100，默认 75（可选）
//   - format: 输出格式 jpeg/png，默认 jpeg（可选）
func CompressImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "请上传图片文件")
		return
	}
	defer file.Close()

	// 读取压缩质量参数
	quality := 75
	if q := c.PostForm("quality"); q != "" {
		if v, err := strconv.Atoi(q); err == nil && v > 0 && v <= 100 {
			quality = v
		}
	}

	// 读取输出格式参数
	outputFormat := c.DefaultPostForm("format", "jpeg")

	// 读取原始图片数据
	imgData, err := io.ReadAll(file)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取图片失败")
		return
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无法解析图片格式，支持 JPEG/PNG/GIF")
		return
	}

	// 压缩编码
	var buf bytes.Buffer
	var contentType string

	switch outputFormat {
	case "png":
		encoder := &png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&buf, img)
		contentType = "image/png"
	default:
		// 默认输出为 JPEG
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		contentType = "image/jpeg"
		outputFormat = "jpeg"
	}

	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "压缩图片失败")
		return
	}

	// 设置响应头：文件下载 + 压缩前后大小信息
	filename := fmt.Sprintf("compressed_%s", header.Filename)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", contentType)
	c.Header("X-Original-Size", strconv.Itoa(len(imgData)))
	c.Header("X-Compressed-Size", strconv.Itoa(buf.Len()))
	c.Header("X-Compression-Ratio", fmt.Sprintf("%.1f%%", float64(buf.Len())/float64(len(imgData))*100))

	c.Data(http.StatusOK, contentType, buf.Bytes())
}
