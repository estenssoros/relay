package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Core struct {
	RelayHome               string `yaml:"relay_home" json:"relay_home"`
	DefaultTimeZone         string `yaml:"default_time_zone" json:"default_time_zone"`
	Executor                string `yaml:"executor" json:"executor"`
	SQLConn                 string `yaml:"sql_conn" json:"sql_conn"`
	Parallelism             int    `yaml:"parallelism" json:"parallelism"`
	DagConcurrency          int    `yaml:"dag_concurrency" json:"dag_concurrency"`
	DagsArePausedAtCreation bool   `yaml:"dags_are_paused_at_creation" json:"dags_are_paused_at_creation"`
	MaxActiveRunsPerDag     int    `yaml:"max_active_runs_per_dag" json:"max_active_runs_per_dag"`
	CipherKey               string `yaml:"cipher_key" json:"cipher_key"`
	DagBagImportTimeout     int    `yaml:"dag_bag_import_timeout" json:"dag_bag_import_timeout"`
	TaskRunner              string `yaml:"task_runner" json:"task_runner"`
}

type DBCreds struct {
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Host     string `yaml:"host" json:"host"`
	Database string `yaml:"database" json:"database"`
	Port     int    `yaml:"port" json:"port"`
	Flavor   string `yaml:"flavor" json:"flavor"`
}

type Webserver struct {
	Host           string `yaml:"host" json:"host"`
	Port           int64  `yaml:"port" json:"port"`
	Workers        int    `yaml:"workers" json:"workers"`
	WorkerClass    string `yaml:"worker_class" json:"worker_class"`
	DagOrientation string `yaml:"dag_orientation" json:"dag_orientation"`
}

type Scheduler struct {
	JobHeartBeatSec               int  `yaml:"job_heart_beat_sec" json:"job_heart_beat_sec"`
	SchedulerHeartBeatSec         int  `yaml:"scheduler_heartbeat_sec" json:"scheduler_heartbeat_sec"`
	NumRuns                       int  `yaml:"num_runs" json:"num_runs"`
	ProcessPollInterval           int  `yaml:"process_poll_interval" json:"process_poll_interval"`
	MinFileProcessInterval        int  `yaml:"min_file_process_interval" json:"min_file_process_interval"`
	DagDirListInterval            int  `yaml:"dag_dir_list_interval" json:"dag_dir_list_interval"`
	SchedulerHealthCheckThreshold int  `yaml:"scheduler_health_check_threshold" json:"scheduler_health_check_threshold"`
	CatchupByDefault              bool `yaml:"catchup_by_default" json:"catchup_by_default"`
}

type Config struct {
	Core      `yaml:"core" json:"core"`
	DBCreds   `yaml:"db_creds" json:"db_creds"`
	Webserver `yaml:"webserver" json:"webserver"`
	Scheduler `yaml:"scheduler" json:"scheduler"`
	Error     error
}

func (c Config) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

var configFile = "relay.yaml"

var DefaultConfig *Config

func init() {
	config, err := Load()
	config.Error = err
	DefaultConfig = config
}

func Load() (*Config, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, errors.Wrap(err, "homedir")
	}
	f, err := os.Open(filepath.Join(homeDir, "relay", configFile))
	if err != nil {
		return nil, errors.Wrap(err, "readfile")
	}
	defer f.Close()
	config := &Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return config, nil
}

func CipherKeyBytes() ([]byte, error) {
	if DefaultConfig.Error != nil {
		return nil, DefaultConfig.Error
	}
	return []byte(DefaultConfig.Core.CipherKey), nil
}
