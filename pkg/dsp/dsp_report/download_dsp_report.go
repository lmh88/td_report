package dsp_report

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"td_report/pkg/logger"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

func DownloadDspReport(location string, ctx context.Context) ([]byte, error) {
	var (
		num = 5
		i   = 0
	)
	for i = 0; i < num; i++ {
		resp, err := requests.Get(location, requests.WithTimeout(time.Hour))
		if err != nil || resp == nil {
			errstr := err.Error()
			if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
				logger.Logger.ErrorWithContext(ctx, "network error")
				time.Sleep(time.Duration((i+1)*20) * time.Second)
				continue
			}

			logger.Logger.ErrorWithContext(ctx, "dsp download error nil", err.Error())
			return nil, err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("dsp报表下载数据结果，StatusCode:%d", resp.StatusCode))
		if resp.StatusCode != http.StatusOK {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"desc": "dsp download error code",
				"err":  err,
				"code": resp.StatusCode,
			})
			return nil, fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
		}

		return resp.Body, nil
	}

	logger.Logger.ErrorWithContext(ctx, "dsp download error try too time")
	return nil, errors.New("dsp download error try too time")
}
