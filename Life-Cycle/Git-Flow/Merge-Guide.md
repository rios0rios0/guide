# Merge Guide

> **TL;DR:** For dependent (chained) branches, merge from the outermost branch inward before merging into `main`. For independent branches, merge one at a time, rebasing each subsequent branch on the updated `main`.

## Overview

This guide covers the correct merge order for two common scenarios: dependent branches and independent branches. Following these procedures ensures a clean, linear commit history.

## Case 1: Dependent Branches

When branches are chained (e.g., `test/2` depends on `test/1`), merge from the **outermost** branch inward before merging into `main`.

![](.assets/case-1-commits-before.png)

### Procedure

1. Merge `test/2` into `test/1`.
2. Merge `test/1` into `main`.

For deeper chains (N branches):

1. Merge `test/N` into `test/N-1`.
2. Merge `test/N-1` into `test/N-2`.
3. Continue until `test/1`.
4. Merge `test/1` into `main`.

### Result

![](.assets/case-1-commits-result.png)

---

## Case 2: Independent Branches

When multiple branches originate independently from `main`, merge them **one at a time**, rebasing each subsequent branch on the updated `main`.

### Procedure

1. Merge `test/1` into `main`.
2. **Update `main` locally** after the merge.
3. Rebase `test/2` with the updated `main`.
4. Merge `test/2` into `main`.
5. Repeat for additional branches.

This produces independent "triangles" in the commit graph, keeping the history clean and traceable.

### Handling Rebase Conflicts

If conflicts occur during the rebase:

1. Resolve the conflicts in your IDE.
2. Stage the resolved files: `git add <resolved-file>`
3. Continue the rebase: `git rebase --continue`
4. Force-push: `git push -f`

![](.assets/case-2-commits-error.png)

Then update `main` locally and proceed with the merge. The resulting graph:

![](.assets/case-2-commits-result.png)
