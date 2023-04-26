package tool

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"log"
	"os/exec"
	"reflect"
)

// IsEmpty 判断一个对象是否为空
func IsEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int:
		return value == 0
	case reflect.String:
		return value == ""
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	default:
		return value == nil
	}
}

func RunCommand(path, name string, arg ...string) (msg string, err error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Dir = path
	err = cmd.Run()
	log.Println(cmd.Args)
	if err != nil {
		msg = fmt.Sprint(err) + ": " + stderr.String()
		err = errors.New(msg)
		log.Println("err", err.Error(), "cmd", cmd.Args)
	}

	g.Log().Println(out.String())
	return
}

