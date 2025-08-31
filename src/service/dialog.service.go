package service

import (
	"context"

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
