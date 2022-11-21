package gcal

import (
	"context"
	"time"

	"github.com/huyungtang/go-lib/google"
	base "google.golang.org/api/calendar/v3"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(opts ...google.Options) (serv Service, err error) {
	cfg := new(google.Option).
		ApplyOptions(opts,
			CalendarIdOption("primary"),
			EventDurationOption(time.Hour),
		)

	var cal *base.Service
	if cal, err = base.NewService(context.Background(), cfg.GetClientOption()); err != nil {
		return
	}

	return &service{cal, cfg}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// service ********************************************************************************************************************************
type service struct {
	*base.Service
	*google.Option
}

// Service
// ****************************************************************************************************************************************
type Service interface {
	AddEvent(string, time.Time, ...google.Options) EventResult
}

// AddEvent()
// ****************************************************************************************************************************************
func (o *service) AddEvent(summary string, tm time.Time, opts ...google.Options) EventResult {
	opt := &google.Option{
		CalId:    o.CalId,
		Duration: o.Duration,
	}
	opt.ApplyOptions(opts,
		EventEndOption(tm.Add(opt.Duration)),
	)

	evt := &base.Event{
		Summary:     summary,
		Description: opt.Desc,
		Recurrence:  opt.Recur,
		Start:       getEventDateTime(tm, opt.TZone, opt.AllDay),
		End:         getEventDateTime(opt.EndTime, opt.TZone, opt.AllDay),
	}

	res := new(result)
	res.Event, res.err = o.Events.Insert(opt.CalId, evt).Do()

	return res
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// getEventDateTime ***********************************************************************************************************************
func getEventDateTime(tm time.Time, tz string, allDay bool) (evt *base.EventDateTime) {
	evt = &base.EventDateTime{
		TimeZone: tz,
	}
	if allDay {
		evt.Date = tm.Format(time.RFC3339)[0:10]
	} else {
		evt.DateTime = tm.Format(time.RFC3339)
	}

	return
}
