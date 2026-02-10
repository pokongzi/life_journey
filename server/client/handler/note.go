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

// ListNotes 获取笔记列表（仅元数据，不加载文件内容）
func ListNotes(c *gin.Context) {
	var notes []model.Note
	query := database.DB.Order("updated_at desc")

	if nbID := c.Query("notebook_id"); nbID != "" {
		query = query.Where("notebook_id = ?", nbID)
	}

	if err := query.Find(&notes).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询笔记失败")
		return
	}
	// 列表场景不返回文件内容，Content 默认为空
	response.Success(c, notes)
}

// GetNote 获取单条笔记（从 .md 文件加载内容）
func GetNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的 ID")
		return
	}

	var note model.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "笔记不存在")
		return
	}

	// 从 Vault 中的 .md 文件读取内容
	content, err := vault.ReadNoteContent(note.FilePath)
	if err != nil {
		// 文件可能被外部删除，返回空内容但不报错
		content = ""
	}
	note.Content = content

	response.Success(c, note)
}

// CreateNote 创建笔记（在 Vault 中生成 .md 文件）
func CreateNote(c *gin.Context) {
	var input struct {
		NotebookID uint   `json:"notebook_id" binding:"required"`
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误：需要 notebook_id 和 title")
		return
	}

	// 查找所属笔记本，确定文件夹路径
	var nb model.Notebook
	if err := database.DB.First(&nb, input.NotebookID).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "笔记本不存在")
		return
	}

	// 生成唯一文件路径
	filePath := vault.UniqueNoteRelPath(nb.Path, input.Title)

	// 写入 .md 文件
	if err := vault.WriteNoteContent(filePath, input.Content); err != nil {
		response.Fail(c, http.StatusInternalServerError, "写入笔记文件失败")
		return
	}

	// 保存元数据到数据库
	note := model.Note{
		NotebookID: input.NotebookID,
		Title:      input.Title,
		FilePath:   filePath,
	}
	if err := database.DB.Create(&note).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建笔记记录失败")
		return
	}

	note.Content = input.Content
	response.Success(c, note)
}

// UpdateNote 更新笔记（标题变更时重命名文件，内容变更时重写文件）
func UpdateNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的 ID")
		return
	}

	var note model.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "笔记不存在")
		return
	}

	var input struct {
		NotebookID *uint   `json:"notebook_id"`
		Title      *string `json:"title"`
		Content    *string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	dbUpdates := map[string]interface{}{}
	needRename := false
	newNotebookDir := ""

	// 确定目标笔记本目录
	if input.NotebookID != nil && *input.NotebookID != note.NotebookID {
		var newNb model.Notebook
		if err := database.DB.First(&newNb, *input.NotebookID).Error; err != nil {
			response.Fail(c, http.StatusBadRequest, "目标笔记本不存在")
			return
		}
		newNotebookDir = newNb.Path
		dbUpdates["notebook_id"] = *input.NotebookID
		needRename = true
	}

	// 标题变更 → 需要重命名文件
	if input.Title != nil && *input.Title != note.Title {
		dbUpdates["title"] = *input.Title
		needRename = true
	}

	// 重命名文件（标题或笔记本变更时）
	if needRename {
		if newNotebookDir == "" {
			// 未换笔记本，取当前笔记本目录
			var currentNb model.Notebook
			database.DB.First(&currentNb, note.NotebookID)
			newNotebookDir = currentNb.Path
		}
		newTitle := note.Title
		if input.Title != nil {
			newTitle = *input.Title
		}
		newFilePath, err := vault.RenameNoteFile(note.FilePath, newNotebookDir, newTitle)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "重命名笔记文件失败")
			return
		}
		dbUpdates["file_path"] = newFilePath
		note.FilePath = newFilePath
	}

	// 内容变更 → 重写文件
	if input.Content != nil {
		currentPath := note.FilePath
		if fp, ok := dbUpdates["file_path"]; ok {
			currentPath = fp.(string)
		}
		if err := vault.WriteNoteContent(currentPath, *input.Content); err != nil {
			response.Fail(c, http.StatusInternalServerError, "写入笔记文件失败")
			return
		}
	}

	// 更新数据库元数据
	if len(dbUpdates) > 0 {
		database.DB.Model(&note).Updates(dbUpdates)
	}

	// 重新加载并附带文件内容返回
	database.DB.First(&note, id)
	content, _ := vault.ReadNoteContent(note.FilePath)
	note.Content = content

	response.Success(c, note)
}

// DeleteNote 删除笔记（同步删除 .md 文件）
func DeleteNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的 ID")
		return
	}

	var note model.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "笔记不存在")
		return
	}

	// 删除 Vault 中的 .md 文件
	vault.DeleteNoteFile(note.FilePath)

	// 删除数据库记录
	database.DB.Delete(&note)

	response.Success(c, nil)
}
