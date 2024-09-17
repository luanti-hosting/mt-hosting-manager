package worker

import (
	"fmt"
	"mt-hosting-manager/api/coinbase"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/core"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Worker struct {
	repos   *db.Repositories
	cfg     *types.Config
	hcc     *hetzner_cloud.HetznerCloudClient
	hdc     *hetzner_dns.HetznerDNSClient
	cbc     *coinbase.CoinbaseClient
	running *atomic.Bool
	core    *core.Core
	wg      *sync.WaitGroup
}

type JobExecutor func(j *types.Job) error

var executors = map[types.JobType]JobExecutor{}

func NewWorker(repos *db.Repositories, cfg *types.Config) *Worker {
	return &Worker{
		repos:   repos,
		cfg:     cfg,
		hcc:     hetzner_cloud.New(cfg.HetznerCloudKey),
		hdc:     hetzner_dns.New(cfg.HetznerApiKey, cfg.HetznerApiZoneID),
		cbc:     coinbase.New(cfg.CoinbaseKey),
		running: &atomic.Bool{},
		core:    core.New(repos, cfg),
		wg:      &sync.WaitGroup{},
	}
}

func (w *Worker) RegisterExecutors() {
	executors[types.JobTypeNodeDestroy] = w.NodeDestroy
	executors[types.JobTypeNodeSetup] = w.NodeProvision
	executors[types.JobTypeServerDestroy] = w.ServerDestroy
	executors[types.JobTypeServerSetup] = w.ServerSetup
	executors[types.JobTypeServerRestore] = w.ServerRestore
	executors[types.JobTypeServerBackup] = w.ServerBackup
}

func (w *Worker) Stop() {
	w.running.Store(false)
	w.wg.Wait()
}

func (w *Worker) Start() {
	if w.running.CompareAndSwap(false, true) {
		go w.Run()
	}
}

func (w *Worker) processNextJob() {
	err := w.repos.Gorm().Transaction(func(tx *gorm.DB) error {
		job, err := w.repos.JobRepo.GetNextJob(tx, types.JobStateRunning, time.Now().Unix())
		if err != nil {
			return fmt.Errorf("get next job error: %v", err)
		}
		if job != nil {
			w.ExecuteJob(tx, job)
		}

		return nil
	})
	if err != nil {
		logrus.WithError(err).Error("job processing error")
	}
}

func (w *Worker) Run() {
	// start collector job
	go w.CollectJob()

	// start housekeeping job
	go w.HousekeepingJob()

	// start transaction update job
	go w.TransactionUpdateJob()

	for w.running.Load() {
		go w.processNextJob()
		time.Sleep(500 * time.Millisecond)
	}
}
