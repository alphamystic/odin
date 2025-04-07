package domain


import (
  "context"
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
  *utils.ErrorLogger
}

func NewDomain(dbs *sql.DB,max,min int) *Domain {
  dbs.SetMaxOpenConns(max)
	dbs.SetMaxIdleConns(min)
  errorFiles := []string{"users_sql", "minions_sql", "auth_sql", "api_sql", "assets_sql", "ms_sql", "recon_sql"}
  errorLogger := utils.NewErrorLogger("./.brain/logs/sql", 066, errorFiles)
  return &Domain {
    Dbs: dbs,
    ErrorLogger: errorLogger,
  }
}


func (d *Domain) GetConnection(ctx context.Context) (*sql.Conn,error){
  return  d.Dbs.Conn(ctx)
}
