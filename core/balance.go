package core

import (
	"fmt"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"

	"github.com/sirupsen/logrus"
)

func (c *Core) SubtractBalance(user_id string, eurocents int64) error {
	before_user, err := c.repos.UserRepo.GetByID(user_id)
	if err != nil {
		return fmt.Errorf("could not fetch before user: %v", err)
	}

	err = c.repos.UserRepo.SubtractBalance(user_id, eurocents)
	if err != nil {
		return fmt.Errorf("could not subtract balance: %v", err)
	}

	after_user, err := c.repos.UserRepo.GetByID(user_id)
	if err != nil {
		return fmt.Errorf("could not fetch after user: %v", err)
	}

	if before_user.Balance >= before_user.WarnBalance && after_user.Balance < before_user.WarnBalance {
		// crossed warning threshold
		c.AddAuditLog(&types.AuditLog{
			Type:   types.AuditLogPaymentWarning,
			UserID: user_id,
			Amount: &after_user.Balance,
		})

		notify.Send(&notify.NtfyNotification{
			Title:    fmt.Sprintf("User %s balance warning (%.2f)", after_user.Mail, float64(after_user.Balance)/100),
			Message:  fmt.Sprintf("User: %s crossed warning threshold: EUR %.2f", after_user.Mail, float64(after_user.Balance)/100),
			Priority: 3,
			Tags:     []string{"credit_card", "warning"},
		}, true)

		err = c.SendBalanceWarning(after_user)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":            err,
				"user_id":        user_id,
				"before_balance": before_user.Balance,
				"after_balance":  after_user.Balance,
				"warn_balance":   before_user.WarnBalance,
			}).Error("could not send balance warning")
		}
	}

	if before_user.Balance >= 0 && after_user.Balance < 0 {
		// crossed zero threshold
		c.AddAuditLog(&types.AuditLog{
			Type:   types.AuditLogPaymentZero,
			UserID: user_id,
			Amount: &after_user.Balance,
		})

		notify.Send(&notify.NtfyNotification{
			Title:    fmt.Sprintf("User %s balance zero warning (%.2f)", after_user.Mail, float64(after_user.Balance)/100),
			Message:  fmt.Sprintf("User: %s crossed zero threshold: EUR %.2f", after_user.Mail, float64(after_user.Balance)/100),
			Priority: 4,
			Tags:     []string{"credit_card", "warning"},
		}, true)

		err = c.SendRemovalNotice(after_user)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":            err,
				"user_id":        user_id,
				"before_balance": before_user.Balance,
				"after_balance":  after_user.Balance,
				"warn_balance":   before_user.WarnBalance,
			}).Error("could not send removal notice")
		}

		nodes, err := c.repos.UserNodeRepo.GetByUserID(before_user.ID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":     err,
				"user_id": before_user.ID,
			}).Error("could not fetch usernodes")
		} else if len(nodes) > 0 {

			for _, node := range nodes {
				j := types.RemoveNodeJob(node)
				err = c.repos.JobRepo.Insert(j)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"err":     err,
						"node_id": node.ID,
					}).Error("could not schedule removal job")
				}
			}
		}

	}

	return nil
}
