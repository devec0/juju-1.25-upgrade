// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package network_test

import (
	"os"
	"runtime"
	"time"

	"github.com/juju/testing"
	"github.com/juju/utils/clock"
	gc "gopkg.in/check.v1"

	"github.com/juju/1.25-upgrade/juju2/network"
)

// A note regarding the use of clock.WallClock in these unit tests.
//
// All the tests pass 0 for a timeout, which means indefinite, and
// therefore no timer/clock is used. There is one test that checks for
// timeout and passes 0.5s as its timeout value. Because of this it's
// not clear why the 'testing clock' would be a better choice.

type BridgeSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&BridgeSuite{})

const echoArgsScript = `
import sys
for arg in sys.argv[1:]: print(arg)
`

func (s *BridgeSuite) SetUpSuite(c *gc.C) {
	s.IsolationSuite.SetUpSuite(c)
	if runtime.GOOS == "windows" {
		c.Skip("skipping BridgeSuite tests on windows")
	}
}

func assertCmdResult(c *gc.C, cmd, expected string) {
	result, err := network.RunCommand(cmd, os.Environ(), clock.WallClock, 0)
	c.Assert(err, gc.IsNil)
	c.Assert(result.Code, gc.Equals, 0)
	c.Assert(string(result.Stdout), gc.Equals, expected)
	c.Assert(string(result.Stderr), gc.Equals, "")
}

func assertENIBridgerError(c *gc.C, devices []network.DeviceToBridge, timeout time.Duration, clock clock.Clock, filename string, dryRun bool, reconfigureDelay int, expected string) {
	bridger := network.NewEtcNetworkInterfacesBridger(clock, timeout, filename, dryRun)
	err := bridger.Bridge(devices, reconfigureDelay)
	c.Assert(err, gc.NotNil)
	c.Assert(err, gc.ErrorMatches, expected)
}

func (*BridgeSuite) TestENIBridgerWithMissingFilenameArgument(c *gc.C) {
	devices := []network.DeviceToBridge{
		network.DeviceToBridge{
			DeviceName: "ens123",
			BridgeName: "br-ens123",
		},
	}
	expected := `bridge activaction error: filename and input is nil`
	assertENIBridgerError(c, devices, 0, clock.WallClock, "", true, 0, expected)
}

func (*BridgeSuite) TestENIBridgerWithEmptyDeviceNamesArgument(c *gc.C) {
	devices := []network.DeviceToBridge{}
	expected := `bridge activaction error: no devices specified`
	assertENIBridgerError(c, devices, 0, clock.WallClock, "testdata/non-existent-filename", true, 0, expected)
}

func (*BridgeSuite) TestENIBridgerWithNonExistentFile(c *gc.C) {
	devices := []network.DeviceToBridge{
		network.DeviceToBridge{
			DeviceName: "ens123",
			BridgeName: "br-ens123",
		},
	}
	expected := `bridge activaction error: open testdata/non-existent-file: no such file or directory`
	assertENIBridgerError(c, devices, 0, clock.WallClock, "testdata/non-existent-file", true, 0, expected)
}

func (*BridgeSuite) TestENIBridgerWithTimeout(c *gc.C) {
	devices := []network.DeviceToBridge{
		network.DeviceToBridge{
			DeviceName: "ens123",
			BridgeName: "br-ens123",
		},
	}
	expected := "bridge activaction error: bridge activation error: command cancelled"
	// 25694 is a magic value that causes the bridging script to sleep
	assertENIBridgerError(c, devices, 500*time.Millisecond, clock.WallClock, "testdata/interfaces", true, 25694, expected)
}

func (*BridgeSuite) TestENIBridgerWithDryRun(c *gc.C) {
	devices := []network.DeviceToBridge{
		network.DeviceToBridge{
			DeviceName: "ens123",
			BridgeName: "br-ens123",
		},
	}
	bridger := network.NewEtcNetworkInterfacesBridger(clock.WallClock, 0, "testdata/interfaces", true)
	err := bridger.Bridge(devices, 0)
	c.Assert(err, gc.IsNil)
}
