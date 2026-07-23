package scheduler

import (
	"context"
	"fmt"
	"github.com/yourname/hunter-system/internal/db"
	"github.com/yourname/hunter-system/internal/quest"
	"github.com/yourname/hunter-system/internal/user"
	"log"
	"time"
)

type Scheduler struct {
	userService  *user.Service
	questService *quest.Service
}

func NewScheduler(userService *user.Service, questService *quest.Service) *Scheduler {
	return &Scheduler{userService: userService, questService: questService}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(3600 * time.Second)

	s.tick(ctx)

	for {
		select {
		case <-ticker.C:
			s.tick(ctx)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	fmt.Println("Tick")
	users, err := s.userService.ListUser(ctx)
	if err != nil {
		fmt.Println("Failed to get users")
	}
	for _, u := range users {
		if needsReset(u, time.Now()) {
			s.resetUser(ctx, u)
		}
	}
}

func (s *Scheduler) resetUser(ctx context.Context, u db.User) {
	err := s.questService.FinishNotCompleteTask(ctx, u.ID)
	if err != nil {
		log.Printf("Failed to finish task: %v", err)
	}
	_, err = s.questService.GenerateDailyQuests(ctx, u.ID)
	if err != nil {
		log.Printf("Failed to generate daily quests: %v", err)
	}
	s.userService.UpdateUserReset(ctx, u)
	fmt.Printf("Generate quests for user %d", u.ID)
}

func needsReset(u db.User, now time.Time) bool {
	loc, err := time.LoadLocation(u.Timezone)
	if err != nil {
		log.Printf("user %d: bad timezone %q, fallback UTC", u.ID, u.Timezone)
		loc = time.UTC
	}
	today := localDate(now, loc)

	return u.LastResetDate.Before(today)
}

func localDate(t time.Time, loc *time.Location) time.Time {
	lt := t.In(loc)
	return time.Date(lt.Year(), lt.Month(), lt.Day(), 0, 0, 0, 0, time.UTC)
}
