// Package client ...
// Copyright 2019 The yeqown Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// https://godoc.org/gopkg.in/src-d/go-git.v4

func main() {

	branchName := "master"
	commitHash := "8a951b33e9cf0e2c13559255d4749d40279125b6"

	// clone repo
	repo, err := git.PlainClone("testdata/", false, &git.CloneOptions{
		URL:           "https://github.com/yeqown/playground",
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
		Progress:      os.Stdout,
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// check out branch and commitHash
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
		// Branch: plumbing.NewBranchReferenceName(branchName), branch 和 hash 是互斥的
		Create: false,
		Force:  false,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
