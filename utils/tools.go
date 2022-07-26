package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
			return true, nil
	}
	if os.IsNotExist(err) {
			return false, nil
	}
	return false, err
}
func DownloaderChooser() string{
	//依次查看系统自带的下载器，如果有则返回，没有则返回空
	///usr/bin/curl,/usr/bin/wget,/usr/bin/aria2c,/usr/bin/axel
	result_curl,_:=PathExists("/usr/bin/curl")
	result_aria2c,_:=PathExists("/usr/bin/aria2c")
	result_wget,_:=PathExists("/usr/bin/wget")
	result_axel,_:=PathExists("/usr/bin/axel")
	if result_axel{
		return "/usr/bin/axel"
	}else if result_aria2c{
		return "/usr/bin/aria2c"
	}else if result_wget{
		return "/usr/bin/wget"
	}else if result_curl{
		return "/usr/bin/curl"
	}else{
		return ""
	}
}
func DownloadCommandWrapper(tools string,url string,path string) string{
	if(tools=="/usr/bin/curl"){
		return tools+" -L "+url+" -o "+path
	}else if(tools=="/usr/bin/wget"){
		return tools+" -O "+path+" "+url
	}else if(tools=="/usr/bin/aria2c"){
		return tools+" "+url+" -d "+path
	}else if(tools=="/usr/bin/axel"){
		return tools+" "+url+" -o "+path
	}else{
		return ""
	}

}
//解压zip到文件夹
func DeCompress(archive string,path string) error{
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return errors.Wrap(err,"无法打开zip文件")
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return errors.Wrap(err,"无法打开zip文件")
		}
		defer rc.Close()
		path := filepath.Join(path, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return errors.Wrap(err,"无法打开zip文件")
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return errors.Wrap(err,"无法打开zip文件")
			}
		}
	}
	return nil
}

