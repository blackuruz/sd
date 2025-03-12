package backend

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Utils struct{}

func NewUtils() *Utils {
	return &Utils{}
}

// SelectFile menggunakan file dialog dari Wails untuk memilih file
func (u *Utils) SelectFile(ctx context.Context) string {
	file, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Pilih File Video",
		Filters: []runtime.FileFilter{
			{DisplayName: "Video Files", Pattern: "*.mp4;*.avi;*.mkv"},
		},
	})
	if err != nil {
		return ""
	}
	return file
}
