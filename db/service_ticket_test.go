package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceTicket(t *testing.T) {
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

	ticket := &types.ServiceTicket{
		UserID: user.ID,
		Title:  "blah",
		State:  types.ServiceTicketOpen,
	}
	assert.NoError(t, repos.ServiceTicketRepo.InsertTicket(ticket))

	msg := &types.ServiceTicketMessage{
		TicketID: ticket.ID,
		Message:  "1234",
		UserID:   user.ID,
	}
	assert.NoError(t, repos.ServiceTicketRepo.InsertMessage(msg))

	list, err := repos.ServiceTicketRepo.SearchTickets(&types.ServiceTicketSearch{
		UserID: &user.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, ticket, list[0])

	msgs, err := repos.ServiceTicketRepo.GetMessagesByTicket(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(msgs))
	assert.Equal(t, msg, msgs[0])

	fakeId := "abcd"
	list, err = repos.ServiceTicketRepo.SearchTickets(&types.ServiceTicketSearch{
		UserID: &fakeId,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))

	msgs, err = repos.ServiceTicketRepo.GetMessagesByTicket(fakeId)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(msgs))
}
