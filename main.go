package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"gocms/application/admin/controller"
	"gocms/application/app_session"
	controller_index "gocms/application/index/controller"
	"gocms/config"
	"strconv"
)

var app *iris.Application

func main() {
	app = iris.Default()

	//app.RegisterView(iris.HTML("./application/admin/view", ".html"))              //注册视图目录下的html模板
	//当注册多个视图环境时，有时会出现随机的找不到视图的BUG，不是必现，不知道为什么
	//app.RegisterView(iris.HTML("./application/admin/view", ".html").Reload(true)) //注册视图目录下的html模板,加上reload方法，静态页面改变时可以实时更新
	//app.RegisterView(iris.HTML("./application/index/view", ".html").Reload(true)) //注册视图目录下的html模板,加上reload方法，静态页面改变时可以实时更新
	app.HandleDir("static", "./static") //定义静态文件的请求路径与实际路径的映射
	//后台管理路由部分
	bossPart := app.Party("/boss", controller.AdminAuth)                               //admin分组,admin下分组都使用adminAuth中间件
	bossPart.RegisterView(iris.HTML("./application/admin/view", ".html").Reload(true)) //注册视图目录下的html模板,加上reload方法，静态页面改变时可以实时更新
	bossPart.Use(iris.Compression)
	index := mvc.New(bossPart.Party("/index")) //IndexController
	index.Handle(new(controller.IndexController))
	public := mvc.New(bossPart.Party("/public")) //PublicController
	public.Handle(new(controller.PublicController))
	node := mvc.New(bossPart.Party("/node")) //NodeController
	node.Handle(new(controller.NodeController))
	role := mvc.New(bossPart.Party("/role")) //RoleController
	role.Handle(new(controller.RoleController))
	user := mvc.New(bossPart.Party("/user")) //UserController
	user.Handle(new(controller.UserController))
	model := mvc.New(bossPart.Party("/model")) //ModelController
	model.Handle(new(controller.ModelController))
	attribute := mvc.New(bossPart.Party("/attribute")) //AttributeController
	attribute.Handle(new(controller.AttributeController))
	system := mvc.New(bossPart.Party("/system")) //SystemController
	system.Handle(new(controller.SystemController))
	category := mvc.New(bossPart.Party("/category")) //CategoryController
	category.Handle(new(controller.CategoryController))
	article := mvc.New(bossPart.Party("/article")) //ArticleController
	article.Handle(new(controller.ArticleController))
	photo := mvc.New(bossPart.Party("/photo")) //PhotoController
	photo.Handle(new(controller.PhotoController))
	static := mvc.New(bossPart.Party("/static")) //StaticController
	static.Handle(new(controller.StaticController))
	message := mvc.New(bossPart.Party("/message")) //MessageController
	message.Handle(new(controller.MessageController))

	//前台展示路由部分
	indexPart := app.Party("/index", controller_index.IndexAuth)
	indexPart.RegisterView(iris.HTML("./application/index/view", ".html").Reload(true)) //注册视图目录下的html模板,加上reload方法，静态页面改变时可以实时更新
	indexPart.Use(iris.Compression)
	indexPart_index := mvc.New(indexPart.Party("/index"))
	indexPart_index.Handle(new(controller_index.IndexController))
	indexPart_about := mvc.New(indexPart.Party("/about"))
	indexPart_about.Handle(new(controller_index.AboutController))
	indexPart_article := mvc.New(indexPart.Party("/article"))
	indexPart_article.Handle(new(controller_index.ArticleController))
	indexPart_product := mvc.New(indexPart.Party("/product"))
	indexPart_product.Handle(new(controller_index.ProductController))
	indexPart_contact := mvc.New(indexPart.Party("/contact"))
	indexPart_contact.Handle(new(controller_index.ContactController))

	iris.RegisterOnInterrupt(func() {
		app_session.Sessdb.Close() //关闭session的本地数据库存储
	})

	app.Listen(":" + strconv.Itoa(config.InitConfig().AppUrlPort))
}
