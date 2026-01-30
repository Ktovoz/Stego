package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"stego/internal/app"
	"stego/internal/config"
	"stego/internal/log"
	"stego/internal/models"
)

type App struct {
	ctx    context.Context
	cfg    *config.Store
	tasks  *app.TaskManager
	logger *log.Store
	info   models.AppInfo
}

func NewApp() *App {
	buildHash := computeBuildHash()

	return &App{
		info: models.AppInfo{
			Name:      "stego",
			Version:   "1.0",
			BuildDate: time.Now().Format("2006-01-02"),
			BuildHash: buildHash,
			Author:    "Ktovoz",
			GitHub:    "https://github.com/Ktovoz",
		},
	}
}

func computeBuildHash() string {
	execPath, err := os.Executable()
	if err != nil {
		return "unknown"
	}

	file, err := os.Open(execPath)
	if err != nil {
		return "unknown"
	}
	defer func() { _ = file.Close() }()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "unknown"
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)[:16]
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dbPath := filepath.Join(".", "data", "config.db")
	cfgStore, err := config.NewStore(dbPath)
	if err != nil {
		runtime.LogErrorf(a.ctx, "config init failed: %v", err)
		cfgStore = config.NewInMemoryStore()
	}
	a.cfg = cfgStore
	a.tasks = app.NewTaskManager()

	logDBPath := filepath.Join(".", "data", "logs.db")
	logger, err := log.NewStore(logDBPath)
	if err != nil {
		runtime.LogErrorf(a.ctx, "logger init failed: %v", err)
	} else {
		a.logger = logger
		_ = a.logger.Add("INFO", "app", "应用程序启动", "Version: "+a.info.Version)
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.logger != nil {
		_ = a.logger.Close()
	}
}

func (a *App) GetConfig() map[string]string {
	if a.cfg == nil {
		return map[string]string{}
	}
	return a.cfg.GetAllWithDefaults()
}

func (a *App) SaveConfig(configMap map[string]string) error {
	if a.cfg == nil {
		return errors.New("config store not initialized")
	}
	err := a.cfg.SaveAll(configMap)
	if err == nil && a.logger != nil {
		_ = a.logger.Add("INFO", "config", "配置已保存", "")
	}
	return err
}

func (a *App) GetAppInfo() models.AppInfo {
	return a.info
}

func (a *App) StartEncrypt(req models.EncryptRequest) string {
	taskID := uuid.NewString()

	if a.logger != nil {
		dataSourceInfo := req.DataSourcePath
		if len(dataSourceInfo) > 50 {
			dataSourceInfo = dataSourceInfo[:47] + "..."
		}
		_ = a.logger.Add("INFO", "encrypt", "开始加密任务",
			fmt.Sprintf("任务ID: %s, 数据源: %s, 输出目录: %s", taskID, dataSourceInfo, req.OutputDir))
	}

	a.tasks.Start(taskID, func(ctx context.Context) error {
		var perf app.PerfLogger
		if a.logger != nil {
			perf = func(module, action, details string) {
				_ = a.logger.Add("INFO", module, action, details)
			}
		}
		err := app.RunEncrypt(ctx, a.cfg.GetAllWithDefaults(), req, func(p models.ProgressEvent) {
			p.TaskID = taskID
			runtime.EventsEmit(a.ctx, "encryptProgress", p)
		}, taskID, perf)

		if err != nil {
			if a.logger != nil {
				_ = a.logger.Add("ERROR", "encrypt", "加密任务失败", fmt.Sprintf("任务ID: %s, 错误: %s", taskID, err.Error()))
			}
		} else {
			if a.logger != nil {
				_ = a.logger.Add("INFO", "encrypt", "加密任务完成", fmt.Sprintf("任务ID: %s, 输出目录: %s", taskID, req.OutputDir))
			}
		}

		return err
	})
	return taskID
}

func (a *App) CancelEncrypt(taskID string) error {
	err := a.tasks.Cancel(taskID)
	if err == nil && a.logger != nil {
		_ = a.logger.Add("WARN", "encrypt", "加密任务已取消", "任务ID: "+taskID)
	}
	return err
}

func (a *App) StartDecrypt(req models.DecryptRequest) string {
	taskID := uuid.NewString()

	if a.logger != nil {
		imagePathInfo := req.ImagePath
		if len(imagePathInfo) > 50 {
			imagePathInfo = imagePathInfo[:47] + "..."
		}
		_ = a.logger.Add("INFO", "decrypt", "开始解密任务",
			fmt.Sprintf("任务ID: %s, 图片路径: %s, 输出目录: %s", taskID, imagePathInfo, req.OutputDir))
	}

	a.tasks.Start(taskID, func(ctx context.Context) error {
		var perf app.PerfLogger
		if a.logger != nil {
			perf = func(module, action, details string) {
				_ = a.logger.Add("INFO", module, action, details)
			}
		}
		err := app.RunDecrypt(ctx, a.cfg.GetAllWithDefaults(), req, func(p models.ProgressEvent) {
			p.TaskID = taskID
			runtime.EventsEmit(a.ctx, "decryptProgress", p)
		}, taskID, perf)

		if err != nil {
			if a.logger != nil {
				_ = a.logger.Add("ERROR", "decrypt", "解密任务失败", fmt.Sprintf("任务ID: %s, 错误: %s", taskID, err.Error()))
			}
		} else {
			if a.logger != nil {
				_ = a.logger.Add("INFO", "decrypt", "解密任务完成", fmt.Sprintf("任务ID: %s, 输出目录: %s", taskID, req.OutputDir))
			}
		}

		return err
	})
	return taskID
}

func (a *App) CancelDecrypt(taskID string) error {
	err := a.tasks.Cancel(taskID)
	if err == nil && a.logger != nil {
		_ = a.logger.Add("WARN", "decrypt", "解密任务已取消", "任务ID: "+taskID)
	}
	return err
}

func (a *App) StartGenerateCarrier(req models.GenerateRequest) string {
	taskID := uuid.NewString()

	if a.logger != nil {
		outputDirInfo := req.OutputDir
		if len(outputDirInfo) > 50 {
			outputDirInfo = outputDirInfo[:47] + "..."
		}
		targetMB := req.TargetBytes / (1024 * 1024)
		_ = a.logger.Add("INFO", "generate", "开始生成载体",
			fmt.Sprintf("任务ID: %s, 输出目录: %s, 目标容量: %dMB, 数量: %d", taskID, outputDirInfo, targetMB, req.Count))
	}

	a.tasks.Start(taskID, func(ctx context.Context) error {
		var perf app.PerfLogger
		if a.logger != nil {
			perf = func(module, action, details string) {
				_ = a.logger.Add("INFO", module, action, details)
			}
		}
		err := app.RunGenerateCarrier(ctx, req, func(p models.ProgressEvent) {
			p.TaskID = taskID
			runtime.EventsEmit(a.ctx, "generateProgress", p)
		}, taskID, perf)

		if err != nil {
			if a.logger != nil {
				_ = a.logger.Add("ERROR", "generate", "生成载体失败", fmt.Sprintf("任务ID: %s, 错误: %s", taskID, err.Error()))
			}
		} else {
			if a.logger != nil {
				_ = a.logger.Add("INFO", "generate", "生成载体完成", fmt.Sprintf("任务ID: %s, 数量: %d", taskID, req.Count))
			}
		}

		return err
	})
	return taskID
}

func (a *App) CancelGenerate(taskID string) error {
	err := a.tasks.Cancel(taskID)
	if err == nil && a.logger != nil {
		_ = a.logger.Add("WARN", "generate", "生成任务已取消", "任务ID: "+taskID)
	}
	return err
}

func (a *App) OpenDirectoryDialog(defaultDir string) string {
	result, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "选择目录",
		DefaultDirectory: defaultDir,
	})
	if err != nil || result == "" {
		return ""
	}
	if a.logger != nil {
		_ = a.logger.Add("INFO", "ui", "选择目录", "路径: "+result)
	}
	return result
}

func (a *App) OpenFileDialog(defaultDir string) string {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "选择文件",
		DefaultDirectory: defaultDir,
	})
	if err != nil || result == "" {
		return ""
	}
	if a.logger != nil {
		_ = a.logger.Add("INFO", "ui", "选择文件", "路径: "+result)
	}
	return result
}

func (a *App) GetLogs(level string, startTime int64, endTime int64, limit int, offset int) []log.Entry {
	if a.logger == nil {
		return []log.Entry{}
	}

	start := time.Unix(0, 0)
	if startTime > 0 {
		start = time.Unix(startTime, 0)
	}
	end := time.Now()
	if endTime > 0 {
		end = time.Unix(endTime, 0)
	}

	var logs []log.Entry
	var err error

	if level == "" || level == "ALL" {
		if startTime > 0 && endTime > 0 {
			logs, err = a.logger.GetByTimeRange(start, end, limit, offset)
		} else {
			logs, err = a.logger.Get(limit, offset)
		}
	} else {
		logs, err = a.logger.GetByLevel(level, limit, offset)
	}

	if err != nil {
		runtime.LogErrorf(a.ctx, "get logs failed: %v", err)
		return []log.Entry{}
	}

	return logs
}

func (a *App) GetLogsCount() int {
	if a.logger == nil {
		return 0
	}
	count, err := a.logger.GetCount()
	if err != nil {
		return 0
	}
	return count
}

func (a *App) ExportLogs(format string, startTime int64, endTime int64) string {
	if a.logger == nil {
		return ""
	}

	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

	var content string
	var err error

	if format == "json" {
		content, err = a.logger.ExportAsJSON(start, end)
	} else {
		content, err = a.logger.ExportAsText(start, end)
	}

	if err != nil {
		runtime.LogErrorf(a.ctx, "export logs failed: %v", err)
		return ""
	}

	if a.logger != nil {
		_ = a.logger.Add("INFO", "logs", "导出日志", fmt.Sprintf("格式: %s, 起始: %s, 结束: %s", format, start.Format("2006-01-02"), end.Format("2006-01-02")))
	}

	return content
}

func (a *App) ExportLogsToFile(format string, startTime int64, endTime int64) (string, error) {
	if a.logger == nil {
		return "", errors.New("logger not initialized")
	}

	start := time.Unix(0, 0)
	if startTime > 0 {
		start = time.Unix(startTime, 0)
	}
	end := time.Now()
	if endTime > 0 {
		end = time.Unix(endTime, 0)
	}

	ext := "txt"
	if format == "json" {
		ext = "json"
	}

	var content string
	var err error
	if format == "json" {
		content, err = a.logger.ExportAsJSON(start, end)
	} else {
		content, err = a.logger.ExportAsText(start, end)
	}
	if err != nil {
		return "", err
	}

	defaultFilename := fmt.Sprintf("stego_logs_%s.%s", time.Now().Format("2006-01-02"), ext)
	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save logs",
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{DisplayName: fmt.Sprintf("Log Files (*.%s)", ext), Pattern: fmt.Sprintf("*.%s", ext)},
		},
	})
	if err != nil {
		return "", err
	}
	if savePath == "" {
		return "", nil
	}
	if filepath.Ext(savePath) == "" {
		savePath = savePath + "." + ext
	}

	if err := os.WriteFile(savePath, []byte(content), 0o644); err != nil {
		return "", err
	}
	if a.logger != nil {
		_ = a.logger.Add("INFO", "logs", "Export logs", "Path: "+savePath)
	}
	return savePath, nil
}

func (a *App) ClearLogs() error {
	if a.logger == nil {
		return errors.New("logger not initialized")
	}

	err := a.logger.Clear()
	if err != nil {
		return err
	}

	_ = a.logger.Add("INFO", "system", "日志已清空", "")
	return nil
}

func (a *App) LogUserAction(module, action, details string) {
	if a.logger == nil {
		return
	}
	_ = a.logger.Add("INFO", module, action, details)
}
