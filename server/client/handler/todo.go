package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"life_journey/database"
	"life_journey/model"
	"life_journey/response"
)

// ListTodos 获取待办列表，支持按 status 过滤
func ListTodos(c *gin.Context) {
	var todos []model.Todo
	query := database.DB.Order("created_at desc")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&todos).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询待办失败")
		return
	}
	response.Success(c, todos)
}

// CreateTodo 创建待办
func CreateTodo(c *gin.Context) {
	var todo model.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if todo.Title == "" {
		response.Fail(c, http.StatusBadRequest, "待办标题不能为空")
		return
	}
	if todo.Status == "" {
		todo.Status = "pending"
	}
	if err := database.DB.Create(&todo).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建待办失败")
		return
	}
	response.Success(c, todo)
}

// UpdateTodo 更新待办
func UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的 ID")
		return
	}

	var todo model.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "待办不存在")
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 防止覆盖不可修改的字段
	delete(input, "id")
	delete(input, "created_at")

	if len(input) > 0 {
		database.DB.Model(&todo).Updates(input)
	}

	// 重新加载获取更新后的数据
	database.DB.First(&todo, id)
	response.Success(c, todo)
}

// DeleteTodo 删除待办
func DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的 ID")
		return
	}

	if err := database.DB.Delete(&model.Todo{}, id).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除待办失败")
		return
	}
	response.Success(c, nil)
}
