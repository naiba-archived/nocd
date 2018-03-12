/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package main

import (
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

func init() {
	raven.SetDSN("https://b2be1d09de6a4765aa1bf2f02c58d156:f1598a174c9441648e09b7d88e29d7a6@sentry.io/301979")
}

func main() {
	raven.CaptureErrorAndWait(errors.New("Hello sentry"), nil)
}
