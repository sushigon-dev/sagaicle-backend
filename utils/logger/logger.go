package logger

import (
	"log"
	"runtime"
)

// エラーと追加情報を受け取り、発生箇所を自動的にログ出力
func LogError(err error, context string) {
	// 呼び出し元の情報を取得
	pc, file, line, ok := runtime.Caller(1)
	funcName := "unknown"
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			funcName = f.Name()
		}
	}
	log.Printf("[ERROR] %s:%d %s() - %s: %v", file, line, funcName, context, err)
}
