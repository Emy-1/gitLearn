package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

var (
	DB *gorm.DB
)

//Todo MODEL
type Todo struct {
	ID     int    `json:"id"`
	Status bool   `json:"status"`
	Title  string `json:"title"`
}

func initMysql() (err error) {
	dsn := "root:951003@(localhost:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}
func main() {
	//创建数据库
	err := initMysql()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	//连接数据库
	//绑定模型
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	//告诉gin模板框架文件引用的静态文件在哪里找
	r.Static("./static", "./static")
	//告诉gin在哪里找模板文件
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	//v1 api
	v1Group := r.Group("v1")

	v1Group.POST("./todo", func(c *gin.Context) {
		var todo Todo
		c.Bind(&todo)
		err := DB.Create(&todo).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, &todo)
		}

	})
	v1Group.GET("./todo", func(c *gin.Context) {
		var todoList []Todo
		c.Bind(&todoList)
		err := DB.Find(&todoList).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, &todoList)
		}
	})
	v1Group.GET("./todo/:id", func(c *gin.Context) {

	})
	v1Group.PUT("./todo/:id", func(c *gin.Context) {
		todoId, ok := c.Params.Get("id")
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"error": "无效id",
			})
		}
		var todo Todo
		err := DB.Where("id=?", todoId).First(&todo).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})

		}

		c.Bind(&todo)
		err = DB.Save(&todo).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, &todo)
		}

	})

	v1Group.DELETE("./todo/:id", func(c *gin.Context) {
		param, ok := c.Params.Get("id")
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"error": "wrong",
			})
		}

		c.Bind(param)
		err := DB.Where("id=?", param).Delete(&Todo{}).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, &Todo{})
		}
	})
	//待办事项
	//添加
	//查看
	//修改
	//删除
	r.Run()
}
