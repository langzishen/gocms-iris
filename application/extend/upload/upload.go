package upload

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/kataras/iris/v12"
	"gocms/config"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

/**
单个文件上传：
	post_file_name：form的文件域name
	返回：上传后url访问路径,上传后文件的绝对路径,错误
*/
func Upload(ctx iris.Context, post_file_name string, UploadPath ...string) (url_path string, Abs_path string, err error) {
	app_conf := config.InitConfig()
	ctx.SetMaxRequestBodySize(app_conf.Upload.UploadTypeConf.MaxSize) //设置最大上传的文件的大小
	file, fileHeader, err := ctx.FormFile(post_file_name)
	defer file.Close()
	if err != nil {
		return "", "", err
	}
	root_path, _ := filepath.Abs(app_conf.Upload.UploadTypeConf.RootPath)
	sub_dir := ""
	if len(UploadPath) == 0 {
		time_now := time.Now()
		date_dir := fmt.Sprintf("%d-%d-%d", time_now.Year(), time_now.Month(), time_now.Day())
		_, err = os.Stat(filepath.Join(root_path, date_dir))
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(filepath.Join(root_path, date_dir), 0777)
				if err != nil {
					return "", "", err
				}
			}
		}
		sub_dir = date_dir
	} else {
		sub_dir = UploadPath[0]
	}
	ext := path.Ext(fileHeader.Filename)
	file_name := filepath.Join(root_path, sub_dir, uniqueId()+ext)
	//file_name := RootPath + date_dir + "/" + uniqueId() + ext
	_, err = ctx.SaveFormFile(fileHeader, file_name)
	if err != nil {
		return "", "", err
	}

	switch app_conf.Upload.FileUploadType {
	case "Local":
		if strings.Contains(sub_dir, "\\") {
			sub_dir = strings.Replace(sub_dir, "\\", "/", -1)
		}
		url_path = app_conf.Upload.UploadTypeConf.RootPath[1:] + sub_dir + "/" + filepath.Base(file_name)
		Abs_path = file_name
		err = nil
		break
	case "AliyunOSS":
		url_path, Abs_path, err = PutToOss(app_conf.Upload.UploadTypeConf, file_name)
		os.Remove(file_name) //删除本地文件
		break
	default:
		if strings.Contains(sub_dir, "\\") {
			sub_dir = strings.Replace(sub_dir, "\\", "/", -1)
		}
		url_path = app_conf.Upload.UploadTypeConf.RootPath[1:] + sub_dir + "/" + filepath.Base(file_name)
		Abs_path = file_name
		err = nil
		break
	}
	return
}

/**
上传到阿里云OSS
*/
func PutToOss(conf config.UploadTypeConf, file string) (url_path string, Abs_path string, err error) {
	client, err := oss.New(conf.Endpoint, conf.AccessKeyId, conf.AccessKeySecret)
	if err != nil {
		return "", "", err
	}
	bucket, err := client.Bucket(conf.Bucket)
	if err != nil {
		return "", "", err
	}
	time_now := time.Now()
	date_dir := fmt.Sprintf("%d-%d-%d", time_now.Year(), time_now.Month(), time_now.Day())
	file_name := "uploads/" + date_dir + "/" + filepath.Base(file)
	err = bucket.PutObjectFromFile(file_name, file)
	if err != nil {
		return "", "", err
	}
	return conf.Requesturl + file_name, conf.RequestUri + file_name, nil
}

func uniqueId() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
