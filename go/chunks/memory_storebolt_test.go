// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package chunks

import (
	//"bytes"
	"testing"

	"github.com/attic-labs/testify/suite"
)

func TestMemoryStoreBoltTestSuite(t *testing.T) {
	suite.Run(t, &MemoryStoreTestSuite{})
}

type MemoryStoreBoltTestSuite struct {
	BoltStoreTestSuite
}

func (suite *MemoryStoreBoltTestSuite) SetupTest() {
	suite.Store = NewMemoryStore()
}

func (suite *MemoryStoreBoltTestSuite) TearDownTest() {
	suite.Store.Close()
}
