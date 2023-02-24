package file

import (
	"cdn/auth"
	"cdn/database"
	"cdn/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
)

func GetAllFiles(ctx *gin.Context) {

	user := auth.ExtractUser(ctx)

	ctx.IndentedJSON(http.StatusOK, utils.Filter(database.FileData, func(file database.StoredFile) bool {
		if file.Public {
			return true
		}
		if user == nil {
			return false
		}
		if (*user).Admin {
			return true
		}
		return file.Owner == (*user).Id
	}))
}

func GetFile(display bool) func(ctx2 *gin.Context) {

	return func(ctx *gin.Context) {
		fileId := ctx.Param("id")

		fileInfo := utils.FindFirst(&database.FileData, func(f database.StoredFile) bool {
			return f.Id == fileId
		})

		if fileInfo == nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if !fileInfo.Public {
			user := auth.ExtractUser(ctx)
			if user == nil {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}
			if !(*user).Admin && (*user).Id != fileInfo.Owner {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}
		}

		targetPath := "./files/" + fileInfo.Id + fileInfo.Ext

		if display {
			ctx.Header("Content-Disposition", "inline")
		} else {
			ctx.Header("Content-Description", "File Transfer")
			ctx.Header("Content-Transfer-Encoding", "binary")
			ctx.Header("Content-Disposition", "attachment; filename="+fileInfo.Name)
			ctx.Header("Content-Type", "application/octet-stream")
		}

		ctx.File(targetPath)
	}
}

func UploadFile(ctx *gin.Context) {
	user := auth.ExtractUser(ctx)
	if user == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fileName := filepath.Base(file.Filename)
	id := uuid.New()

	ext := filepath.Ext(file.Filename)

	newFile := database.StoredFile{
		Id:        id.String(),
		Ext:       ext,
		Name:      fileName,
		Public:    true,
		Owner:     user.Id,
		Encrypted: false,
	}

	database.FileData = append(database.FileData, newFile)

	err = ctx.SaveUploadedFile(file, "./files/"+id.String()+ext)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.IndentedJSON(http.StatusOK, newFile)
}
