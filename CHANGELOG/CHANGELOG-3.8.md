## [v3.8.3-patch.14](https://github.com/openimsdk/openim-sdk-core/releases/tag/v3.8.3-patch.14) 	(2026-06-10)

### New Features
* feat: sorted send msg [#1065](https://github.com/openimsdk/openim-sdk-core/pull/1065)
* feat: wasm support application layer ping (#936) [#1103](https://github.com/openimsdk/openim-sdk-core/pull/1103)

### Bug Fixes
* fix: add safe submodule in workflows in v3.8.3-patch. [#1009](https://github.com/openimsdk/openim-sdk-core/pull/1009)
* fix: resolve SQL injection in SearchLocalMessages in v3.8.3-patch branch. [#1023](https://github.com/openimsdk/openim-sdk-core/pull/1023)
* fix: Resolved the issue where the delete message method only performs local deletion when the seq is 0 in v3.8.3-patch branch. [#1027](https://github.com/openimsdk/openim-sdk-core/pull/1027)
* fix: test [#1038](https://github.com/openimsdk/openim-sdk-core/pull/1038)
* fix: create group group [#1039](https://github.com/openimsdk/openim-sdk-core/pull/1039)
* fix: solve incorrect conversationID generate and solve incorrect error handle in v3.8.3-patch [#1041](https://github.com/openimsdk/openim-sdk-core/pull/1041)
* fix: delete sending failed msg [#1055](https://github.com/openimsdk/openim-sdk-core/pull/1055)
* fix: copier.Copy missing field (#1066) [#1067](https://github.com/openimsdk/openim-sdk-core/pull/1067)
* fix: infinite loop caused [#1072](https://github.com/openimsdk/openim-sdk-core/pull/1072)
* fix: refactor message handling and improve attached info parsing [#1073](https://github.com/openimsdk/openim-sdk-core/pull/1073)
* fix: optimize conversation message handling and improve message synchronization [#1078](https://github.com/openimsdk/openim-sdk-core/pull/1078)
* fix: optimize conversation message handling and improve message synchronization  [#1080](https://github.com/openimsdk/openim-sdk-core/pull/1080)
* fix: can not receive OnJoinedGroupAdded notification [#1081](https://github.com/openimsdk/openim-sdk-core/pull/1081)
* fix the problem that  abnormal groups are still synced with error, then failed to sync [#1086](https://github.com/openimsdk/openim-sdk-core/pull/1086)
* fix: handle record not found errors in group and friend list retrieval [#1100](https://github.com/openimsdk/openim-sdk-core/pull/1100)
* fix(chat_log): repair abnormal messages with sendTime=0 to restore no… [#1106](https://github.com/openimsdk/openim-sdk-core/pull/1106)

### Others
* V3.8.3 patch opti [#1088](https://github.com/openimsdk/openim-sdk-core/pull/1088)
* bugfix(arr): ensure 16KB page size support on v3.8.3-patch [#1104](https://github.com/openimsdk/openim-sdk-core/pull/1104)
* @dsx137 made their first contribution in https://github.com/openimsdk/openim-sdk-core/pull/1103 [#1103](https://github.com/openimsdk/openim-sdk-core/pull/1103)

**Full Changelog**: [v3.8.3-patch.10...v3.8.3-patch.14](https://github.com/openimsdk/openim-sdk-core/compare/v3.8.3-patch.10...v3.8.3-patch.14)

