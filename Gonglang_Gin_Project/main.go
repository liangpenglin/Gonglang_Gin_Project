package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

var (
	DB *gorm.DB
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func initMysql() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return err
}

func main() {
	// 创建数据库

	// 连接数据库
	err := initMysql()
	if err != nil {
		panic(err)
	}
	// 模型绑定
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	// 告诉 gin 框架模板文件引用的静态文件去哪里找
	r.Static("/static", "static")
	// 告诉 gin 框架去哪里找模板文件
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// v1
	v1Group := r.Group("v1")
	{

		// 待办事项
		// 添加
		v1Group.POST("/todo", func(context *gin.Context) {
			// 前端页面填写待办事项 点击提交 会发请求到这里
			// 1.从请求中把数据拿出来
			var todo Todo
			context.BindJSON(&todo)

			// 2.存入数据库
			//err = DB.Create(&todo).Error
			//if err != nil {
			//}
			if err = DB.Create(&todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, todo)
				//context.JSON(http.StatusOK,gin.H{
				//	"code":2000,
				//	"data":todo,
				//})
			}
		})
		// 查看所有待办事项
		v1Group.GET("/todo", func(context *gin.Context) {
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				context.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, todoList)
			}
		})
		// 修改某一个待办事项
		v1Group.PUT("/todo/:id", func(context *gin.Context) {
			// 根据传入的 id 修改状态
			id, ok := context.Params.Get("id")
			if !ok {
				context.JSON(http.StatusOK, gin.H{"error": "id 不存在"})
			}
			var todo Todo
			if err = DB.Where("id=?", id).First(&todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			}
			context.BindJSON(&todo)
			if err = DB.Save(&todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, todo)
			}
		})
		// 删除某一个待办事项
		v1Group.DELETE("/todo/:id", func(context *gin.Context) {
			id, ok := context.Params.Get("id")
			if !ok {
				context.JSON(http.StatusOK, gin.H{
					"error": "无效的id",
				})
				return
			}

			if err = DB.Where("id=?", id).Delete(&Todo{}).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					id: "deleted",
				})
			}

		})
	}
	r.Run()
}
