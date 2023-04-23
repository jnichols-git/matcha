# Matcha and Version Control

One of the priorities of `matcha` is *backwards compatibility*; that is, updates to the minor version of the library shouldn't break pre-existing codebases using it. To clarify what that means for developers, this document lays out how backwards compatibility is maintained, and under what circumstances that policy might be violated, and what happens when those circumstances occur.

## Maintaining Backwards Compatibility

Most of `matcha`'s developer-facing capabilities are behind interfaces, collections of types that can do a certain set of clearly-defined things. Some are defined by us, while some are defined by the Go standard library, like `context.Context`. What this means in practice is that when changes are made to the underlying mechanisms of routers, routes, and context, the code that you use to invoke those things doesn't change. The result is an easy-to-read definition of what functionality will always be available to you, the developer, throughout versions.

Any changes to public/exported packages will not remove features until version 2.0, even if that functionality ceases use internally.

## Problem Cases

There are limited cases when we may break backwards compatibility, such as severe bugs, where maintaining stability takes priority over compatibility. This won't be done unless the current available functionality compromises the library in some way. If this occurs, an announcement will be posted at the top of the README, as well as on our communication board on [Discord](https://discord.gg/CYPnDDCG).
