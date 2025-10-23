package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

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

	err = a.repos.ServiceTicketRepo.InsertTicket(ticket)
	Send(w, ticket, err)
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
