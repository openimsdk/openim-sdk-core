## [v3.8.3-patch.10](https://github.com/openimsdk/openim-sdk-core/releases/tag/v3.8.3-patch.10) 	(2025-07-16)

### New Features
* feat: add sdk releaser workflows in v3.8.3-patch. [#982](https://github.com/openimsdk/openim-sdk-core/pull/982)

### Bug Fixes
* fix: use a custom close code for javascript websocket active disconnection to prevent close failures, according to the WebSocket RFC documentation. [#923](https://github.com/openimsdk/openim-sdk-core/pull/923)
* fix: recycle javascript promise function manually to prevent memory l… [#922](https://github.com/openimsdk/openim-sdk-core/pull/922)
* fix: SyncLoginUserInfo [#929](https://github.com/openimsdk/openim-sdk-core/pull/929)
* fix: repeat notification trigger [#948](https://github.com/openimsdk/openim-sdk-core/pull/948)
* fix: group request trigger [#950](https://github.com/openimsdk/openim-sdk-core/pull/950)
* fix: group request trigger [#953](https://github.com/openimsdk/openim-sdk-core/pull/953)
* fix: modify the app's foreground and background status switch, first modify it locally to prevent the SDK from triggering offline messages and not triggering the normal newMessage. [#972](https://github.com/openimsdk/openim-sdk-core/pull/972)
* fix: update version file in sdk releaser in v3.8.3-patch branch. [#985](https://github.com/openimsdk/openim-sdk-core/pull/985)
* fix: remove new line in update version file in v3.8.3-patch. [#989](https://github.com/openimsdk/openim-sdk-core/pull/989)
* fix: update wasm archive have wasm_exec.js in v3.8.3-patch branch. [#995](https://github.com/openimsdk/openim-sdk-core/pull/995)
* fix: solve incorrect distinct in IncrSync and remove goroutine in syncAndTriggerReinstallMsgs func in v3.8.3-patch. [#1000](https://github.com/openimsdk/openim-sdk-core/pull/1000)

**Full Changelog**: [v3.8.3...v3.8.3-patch.10](https://github.com/openimsdk/openim-sdk-core/compare/v3.8.3...v3.8.3-patch.10)

