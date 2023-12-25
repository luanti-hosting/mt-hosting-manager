package core

import (
	"fmt"
	"mt-hosting-manager/types"
	"time"

	"github.com/minetest-go/oauth"

	"github.com/dchest/captcha"
	"golang.org/x/crypto/bcrypt"
)

func (c *Core) RegisterOauth(user_info *oauth.OauthUserInfo) (*types.User, error) {
	new_name := user_info.Name
	user, err := c.repos.UserRepo.GetByName(new_name)
	if err != nil {
		return nil, err
	}
	if user != nil {
		// a user with that name and different auth provider already exists
		// add a suffix
		i := 2
		for {
			new_name = fmt.Sprintf("%s_%d", user_info.Name, i)
			user, err = c.repos.UserRepo.GetByName(new_name)
			if err != nil {
				return nil, err
			}
			if user == nil {
				break
			}
			i++
			if i > 50 {
				return nil, fmt.Errorf("username register iterations exceeded %d tries, aborting", i)
			}
		}
	}

	user = &types.User{
		Created:    time.Now().Unix() * 1000,
		Name:       new_name,
		Type:       types.UserType(user_info.Provider),
		Role:       types.UserRoleUser,
		ExternalID: user_info.ExternalID,
	}
	err = c.repos.UserRepo.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Core) RegisterLocal(rr *types.RegisterRequest) (*types.User, *types.CheckRegisterResponse, error) {
	resp, err := c.CheckRegisterLocal(rr)
	if err != nil {
		return nil, nil, err
	}
	if !resp.Success {
		return nil, resp, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user := &types.User{
		Created: time.Now().Unix() * 1000,
		Name:    rr.Name,
		Type:    types.UserTypeLocal,
		Role:    types.UserRoleUser,
		Hash:    string(hash),
	}
	err = c.repos.UserRepo.Insert(user)
	if err != nil {
		return nil, nil, err
	}

	return user, resp, nil
}

func (c *Core) CheckRegisterLocal(rr *types.RegisterRequest) (*types.CheckRegisterResponse, error) {
	resp := &types.CheckRegisterResponse{
		Success: true,
	}

	if !ValidateName(rr.Name) || rr.Name == "" {
		resp.ErrInvalidUsername = true
		resp.Success = false
	}

	existing_user, err := c.repos.UserRepo.GetByName(rr.Name)
	if err != nil {
		return nil, err
	}
	if existing_user != nil {
		resp.ErrUsernameTaken = true
		resp.Success = false
	}

	if len(rr.Password) < 6 {
		resp.ErrPasswordTooShort = true
		resp.Success = false
	}

	if !captcha.VerifyString(rr.CaptchaID, rr.CaptchaAnswer) {
		resp.ErrCaptcha = true
		resp.Success = false
	}

	return resp, nil
}
