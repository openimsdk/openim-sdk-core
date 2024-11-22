## [v3.8.2](https://github.com/openimsdk/openim-sdk-core/releases/tag/v3.8.2) 	(2024-11-22)

### New Features
* feat: Support FetchSurroundingMessages [#741](https://github.com/openimsdk/openim-sdk-core/pull/741)
* feat: mark all conversation as read [#743](https://github.com/openimsdk/openim-sdk-core/pull/743)
* feat: implement default logger when no init. [#755](https://github.com/openimsdk/openim-sdk-core/pull/755)
* feat: implement error stack print. [#733](https://github.com/openimsdk/openim-sdk-core/pull/733)
* feat: support stream message [#770](https://github.com/openimsdk/openim-sdk-core/pull/770)

### Bug Fixes
* fix: change data convert to Values [#731](https://github.com/openimsdk/openim-sdk-core/pull/731)
* fix: SearchLocalMessages no such table [#737](https://github.com/openimsdk/openim-sdk-core/pull/737)
* fix: remove version space [#750](https://github.com/openimsdk/openim-sdk-core/pull/750)
* fix: update the latest message when group member's changed. [#752](https://github.com/openimsdk/openim-sdk-core/pull/752)
* fix: remove duplicate License. [#747](https://github.com/openimsdk/openim-sdk-core/pull/747)
* fix: After being removed as a group admin, delete group requests [#754](https://github.com/openimsdk/openim-sdk-core/pull/754)
* fix: improve batchUserFaceURLandName logic. [#756](https://github.com/openimsdk/openim-sdk-core/pull/756)
* fix: escape table names to avoid the sqlite error: near "-": syntax e… [#762](https://github.com/openimsdk/openim-sdk-core/pull/762)
* Fix local cache: user cache and group member cache [#765](https://github.com/openimsdk/openim-sdk-core/pull/765)
* fix: fix temp file don't remove when upload file. [#764](https://github.com/openimsdk/openim-sdk-core/pull/764)
* fix: GetGroupMembersInfoFunc return key [#767](https://github.com/openimsdk/openim-sdk-core/pull/767)
* fix: Change check reinstall logic [#766](https://github.com/openimsdk/openim-sdk-core/pull/766)
* fix: deleting the last message in a conversation will prompt failure [#771](https://github.com/openimsdk/openim-sdk-core/pull/771)
* fix: the bug where isEnd for fetching message history is not working correctly. [#773](https://github.com/openimsdk/openim-sdk-core/pull/773)
* fix: solve uncorrect log in revoke handle. [#777](https://github.com/openimsdk/openim-sdk-core/pull/777)
* fix: solve uncorrect delete temp file. [#784](https://github.com/openimsdk/openim-sdk-core/pull/784)
* Fix：Change check reinstall logic [#789](https://github.com/openimsdk/openim-sdk-core/pull/789)

### Refactors
* refactor: change log and avoid nil array. [#728](https://github.com/openimsdk/openim-sdk-core/pull/728)
* refactor: remove batchListener. [#729](https://github.com/openimsdk/openim-sdk-core/pull/729)

### Builds
* build: implement changelog generate. [#748](https://github.com/openimsdk/openim-sdk-core/pull/748)
* build: implement create Pre-release PR from Milestone. [#746](https://github.com/openimsdk/openim-sdk-core/pull/746)
* build: remove uncorrect schedule. [#782](https://github.com/openimsdk/openim-sdk-core/pull/782)
* build: add bot PR merged filter. [#788](https://github.com/openimsdk/openim-sdk-core/pull/788)

### Others
* @qmarliu made their first contribution in https://github.com/openimsdk/openim-sdk-core/pull/771 [#771](https://github.com/openimsdk/openim-sdk-core/pull/771)

**Full Changelog**: [v3.8.1...v3.8.2](https://github.com/openimsdk/openim-sdk-core/compare/v3.8.1...v3.8.2)

