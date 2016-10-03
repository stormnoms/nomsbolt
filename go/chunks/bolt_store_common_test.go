// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package chunks

import (
	"github.com/attic-labs/testify/suite"

	//"github.com/attic-labs/noms/go/constants"
	//"github.com/attic-labs/noms/go/hash"
)

type BoltStoreTestSuite struct {
	suite.Suite
	Store      ChunkStore
	putCountFn func() int
}

func (suite *BoltStoreTestSuite) TestChunkBoltStorePut() {
	input := "abc"
	c := NewChunk([]byte(input))
	suite.Store.Put(c)
	h := c.Hash()

	// See http://www.di-mgt.com.au/sha_testvectors.html
	suite.Equal("rmnjb8cjc5tblj21ed4qs821649eduie", h.String())

	suite.Store.UpdateRoot(h, suite.Store.Root()) // Commit writes

	// And reading it via the API should work...
	assertInputInStore(input, h, suite.Store, suite.Assert())
	if suite.putCountFn != nil {
		suite.Equal(1, suite.putCountFn())
	}

	// Re-writing the same data should cause a second put
	c = NewChunk([]byte(input))
	suite.Store.Put(c)
	suite.Equal(h, c.Hash())
	assertInputInStore(input, h, suite.Store, suite.Assert())
	suite.Store.UpdateRoot(h, suite.Store.Root()) // Commit writes

	if suite.putCountFn != nil {
		suite.Equal(2, suite.putCountFn())
	}
}
