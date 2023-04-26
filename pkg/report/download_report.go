package report

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"td_report/app/bean"
	"td_report/pkg/logger"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

func DownloadReport(profileToken *bean.ProfileToken, location string, ctx context.Context) ([]byte, error) {

	headers, err := fakeHeaders(profileToken.ProfileId, profileToken.RefreshToken, profileToken.ClientId, profileToken.ClientSecret)
	if err != nil {
		return nil, err
	}

	num := 5
	for i := 0; i < num; i++ {
		resp, err := requests.Get(location, requests.WithHeaders(headers), requests.WithTimeout(time.Hour))
		if err != nil || resp == nil {
			errstr := err.Error()
			if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
				logger.Logger.ErrorWithContext(ctx, "network error")
				time.Sleep(time.Duration((i+1)*20) * time.Second)
				continue
			}

			logger.Logger.ErrorWithContext(ctx, "get download error:", err.Error())
			return nil, err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("报表下载数据结果，statusCode:%d", resp.StatusCode))

		if resp.StatusCode != http.StatusOK {
			logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("download error status,code:%s, body:%s", resp.StatusCode, string(resp.Body)))
			if resp.StatusCode >= 500 {
				time.Sleep(time.Duration((i+1)*20) * time.Second)
				continue
			}
			return nil, fmt.Errorf("download report resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
		}

		return resp.Body, nil
	}

	logger.Logger.ErrorWithContext(ctx, "get download error try 5 time:")
	return nil, errors.New("download file error")
}

func DownloadReportWithoutAuth(location string, ctx context.Context) ([]byte, error) {
	num := 5
	for i := 0; i < num; i++ {
		resp, err := requests.Get(location, requests.WithTimeout(time.Hour))
		if err != nil || resp == nil {
			if strings.Contains(err.Error(), vars.Timeouttxt) || strings.Contains(err.Error(), vars.Timeouttxt1) {
				logger.Logger.ErrorWithContext(ctx, "network error")
				time.Sleep(time.Duration((i+1)*20) * time.Second)
				continue
			}

			logger.Logger.ErrorWithContext(ctx, "get download error:", err.Error())
			time.Sleep(time.Duration((i+1)*20) * time.Second)
			return nil, err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("报表下载数据结果，statusCode:%d", resp.StatusCode))
		if resp.StatusCode != http.StatusOK {
			logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("download error status,code:%s, body:%s", resp.StatusCode, string(resp.Body)))
			if resp.StatusCode >= 500 {
				time.Sleep(time.Duration((i+1)*20) * time.Second)
				continue
			}
			return nil, fmt.Errorf("download report resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
		}

		return resp.Body, nil
	}

	logger.Logger.ErrorWithContext(ctx, "get download error try 5 time:")
	return nil, errors.New("download error ")

}
