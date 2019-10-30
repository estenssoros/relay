package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Core config for relay
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
	TaskRunner              string `yaml:"task_runner" json:"task_runner"`
}

// DBCreds creds for database
type DBCreds struct {
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Host     string `yaml:"host" json:"host"`
	Database string `yaml:"database" json:"database"`
	Port     int    `yaml:"port" json:"port"`
	Flavor   string `yaml:"flavor" json:"flavor"`
}

// Webserver webserer config
type Webserver struct {
	Port           int64  `yaml:"port" json:"port"`
	Workers        int    `yaml:"workers" json:"workers"`
	WorkerClass    string `yaml:"worker_class" json:"worker_class"`
	DagOrientation string `yaml:"dag_orientation" json:"dag_orientation"`
}

// Scheduler scheduler config
type Scheduler struct {
	JobHeartBeatSec       int `yaml:"job_heart_beat_sec" json:"job_heart_beat_sec"`
	SchedulerHeartBeatSec int `yaml:"scheduler_heartbeat_sec" json:"scheduler_heartbeat_sec"`
	NumRuns               int `yaml:"num_runs" json:"num_runs"`
}

// Config holds all configs
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

// DefaultConfig to be used by relay process
var DefaultConfig *Config

func init() {
	config, err := Load()
	config.Error = err
	DefaultConfig = config
}

// Load reads in a config
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

// CipherKeyBytes returns bytes of cipher key
func CipherKeyBytes() ([]byte, error) {
	if DefaultConfig.Error != nil {
		return nil, DefaultConfig.Error
	}
	return []byte(DefaultConfig.Core.CipherKey), nil
}

func newCipherKey() string {
	return ``
}

func createConfig(path string) error {
	config := &Config{
		Core: Core{
			RelayHome:               path,
			DefaultTimeZone:         "utc",
			Executor:                defaultExecutor,
			SQLConn:                 defaultSQLConn,
			Parallelism:             defaultParrallelism,
			DagConcurrency:          defaultDagConcurrency,
			DagsArePausedAtCreation: false,
			MaxActiveRunsPerDag:     defaultMaxActiveRunsPerDag,
			CipherKey:               newCipherKey(),
			TaskRunner:              defaultTaskRunner,
		},
		DBCreds: DBCreds{
			User:     "",
			Password: "",
			Host:     "",
			Database: "",
			Port:     0,
			Flavor:   "",
		},
		Webserver: Webserver{
			Port:           int64(defaultPort),
			Workers:        defaultWorkers,
			WorkerClass:    defaultWorkerClass,
			DagOrientation: defaultDagOrientation,
		},
		Scheduler: Scheduler{
			JobHeartBeatSec:       defaultJobHeartBeatSec,
			SchedulerHeartBeatSec: defaultSchedulerheartBeatSec,
			NumRuns:               defaultNumRuns,
		},
		Error: nil,
	}
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "create file")
	}
	defer f.Close()
	ym, _ := yaml.Marshal(config)
	if _, err := f.Write(ym); err != nil {
		return errors.Wrap(err, "file write")
	}
	return nil
}

// CreateIfNotExists creates a config if it doesn't exist
func CreateIfNotExists() error {
	homeDir, err := homedir.Dir()
	if err != nil {
		return errors.Wrap(err, "homedir")
	}
	dir := filepath.Join(homeDir, "relay")
	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, os.ModeDir); err != nil {
			return errors.Wrap(err, "mkdir")
		}
	}
	pathToConfig := filepath.Join(dir, configFile)
	if _, err := os.Stat(pathToConfig); err != nil {
		if err := createConfig(pathToConfig); err != nil {
			return errors.Wrap(err, "create config")
		}
	}
	return nil
}
