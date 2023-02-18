package config

type SQLConnect struct {
	DriverName     string
	DataSourceName string
}

func NewSqlConnect(driverName string, dataSourceName string) SQLConnect {
	return SQLConnect{
		DriverName:     driverName,
		DataSourceName: dataSourceName,
	}
}

var SqlConnect = NewSqlConnect("sqlite3", "internal/repository/files.db")
