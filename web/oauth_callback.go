package web

import (
	"net/http"

	"github.com/minetest-go/oauth"
)

func (api *Api) OauthCallback(w http.ResponseWriter, r *http.Request, user_info *oauth.OauthUserInfo) error {

	user, err := api.repos.UserRepo.GetByNameAndExternalID(user_info.Name, user_info.ExternalID)
	if err != nil {
		return err
	}

	if user == nil {
		user, err = api.core.RegisterOauth(user_info)
		if err != nil {
			return err
		}
	}

	err = api.loginUser(w, user)
	if err != nil {
		return err
	}

	target := api.cfg.BaseURL + "/#/profile"
	http.Redirect(w, r, target, http.StatusSeeOther)

	return nil
}
