package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"life_journey/database"
	"life_journey/model"
	"life_journey/response"
	"life_journey/vault"
)

// ListNotebooks 获取笔记本列表
func ListNotebooks(c *gin.Context) {
	var notebooks []model.Notebook
	if err := database.DB.Order("created_at desc").Find(&notebooks).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询笔记本失败")
		return
	}
	response.Success(c, notebooks)
}

// CreateNotebook 创建笔记本（同步创建 Vault 中的文件夹）
func CreateNotebook(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if input.Name == "" {
		response.Fail(c, http.StatusBadRequest, "笔记本名称不能为空")
		return
	}

	// 目录名 = 清理后的笔记本名
	dirName := vault.SanitizeName(input.Name)

	// 在 Vault 中创建文件夹
	if err := vault.CreateNotebookDir(dirName); err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建笔记本文件夹失败")
		return
	}

	nb := model.Notebook{
		Name:        input.Name,
		Description: input.Description,
		Path:        dirName,
	}
	if err := database.DB.Create(&nb).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建笔记本失败")
		return
	}
	response.Success(c, nb)
}

// UpdateNotebook 更新笔记本（名称变更时同步重命名文件夹，并更新所有笔记路径）
func UpdateNotebook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的 ID")
		return
	}

	var nb model.Notebook
	if err := database.DB.First(&nb, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "笔记本不存在")
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if input.Description != nil {
		updates["description"] = *input.Description
	}

	// 名称变更 → 重命名文件夹 + 更新笔记文件路径
	if input.Name != nil && *input.Name != nb.Name {
		newDir := vault.SanitizeName(*input.Name)
		oldDir := nb.Path

		if newDir != oldDir {
			// 重命名 Vault 中的文件夹
			if err := vault.RenameNotebookDir(oldDir, newDir); err != nil {
				response.Fail(c, http.StatusInternalServerError, "重命名文件夹失败")
				return
			}

			// 批量更新该笔记本下所有笔记的 file_path
			var notes []model.Note
			database.DB.Where("notebook_id = ?", id).Find(&notes)
			for _, note := range notes {
				// 替换路径前缀: oldDir/xxx.md → newDir/xxx.md
				newFilePath := vault.NoteRelPath(newDir, note.Title)
				database.DB.Model(&note).Update("file_path", newFilePath)
			}

			updates["path"] = newDir
		}
		updates["name"] = *input.Name
	}

	if len(updates) > 0 {
		database.DB.Model(&nb).Updates(updates)
	}

	database.DB.First(&nb, id)
	response.Success(c, nb)
}

// DeleteNotebook 删除笔记本（同步删除 Vault 中的文件夹及其所有 .md 文件）
func DeleteNotebook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的 ID")
		return
	}

	var nb model.Notebook
	if err := database.DB.First(&nb, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "笔记本不存在")
		return
	}

	// 删除 Vault 中的文件夹
	if nb.Path != "" {
		vault.DeleteNotebookDir(nb.Path)
	}

	// 删除数据库中的笔记本及其所有笔记记录
	database.DB.Where("notebook_id = ?", id).Delete(&model.Note{})
	database.DB.Delete(&nb)

	response.Success(c, nil)
}
