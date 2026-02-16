package db

import (
	"mt-hosting-manager/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceTicketRepository struct {
	g *gorm.DB
}

func (r *ServiceTicketRepository) InsertTicket(n *types.ServiceTicket) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	if n.Created == 0 {
		n.Created = time.Now().Unix()
	}
	return r.g.Create(n).Error
}

func (r *ServiceTicketRepository) UpdateTicket(n *types.ServiceTicket) error {
	return r.g.Model(n).Select("*").Updates(n).Error
}

func (r *ServiceTicketRepository) InsertMessage(n *types.ServiceTicketMessage) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	if n.Timestamp == 0 {
		n.Timestamp = time.Now().Unix()
	}
	return r.g.Create(n).Error
}

func (r *ServiceTicketRepository) SearchTickets(search *types.ServiceTicketSearch) ([]*types.ServiceTicket, error) {
	g := r.g.Model(types.ServiceTicket{})

	if search.TicketID != nil {
		g = g.Where(types.ServiceTicket{ID: *search.TicketID})
	}
	if search.UserID != nil {
		g = g.Where(types.ServiceTicket{UserID: *search.UserID})
	}
	if search.UserNodeID != nil {
		g = g.Where(types.ServiceTicket{UserNodeID: search.UserNodeID})
	}
	if search.MinetestServerID != nil {
		g = g.Where(types.ServiceTicket{MinetestServerID: search.MinetestServerID})
	}
	if search.BackupID != nil {
		g = g.Where(types.ServiceTicket{BackupID: search.BackupID})
	}
	if search.State != nil {
		g = g.Where(types.ServiceTicket{State: *search.State})
	}

	g = g.Order("created desc")
	g = g.Limit(100)

	return FindMulti[types.ServiceTicket](g)
}

func (r *ServiceTicketRepository) GetMessagesByTicket(ticketID string) ([]*types.ServiceTicketMessage, error) {
	return FindMulti[types.ServiceTicketMessage](r.g.Where(types.ServiceTicketMessage{TicketID: ticketID}).Order("timestamp asc"))
}
