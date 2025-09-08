package service

import (
	"context"
	"yotudo/src/lib/logger"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DialogService struct {
	ctx *context.Context
}

func NewDialogService(ctx *context.Context) *DialogService {
	return &DialogService{ctx}
}

func (s *DialogService) OpenFileDialog() (string, error) {
	return runtime.OpenFileDialog(*s.ctx, runtime.OpenDialogOptions{
		Title:           "Image Selector",
		ShowHiddenFiles: false,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Összes Kép",
				Pattern:     "*.jpeg;*.jpg;*.png;*.webp",
			},
			{
				DisplayName: "Képek (.png; .jpg; .jpeg)",
				Pattern:     "*.jpeg;*.jpg;*.png",
			},
			{
				DisplayName: "Web Képek (.webp)",
				Pattern:     "*.webp",
			},
		},
		CanCreateDirectories: false,
	})
}

func (s *DialogService) OpenConfirmationDialog(title string, message string) bool {
	res, err := runtime.MessageDialog(*s.ctx, runtime.MessageDialogOptions{
		Title:         title,
		Message:       message,
		Type:          runtime.QuestionDialog,
		DefaultButton: "No",
	})
	if err != nil {
		logger.Error(err)
		return false
	}

	switch res {
	case "Yes":
		return true
	case "No":
		fallthrough
	default:
		return false
	}
}
