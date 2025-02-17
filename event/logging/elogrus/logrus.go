// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !disable_events

// Package elogrus provides a logrus Formatter for events.
// To use for the global logger:
//   logrus.SetFormatter(elogrus.NewFormatter(exporter))
//   logrus.SetOutput(io.Discard)
// and for a Logger instance:
//   logger.SetFormatter(elogrus.NewFormatter(exporter))
//   logger.SetOutput(io.Discard)
//
// If you call elogging.SetExporter, then you can pass nil
// for the exporter above and it will use the global one.
package elogrus

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/event"
	"golang.org/x/exp/event/keys"
	"golang.org/x/exp/event/logging/internal"
)

type formatter struct{}

func NewFormatter() logrus.Formatter {
	return &formatter{}
}

var _ logrus.Formatter = (*formatter)(nil)

// Format writes an entry to an event.Exporter. It always returns nil (see below).
// If e.Context is non-nil, Format gets the exporter from the context. Otherwise
// it uses the default exporter.
//
// Logrus first calls the Formatter to get a []byte, then writes that to the
// output. That doesn't work for events, so we subvert it by having the
// Formatter export the event (and thereby write it). That is why the logrus
// Output io.Writer should be set to io.Discard.
func (f *formatter) Format(e *logrus.Entry) ([]byte, error) {
	ctx := e.Context
	if ctx == nil {
		ctx = context.Background()
	}
	b := event.To(ctx).At(e.Time)
	b.With(internal.LevelKey.Of(int(e.Level))) // TODO: convert level
	for k, v := range e.Data {
		b.With(keys.Value(k).Of(v))
	}
	b.Log(e.Message)
	return nil, nil
}
