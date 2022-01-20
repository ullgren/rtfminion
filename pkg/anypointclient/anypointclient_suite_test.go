package anypointclient

/*
Copyright © 2021 Pontus Ullgren

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	client     AnypointClient
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

var _ = BeforeSuite(func() {
	client = NewAnypointClientWithCredentials(
		REGION_US,
		"user",
		"password",
	)
	// block all HTTP requests
	httpmock.ActivateNonDefault(client.HTTPClient)
})

var _ = BeforeEach(func() {
	// remove any mocks
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

func TestApClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Anypoint Client Test Suite")
}
