package config

type PGXDbConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	DefaultDb string `json:"default_db"`
	MaxConn   int    `json:"max_conn"`
}

type PQDbConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	DefaultDb string `json:"default_db"`
}

type SQLXDbConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	DefaultDb string `json:"default_db"`
}
