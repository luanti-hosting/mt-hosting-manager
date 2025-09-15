package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserNotification(t *testing.T) {
	repos := SetupRepos(t)

	user := &types.User{
		Name:       "somedude",
		State:      types.UserStateActive,
		Created:    time.Now().Unix(),
		ExternalID: "abc",
		Type:       types.UserTypeGithub,
		Role:       types.UserRoleUser,
	}
	assert.NoError(t, repos.UserRepo.Insert(user))

	un := &types.UserNotification{
		UserID:    user.ID,
		Timestamp: time.Now().Unix(),
		Title:     "Stuff",
		Message:   "Stuff happened!",
	}
	assert.NoError(t, repos.UserNotificationRepo.Insert(un))

	list, err := repos.UserNotificationRepo.GetUnread(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))

	assert.NoError(t, repos.UserNotificationRepo.MarkRead(list[0].ID))

	un.Seen = true
	un1, err := repos.UserNotificationRepo.GetByID(list[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, un, un1)

	list, err = repos.UserNotificationRepo.GetUnread(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))

}
