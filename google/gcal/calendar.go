package gcal

import (
	"context"
	"time"

	"github.com/huyungtang/go-lib/google"
	"github.com/huyungtang/go-lib/times"
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
	AddEvent(string, ...google.Options) google.Event
	// DelEvent(eventId, CalendarId)
	DelEvent(string, ...google.Options) google.Event
	ListEvent(...google.Options) google.Events
}

// AddEvent()
// ****************************************************************************************************************************************
func (o *service) AddEvent(summary string, opts ...google.Options) google.Event {
	opt := &google.Option{
		CalendarId:   o.CalendarId,
		Transparency: "transparent",
		StartTime:    time.Now(),
	}
	opt.ApplyOptions(opts)

	evt := &base.Event{
		Summary:      summary,
		Description:  opt.Description,
		Recurrence:   opt.Recurrency,
		Start:        getEventDateTime(opt.StartTime, opt.Timezone, opt.AllDay),
		End:          getEventDateTime(opt.EndTime, opt.Timezone, opt.AllDay),
		Transparency: opt.Transparency,
	}

	res := new(result)
	res.Event, res.err = o.Events.Insert(opt.CalendarId, evt).Do()

	return res
}

// DelEvent
// ****************************************************************************************************************************************
func (o *service) DelEvent(evtId string, opts ...google.Options) google.Event {
	opt := (&google.Option{
		CalendarId: o.CalendarId,
	}).ApplyOptions(opts)

	return &result{
		err: o.Events.Delete(opt.CalendarId, evtId).Do(),
	}
}

// ListEvent
// ****************************************************************************************************************************************
func (o *service) ListEvent(opts ...google.Options) google.Events {
	opt := (&google.Option{
		CalendarId: o.CalendarId,
		MaxResult:  100,
		StartTime:  times.Today().Add(0, -1).Time(),
		EndTime:    times.Today().Add(0, 0, 1).Time(),
	}).ApplyOptions(opts)

	res := new(result)
	if res.Events, res.err = o.Events.
		List(opt.CalendarId).
		Q(opt.Description).
		MaxResults(opt.MaxResult).
		TimeMin(getEventDateTime(opt.StartTime, o.Timezone, false).DateTime).
		TimeMax(getEventDateTime(opt.EndTime, o.Timezone, false).DateTime).
		PageToken(opt.NextPage).
		Do(); res.err == nil {
		res.evts = len(res.Events.Items)
	}

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
