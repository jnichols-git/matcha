# Contributing to Matcha

- [Contributing to Matcha](#contributing-to-matcha)
  - [Getting Started](#getting-started)
  - [Making Changes](#making-changes)
    - [Picking an Issue](#picking-an-issue)
    - [Creating a Branch for Changes](#creating-a-branch-for-changes)
    - [Submitting Changes](#submitting-changes)
  - [Submission Standards](#submission-standards)
  - [Versioning Policy](#versioning-policy)
    - [Deprecated Features](#deprecated-features)

We welcome community contributions to Matcha!

## Getting Started

1. Ensure you have [installed Golang](https://go.dev/dl/) on your development machine. Matcha is currently on version `1.20.2`.
2. Create a fork of Matcha to your personal GitHub account. Direct branches or pushes to the `CloudRETIC` repository are not accepted.
3. Clone your personal fork to your development machine, and enter the directory.
4. Add the upstream repository: `git remote add upstream [http-or-ssh-address]`

## Making Changes

### Picking an Issue

We do our best to keep issues up to date and appropriately tagged. Here's a quick rundown on how to pick an issue based on tags.

1. If you're a first time contributor, find one tagged `good first issue` or `patch`. These tend to be short, accessible tasks to help you get to know a specific part of the codebase, like middleware.
2. Once you're more comfortable with the structure of Matcha, pick up `patch`, `minor`, and `bugfix` issues that interest you.
3. If you're feeling a non-code task, there's usually a `documentation` issue or two avaliable for assignment.

Once you have decided on an issue, just leave a comment asking to have it assigned to you. Issues are generally first-come first-serve. If an issue has been inactive for 3 weeks, we will consider reassignment.

### Creating a Branch for Changes

1. Identify the version branch you would like to make a change to. *These branches are specified in Versioning*.
2. Ensure that your local repository is up to date:

    ```bash
    git fetch upstream
    git checkout [version-branch]
    git rebase upstream/[version-branch]
    ```

3. Create a new branch to make changes: `git checkout -b [new-branch-name]`

### Submitting Changes

1. Push your changes to your personal fork: `git push origin [branch-name]`
2. Make a pull request to the development branch, and request review from a maintainer.
3. Revise your changes based on maintainer feedback and test results.

Once your changes are approved, they'll be squash-and-merged into the feature branch.

## Submission Standards

Currently, new submissions to Matcha are subject to the following criteria:

1. **Performance**: Changes must not significantly decrease performance unless they are urgent bugfixes. Benchmarks for routers and routes, as well as a more comprehensive benchmark on the GitHub API (courtesy of [julienschmidt](https://github.com/julienschmidt/go-http-routing-benchmark)) are provided to help evaluate this as you work. Run benchmarks before and after making changes to most accurately assess impact.
2. **Testing**: Test coverage should stay above 95%. New behavior is expected to have associated unit tests, and PRs that drop coverage by more than 2%, or below 90%, will be automatically rejected.
3. **Documentation**: CloudRETIC strongly encourages detailed documentation of code, and pull requests will be evaluated on quality of comments and external documentation in `/docs`.
4. **Style**: Follow good Go code practices. We use [this list](https://github.com/golang/go/wiki/CodeReviewComments#gofmt) to evaluate style.
5. **Zero-Dependency**: Matcha does not use any external libraries. Changes with dependencies will be rejected.

We additionally ask that you avoid the use of AI tools like ChatGPT and GitHub Copilot in your contributions to Matcha. CloudRETIC has concerns regarding the ethics and copyright implications of scraping open-source code for training data, and would prefer that any work submitted be attributable to the users that directly contribute to it.

## Versioning Policy

Matcha's latest stable release is `v1.1.1` (major version 1, minor version 1, patch 0). Releases are tagged in the GitHub repository. We follow these guidelines when deciding if a change is major, minor, or patch:

- Major: The change is non-essential and breaks the existing API. Major changes are not currently being accepted.
- Minor: The change doesn't break the existing API, but changes large portions of the internals of the library, or adds significant functionality to the API. Minor changes should relate to branch `v1.2.0`.
- Patch: The change is a bugfix, minor performance improvement, minor API change, or auxilary component (like middleware). Patches should relate to branch `main`.

### Deprecated Features

To maintain the long-term health of the project, you/we may elect to *deprecate* features, meaning that we won't support them going forward. If you have a change that overrides old functionality, it's preferred that you mark the old as deprecated and implement alongside it, rather than delete it from the project entirely; we generally prefer to keep behavior the same between versions when it's not a bug problem.
