package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

var (
	emptyErr = errors.Wrap(sql.ErrNoRows, "empty error")
	sqlErr   = errors.Wrap(sql.ErrNoRows, "sql error")
)

// dao层的sqlErr 需要wrap上抛给application
// 理由：因为sqlErr本身只是sql包中一个sentinel error
// 如果直接将此error抛给application，会有以下几点问题:
// 1. 堆栈信息没有携带application dao层的错误堆栈, 无法获取更多有用的上下文信息
// 2. 对于application而言应该定义包含其自身业务含义的error, sql error只能作为底层的cause error, 在做application全局异常处理时
//    可以对dao层上抛到application的error作类型划分, 通过errors.Is() 根据不同的error类型来做不同的业务处理
//   一般思路是使用pkg/errors wrap sql.ErrNoRows, 然后抛到application做全局异常处理
func dao() error {
	return errors.Wrap(emptyErr, "emptyErr")
}

func application() {
	if e := dao(); e != nil {
		panic(e)
	}
}

func main() {
	defer func() {
		// global recover
		if es := recover(); es != nil {
			switch {
			case errors.Is(es.(error), emptyErr):
				// TODO do something
				fmt.Printf("%+v\n", sql.ErrNoRows)
				fmt.Println("====================================")
				fmt.Printf("%+v\n", es)
				fmt.Println("====================================")
				fmt.Println("do something")
			}
		}
	}()
	application()
}
