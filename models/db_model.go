package models

type SqlConnection struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type MongoConnection struct {
	URI      string
	Database string
}

type RedisConnection struct {
	Host     string
	Port     string
	Password string
	Database int
}

type CassandraConnection struct {
	Keyspace        string
	Host            string
	ProtocolVersion int
}

type ElasticConnection struct{}
