// Copyright (C) 2020-2021 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package generic

import (
	"fmt"

	"github.com/test-network-function/test-network-function/pkg/tnf"
	"github.com/test-network-function/test-network-function/pkg/tnf/handlers/base/redhat"
	"github.com/test-network-function/test-network-function/pkg/tnf/reel"
	"github.com/test-network-function/test-network-function/pkg/tnf/testcases"

	"github.com/test-network-function/test-network-function/test-network-function/common"

	"github.com/onsi/ginkgo"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	"github.com/onsi/gomega"
)

const (
	testsKey = "generic"
)

//
// All actual test code belongs below here.  Utilities belong above.
//

// Runs the "generic" CNF test cases.
var _ = ginkgo.Describe(testsKey, func() {
	if testcases.IsInFocus(ginkgoconfig.GinkgoConfig.FocusStrings, testsKey) {
		configData := common.ConfigurationData{}
		configData.SetNeedsRefresh()
		ginkgo.BeforeEach(func() {
			common.ReloadConfiguration(&configData)
		})

		testIsRedHatRelease(&configData)
	}
})

// testIsRedHatRelease fetch the configuration and test containers attached to oc is Red Hat based.
func testIsRedHatRelease(configData *common.ConfigurationData) {
	ginkgo.It("Should report a proper Red Hat version", func() {
		for _, cut := range configData.ContainersUnderTest {
			testContainerIsRedHatRelease(cut)
		}
		testContainerIsRedHatRelease(configData.TestOrchestrator)
	})
}

// testContainerIsRedHatRelease tests whether the container attached to oc is Red Hat based.
func testContainerIsRedHatRelease(cut *common.Container) {
	podName := cut.Oc.GetPodName()
	containerName := cut.Oc.GetPodContainerName()
	context := cut.Oc
	ginkgo.By(fmt.Sprintf("%s(%s) is checked for Red Hat version", podName, containerName))
	versionTester := redhat.NewRelease(common.DefaultTimeout)
	test, err := tnf.NewTest(context.GetExpecter(), versionTester, []reel.Handler{versionTester}, context.GetErrorChannel())
	gomega.Expect(err).To(gomega.BeNil())
	testResult, err := test.Run()
	gomega.Expect(testResult).To(gomega.Equal(tnf.SUCCESS))
	gomega.Expect(err).To(gomega.BeNil())
}
