package logger

import "fmt"

type Logger struct{}

func (l *Logger) LogErr(args ...interface{}) {
	fmt.Println("ERROR ", args)
}
