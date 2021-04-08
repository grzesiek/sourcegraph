package window

import (
	"math"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/sourcegraph/sourcegraph/schema"
)

// Configuration represents the rollout windows configured on the site.
type Configuration struct {
	windows []Window
}

// NewConfiguration constructs a Configuration based on the current site
// configuration, including watching the configuration for updates.
func NewConfiguration(raw *[]*schema.BatchChangeRolloutWindow) (*Configuration, error) {
	windows, err := parseConfiguration(raw)
	if err != nil {
		return nil, err
	}

	return &Configuration{windows: windows}, nil
}

// ValidateConfiguration validates the given site configuration.
func ValidateConfiguration(raw *[]*schema.BatchChangeRolloutWindow) error {
	_, err := parseConfiguration(raw)
	return err
}

// Estimate attempts to estimate when the given entry in a queue of changesets
// to be reconciled would be reconciled. nil indicates that there is no
// reasonable estimate, either because all windows are zero or the estimate is
// too far in the future to be reliable.
func (cfg *Configuration) Estimate(now time.Time, n int) *time.Time {
	if !cfg.HasRolloutWindows() {
		return &now
	}

	// Roughly speaking, we iterate over schedules until we reach the one that
	// would include the given entry.
	rem := n
	at := now
	until := at.Add(7 * 24 * time.Hour)
	for at.Before(until) {
		schedule := cfg.scheduleAt(at, false)

		// An unlimited schedule means that the reconciliation will happen
		// immediately.
		if schedule.total() == -1 {
			return &at
		}

		total := schedule.total()
		if total == 0 {
			at = schedule.ValidUntil()
			continue
		}

		rem -= total
		if rem < 0 {
			// We know how many extra reconciliations will occur within this
			// schedule, so we can use that calculate what percentage of the way
			// into the window our target will be reconciled, then we can
			// multiple the schedule duration by that to get the approximate
			// time.
			perc := 1.0 - math.Abs(float64(rem))/float64(total)
			duration := time.Duration(float64(schedule.ValidUntil().Sub(at)) * perc)
			at = at.Add(duration)
			return &at
		} else if rem == 0 {
			// Special case: this will be the very last entry to be reconciled.
			at = schedule.ValidUntil()
			return &at
		}

		at = schedule.ValidUntil()
	}

	return nil
}

// HasRolloutWindows returns true if one or more windows have been defined.
func (cfg *Configuration) HasRolloutWindows() bool {
	return len(cfg.windows) != 0
}

// Schedule calculates a schedule for the near future (most likely about a
// minute), based on the history.
func (cfg *Configuration) Schedule() Schedule {
	// If there are no rollout windows, then we return an unlimited schedule and
	// have the scheduler check back in periodically in case the configuration
	// updated. Ten minutes is probably safe enough.
	if !cfg.HasRolloutWindows() {
		return newSchedule(time.Now(), 10*time.Minute, rate{n: -1})
	}

	return cfg.scheduleAt(time.Now(), true)
}

// currentFor returns the current rollout window, if any, and the duration for
// which that window applies from now. The duration will be nil if the current
// window is applied forever.
//
// As this is an internal function, no locking is applied here.
func (cfg *Configuration) currentFor(now time.Time) (*Window, *time.Duration) {
	// If there are no rollout windows, there's no current window. This should
	// be checked before entry, but let's at least not panic here.
	if len(cfg.windows) == 0 {
		return nil, nil
	}

	// Find the last matching window that is currently active.
	index := -1
	for i := range cfg.windows {
		if cfg.windows[i].IsOpen(now) {
			index = i
		}
	}
	if index == -1 {
		// No matching window, so let's figure out when the next window would
		// open and return a nil window.
		var next *time.Time
		for i := range cfg.windows {
			at := cfg.windows[i].NextOpenAfter(now)
			if next == nil || at.Before(*next) {
				next = &at
			}
		}

		// If we never saw a time, that's weird, since this scenario shouldn't
		// occur if there are windows defined, but let's just say nothing can
		// happen forever for now.
		if next == nil {
			return nil, nil
		}

		duration := next.Sub(now)
		return nil, &duration
	}
	window := &cfg.windows[index]

	// Calculate when this window closes. This may not be the end time on the
	// window: if there's a later window that starts before the end time of this
	// window, that will end up taking precedence.
	var end *time.Time

	if window.end == nil {
		// There may still be a weekday restriction, so we should figure that out.
		if !window.days.all() {
			start := now.Truncate(24 * time.Hour)
			for {
				start = start.Add(24 * time.Hour)
				if !window.days.includes(start.Weekday()) {
					start.Add(-1 * time.Second)
					end = &start
					break
				} else if start.After(now.Add(7 * 24 * time.Hour)) {
					panic("could not find end of a day-limited window in the next week")
				}
			}
		}
	} else {
		// We have a concrete end time for this window, so we can set end to
		// that.
		windowEnd := time.Date(now.Year(), now.Month(), now.Day(), int(window.end.hour), int(window.end.minute), 0, 0, time.UTC)
		end = &windowEnd
	}

	// Now we iterate over the subsequent windows in the configuration and see
	// if any of them would start before the existing end time, which would make
	// them active (since they're subsequent). Note that we're using a C style
	// for loop here instead of slicing: we'd have to check the bounds of the
	// cfg.windows slice before being able to subslice, and this feels more
	// readable.
	for i := index + 1; i < len(cfg.windows); i++ {
		nextActive := cfg.windows[i].NextOpenAfter(now)
		if end != nil {
			if nextActive.Before(*end) {
				end = &nextActive
			}
		} else {
			end = &nextActive
		}
	}

	// If we still don't have an end time, then this window remains open forever
	// and cannot be overridden. Cool.
	if end == nil {
		return window, nil
	}

	// Otherwise, let's calculate how long we have until this window closes, and
	// return that.
	d := end.Sub(now)
	return window, &d
}

func (cfg *Configuration) scheduleAt(at time.Time, minimal bool) Schedule {
	window, validity := cfg.currentFor(at)

	// No window means a zero schedule should be returned until the next window
	// change.
	if window == nil {
		if validity != nil {
			if *validity >= 1*time.Minute && minimal {
				return newSchedule(at, 1*time.Minute, rate{n: 0})
			}
			return newSchedule(at, *validity, rate{n: 0})
		}
		// We should always have a validity in this case, but let's be defensive
		// if we don't for some reason. The scheduler can check back in a
		// minute.
		return newSchedule(at, 1*time.Minute, rate{n: 0})
	}

	// OK, so we have a rollout window. It may or may not have an expiry. Either
	// way, let's calculate how long we'd want to schedule in an ideal world.
	if validity == nil {
		if window.rate.unit == ratePerHour {
			return newSchedule(at, 1*time.Hour, window.rate)
		}
		// We never really want to have less than a minute's worth of schedule
		// in this case: we want to check occasionally for updated
		// configurations, but don't want to calculate every time the scheduler
		// needs a changeset.
		return newSchedule(at, 1*time.Minute, window.rate)
	}

	// TODO: minimise if needed.
	return newSchedule(at, *validity, window.rate)
}

func parseConfiguration(raw *[]*schema.BatchChangeRolloutWindow) ([]Window, error) {
	// Ensure we always start with an empty window slice.
	windows := []Window{}

	// If there's no window configuration, there are no windows, and we can just
	// return here.
	if raw == nil {
		return windows, nil
	}

	var errs *multierror.Error
	for i, rawWindow := range *raw {
		if window, err := parseWindow(rawWindow); err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "window %d", i))
		} else {
			windows = append(windows, window)
		}
	}

	return windows, errs.ErrorOrNil()
}