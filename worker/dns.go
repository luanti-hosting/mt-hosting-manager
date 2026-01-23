package worker

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/sirupsen/logrus"
)

func (w *Worker) CreateOrUpdateDNSRecord(t hcloud.ZoneRRSetType, name, value string) (*hcloud.ZoneRRSet, error) {

	existing_record, _, err := w.hc.Zone.GetRRSetByNameAndType(context.Background(), &hcloud.Zone{Name: w.cfg.HetznerZoneName}, name, t)
	if err != nil {
		return nil, fmt.Errorf("get record error type=%s, name=%s, value=%s: %v", t, name, value, err)
	}

	if existing_record == nil {
		// create new record
		logrus.WithFields(logrus.Fields{
			"value": value,
			"name":  name,
			"type":  t,
		}).Info("Creating Record")

		record, _, err := w.hc.Zone.CreateRRSet(context.Background(), &hcloud.Zone{Name: w.cfg.HetznerZoneName}, hcloud.ZoneRRSetCreateOpts{
			Name: name,
			Type: t,
			TTL:  hcloud.Ptr(300),
			Records: []hcloud.ZoneRRSetRecord{
				{Value: value},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("create record error type=%s, name=%s, value=%s: %v", t, name, value, err)
		}
		return record.RRSet, nil
	} else {
		// update existing
		if len(existing_record.Records) != 1 {
			return nil, fmt.Errorf("invalid record count: %d", len(existing_record.Records))
		}
		existing_record.Records[0].Value = value
		_, _, err = w.hc.Zone.UpdateRRSet(context.Background(), existing_record, hcloud.ZoneRRSetUpdateOpts{})
		if err != nil {
			return nil, fmt.Errorf("update record error type=%s, name=%s, value=%s: %v", t, name, value, err)
		}
		return existing_record, nil
	}
}
