<h1 align="center" style="border-bottom: none">
    <b>
        <a href="https://doc.rentsoft.cn/">openim-sdk-core</a><br>
    </b>
</h1>
<h3 align="center" style="border-bottom: none">
      ⭐️  Used in IOS, Android, PC and other platforms  ⭐️ <br>
<h3>

<p align=center>
<a href="https://goreportcard.com/report/github.com/OpenIMSDK/openim-sdk-core"><img src="https://goreportcard.com/badge/github.com/OpenIMSDK/openim-sdk-core" alt="A+"></a>
<a href="https://github.com/OpenIMSDK/openim-sdk-core/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22"><img src="https://img.shields.io/github/issues/OpenIMSDK/Open-IM-Server/good%20first%20issue?logo=%22github%22" alt="good first"></a>
<a href="https://github.com/OpenIMSDK/openim-sdk-core"><img src="https://img.shields.io/github/stars/OpenIMSDK/openim-sdk-core.svg?style=flat&logo=github&colorB=deeppink&label=stars"></a>
<a href="https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg"><img src="https://img.shields.io/badge/Slack-100%2B-blueviolet?logo=slack&amp;logoColor=white"></a>
<a href="https://github.com/OpenIMSDK/openim-sdk-core/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/badge/Language-Go-blue.svg"></a>
</p>

</p>

<p align="center">
    <a href="./README.md"><b>English</b></a> •
    <a href="./README_zh-CN.md"><b>中文</b></a>
</p>

</p>

----

## 🧩 Awesome features

OpenIM-SDK-core is a core SDK of OpenIM. 

1. Manage WebSocket long connections, responsible for creating, closing and reconnecting connections. 
2. Encoding and decoding. Encode and decode messages in binary format to achieve cross-language compatibility.
3. Implement basic protocols of OpenIM, such as login, push, etc. 
4. Provide an event handling mechanism to convert received messages into corresponding events and pass them to upper layer applications for processing.
5. Cache management. Manage user, group, blacklists, and other cache information. 
6. Provide basic IM function APIs such as sending messages, creating groups, etc. Hide the underlying implementation details from the upper layer application.


## Quickstart

> **Note**: You can get started quickly with openim-sdk-core.

<details>
  <summary>Work with Makefile</summary>

```bash
❯ make help    # show help
❯ make build   # build binary
```

</details>
<details>
  <summary>Work with actions</summary>

Actions provide handling of PR and issue.
We used the bot @kubbot, It can detect issues in Chinese and translate them to English, and you can interact with it using the command `/comment`.

Comment in an issue:

```bash
❯ /intive
```

</details>
<details>
  <summary>Work with Tools</summary>

```bash
❯ make tools
```

</details>
<details>
  <summary>Work with Docker</summary>

```bash
$ make deploy
```

</details>


## Contributing & Development

OpenIM Our goal is to build a top-level open source community. We have a set of standards, in the [Community repository](https://github.com/OpenIMSDK/community).

If you'd like to contribute to this openim-sdk-core repository, please read our [contributor documentation](https://github.com/OpenIMSDK/openim-sdk-core/blob/main/CONTRIBUTING.md).

## community meeting

We welcome everyone to join us and contribute to openim-sdk-core, whether you are new to open source or professional. We are committed to promoting an open source culture, so we offer community members neighborhood prizes and reward money in recognition of their contributions. We believe that by working together, we can build a strong community and make valuable open source tools and resources available to more people. So if you are interested in openim-sdk-core, please join our community and start contributing your ideas and skills!

We take notes of each [biweekly meeting](https://github.com/OpenIMSDK/Open-IM-Server/issues/381) in [GitHub discussions](https://github.com/OpenIMSDK/Open-IM-Server/discussions/categories/meeting), and our minutes are written in [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

openim-sdk-core maintains a [public roadmap](https://github.com/OpenIMSDK/community/tree/main/roadmaps). It gives a a high-level view of the main priorities for the project, the maturity of different features and projects, and how to influence the project direction.

## about OpenIM

### common

+ https://github.com/OpenIMSDK/automation: OpenIM Automation, cicd, and actions, Robotics.
+ https://github.com/OpenIMSDK/community: Community Management for OpenIM

### OpenIM **Links**

Contains some common parts of the OpenIM community.

+ https://github.com/OpenIMSDK/automation: OpenIM Automation, cicd, and actions, Robotics.
+ https://github.com/OpenIMSDK/openim-sdk-core: The IMSDK implemented by golang can be used in IOS, Android, PC and other platforms.
+ https://github.com/OpenIMSDK/openim-sdk-core: Instant messaging IM server.
+ https://github.com/OpenIMSDK/community: Community Management for OpenIM.

### SDKs

+ [openim-sdk-core](https://github.com/OpenIMSDK/openim-sdk-core): A cross-platform SDK implemented in golang that can be used in IOS, Android, PC, and other platforms.
+ [Open-IM-SDK-iOS](https://github.com/OpenIMSDK/Open-IM-SDK-iOS): An iOS SDK generated based on openim-sdk-core, available for developers to reference.
+ [Open-IM-SDK-Android](https://github.com/OpenIMSDK/Open-IM-SDK-Android): An Android SDK generated based on openim-sdk-core, available for developers to reference.
+ [Open-IM-SDK-Flutter](https://github.com/OpenIMSDK/Open-IM-SDK-Flutter): A Flutter SDK generated based on Open-IM-SDK-iOS and Open-IM-SDK-Android, available for developers to reference.
+ [Open-IM-SDK-Uniapp](https://github.com/OpenIMSDK/Open-IM-SDK-Uniapp): A uni-app SDK generated based on Open-IM-SDK-iOS and Open-IM-SDK-Android, available for developers to reference.

### Demos

+ [Open-IM-iOS-Demo](https://github.com/OpenIMSDK/Open-IM-iOS-Demo): An iOS demo based on Open-IM-SDK-iOS, available for developers to reference.
+ [Open-IM-Android-Demo](https://github.com/OpenIMSDK/Open-IM-Android-Demo): An Android demo based on Open-IM-SDK-Android, available for developers to reference.
+ [Open-IM-Flutter-Demo](https://github.com/OpenIMSDK/Open-IM-Flutter-Demo): A Flutter demo based on Open-IM-SDK-Flutter, available for developers to reference.

## Used By

OpenIM is used by the following companies ,let's write it down in [ADOPTER](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md).

Please leave your use cases in the comments [here](https://github.com/OpenIMSDK/Open-IM-Server/issues/379).

## License

[openim-sdk-core](https://github.com/OpenIMSDK/openim-sdk-core) is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/OpenIMSDK/openim-sdk-core/tree/main/LICENSE) for the full license text.

## Thanks to our contributors!

<a href="https://github.com/OpenIMSDK/openim-sdk-core/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=OpenIMSDK/openim-sdk-core" />
</a>
