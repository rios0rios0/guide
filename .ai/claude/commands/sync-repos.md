Sync all git repositories under `$HOME/Development`. For each repo at depth 2 (`<provider>/<org-or-user>/<repo>`):

1. If the repo has uncommitted changes (staged, unstaged, or untracked), create a WIP branch (`wip/auto-stash-<timestamp>`) and commit all changes with message `wip: auto-stash uncommitted changes`.
2. Checkout the default branch (detected via `git symbolic-ref refs/remotes/origin/HEAD`, fallback to `main`).
3. Run `git fetch --all --prune` and `git pull --rebase`.
4. If pull/rebase fails, abort the rebase and restore the previous branch.
5. Restore the original branch (or WIP branch if changes were stashed).

Run this by executing the shell function:

```bash
source ~/.scripts/linux-engineering-git-sync-repos.sh && git-sync-repos
```

If the user provides a specific subdirectory (e.g. `dev.azure.com/ZestSecurity`), pass it as argument:

```bash
source ~/.scripts/linux-engineering-git-sync-repos.sh && git-sync-repos "$HOME/Development/dev.azure.com/ZestSecurity"
```

Print the summary at the end showing total repos, synced, WIP commits, and failures.
