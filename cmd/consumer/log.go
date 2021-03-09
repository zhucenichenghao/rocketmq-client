package main

var XRocketMQLog = &defaultLogger{}

type defaultLogger struct {
}

func (l *defaultLogger) Debug(msg string, fields map[string]interface{}) {

}

func (l *defaultLogger) Info(msg string, fields map[string]interface{}) {

}

func (l *defaultLogger) Warning(msg string, fields map[string]interface{}) {

}

func (l *defaultLogger) Error(msg string, fields map[string]interface{}) {

}

func (l *defaultLogger) Fatal(msg string, fields map[string]interface{}) {

}
