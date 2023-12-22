package definers

import (
  "fmt"
  "database/sql"
  "github.com/go-sql-driver/mysql"
)

// Initialize database connection for the given domain
type DBConfig struct {
  Username string
  Password string
  DBName string
  Host string
}

// Initiate a new MysqlDB COnnector
func NewMySQLConnector(db_config *DBConfig) (*sql.DB,error) {
  db,err := sql.Open("mysql",fmt.Sprintf("%s:%s@tcp(%s)/%s",&dbconfig.Username,&dbconfig.Password,&dbconfig.Host,&dbconfig.DBName))
  if err != nil{
    return nil,fmt.Errorf("Error creating new Connector: %v",err)
  }
  return db,nil
}

// Initialize a db configurator
func IntitializeConnector(username,pass,host,dbname string)*DBConfig{
  return &DBConfig{
    UserName:username,
    Password: pass,
    DBName: dbname,
    host: host.
  }
}
