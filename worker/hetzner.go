package worker

import (
	"fmt"
	"mt-hosting-manager/notify"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) HetznerServerTypeCheckJob() {
	for w.running.Load() {
		err := w.core.CheckHetznerServerTypes()
		if err != nil {
			logrus.WithError(err).Error("hetzner server type check failed")
			notify.Send(&notify.NtfyNotification{
				Title:    fmt.Sprintf("Hetzner server type check failed: %v", err),
				Message:  fmt.Sprintf("Hetzner server type check failed: %v", err),
				Priority: 4,
				Tags:     []string{"computer", "warning"},
			}, true)
		} else {
			logrus.Info("hetzner server type check succeeded")
		}

		time.Sleep(time.Hour * 6)
	}
}
