package handler

import (
	"archive/zip"
	"bytes"
	"io"
	"logger"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetFidList(ctx *gin.Context) {
	fidList := ctx.QueryArray("fid")
	if len(fidList) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		logger.Logger.Println("Invalid request: fid is required")
		ctx.Abort()
		return
	}

	ctx.Set("fidList", fidList)
	ctx.Next()
}

func VerifyFidList(ctx *gin.Context) {
	fidList := ctx.MustGet("fidList").([]string)

	nonExtFidList := make([][]string, 0)
	// 验证fid是否符合要求：fid只允许包含一个点分隔符
	for _, fid := range fidList {
		nonExtFid := strings.Split(fid, ".")
		if len(nonExtFid) != 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			logger.Logger.Println("Invalid request: fid allows only one point. where=" + fid)
			ctx.Abort()
			return
		}
	}

	// 检查是否存在对应文件
	for _, nonExtFid := range nonExtFidList {
		if ok := utils.CheckSum256InSum256Map(nonExtFid[0]); !ok {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "File not found",
			})
			logger.Logger.Println("File not found. where=" + nonExtFid[0])
			ctx.Abort()
			return
		}
	}

	ctx.Set("nonExtFidList", nonExtFidList)
	ctx.Next()
}

// Invoke-WebRequest http://localhost:8080/api/v1/static/delete?fid -Method DELETE
func RemoveImages(ctx *gin.Context) {
	fidList := ctx.MustGet("fidList").([]string)
	nonExtFidList := ctx.MustGet("nonExtFidList").([][]string)

	for _, nonExtFid := range nonExtFidList {
		utils.RemoveSum256Map(nonExtFid[0])
	}
	utils.RemoveImage(fidList...)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Files deleted successfully",
	})
	logger.Logger.Println("Files deleted successfully. filenames: " + strings.Join(fidList, ", "))
	ctx.Next()
}

// Invoke-WebRequest http://localhost:8080/api/v1/static/:controller/get?fid
func GetImages(ctx *gin.Context) {
	fidList := ctx.MustGet("fidList").([]string)
	controller := ctx.Param("controller")
	var prefix string

	switch controller {
	case "original":
		prefix = viper.GetString("originalPrefix")
	case "thumbnail":
		prefix = viper.GetString("thumbnailPrefix")
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		logger.Logger.Println("Invalid request: controller must be either 'original' or 'thumbnail'")
		ctx.Abort()
		return
	}

	errStrRet := ""
	errStrLog := "Some files not found. {\n"
	baseImageRet := make([]string, 0)
	for _, fid := range fidList {
		if baseImage, ok := utils.GetImageUrlBase64(prefix + fid); !ok {
			errStrRet = "Some files not found"
			errStrLog += "\twhere=" + fid + "\n"
		} else {
			baseImageRet = append(baseImageRet, baseImage)
		}
	}
	errStrLog += "}"

	if errStrRet != "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": errStrRet,
			"data": gin.H{
				controller: baseImageRet,
			},
		})
		logger.Logger.Println(errStrLog)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get files successfully",
		"data": gin.H{
			controller: baseImageRet,
		},
	})
	ctx.Next()
}

// 将图片上传到服务器，并保存为original和thumbnail两张
func AddImages(ctx *gin.Context) {
	ctx.Request.ParseForm()

	uploadFile, handle, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		logger.Logger.Println("Invalid request: image is required")
		ctx.Abort()
		return
	}

	// 将图片内容缓存到[]byte
	fileBytesBuffer := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(fileBytesBuffer, uploadFile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		logger.Logger.Println("Internal server error: Failed to copy image data")
		logger.Logger.Println(err)
		ctx.Abort()
		return
	}
	fileBytes := fileBytesBuffer.Bytes()
	uploadFile.Close()

	// 计算图片的sha256校验和
	sum256 := utils.GetSHA256(fileBytes)

	// 检查sum256是否重复
	if ok := utils.CheckSum256InSum256Map(sum256); ok {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "File already exists",
		})
		logger.Logger.Println("File already exists. where=" + handle.Filename)
		ctx.Abort()
		return
	}

	tmp := strings.Split(handle.Filename, ".")
	ext := tmp[len(tmp)-1]

	// 保存original和thumbnail
	if ok := utils.ImageCompress(ext, fileBytes, sum256); !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		logger.Logger.Println("Internal server error: Failed to compress image")
		ctx.Abort()
		return
	}

	// 将sum256写入sum256Map
	utils.AddSum256Map(sum256)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
	})
	ctx.Next()
}

func GetFidListPagingQuery(ctx *gin.Context) {
	pageNumStr := ctx.Query("page")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		logger.Logger.Println("Invalid request: pageNum is required and must be a number")
		ctx.Abort()
		return
	}
	pageSize := viper.GetInt("indexPageSize")

	fidList := utils.GetSum256ListPagingQuery(pageNum, pageSize)
	totalFidList := utils.GetTotalPage(pageSize)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get fid list successfully",
		"data": gin.H{
			"fidList": fidList,
			"total":   totalFidList,
		},
	})
	ctx.Next()
}

func AddZippedImages(ctx *gin.Context) {
	uploadFileHandler, err := ctx.FormFile("imageZip")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		logger.Logger.Println("Invalid request: imageZip is required")
		ctx.Abort()
		return
	}

	uploadFile, err := uploadFileHandler.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		logger.Logger.Println("Internal server error: Failed to open imageZip file")
		ctx.Abort()
		return
	}

	tmpFilename := "./tmpZip" + strconv.Itoa(time.Now().Nanosecond()) + ".zip"
	zipfile, _ := os.Create(tmpFilename)
	io.Copy(zipfile, uploadFile)
	zipfile.Close()
	uploadFile.Close()

	zipReader, err := zip.OpenReader(tmpFilename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		logger.Logger.Println("Internal server error: Failed to open zip file")
		ctx.Abort()
		return
	}
	defer zipReader.Close()

	ch := make(chan bool, viper.GetInt("maxUnzipThreadNum"))
	wg := &sync.WaitGroup{}

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		ch <- true
		wg.Add(1)
		go func(file *zip.File) {
			// 将图片内容缓存到[]byte
			fileReader, err := file.Open()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				logger.Logger.Println("Internal server error: Failed to open image file")
				<-ch
				wg.Done()
				return
			}
			fileBytesBuffer := bytes.NewBuffer([]byte{})
			if _, err := io.Copy(fileBytesBuffer, fileReader); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				logger.Logger.Println("Internal server error: Failed to copy image data")
				<-ch
				wg.Done()
				return
			}
			fileBytes := fileBytesBuffer.Bytes()
			fileReader.Close()

			// 计算图片的sha256校验和
			sum256 := utils.GetSHA256(fileBytes)

			// 检查sum256是否重复
			if ok := utils.CheckSum256InSum256Map(sum256); ok {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "File already exists",
				})
				logger.Logger.Println("File already exists. where=" + file.Name)
				<-ch
				wg.Done()
				return
			}

			tmp := strings.Split(file.Name, ".")
			ext := strings.ToLower(tmp[len(tmp)-1])

			// 保存original和thumbnail
			if ok := utils.ImageCompress(ext, fileBytes, sum256); !ok {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				logger.Logger.Println("Internal server error: Failed to compress image")
				<-ch
				wg.Done()
				return
			}

			// 将sum256写入sum256Map
			utils.AddSum256Map(sum256)

			<-ch
			wg.Done()
		}(file)
	}

	wg.Wait()

	os.Remove(tmpFilename) // 删除zip缓存文件
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Zip uploaded and unzipped successfully",
	})
	ctx.Next()
}
