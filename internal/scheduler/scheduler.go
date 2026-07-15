package scheduler

import (
	"context"
	"fmt"
	"github.com/yourname/hunter-system/internal/db"
	"github.com/yourname/hunter-system/internal/quest"
	"github.com/yourname/hunter-system/internal/user"
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
	ticker := time.NewTicker(15 * time.Second)

	s.tick(ctx)

	for {
		select {
		case <-ticker.C:
			s.tick(ctx)
		case <-ctx.Done():
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	fmt.Println("Tick")
	users, err := s.userService.ListUser(ctx)
	if err != nil {

	}
	for _, u := range users {
		go func() {
			if u.LastResetDate.Day() < time.Now().Day() && u.LastResetDate.Month() <= time.Now().Month() {
				s.resetUser(ctx, u)
			}
		}()
	}
}

func (s *Scheduler) resetUser(ctx context.Context, u db.User) {
	s.questService.GenerateDailyQuests(ctx, u.ID)
	s.userService.UpdateUserReset(ctx, u)
	fmt.Printf("Generate quests for user %d", u.ID)
}
