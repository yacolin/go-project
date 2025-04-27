package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
 * @description: 获取所有专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /books [get]
func GetAllBooks(c *gin.Context) {
	// 1. 参数解析与校验
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return // 直接终止
	}

	// 2. 数据库操作
	var (
		books []models.Book
		count int64
	)

	baseQuery := configs.DB.Model(&models.Book{})

	// 获取数据总数
	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_books"},
			fmt.Errorf("查询总计失败：%w", err),
		))
	}

	// 获取分页数据
	if err := baseQuery.Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_books"},
			fmt.Errorf("查询失败：%w", err),
		))
	}

	// 3. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  books,
		Total: count,
	})
}

/**
 * @description: 获取单个专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /book/:id [get]
func GetBookByID(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := utils.FindByID(
		c,
		configs.DB,
		id,
		&book,
		utils.QueryOptions{ResourceName: "book"},
	); err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, book)
}

/**
 * @description: 创建专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /books [post]
func CreateBook(c *gin.Context) {
	// 绑定请求数据
	var createReq models.BookForm
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err)},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	newBook := models.Book{
		ISBN:        createReq.ISBN,
		Title:       createReq.Title,
		Author:      createReq.Author,
		Stock:       createReq.Stock,
		Publisher:   createReq.Publisher,
		PublishDate: createReq.PublishDate,
	}

	// 写入数据库
	if err := configs.DB.Create(&newBook).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_book"},
			fmt.Errorf("book创建失败：%w", err),
		))
		return
	}

	// 返回创建结果
	utils.Success(
		c,
		http.StatusCreated,
		utils.Created,
		newBook,
	)
}

/**
 * @description: 删除专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /book/:id [delete]
func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := utils.FindByID(
		c,
		configs.DB,
		id,
		&book,
		utils.QueryOptions{ResourceName: "book"},
	); err != nil {
		c.Error(err)
		return
	}

	if err := configs.DB.Delete(&models.Book{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_book"},
			fmt.Errorf("删除失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"id": id})
}

/**
 * @description: 更新专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /book/:id [put]
func UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := utils.FindByID(
		c,
		configs.DB,
		id,
		&book,
		utils.QueryOptions{ResourceName: "book"},
	); err != nil {
		c.Error(err)
		return
	}

	// 绑定请求数据
	var updateReq models.BookForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err)},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 执行更新
	if err := configs.DB.Model(&book).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_book"},
			fmt.Errorf("book更新失败：%w", err),
		))
		return
	}

	// 返回创建结果
	utils.Success(
		c,
		http.StatusOK,
		utils.Updated,
		nil,
	)
}

// @router /books/search [get]
func SearchBooks(c *gin.Context) {
	// 1. 获取查询参数
	author := strings.TrimSpace(c.Query("author"))
	title := strings.TrimSpace(c.Query("title"))
	isbn := strings.TrimSpace(c.Query("isbn"))

	// 2. 检查是否至少提供了一个查询参数
	if author == "" && title == "" && isbn == "" {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"error": "至少提供一个查询参数"},
			fmt.Errorf("查询参数缺失"),
		))
		return
	}

	// 3. 构建查询条件
	query := configs.DB.Model(&models.Book{})
	if author != "" {
		query = query.Or("author LIKE ?", "%"+author+"%")
	}
	if title != "" {
		query = query.Or("title LIKE ?", "%"+title+"%")
	}
	if isbn != "" {
		query = query.Or("isbn LIKE ?", "%"+isbn+"%")
	}

	// 4. 查询数据库
	var books []models.Book
	if err := query.Find(&books).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "search_books"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 5. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, books)
}
