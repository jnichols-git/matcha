# Contributing to Router

> ***The current development branch is `v1.1`.***

We welcome community contributions to `router`!

## Getting Started

1. Ensure you have [installed Golang](https://go.dev/dl/) on your development machine. `router` is currently on version `1.20.2`.
2. Create a fork of `router` to your personal GitHub account. Direct branches or pushes to the `CloudRETIC` repository are not accepted.
3. Clone your personal fork to your development machine, and enter the directory.
4. Add the upstream repository: `git remote add upstream [http-or-ssh-address]`

## Creating a Branch for Changes

1. Identify the version number you would like to make a change to. *The current active development branch is noted at the top of this document.*
2. Ensure that your local repository is up to date:

    ```bash
    git fetch upstream
    git checkout [version-branch]
    git rebase upstream/[version-branch]
    ```

3. Create a new branch to make changes: `git checkout -b [new-branch-name]`

## Submitting Changes

1. Push your changes to your personal fork: `git push origin [branch-name]`
2. Make a pull request to the development branch, and request review from a maintainer.
3. Revise your changes based on maintainer feedback and test results.

Once your changes are approved, they'll be squash-and-merged into the feature branch.

## Guidelines

Currently, new submissions to `router` are subject to the following criteria:

1. **Performance**: Changes must not significantly decrease performance. Benchmarks for routers and routes, as well as a more comprehensive benchmark on the GitHub API (courtesy of [julienschmidt](https://github.com/julienschmidt/go-http-routing-benchmark)) are provided to help evaluate this as you work. Run benchmarks before and after making changes to most accurately assess impact.
2. **Testing**: Test coverage should stay above 95%. New behavior is expected to have associated unit tests
3. **Documentation**: CloudRETIC strongly encourages detailed documentation of code, and pull requests will be evaluated on quality of comments and external documentation in `/docs`.
4. **Style**: While style is massively subjective, you should follow good Go code practices. We use [this list](https://github.com/golang/go/wiki/CodeReviewComments#gofmt) to evaluate style.
