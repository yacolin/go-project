package controllers

import (
	"encoding/json"
	"fmt"
	"go-project/configs"
	"go-project/constants"
	"go-project/models"
	"go-project/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// BookCache 图书相关的缓存键
	BookCache = utils.NewCacheKeys(constants.BOOK)
)

// @Summary 获取所有图书信息
// @Description 获取所有图书信息，支持分页
// @Tags books
// @Produce json
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ApiRes
// @Failure 500 {object} utils.ApiRes
// @Router /books [get]
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

	// 3.1 获取数据总数缓存
	totalStr, err := configs.GetCache(BookCache.TotalKey)
	if err == nil {
		// 缓存命中，解析总数
		count, _ = strconv.ParseInt(totalStr, 10, 64)
	} else {
		if err := baseQuery.Count(&count).Error; err != nil {
			c.Error(utils.NewBusinessError(
				utils.DBQuery,
				http.StatusInternalServerError,
				gin.H{"operation": "query_books"},
				fmt.Errorf("查询总计失败：%w", err),
			))
			return
		}
		// 设置总数缓存，过期时间5分钟
		_ = configs.SetCache(BookCache.TotalKey, fmt.Sprintf("%d", count), utils.DefaultCacheTime)

	}

	// 3.2 尝试获取分页数据缓存
	listCacheKey := utils.GenListCacheKey(BookCache.ListPrefix, limit, offset)
	listCache, err := configs.GetCache(listCacheKey)
	if err == nil {
		// 缓存命中，解析数据
		if err := json.Unmarshal([]byte(listCache), &books); err == nil {
			utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
				List:  books,
				Total: count,
			})
			return
		}
	}

	// 3.3 缓存未命中或解析失败，从数据库查询
	if err := baseQuery.Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_books"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 3.4 设置分页数据缓存，过期时间5分钟
	if listData, err := json.Marshal(books); err == nil {
		_ = configs.SetCache(listCacheKey, string(listData), utils.DefaultCacheTime)
	}

	// 3. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  books,
		Total: count,
	})
}

// @Summary 获取单个图书信息
// @Description 根据ID获取图书详情
// @Tags books
// @Produce json
// @Param id path int true "图书ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} utils.BusinessError
// @Router /book/{id} [get]
func GetBookByID(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := utils.FindByID(
		c,
		configs.DB,
		id,
		&book,
		utils.QueryOptions{ResourceName: constants.BOOK},
	); err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, book)
}

// @Summary 创建图书信息
// @Description 创建一本新图书
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.BookForm true "图书信息"
// @Success 201 {object} models.Book
// @Failure 400 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /books [post]
func CreateBook(c *gin.Context) {
	// 绑定请求数据
	var createReq models.BookForm
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig(constants.BOOK))},
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
			utils.DBCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_book"},
			fmt.Errorf("book创建失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(BookCache)

	// 返回创建结果
	utils.Success(
		c,
		http.StatusCreated,
		utils.Created,
		newBook,
	)
}

// @Summary 删除图书信息
// @Description 根据ID删除图书
// @Tags books
// @Produce json
// @Param id path int true "图书ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /book/{id} [delete]
func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := utils.FindByID(
		c,
		configs.DB,
		id,
		&book,
		utils.QueryOptions{ResourceName: constants.BOOK},
	); err != nil {
		c.Error(err)
		return
	}

	if err := configs.DB.Delete(&models.Book{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_book"},
			fmt.Errorf("删除失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(BookCache)

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"id": id})
}

// @Summary 更新图书信息
// @Description 根据ID更新图书信息
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "图书ID"
// @Param book body models.BookForm true "图书信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.BusinessError
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /book/{id} [put]
func UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := utils.FindByID(
		c,
		configs.DB,
		id,
		&book,
		utils.QueryOptions{ResourceName: constants.BOOK},
	); err != nil {
		c.Error(err)
		return
	}

	// 绑定请求数据
	var updateReq models.BookForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig(constants.BOOK))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 执行更新
	if err := configs.DB.Model(&book).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_book"},
			fmt.Errorf("book更新失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(BookCache)

	// 返回创建结果
	utils.Success(
		c,
		http.StatusOK,
		utils.Updated,
		nil,
	)
}

// @Summary 搜索图书
// @Description 根据作者、标题或ISBN模糊搜索图书
// @Tags books
// @Produce json
// @Param author query string false "作者"
// @Param title query string false "标题"
// @Param isbn query string false "ISBN"
// @Success 200 {array} models.Book
// @Failure 400 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /books/search [get]
func SearchBooks(c *gin.Context) {
	// 1. 获取查询参数
	author := strings.TrimSpace(c.Query("author"))
	title := strings.TrimSpace(c.Query("title"))
	isbn := strings.TrimSpace(c.Query("isbn"))

	// 2. 检查是否至少提供了一个查询参数
	if author == "" && title == "" && isbn == "" {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
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
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "search_books"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 5. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, books)
}
