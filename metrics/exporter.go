package metrics

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"time"
)

func (r *Registry) InfluxDb(interval time.Duration, cfg InfluxDbConfig) {
	var c client.Client
	sendTicker := time.Tick(interval)
	pingTicker := time.Tick(time.Second * 5)
	var err error
	for c, err = Connect(cfg); err != nil; {
		log.Printf("metrics: %v %v", time.Now(), err)
		time.Sleep(interval)
	}
	for {
		select {
		case <-sendTicker:
			log.Printf("metrics: export %v", time.Now())
			err = Export(r, c, cfg)
		case <-pingTicker:
			log.Printf("metrics: tick %v", time.Now())
			_, _, err = c.Ping(time.Second * 1)
		}
		if err != nil {
			log.Printf("metrics: %v %v", time.Now(), err)
			for c, err = Connect(cfg); err != nil; {
				log.Printf("metrics: %v %v", time.Now(), err)
			}
		}
	}
}

func Connect(cfg InfluxDbConfig) (result client.Client, err error) {
	conf := client.HTTPConfig{
		Addr:     cfg.Host,
		Username: cfg.Username,
		Password: cfg.Password,
	}
	result, err = client.NewHTTPClient(conf)
	return
}

func Export(r *Registry, c client.Client, cfg InfluxDbConfig) (err error) {
	t := time.Now()
	bpc := client.BatchPointsConfig{
		Database:  cfg.Database,
		Precision: "s",
	}
	var bps client.BatchPoints
	if bps, err = client.NewBatchPoints(bpc); err != nil {
		return
	}
	var p *client.Point
	for _, m := range r.Measures {
		fields := make(map[string]interface{})
		for field, c := range m.Counters {
			fields[field] = c.Value()
		}
		if p, err = client.NewPoint(m.Name, m.Tags, fields, t); err != nil {
			return
		}
		bps.AddPoint(p)
	}
	err = c.Write(bps)
	return
}
