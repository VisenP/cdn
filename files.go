package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
)

func getAllFiles(ctx *gin.Context) {

	user := extractUser(ctx)

	ctx.IndentedJSON(http.StatusOK, filter(fileData, func(file storedFile) bool {
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

func getFile(ctx *gin.Context) {

	fileId := ctx.Param("id")

	fileInfo := findFirst(&fileData, func(f storedFile) bool {
		return f.Id == fileId
	})

	if fileInfo == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if !fileInfo.Public {
		user := extractUser(ctx)
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

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+fileInfo.Name)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(targetPath)
}

func uploadFile(ctx *gin.Context) {
	user := extractUser(ctx)
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

	newFile := storedFile{
		Id:        id.String(),
		Ext:       ext,
		Name:      fileName,
		Public:    true,
		Owner:     user.Id,
		Encrypted: false,
	}

	fileData = append(fileData, newFile)

	err = ctx.SaveUploadedFile(file, "./files/"+id.String()+ext)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.IndentedJSON(http.StatusOK, newFile)
}
