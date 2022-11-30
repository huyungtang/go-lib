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
	// AddEvent(summary, event_time, CalendarId, Description, Recurrency, Timezone, AllDay, EndTime, Busy)
	AddEvent(string, time.Time, ...google.Options) google.EventResult
	// DelEvent(eventId, CalendarId)
	DelEvent(string, ...google.Options) google.EventResult
}

// AddEvent()
// ****************************************************************************************************************************************
func (o *service) AddEvent(summary string, tm time.Time, opts ...google.Options) google.EventResult {
	opt := &google.Option{
		CalendarId:   o.CalendarId,
		Duration:     o.Duration,
		Transparency: "transparent",
	}
	opt.ApplyOptions(opts,
		EventEndOption(tm.Add(opt.Duration)),
	)

	evt := &base.Event{
		Summary:      summary,
		Description:  opt.Description,
		Recurrence:   opt.Recurrency,
		Start:        getEventDateTime(tm, opt.Timezone, opt.AllDay),
		End:          getEventDateTime(opt.EndTime, opt.Timezone, opt.AllDay),
		Transparency: opt.Transparency,
	}

	res := new(result)
	res.Event, res.err = o.Events.Insert(opt.CalendarId, evt).Do()

	return res
}

// DelEvent
// ****************************************************************************************************************************************
func (o *service) DelEvent(evtId string, opts ...google.Options) google.EventResult {
	opt := (&google.Option{
		CalendarId: o.CalendarId,
	}).ApplyOptions(opts)

	return &result{
		err: o.Events.Delete(opt.CalendarId, evtId).Do(),
	}
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
