<h1 align="center" style="border-bottom: none">
    <b>
        <a href="https://doc.rentsoft.cn/sdks/introduction">openim-sdk-core</a><br>
    </b>
</h1>
<h3 align="center" style="border-bottom: none">
      ⭐️  适用于 iOS、Android、PC、Web（WebAssembly）及其他平台  ⭐️ <br>
<h3>

<p align=center>
<a href="https://goreportcard.com/report/github.com/openimsdk/openim-sdk-core"><img src="https://goreportcard.com/badge/github.com/openimsdk/openim-sdk-core" alt="A+"></a>
<a href="https://github.com/openimsdk/openim-sdk-core/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22"><img src="https://img.shields.io/github/issues/OpenIMSDK/Open-IM-Server/good%20first%20issue?logo=%22github%22" alt="good first"></a>
<a href="https://github.com/openimsdk/openim-sdk-core"><img src="https://img.shields.io/github/stars/OpenIMSDK/openim-sdk-core.svg?style=flat&logo=github&colorB=deeppink&label=stars"></a>
<a href="https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg"><img src="https://img.shields.io/badge/Slack-100%2B-blueviolet?logo=slack&amp;logoColor=white"></a>
<a href="https://github.com/openimsdk/openim-sdk-core/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/badge/Language-Go-blue.svg"></a>
</p>

</p>

<p align="center">
    <a href="./README.md"><b>English</b></a> •
    <a href="./README_zh-CN.md"><b>中文</b></a>
</p>

</p>

----


## 🧩 Features
<!--BEGIN_DESCRIPTION-->
OpenIM-SDK-core is the core SDK of OpenIM, serving as the cross-platform foundation for all open-source OpenIM SDKs (excluding mini web).
All open-source OpenIM SDKs (except mini web) are built upon this core layer, ensuring consistency, stability, and seamless cross-platform integration.
<!--END_DESCRIPTION-->

- [x] Network management with intelligent heartbeat
- [x] Message encoding and decoding
- [x] Local message storage
- [x] Relationship data synchronization
- [x] IM message synchronization
- [x] Cross-platform communication and callback management

 - Supported Platforms
   - [x] Windows
   - [x] MacOS
   - [x] Linux
   - [x] iOS
   - [x] Android
   - [x] Web (WebAssembly)
   - [ ] Mini Web


## Quickstart

> **Note**: This section guides you on how to quickly connect to the server and get OpenIM-SDK-core running.

### 🚀 Connect to the Server and Run

Follow these steps to quickly set up and run OpenIM-SDK-core by simulating an app environment using test files.

1. **Enter the `test` directory**
   ```bash
    # This folder contains unit test files for all interface functions of OpenIM-SDK-core,  
    # used to simulate an app connecting to the server for login testing.  
   cd test
   ```
2. **Modify the configuration file**
   > [Set up your own server beforehand.](https://github.com/openimsdk/open-im-server.git)
 - Open the config file in the test directory.
 - Update the following fields with your server information:
   ```json
    {
    "APIADDR": "http://your-server-api-address",
    "WSADDR": "ws://your-server-websocket-address",
    "UserID": "your-test-user-id"
    }

   ```
3. **Run test functions to simulate an app using the SDK**
 - Identify the test function you want to execute (The `init` file has already completed the SDK initialization and login logic. You can now call other functions).
   ```bash
   go test -run TestFunctionName
   ```
 - Example: Running the login test
    ```bash
    go test -run Test_GetAllConversationList
    ```
Now, you can use the test cases to simulate real SDK usage, just like an actual app.


## 📦 Build and Package for Different Platforms

Once the SDK is tested successfully, you can build and package it for various platforms:

- **Android/iOS**

Refer to [this guide](./docs/CHANGELOG.md) for detailed instructions on building and packaging for Android and iOS.
- **WebAssembly**

Navigate to the `wasm/cmd` directory and run the following command to build the WebAssembly package:
  ```bash
  make wasm  # Ensure Go is installed
  ```
If you are on Windows, use the following command instead:
  ```bash
  mingw32-make wasm  # Ensure MinGW64 is installed
  ```
- **Windows, MacOS, Linux**

Refer to [this repository](https://github.com/openimsdk/openim-sdk-cpp.git) for platform-specific build instructions.

## Contributing & Development

OpenIM Our goal is to build a top-level open source community. We have a set of standards, in the [Community repository](https://github.com/openimsdk/community).

If you'd like to contribute to this openim-sdk-core repository, please read our [contributor documentation](https://github.com/openimsdk/openim-sdk-core/blob/main/CONTRIBUTING.md).

## community meeting

We welcome everyone to join us and contribute to openim-sdk-core, whether you are new to open source or professional. We are committed to promoting an open source culture, so we offer community members neighborhood prizes and reward money in recognition of their contributions. We believe that by working together, we can build a strong community and make valuable open source tools and resources available to more people. So if you are interested in openim-sdk-core, please join our community and start contributing your ideas and skills!

We take notes of each [biweekly meeting](https://github.com/openimsdk/Open-IM-Server/issues/381) in [GitHub discussions](https://github.com/openimsdk/Open-IM-Server/discussions/categories/meeting), and our minutes are written in [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

openim-sdk-core maintains a [public roadmap](https://github.com/openimsdk/community/tree/main/roadmaps). It gives a a high-level view of the main priorities for the project, the maturity of different features and projects, and how to influence the project direction.

## about OpenIM

### common

+ https://github.com/openimsdk/automation: OpenIM Automation, cicd, and actions, Robotics.
+ https://github.com/openimsdk/community: Community Management for OpenIM

### OpenIM **Links**

Contains some common parts of the OpenIM community.

+ https://github.com/openimsdk/automation: OpenIM Automation, cicd, and actions, Robotics.
+ https://github.com/openimsdk/openim-sdk-core: The IMSDK implemented by golang can be used in iOS, Android, PC and other platforms.
+ https://github.com/openimsdk/openim-sdk-core: Instant messaging IM server.
+ https://github.com/openimsdk/community: Community Management for OpenIM.

### SDKs

+ [openim-sdk-core](https://github.com/openimsdk/openim-sdk-core): A cross-platform SDK implemented in golang that can be used in iOS, Android, PC, and other platforms.
+ [Open-IM-SDK-iOS](https://github.com/openimsdk/Open-IM-SDK-iOS): An iOS SDK generated based on openim-sdk-core, available for developers to reference.
+ [Open-IM-SDK-Android](https://github.com/openimsdk/Open-IM-SDK-Android): An Android SDK generated based on openim-sdk-core, available for developers to reference.
+ [Open-IM-SDK-Flutter](https://github.com/openimsdk/Open-IM-SDK-Flutter): A Flutter SDK generated based on Open-IM-SDK-iOS and Open-IM-SDK-Android, available for developers to reference.
+ [Open-IM-SDK-Uniapp](https://github.com/openimsdk/Open-IM-SDK-Uniapp): A uni-app SDK generated based on Open-IM-SDK-iOS and Open-IM-SDK-Android, available for developers to reference.

### Demos

+ [Open-IM-iOS-Demo](https://github.com/openimsdk/Open-IM-iOS-Demo): An iOS demo based on Open-IM-SDK-iOS, available for developers to reference.
+ [Open-IM-Android-Demo](https://github.com/openimsdk/Open-IM-Android-Demo): An Android demo based on Open-IM-SDK-Android, available for developers to reference.
+ [Open-IM-Flutter-Demo](https://github.com/openimsdk/Open-IM-Flutter-Demo): A Flutter demo based on Open-IM-SDK-Flutter, available for developers to reference.

## Used By

OpenIM is used by the following companies ,let's write it down in [ADOPTER](https://github.com/openimsdk/community/blob/main/ADOPTERS.md).

Please leave your use cases in the comments [here](https://github.com/openimsdk/Open-IM-Server/issues/379).

## License

For more details, please refer to [here](./LICENSE).


## Thanks to our contributors!

<a href="https://github.com/openimsdk/openim-sdk-core/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=OpenIMSDK/openim-sdk-core" />
</a>
