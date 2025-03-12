package backend

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
)

type StreamManager struct {
	mu        sync.Mutex
	processes map[string]*exec.Cmd
}

func NewStreamManager() *StreamManager {
	return &StreamManager{processes: make(map[string]*exec.Cmd)}
}

func (sm *StreamManager) StartStream(ctx context.Context, config StreamConfig) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.processes[config.StreamKey]; exists {
		return fmt.Errorf("stream already running")
	}

	var args []string
	if config.Quality == "Default" {
		args = []string{
			"-re", "-stream_loop", "-1",
			"-i", config.VideoFile,
			"-c", "copy", "-preset", "ultrafast",
			"-g", "120",
			"-f", "flv", "rtmp://a.rtmp.youtube.com/live2/" + config.StreamKey,
			"-err_detect", "ignore_err", "-xerror",
		}
	} else {
		args = []string{
			"-re", "-stream_loop", "-1",
			"-i", config.VideoFile,
			"-s", "1280x720", "-c:v", "libx264", "-preset", "veryfast",
			"-maxrate", "2500k", "-bufsize", "5000k", "-b:v", "1500k",
			"-pix_fmt", "yuv420p", "-g", "60",
			"-c:a", "aac", "-b:a", "128k", "-ar", "44100",
			"-f", "flv", "rtmp://a.rtmp.youtube.com/live2/" + config.StreamKey,
			"-err_detect", "ignore_err", "-xerror",
		}
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	sm.processes[config.StreamKey] = cmd

	return cmd.Start()
}

func (sm *StreamManager) StopStream(ctx context.Context, streamKey string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if cmd, exists := sm.processes[streamKey]; exists {
		err := cmd.Process.Kill()
		delete(sm.processes, streamKey)
		return err
	}
	return fmt.Errorf("stream not found")
}
