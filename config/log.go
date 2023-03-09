package config

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

// SelfLogConfig 日志配置
type LogConfig struct {
	LogPattern string `yaml:"log_pattern" mapstructure:"log_pattern"` // 日志输出标准，终端输出/文件输出
	LogPath    string `yaml:"log_path" mapstructure:"log_path"`       // 日志路径
	SaveDays   uint   `yaml:"save_days" mapstructure:"save_days"`     // 日志保存天数
	Level      string `yaml:"level" mapstructure:"level"`             // 日志级别
	PidPath    string `yaml:"pid_path" mapstructure:"pid_path"`
}

// logFormatter 日志数据格式化
type logFormatter struct {
	log.TextFormatter
}

// Format 自定义日志输出格式
func (c *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	prettyCaller := func(frame *runtime.Frame) string {
		_, fileName := filepath.Split(frame.File)
		return fmt.Sprintf("%s:%d", fileName, frame.Line)
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("[%s] %s", entry.Time.Format(c.TimestampFormat), // 输出日志时间
		strings.ToUpper(entry.Level.String())))
	if entry.HasCaller() {
		b.WriteString(fmt.Sprintf("[%s]", prettyCaller(entry.Caller))) // 输出日志所在文件，行数位置
	}
	b.WriteString(fmt.Sprintf(" %s\n", entry.Message)) // 输出日志内容
	return b.Bytes(), nil
}
