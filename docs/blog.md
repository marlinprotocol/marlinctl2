# Improved process management for marlin applications
We understood the importance of seamless experience that needs to be provided to end users for managing nodes running marlin applications. With this, we had created marlinctl. It solved a variety of issues with node management - by providing quick onboarding, logical groups, fully discoverable command tree among other things.

However, there was a lot of capabilities that were desired with marlinctl's use case. Some of the feedback and use cases we discovered were:
1. Support for runtimes other than supervisor - such as support for systemd
2. Support for multiple intances of same application able to run on the same machine with varied configurations
3. Support for platforms other than linux-amd64
4. Ability to freeze a version of application for production use so it does not update automatically
5. Integrity verification for binaries - so one can be assured that binaries are compiled by marlin team
6. Ability to specify update policies for applications so that programs don't get updated with breaking changes
7. Better log and debugging for running application
8. Unstable channel support for people who would like to test release candidates and support releases before public release.

among many others

Apart from this, we also understand that a strong release pipeline from standpoint of development is critical to bring better experience to the end user. Dev teams requirements from the marlinctl and hence marlinctl was supposed to have a few more technical prowess up it's sleeve:

1. Ability to play well with semantic versioning.
2. Ability to play well with release channels - public, beta, alpha and dev channels so that marlinctl is as relevant to a geek as it is to a user with simple use case.
3. Resource management to keep tabs on how multiple instances are running (IaaC like resources)
4. Multiple runners, platforms and runtimes support - so that applications play just as well in macos with plists as it does on linux-amd64 with supervisor

All these demanded a complete architectural rework of marlinctl, and we are pleased to inform that we have made it work. We introduce marlinctl 2.0.0.

With marlinctl 2.0.0 you can expect all the features of marlinctl plus a few more:
1. Future releases for projects to be additive for most part, breaking changes will be minimised
2. Advanced users can freeze their versions to a set state
3. Logging, versions, multiple release channel subscription, semantic versioning along with integrity checks on binary makes for a more predictible and secure experience
4. Support for multiple runtimes, platforms etc.

marlinctl2 is ready with a lot up its sleeve. However, we want to make sure we have got all the rough edges cleaned up before we make it available to you. If you are reading this, marlinctl2 is made available to you already. Read through our docs on how you can set it up. Support for various other chains - their gateways, relays, configuration management endpoints among other things will be rolled out soon as minor updates which marlinctl2 should automatically update to once they are available.

As always, since we might not be able to prioritize something you want, we're always open to community contributions and PRs. And if enough people would find it useful, reach out to us and we might be able to provide dev grants for the same. If you have a suggestion, a feature request or you came across a bug, let us know on our discord channel.

If you have any ideas on how this can evolve, we're all ears!

Follow our official social media channels to get the latest updates as and when they come out! 

