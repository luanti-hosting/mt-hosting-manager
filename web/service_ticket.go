package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) CreateTicket(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	ticket := &types.ServiceTicket{}
	err := json.NewDecoder(r.Body).Decode(ticket)
	if err != nil {
		SendError(w, http.StatusInternalServerError, fmt.Errorf("json decode error: %v", err))
		return
	}
	ticket.UserID = c.UserID
	ticket.ID = uuid.NewString()
	ticket.Created = time.Now().Unix()
	ticket.State = types.ServiceTicketOpen

	err = a.repos.ServiceTicketRepo.InsertTicket(ticket)
	Send(w, ticket, err)
}

func (a *Api) UpdateTicket(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	ticketID := vars["ticket_id"]

	ticket := &types.ServiceTicket{}
	err := json.NewDecoder(r.Body).Decode(ticket)
	if err != nil {
		SendError(w, http.StatusInternalServerError, fmt.Errorf("json decode error: %v", err))
		return
	}

	tickets, err := a.repos.ServiceTicketRepo.SearchTickets(&types.ServiceTicketSearch{
		TicketID: &ticketID,
	})
	if err != nil {
		SendError(w, http.StatusInternalServerError, fmt.Errorf("ticket search error: %v", err))
		return
	}
	if len(tickets) != 1 {
		SendError(w, http.StatusNotFound, fmt.Errorf("ticket not found"))
		return
	}
	saved_ticket := tickets[0]
	if c.Role != types.UserRoleAdmin && ticket.UserID != c.UserID {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	if ticket.State != saved_ticket.State {
		// State change, notify
		click_url := fmt.Sprintf("%s#/ticket/%s", a.cfg.BaseURL, ticket.ID)
		notify.Send(&notify.NtfyNotification{
			Title:    fmt.Sprintf("Service ticket state change: %s", ticket.ID),
			Message:  fmt.Sprintf("Ticket changed from %s to %s", saved_ticket.State, ticket.State),
			Click:    &click_url,
			Priority: 4,
			Tags:     []string{"ticket"},
		}, true)
	}

	saved_ticket.State = ticket.State
	err = a.repos.ServiceTicketRepo.UpdateTicket(saved_ticket)
	Send(w, saved_ticket, err)
}

func (a *Api) CreateTicketMessage(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	msg := &types.ServiceTicketMessage{}
	err := json.NewDecoder(r.Body).Decode(msg)
	if err != nil {
		SendError(w, http.StatusInternalServerError, fmt.Errorf("json decode error: %v", err))
		return
	}
	tickets, err := a.repos.ServiceTicketRepo.SearchTickets(&types.ServiceTicketSearch{
		TicketID: &msg.TicketID,
	})
	if err != nil {
		SendError(w, http.StatusInternalServerError, fmt.Errorf("ticket search error: %v", err))
		return
	}
	if len(tickets) != 1 {
		SendError(w, http.StatusNotFound, fmt.Errorf("ticket not found"))
		return
	}
	ticket := tickets[0]
	if c.Role != types.UserRoleAdmin && ticket.UserID != c.UserID {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	click_url := fmt.Sprintf("%s#/ticket/%s", a.cfg.BaseURL, ticket.ID)
	notify.Send(&notify.NtfyNotification{
		Title:    fmt.Sprintf("Service ticket created: %s", ticket.ID),
		Message:  fmt.Sprintf("Title: '%s' from user: %s", ticket.Title, ticket.UserID),
		Click:    &click_url,
		Priority: 4,
		Tags:     []string{"ticket", "heavy_plus_sign"},
	}, true)

	msg.ID = uuid.NewString()
	msg.Timestamp = time.Now().Unix()
	msg.UserID = c.UserID
	err = a.repos.ServiceTicketRepo.InsertMessage(msg)
	Send(w, msg, err)
}

func (a *Api) SearchTickets(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	search := &types.ServiceTicketSearch{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, http.StatusInternalServerError, fmt.Errorf("json decode error: %v", err))
		return
	}
	if c.Role != types.UserRoleAdmin {
		search.UserID = &c.UserID
	}
	list, err := a.repos.ServiceTicketRepo.SearchTickets(search)
	Send(w, list, err)
}

func (a *Api) GetTicketMessages(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	ticketID := vars["ticket_id"]

	tickets, err := a.repos.ServiceTicketRepo.SearchTickets(&types.ServiceTicketSearch{
		TicketID: &ticketID,
	})
	if err != nil {
		SendError(w, http.StatusInternalServerError, fmt.Errorf("ticket search error: %v", err))
		return
	}
	if len(tickets) != 1 {
		SendError(w, http.StatusNotFound, fmt.Errorf("ticket not found"))
		return
	}
	ticket := tickets[0]
	if c.Role != types.UserRoleAdmin && ticket.UserID != c.UserID {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	list, err := a.repos.ServiceTicketRepo.GetMessagesByTicket(ticketID)
	Send(w, list, err)
}
