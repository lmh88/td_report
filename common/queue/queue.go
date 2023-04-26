package queue

import (
	"errors"
	"td_report/vars"
)

func GetQueue(reportType string, common bool) (string, error) {
	if common == false {
		queue := map[string]string{
			vars.SB:  vars.QueueMessageSb,
			vars.SD:  vars.QueueMessageSd,
			vars.SP:  vars.QueueMessageSp,
			vars.DSP: vars.QueueMessageDsp,
		}
		if value, ok := queue[reportType]; ok {
			return value, nil
		} else {
			return "", errors.New("paramas error")
		}
	} else {

		queue := map[string]string{
			vars.SB:  vars.QueueMessageSbCommon,
			vars.SD:  vars.QueueMessageSdCommon,
			vars.SP:  vars.QueueMessageSpCommon,
			vars.DSP: vars.QueueMessageDspCommon,
		}
		if value, ok := queue[reportType]; ok {
			return value, nil
		} else {
			return "", errors.New("paramas error")
		}
	}

}
