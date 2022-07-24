package bedrockc

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"

	"github.com/PuerkitoBio/goquery"
)
const(
	bedrockDownloadWebsite="https://www.minecraft.net/en-us/download/server/bedrock"
)
type bedrockS struct {
	path string
	version string
	config *Config
}
func (b *bedrockS) Init(path string,version string,config *Config) {
	b.path=config.Require("BedrockPath","./BDS")
	//查看path是否存在，否则创建
	if _, err := os.Stat(b.path); os.IsNotExist(err) {
		os.Mkdir(b.path, os.ModePerm)
		var downloadPath string;
		if(config.Require("DownloadFromOfficial","true")=="true"){
			//分析网页，获取下载地址
			client:=http.Client{};
			req,err:=http.NewRequest(http.MethodGet,bedrockDownloadWebsite,nil);
			if err!=nil{
				errors.Wrap(err,"访问Bedrock Server官方网页失败")
			}
			req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
			resp,err:=client.Do(req);
			defer resp.Body.Close()
			doc,err:=goquery.NewDocumentFromReader(resp.Body)
			if err!=nil{
				errors.Wrap(err,"无法解析网页数据")
			}
			//找到一个a标签，data-platform="serverBedrockWindows"，内容是Download
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				if s.AttrOr("data-platform","")=="serverBedrockWindows" && s.Text()=="Download"{
					downloadPath,_=s.Attr("href")
				}
			})
			//下载游戏
			err=b.downloadGame(downloadPath,b.path,true)
			if err!=nil{
				errors.Wrap(err,"下载游戏失败")
			}
		}else{
			downloadPath=config.Require("DownloadLink","");
		}
		
		b.downloadGame(downloadPath,config.Require("BedrockPath","./BDS")+"/"+version,true)
	}
}
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
func downloaderChooser() string{
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
func downloadCommandWrapper(tools string,url string,path string) string{
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
func deCompress(archive string,path string) error{
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

func (b *bedrockS) downloadGame(url string,path string,tell bool) error {
	re:=regexp.MustCompile(`bedrock-server-([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\.zip`)
	version:=re.FindStringSubmatch(url)[1]
	if(tell){
		println("下载游戏,地址:"+url+",路径:"+path+"版本:"+version)
	}
	downloader:=downloaderChooser()
	if(downloader==""){
		return errors.New("没有找到可用的下载器")
	}
	//如果bedrock.zip已经存在就删除
	if _, err := os.Stat(path+"/bedrock.zip"); err == nil {
		os.Remove(path+"/bedrock.zip")
	}
	command:=downloadCommandWrapper(downloader,url,"./bedrock.zip");
	//执行下载命令
	cmd:=exec.Command("/bin/bash","-c",command)
	err:=cmd.Run()
	if err!=nil{
		return errors.Wrap(err,"下载游戏失败")
	}
	//解压游戏
	err=deCompress("./bedrock.zip",path)
	if err!=nil{
		return errors.Wrap(err,"解压游戏失败")
	}
	//删除游戏包
	os.Remove(path+"/bedrock.zip")
	return nil
}