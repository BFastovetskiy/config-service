package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	about "config-service/about"

	"github.com/gin-gonic/gin"
)

func (c Controller) Backup(ctx *gin.Context) {
	backupFileName := filepath.Join(c.app.GetWorkDir(), about.Database_Directory, "backup.db")
	backupFile, err := os.Create(backupFileName)
	var len int64 = 0

	if err != nil {
		ctx.JSON(500, gin.H{"Error": err})
		return
	}

	db := c.app.GetDatabase()
	err = db.Backup(backupFile)
	backupFile.Close()
	if err != nil {
		ctx.JSON(500, gin.H{"Error": err})
		return
	}

	backupFile, _ = os.Open(backupFileName)
	data, err := io.ReadAll(backupFile)
	if err != nil {
		ctx.JSON(500, gin.H{"Error": err})
		return
	}
	disp := fmt.Sprintf(`attachment; filename="backup-%v.db"`, time.Now().UTC())
	ctx.Header("Content-Disposition", disp)
	ctx.Header("Content-Length", strconv.Itoa(int(len)))
	ctx.Data(http.StatusOK, "application/octet-stream", data)
}
