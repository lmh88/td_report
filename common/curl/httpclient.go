package curl

import (
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/guonaihong/gout/filter"
	"time"
)

type Method string

const (
	GET    Method = "GET"
	POST          = "POST"
	PUT           = "PUT"
	DELETE        = "DELETE"
)

// MaxWaitInterval 最大等待时间单位s
const MaxWaitInterval = 300

type Httpclient struct {
	apiUrl     string
	header     map[string]interface{}
	method     Method
	isJson     bool
	retryTimes int           // 重试次数
	timeout    time.Duration // 超时时间单位秒
	http       *dataflow.Gout
	Code       int
	isDebug    bool
	result     interface{}
}

func NewHttpclient() *Httpclient {
	client := &Httpclient{
		timeout:    900 * time.Second,
		retryTimes: 8,
		http:       gout.New(),
		isDebug:    false,
		header:     make(map[string]interface{}, 0),
	}

	client.http.SetTimeout(client.timeout)
	return client

}

func (c *Httpclient) SetIsDebug(isDebug bool) *Httpclient {
	c.isDebug = isDebug
	return c
}

func (c *Httpclient) SetTimeOut(timeout time.Duration) *Httpclient {
	c.timeout = timeout
	c.http.SetTimeout(timeout)
	return c
}

func (c *Httpclient) SetRetryTimes(retryTimes int) *Httpclient {
	c.retryTimes = retryTimes
	return c
}

func (c *Httpclient) SetUrl(url string) *Httpclient {
	c.apiUrl = url
	return c
}

func (c *Httpclient) SetHeaders(headers map[string]string) *Httpclient {
	// 添加header头信息
	for name, value := range headers {
		c.header[name] = value
	}
	return c
}

func (c *Httpclient) SetHeader(name string, value string) *Httpclient {
	c.header[name] = value
	return c
}

func (c *Httpclient) GetHeaders() map[string]interface{} {
	return c.header
}

func (c *Httpclient) IsJson(isjson bool) *Httpclient {
	c.isJson = isjson
	return c
}

func (c *Httpclient) SetMethod(method Method) *Httpclient {
	c.method = method
	return c
}

func (c *Httpclient) GetDataFlow() *dataflow.DataFlow {
	var (
		request *dataflow.DataFlow
	)

	if c.method == GET {
		request = c.http.GET(c.apiUrl)
	} else if c.method == POST {
		request = c.http.POST(c.apiUrl)
	} else if c.method == PUT {
		request = c.http.PUT(c.apiUrl)
	} else if c.method == DELETE {
		request = c.http.DELETE(c.apiUrl)
	} else {
		return nil
	}

	request.Code(&c.Code).Debug(c.isDebug)
	return request
}

func (c *Httpclient) Request(repos *dataflow.DataFlow) error {
	var (
		err error
	)

	err = repos.
		F().
		Retry().
		Attempt(c.retryTimes).
		WaitTime(c.timeout).
		MaxWaitTime(MaxWaitInterval).Func(func(cg *gout.Context) error {
		if cg.Error != nil || cg.Code != 200 {
			// 是否需要一直重试待定
			if cg.Code == 429 || cg.Code == 401 {
				return filter.ErrRetry
			} else if cg.Code >= 500 {
				return filter.ErrRetry
			} else {
				return nil
			}
		}

		return nil
	}).Do()

	return err
}
