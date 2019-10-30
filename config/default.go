package config

var (
	defaultHome                  string
	defaultTimeZone              = "utc"
	defaultExecutor              = "SequentialExecutor"
	defaultSQLConn               string
	defaultParrallelism          = 32
	defaultDagConcurrency        = 16
	defaultMaxActiveRunsPerDag   = 16
	defaultTaskRunner            = "StandardTaksRunner"
	defaultPort                  = 3000
	defaultWorkers               = 4
	defaultWorkerClass           = "sync"
	defaultDagOrientation        = "LR"
	defaultJobHeartBeatSec       = 5
	defaultSchedulerheartBeatSec = 5
	defaultNumRuns               = -1
)
