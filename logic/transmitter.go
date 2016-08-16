package logic

import "receiver/data"

type Transiver interface {
	Send(d *data.Configuration, r *data.Record)
}