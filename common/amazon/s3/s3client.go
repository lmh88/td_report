package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/frame/g"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"td_report/pkg/logger"
)

const (
	// 开发模式
	developModel = "prod"
	// dev下的路劲前缀
	devpath = "raw_dev/amazon/"
	// prod 下的路劲前缀
	prodPath = "raw/amazon/"
)

type Client struct {
	Bucket   string
	S3client *s3.Client
}

func NewClient(Bucket string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(g.Cfg().GetString("s3.regin")),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(g.Cfg().GetString("s3.key"), g.Cfg().GetString("s3.secret"), "")),
	)

	if err != nil {
		logger.Logger.Error("Failed to load configuration")
		return nil, err
	}

	return &Client{
		Bucket:   Bucket,
		S3client: s3.NewFromConfig(cfg),
	}, nil
}

// ListFilesWithPrefix 罗列s3的文件
func (t *Client) ListFilesWithPrefix(prefix string) error {
	listObjsResponse, err := t.S3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(t.Bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		logger.Logger.Error("Couldn't list bucket contents", err.Error())
		return err
	}

	fmt.Println(len(listObjsResponse.Contents))
	// 只能罗列1000条数据
	for _, object := range listObjsResponse.Contents {
		fmt.Printf("%s (%d bytes, class %v) \n", *object.Key, object.Size, object.StorageClass)
	}

	return nil
}

// UploadDir 上传指定目录， 返回error 和bool ,表示是否为空目录
func (t *Client) UploadDir(ctx context.Context, dirpath string, suffix string, storePath string) (error, bool) {
	if filesMap, err := t.getFileList(ctx, dirpath, suffix, storePath); err != nil {
		logger.Logger.ErrorWithContext(ctx, " upload dir error:", err.Error())
		return err, false
	} else {

		if len(filesMap) == 0 { // 空目录或者没有匹配的后缀的文件
			return nil, true
		} else {
			err = t.uploadMutil(ctx, filesMap)
			return err, false
		}
	}
}

func (t *Client) getFileList(ctx context.Context, dirpath string, suffix string, storePath string) (files map[string]string, err error) {
	files = make(map[string]string, 0)
	dir, err := ioutil.ReadDir(dirpath)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "get dir file error:", err.Error())
		return files, err
	}

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			key := fmt.Sprintf(storePath + "/" + fi.Name())
			files[key] = dirpath + PthSep + fi.Name()
		}
	}

	return files, nil
}

func (t *Client) uploadMutil(ctx context.Context, fileList map[string]string) error {
	var (
		amazonkey string
	)

	prefixPath := t.getPrefixKey()
	for key, filepath := range fileList {
		stat, err := os.Stat(filepath)
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, "Couldn't stat file: "+err.Error())
			return err
		}

		file, err := os.Open(filepath)

		if err != nil {
			logger.Logger.ErrorWithContext(ctx, "Couldn't open local file")
			return err
		}

		amazonkey = fmt.Sprintf("%s%s", prefixPath, key)
		_, err = t.S3client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:        aws.String(t.Bucket),
			Key:           aws.String(amazonkey),
			Body:          file,
			ContentLength: stat.Size(),
		})

		file.Close()
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, "Couldn't upload file: "+err.Error())
			return err
		}
	}

	return nil
}

func (t *Client) UploadSingle(ctx context.Context, filepath string, key string) error {
	var (
		amazonkey string
		n         = 3
		doflag    = false
	)

	stat, err := os.Stat(filepath)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "Couldn't stat file: "+err.Error())
		return err
	}

	file, err := os.Open(filepath)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "Couldn't open local file")
		return err
	}

	amazonkey = t.getPrefiFullKey(key)

	// 重试3次
	for i := 0; i < n; i++ {
		_, err = t.S3client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:        aws.String(t.Bucket),
			Key:           aws.String(amazonkey),
			Body:          file,
			ContentLength: stat.Size(),
		})

		if err == nil {
			doflag = true
			break
		} else {
			logger.Logger.ErrorWithContext(ctx, "upload file error, retry")
		}
	}

	file.Close()
	if doflag != true {
		logger.Logger.ErrorWithContext(ctx, "Couldn't upload file: "+err.Error())
		return err
	}

	return nil
}

func (t *Client) getPrefixKey() string {
	var amazonkey string
	devmod := g.Cfg().GetString("server.Env")
	if devmod == "" {
		devmod = "dev"
	}
	if devmod == developModel {
		amazonkey = prodPath
	} else {
		amazonkey = devpath
	}

	return amazonkey
}

func (t *Client) getPrefiFullKey(key string) string {
	var amazonkey string
	devmod := g.Cfg().GetString("server.Env")
	if devmod == developModel {
		amazonkey = fmt.Sprintf("%s%s", prodPath, key)
	} else {
		amazonkey = fmt.Sprintf("%s%s", devpath, key)
	}

	return amazonkey
}

// DownloadSingle 下载单个文件
func (t *Client) DownloadSingle(key string, filename string) error {
	var (
		file      *os.File
		amazonkey string
	)

	amazonkey = t.getPrefiFullKey(key)
	getObjectResponse, err := t.S3client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(t.Bucket),
		Key:    aws.String(amazonkey),
	})

	if err == nil {
		file, err = os.Create(filename)
		if err != nil {
			logger.Logger.Error("didnt open file to write: " + err.Error())
			return err
		}

		written, err := io.Copy(file, getObjectResponse.Body)
		if err != nil {
			logger.Logger.Error("Failed to write file contents! " + err.Error())
			return err
		} else if written != getObjectResponse.ContentLength {
			logger.Logger.Error("wrote a different size than was given to us")
			return err
		}

		file.Close()

	} else {
		logger.Logger.Error("Couldn't download object")
		return err
	}

	return nil
}

// GetFileNum 获取指定路劲下文件的数量
func (t *Client) GetFileNum(key string) (int, error) {
	path := t.getPrefiFullKey(key)
	listObjsResponse, err := t.S3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(t.Bucket),
		Prefix: aws.String(path),
	})

	if err != nil {
		logger.Logger.Error("Couldn't list bucket contents")
		return 0, err
	}

	return len(listObjsResponse.Contents), nil
}
