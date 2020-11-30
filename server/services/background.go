package services

import (
	"log"
	"time"

	"app-pointment/server/models"
)

type saver interface {
	save() error
}

type BackgroundSaver struct {
	ticker  *time.Ticker
	service saver
}

func NewSaver(service saver) *BackgroundSaver {
	ticker := time.NewTicker(30 * time.Second)
	return &BackgroundSaver{
		ticker:  ticker,
		service: service,
	}
}

func (s *BackgroundSaver) Start() {
	log.Println("background saver started")
	for {
		select {
		case <-s.ticker.C:
			err := s.service.save()
			if err != nil {
				log.Printf("could not save records in background: %v", err)
			}
		}
	}
}

func (s *BackgroundSaver) Stop() error {
	s.ticker.Stop()
	err := s.service.save()
	if err != nil {
		return err
	}
	log.Println("background saver stopped")
	return nil
}

type HTTPNotifierClient interface {
	Notify(reminder models.Reminder) (NotificationResponse, error)
}

type snapshotManager interface {
	snapshot() Snapshot
	snapshotGrooming(notifiedReminders ...models.Reminder)
	retry(reminder models.Reminder, duration time.Duration)
}

type BackgroundNotifier struct {
	ticker    *time.Ticker
	service   snapshotManager
	completed chan models.Reminder
	Client    HTTPNotifierClient
}

func NewNotifier(notifierURI string, service snapshotManager) *BackgroundNotifier {
	ticker := time.NewTicker(1 * time.Second)
	httpClient := NewHTTPClient(notifierURI)
	return &BackgroundNotifier{
		ticker:    ticker,
		service:   service,
		completed: make(chan models.Reminder),
		Client:    httpClient,
	}
}

func (s *BackgroundNotifier) Start() {
	log.Println("background notifier started")
	for {
		select {
		case <-s.ticker.C:
			snapshot := s.service.snapshot()
			for id := range snapshot.UnCompleted {
				_, reminder := snapshot.UnCompleted.flatten(id)
				reminderTick := reminder.ModifiedAt.Add(reminder.Duration).UnixNano()
				nowTick := time.Now().UnixNano()
				deltaTick := time.Now().Add(time.Second).UnixNano()
				if reminderTick > nowTick && reminderTick < deltaTick {
					go s.notify(reminder)
				}
			}
		case r := <-s.completed:
			log.Printf("reminder with with: %d was completed\n", r.ID)
		}
	}
}

func (s *BackgroundNotifier) notify(r models.Reminder) {
	res, err := s.Client.Notify(r)
	if err != nil {
		log.Printf("could not notify reminder with id %d\n", r.ID)
		log.Printf("background http client error: %v\n", err)

	} else if res.completed {
		s.service.snapshotGrooming(r)
		s.completed <- r
		return
	}
	s.service.retry(r, res.duration)
}

func (s *BackgroundNotifier) Stop() error {
	s.ticker.Stop()
	log.Println("background notifier stopped")
	return nil
}
