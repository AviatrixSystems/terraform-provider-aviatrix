package goaviatrix

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"path/filepath"
	"runtime"
	"time"
)

func LogDebug(ctx context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	fileInfo := "unknown"
	if ok {
		fileInfo = filepath.Base(file)
	}

	message := fmt.Sprintf(format, args...)

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	fullMessage := fmt.Sprintf("[AVXDBG] %s %s:%d - %s", timestamp, fileInfo, line, message)

	tflog.Debug(ctx, fullMessage)
}
