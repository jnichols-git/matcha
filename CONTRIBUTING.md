# Contributing to Router

We welcome community contributions to Matcha!

## Versioning

Matcha's latest stable release is `v1.1.0` (major version 1, minor version 1, patch 0). Releases are tagged in the GitHub repository. We follow these guidelines when deciding if a change is major, minor, or patch:

- Major: The change is non-essential and breaks the existing API. See our backwards compatiblility policy [here](docs/versioning.md). Major changes are not currently being accepted.
- Minor: The change doesn't break the existing API, but changes large portions of the internals of the library, or adds significant functionality to the API. Minor changes should relate to branch `v1.2.0`.
- Patch: The change is a bugfix, minor performance improvement, minor API change, or auxilary component (like middleware). Patches should relate to branch `main`.

Release eligibility will be evaluated biweekly, with the next evaluation being on May 15, 2023.

## Getting Started

1. Ensure you have [installed Golang](https://go.dev/dl/) on your development machine. Matcha is currently on version `1.20.2`.
2. Create a fork of Matcha to your personal GitHub account. Direct branches or pushes to the `CloudRETIC` repository are not accepted.
3. Clone your personal fork to your development machine, and enter the directory.
4. Add the upstream repository: `git remote add upstream [http-or-ssh-address]`

## Creating a Branch for Changes

1. Identify the version branch you would like to make a change to. *These branches are specified in Versioning*.
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

Currently, new submissions to Matcha are subject to the following criteria:

1. **Performance**: Changes must not significantly decrease performance unless they are urgent bugfixes. Benchmarks for routers and routes, as well as a more comprehensive benchmark on the GitHub API (courtesy of [julienschmidt](https://github.com/julienschmidt/go-http-routing-benchmark)) are provided to help evaluate this as you work. Run benchmarks before and after making changes to most accurately assess impact.
2. **Testing**: Test coverage should stay above 95%. New behavior is expected to have associated unit tests
3. **Documentation**: CloudRETIC strongly encourages detailed documentation of code, and pull requests will be evaluated on quality of comments and external documentation in `/docs`.
4. **Style**: While style is massively subjective, you should follow good Go code practices. We use [this list](https://github.com/golang/go/wiki/CodeReviewComments#gofmt) to evaluate style.
