package poll

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

type Handler func(ctx context.Context, filename string) error

type CSVPoll struct {
	location string
}

func NewCSVPoll(location string) *CSVPoll {
	return &CSVPoll{
		location: location,
	}
}

func (c *CSVPoll) Poll(ctx context.Context, handler Handler) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() { _ = watcher.Close() }()

	if err = watcher.Add(c.location); err != nil {
		logrus.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				logrus.Fatal("watcher error occurred")
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				logrus.Infof("The new file was added: %s", event.Name)
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := handler(ctx, event.Name); err != nil {
					logrus.Error(err.Error())
				}
			}

			// we do not care about any other events
		case err, ok := <-watcher.Errors:
			if !ok {
				logrus.Fatal("watcher error occurred")
			}

			logrus.Error(err)
		}
	}
}
