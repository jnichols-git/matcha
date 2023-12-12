# Contributing to Matcha

We welcome community contributions to Matcha!

## Getting Started

1. Ensure you have [installed Golang](https://go.dev/dl/) on your development machine. Matcha requires at least version 1.20.
2. `git clone git@github.com:jnichols-git/matcha/v2.git`

## Making Changes

1. Check out the correct development branch (usually `main`).
2. Pull in remote changes by doing a fetch/rebase: `git pull --rebase`
3. Create your own development branch: `git branch [my-branch]`
4. Check out your branch: `git checkout [my-branch]`
5. Make your changes!
6. Push your changes to remote and make a pull request.

Please don't make changes on branch `main` or any semver (`vX.X.X`) branch. It will only make your life harder. If you accidentally commit to one of these, this guide may help (external link): <https://dev.to/projectpage/how-to-move-a-commit-to-another-branch-in-git-4lj4>

## Code Quality

Please run `gofmt` and ensure that your code passes tests before you make a PR.
