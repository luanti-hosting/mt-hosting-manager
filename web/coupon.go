package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type RedeemStatus struct {
	Success bool `json:"success"`
}

func (a *Api) RedeemCoupon(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	code := vars["code"]

	coupon, err := a.repos.CouponRepo.GetByCode(code)
	if err != nil {
		SendError(w, 500, fmt.Errorf("get coupon error: %v", err))
		return
	}
	if coupon == nil {
		// no such code
		Send(w, RedeemStatus{Success: false}, nil)
		return
	}
	if coupon.UseCount >= coupon.UseMax {
		// max reached
		Send(w, RedeemStatus{Success: false}, nil)
		return
	}
	now := time.Now().Unix()
	if now > coupon.ValidUntil || now < coupon.ValidFrom {
		// not valid yet or expired
		Send(w, RedeemStatus{Success: false}, nil)
		return
	}

	already_redeemed, err := a.repos.CouponRepo.IsRedeemed(coupon.ID, c.UserID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("get redeemed_coupon error: %v", err))
		return
	}
	if already_redeemed {
		// user already redeemed coupon
		Send(w, RedeemStatus{Success: false}, nil)
		return
	}

	// everything ok here!
	err = a.repos.UserRepo.AddBalance(c.UserID, coupon.Value)
	if err != nil {
		SendError(w, 500, fmt.Errorf("addbalance error: %v", err))
		return
	}

	// set redeemed
	err = a.repos.CouponRepo.Redeem(coupon.ID, c.UserID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("redeem error: %v", err))
		return
	}

	// increment use count
	coupon.UseCount++
	err = a.repos.CouponRepo.Update(coupon)
	if err != nil {
		SendError(w, 500, fmt.Errorf("update usecount error: %v", err))
		return
	}

	a.core.AddAuditLog(&types.AuditLog{
		Type:   types.AuditLogCouponRedeemed,
		UserID: c.UserID,
		Amount: &coupon.Value,
	})

	// all done
	Send(w, RedeemStatus{Success: true}, nil)

	notify.Send(&notify.NtfyNotification{
		Title:    fmt.Sprintf("Coupon '%s' used by %s (%.2f)", coupon.Code, c.Name, float64(coupon.Value)/100),
		Message:  fmt.Sprintf("User: %s\nEUR %.2f\nCode: %s\nUses left: %d", c.Name, float64(coupon.Value)/100, coupon.Code, coupon.UseMax-coupon.UseCount),
		Priority: 3,
		Tags:     []string{"ticket"},
	}, true)
}

func (a *Api) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	coupon := &types.Coupon{}
	err := json.NewDecoder(r.Body).Decode(coupon)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json error: %v", err))
		return
	}

	err = a.repos.CouponRepo.Insert(coupon)
	Send(w, coupon, err)
}

func (a *Api) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	coupon, err := a.repos.CouponRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, fmt.Errorf("get coupon error: %v", err))
		return
	}
	if coupon == nil {
		SendError(w, http.StatusNotFound, fmt.Errorf("coupon not found: '%s'", id))
		return
	}

	updated_coupon := &types.Coupon{}
	err = json.NewDecoder(r.Body).Decode(updated_coupon)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json error: %v", err))
		return
	}

	// apply modifiable fields
	coupon.Name = updated_coupon.Name
	coupon.ValidUntil = updated_coupon.ValidUntil
	coupon.ValidFrom = updated_coupon.ValidFrom
	coupon.UseMax = updated_coupon.UseMax

	err = a.repos.CouponRepo.Update(coupon)
	Send(w, coupon, err)
}

func (a *Api) GetCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	coupon, err := a.repos.CouponRepo.GetByID(id)
	Send(w, coupon, err)
}

func (a *Api) GetCouponUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	redeemed_coupons, err := a.repos.CouponRepo.GetRedeemedCoupons(id)
	if err != nil {
		SendError(w, 500, fmt.Errorf("get redeemed coupons error: %v", err))
		return
	}

	user_list := []*types.User{}
	for _, rc := range redeemed_coupons {
		user, err := a.repos.UserRepo.GetByID(rc.UserID)
		if err != nil {
			SendError(w, 500, fmt.Errorf("get redeemed coupons error: %v", err))
			return
		}

		user.RemoveSensitiveFields()
		user_list = append(user_list, user)
	}

	Send(w, user_list, err)
}

func (a *Api) GetCoupons(w http.ResponseWriter, r *http.Request) {
	list, err := a.repos.CouponRepo.GetAll()
	Send(w, list, err)
}
