package helper

import "time"

func IsUTC(t time.Time) bool { _, off := t.Zone(); return off == 0 }
