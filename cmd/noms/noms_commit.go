// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/stormasm/nomsbolt/cmd/util"
	"github.com/stormasm/nomsbolt/go/config"
	"github.com/stormasm/nomsbolt/go/d"
	"github.com/stormasm/nomsbolt/go/datas"
	"github.com/stormasm/nomsbolt/go/spec"
	"github.com/stormasm/nomsbolt/go/util/verbose"
	flag "github.com/juju/gnuflag"
)

var allowDupe bool

var nomsCommit = &util.Command{
	Run:       runCommit,
	UsageLine: "commit [options] [absolute-path] <dataset>",
	Short:     "Commits a specified value as head of the dataset",
	Long:      "If absolute-path is not provided, then it is read from stdin. See Spelling Objects at https://github.com/stormasm/nomsbolt/blob/master/doc/spelling.md for details on the dataset and absolute-path arguments.",
	Flags:     setupCommitFlags,
	Nargs:     1, // if absolute-path not present we read it from stdin
}

func setupCommitFlags() *flag.FlagSet {
	commitFlagSet := flag.NewFlagSet("commit", flag.ExitOnError)
	commitFlagSet.BoolVar(&allowDupe, "allow-dupe", false, "creates a new commit, even if it would be identical (modulo metadata and parents) to the existing HEAD.")
	spec.RegisterCommitMetaFlags(commitFlagSet)
	verbose.RegisterVerboseFlags(commitFlagSet)
	return commitFlagSet
}

func runCommit(args []string) int {
	cfg := config.NewResolver()
	db, ds, err := cfg.GetDataset(args[len(args)-1])
	d.CheckError(err)
	defer db.Close()

	var path string
	if len(args) == 2 {
		path = args[0]
	} else {
		readPath, _, err := bufio.NewReader(os.Stdin).ReadLine()
		d.CheckError(err)
		path = string(readPath)
	}
	absPath, err := spec.NewAbsolutePath(path)
	d.CheckError(err)

	value := absPath.Resolve(db)
	if value == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", path)))
	}

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		head := ds.HeadValue()
		if head.Hash() == value.Hash() && !allowDupe {
			fmt.Fprintf(os.Stdout, "Commit aborted - allow-dupe is set to off and this commit would create a duplicate\n")
			return 0
		}
	}

	meta, err := spec.CreateCommitMetaStruct(db, "", "", nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, value, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	if oldCommitExists {
		fmt.Fprintf(os.Stdout, "New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Fprintf(os.Stdout, "New head #%v\n", ds.HeadRef().TargetHash().String())
	}
	return 0
}
