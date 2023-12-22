package domain


import (
  "database/sql"
  "github.com/alphamystic/odin/lib/utils"
)

/*
  * Defines a domain for any service that exports all DBFunctions at a go
    This allows us to call th functions from anywhere i.e:
      1. You can leverage each service to do it's own DB Functions
      2. You can leverage each handler.
    Always remember to initiate your own domain
*/

type Domain struct {
  Dbs *sql.DB
  LogChan chan *utils.Logger
}

func NewDomain(dbs *sql.DB, logChan *utils.Logger,max,min int) *Domain {
  dbs.SetMaxOpenConns(max)
	dbs.SetMaxIdleConns(min)
  return &Domain {
    Dbs: dbs,
    Logger: logChan,
  }
}


func (d *Domain) GetConnection(ctx context.Context) (*sql.Conn,error){
  return  conn.Conn(ctx)
}
