package controller

import (
	"github.com/kataras/iris/v12"
	"gocms/application/extend/upload"
	"gocms/config"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type StaticController struct {
	BaseController
}

func (sc *StaticController) GetIndex(ctx iris.Context) {
	ctx.View(RequestController + "/index.html")
}

func (sc *StaticController) PostList(ctx iris.Context) {
	conf := config.InitConfig()
	static_root := conf.Upload.UploadTypeConf.RootPath
	dir := ""
	search := ""
	keyPrefix := ctx.PostValue("keyPrefix")
	if keyPrefix != "" {
		dir = filepath.Join(static_root, keyPrefix)
	} else {
		dir = static_root
	}
	if ctx.PostValue("search") != "" {
		search = ctx.PostValue("search")
	}
	static_dir_list, err := ioutil.ReadDir(dir)
	if err != nil {
		if len(static_dir_list) != 0 {
			sc.TopjuiError(ctx, "目录错误："+err.Error())
		}
	}
	list := []map[string]interface{}{}
	for _, static_dir_vo := range static_dir_list {
		file_map := make(map[string]interface{})
		if static_dir_vo.IsDir() {
			file_map["isFile"] = false
			file_map["mode"] = static_dir_vo.Mode()
			if keyPrefix != "" {
				file_map["uuid"] = keyPrefix + static_dir_vo.Name() + "/"
				file_map["objectName"] = keyPrefix + static_dir_vo.Name() + "/"
			} else {
				file_map["uuid"] = static_dir_vo.Name() + "/"
				file_map["objectName"] = static_dir_vo.Name() + "/"
			}
		} else {
			file_map["isFile"] = true
			file_map["size"] = static_dir_vo.Size()
			file_map["ext"] = filepath.Ext(static_dir_vo.Name())
			if keyPrefix != "" {
				file_map["uuid"] = keyPrefix + static_dir_vo.Name()
				file_map["objectName"] = keyPrefix + static_dir_vo.Name()
			} else {
				file_map["uuid"] = static_dir_vo.Name()
				file_map["objectName"] = static_dir_vo.Name()
			}
		}
		file_map["updateTime"] = static_dir_vo.ModTime()
		if search != "" {
			if strings.Contains(static_dir_vo.Name(), search) {
				list = append(list, file_map)
			}
		} else {
			list = append(list, file_map)
		}
	}
	ctx.JSON(map[string]interface{}{"maxKeys": 100, "keyPrefix": ctx.PostValue("keyPrefix"), "rows": list})
}

func (sc *StaticController) PostAdddir(ctx iris.Context) {
	ctx.View("static/adddir.html")
}

//创建文件夹
func (sc *StaticController) PostDoadddir(ctx iris.Context) {
	dir_name := ctx.PostValue("dir_name")
	keyPrefix := ctx.PostValue("keyPrefix")
	dir_s := strings.Split(dir_name, "/")
	if dir_name == "" {
		sc.TopjuiError(ctx, "目录名不能为空")
	}
	if strings.Contains(dir_name, ".") || strings.Contains(dir_name, "..") {
		sc.TopjuiError(ctx, "目录名不能饱含 '.'或'..'")
	}
	if dir_s[0] == "" {
		sc.TopjuiError(ctx, "目录不能以 ‘/’ 开头")
	}
	conf := config.InitConfig()
	static_root := conf.Upload.UploadTypeConf.RootPath
	dir := ""
	if keyPrefix != "" {
		dir = filepath.Join(static_root, keyPrefix)
	} else {
		dir = static_root
	}
	_, err := os.Stat(dir)
	if err != nil {
		sc.TopjuiError(ctx, "根目录不存在:"+err.Error())
	}
	sub_dir := filepath.Join(dir, dir_name)
	if len(dir_s) > 1 { //多级目录
		os.MkdirAll(sub_dir, 0777)
	} else {
		os.Mkdir(sub_dir, 0777)
	}
	os.Chmod(sub_dir, 0777)
	sc.TopjuiSucess(ctx, "创建成功")
}

/**
文件上传
*/
func (sc *StaticController) PostMultiupload(ctx iris.Context) {
	keyPrefix := ctx.GetCookie("keyPrefix")
	var err error
	if keyPrefix == "" {
		_, _, err = upload.Upload(ctx, "file") //topjui的单文件上传的文件名默认为file，与input的name无关
	} else {
		dir_s := strings.Split(keyPrefix, "/")
		var sub_dir string
		if len(dir_s) == 2 {
			sub_dir = dir_s[0]
		} else {
			sub_dir = filepath.Join(dir_s[0 : len(dir_s)-1]...)
		}
		_, _, err = upload.Upload(ctx, "file", sub_dir)
	}
	if err != nil {
		sc.TopjuiError(ctx, "上传失败："+err.Error())
	} else {
		sc.TopjuiSucess(ctx, "上传成功")
	}
}

/**
删除文件或空文件夹
*/
func (sc *StaticController) PostDelfile(ctx iris.Context) {
	uuid := ctx.PostValue("uuids")
	uuid_s := strings.Split(uuid, ",")
	if len(uuid_s) == 0 {
		return
	}
	conf := config.InitConfig()
	static_root := conf.Upload.UploadTypeConf.RootPath
	dir := ""
	for _, uuid_vo := range uuid_s {
		dir = filepath.Join(static_root, uuid_vo)
		file, _ := os.Stat(dir)
		if file.IsDir() { //目录
			if rd_dir, _ := ioutil.ReadDir(dir); len(rd_dir) != 0 {
				sc.TopjuiError(ctx, "删除失败：选项中饱含非空目录,请依次删除非空目录后再执行删除")
				return
			}
		}
	}
	for _, uuid_vo := range uuid_s {
		dir = filepath.Join(static_root, uuid_vo)
		file, _ := os.Stat(dir)
		if file.IsDir() {
			res := os.RemoveAll(dir)
			if res != nil {
				sc.TopjuiError(ctx, "删除失败："+res.Error())
			}
		} else {
			res := os.Remove(dir)
			if res != nil {
				sc.TopjuiError(ctx, "删除失败："+res.Error())
			}
		}
	}
	sc.TopjuiSucess(ctx, "删除成功：共删除"+strconv.Itoa(len(uuid_s))+"个文件或目录")
}

/**
下载
*/
func (sc *StaticController) GetDownload(ctx iris.Context) {
	uuid := ctx.URLParam("uuid")
	if uuid == "" {
		sc.TopjuiError(ctx, "参数uuid丢失")
	}
	conf := config.InitConfig()
	static_root := conf.Upload.UploadTypeConf.RootPath
	dir := ""
	dir = filepath.Join(static_root, uuid)
	_, err := os.Stat(dir)
	if err != nil {
		sc.TopjuiError(ctx, "文件打开失败:"+err.Error())
	}
	file, _ := os.ReadFile(dir)
	ctx.Header("content-type", "application/octet-stream") //强制下载
	file_name := filepath.Base(dir)
	file_name = html.EscapeString(file_name)
	ctx.Header("content-disposition", "attachment;filename="+file_name) //默认下载名
	ctx.Write(file)
}
