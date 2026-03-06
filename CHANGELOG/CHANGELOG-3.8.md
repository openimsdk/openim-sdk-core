## [v3.8.3-patch.3](https://github.com/openimsdk/openim-sdk-core/releases/tag/v3.8.3-patch.3) 	(2025-03-07)

### Others
* 9d105d673f90ed87c6be1098672ca3cb89b31655 fix: get group member info maybe failed. (#881)
* afebca309b46078ae84a929ddeef64c6e7dc9598 fix: directly deduplicate the messages pulled from the server. (#874)
* 5500f4ddf663295354ba032208cac2162279c044 fix: sync self conversation's avatar when user's info changed. (#871)
* 85d51fd5fbe7ce70d3bf1a57a05f1ac3581088fc fix: add a manually triggered IM message synchronization mechanism to prevent message recall failure due to seq=0. (#869)
* 8f704b30300ce3f2e2e557a40090671c095d0544 fix: modify the historical message retrieval interface to address the message gap problem caused by server crashes or redis seq cache expired. (#857)
* d24833b73486701fdd9ca19b853a1f190dbd2020 optimize the freeze caused by too many friends and group applications && fix: GetConversationIDBySessionType returns a string with escape characters. (#853)
* b01a5517f78c32113811e93e93387832e14e0a3a fix: refine exception message handling to prevent duplicate messages … (#841)
* d7749bd42f78946c3112156e3e2ac65dc81c504b build: improve workflows contents. (#843)


