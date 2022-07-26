package bedrock

import (
	"BedrockC/config"
	"BedrockC/logger"
	"BedrockC/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	bedrockDownloadWebsite = "https://www.minecraft.net/en-us/download/server/bedrock"
)

type bedrockHelper struct {
	bedrockPath     string //用于放置bedrock各版本的路径
	bedrockVersions []string
}

func InitBedrockServer(executable string) (error, BedrockS) {
	temp := BedrockS{}
	temp.path = filepath.Dir(executable)
	temp.executable = executable
	return nil, temp
}
func (b *bedrockHelper) InitByVersion(version string) BedrockS {
	var temp string
	if version == "latest" {
		var versions [][]int
		for _, v := range b.bedrockVersions {
			temp := strings.Split(v, ".")[:]
			tempInt := make([]int, len(temp))
			for j, v2 := range temp {
				tempInt[j], _ = strconv.Atoi(v2)
			}
			versions = append(versions, tempInt)
		}
		//在versions中找到最大的版本
		var maxVersion []int = versions[0]
		for _, v := range versions {
			for i, v2 := range v {
				if v2 <= maxVersion[i] {
					break
				}
				if v2 == 3 {
					maxVersion = v
				}
			}
		}
		for i, v := range maxVersion {
			temp += strconv.Itoa(v)
			if i < len(maxVersion)-1 {
				temp += "."
			}
		}
	} else {
		temp = version
	}
	logger.DefaultLogger().Message("准备启动实例"+temp, "bedrockHelper")
	err, bs := InitBedrockServer(filepath.Join(b.bedrockPath, temp, "bedrock"))
	if err != nil {
		logger.DefaultLogger().Error(err, "无法启动实例", "bedrockHelper")
		return BedrockS{}
	}
	return bs
}

//从官方下载最新的游戏
func (b *bedrockHelper) UpdateGame(path string) error {
	link := GetLatestLink()
	re := regexp.MustCompile(`bedrock-server-([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\.zip`)
	version := re.FindStringSubmatch(link)[1]
	logger.DefaultLogger().Message("开始下载Bedrock Server "+version, "bedrockHelper")
	downloader := utils.DownloaderChooser()
	if downloader == "" {
		return errors.New("没有可用下载器")
	}
	if _, err := os.Stat("./bedrock.zip"); err != nil {
		os.Remove("./bedrock.zip")
		logger.DefaultLogger().Warn("下载游戏时发现已存在的下载文件，删除", "bedrockHelper")
	}
	command := utils.DownloadCommandWrapper(downloader, link, "./bedrock.zip")
	cmd := exec.Command(command)
	err := cmd.Run()
	if err != nil {
		logger.DefaultLogger().Error(errors.Wrap(err, "下载失败"), "下载失败", "bedrockHelper")
		return err
	}
	err = utils.DeCompress("./bedrock.zip", filepath.Join(path, version))
	if err != nil {
		logger.DefaultLogger().Error(errors.Wrap(err, "解压失败"), "解压失败", "bedrockHelper")
		return err
	}
	os.Remove("./bedrock.zip")
	logger.DefaultLogger().Message("下载完成", "bedrockHelper")
	//给config添加
	b.bedrockVersions = append(b.bedrockVersions, version)
	return nil
}
func GetLatestLink() string {
	client := http.Client{}
	var downloadPath string
	req, err := http.NewRequest(http.MethodGet, bedrockDownloadWebsite, nil)
	if err != nil {
		errors.Wrap(err, "访问Bedrock Server官方网页失败")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		errors.Wrap(err, "访问Bedrock Server官方网页失败")
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		errors.Wrap(err, "无法解析网页数据")
	}
	//找到一个a标签，data-platform="serverBedrockWindows"，内容是Download
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.AttrOr("data-platform", "") == "serverBedrockWindows" && s.Text() == "Download" {
			downloadPath, _ = s.Attr("href")
		}
	})
	return downloadPath
}

//在最后的时候调用，用于保存配置
func (b *bedrockHelper) UpdateConfig(c *config.Config) {
	c.Set("bedrockPath", b.bedrockPath)
	c.Set("bedrockServers", b.bedrockVersions)
}
func (b *bedrockHelper) Init(c *config.Config) error {
	b.bedrockPath = c.Require("bedrockPath", "./BDS").(string)
	b.bedrockVersions = c.Require("bedrockServers", []string{}).([]string)
	if len(b.bedrockVersions) == 0 {
		//认定为新目录，删除bedrockPath下的所有文件
		err := os.RemoveAll(b.bedrockPath)
		if err != nil {
			return errors.Wrap(err, "无法删除bedrock目录")
		}
	} else {
		//防止下载到一半，删除bedrockPath下没有记录的文件夹
		for _, v := range b.bedrockVersions {
			temp, _ := utils.PathExists(filepath.Join(b.bedrockPath, v))
			if !temp {
				err := os.RemoveAll(filepath.Join(b.bedrockPath, v))
				if err != nil {
					return errors.Wrap(err, "无法删除bedrock目录")
				}
			}
		}
	}
	return nil
}
