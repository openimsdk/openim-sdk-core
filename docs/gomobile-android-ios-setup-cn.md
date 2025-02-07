# 使用 gomobile 编译 openim-sdk-core 以供 Android/iOS 使用
</p>

<p align="center">
    <a href="./README.md"><b>English</b></a> •
    <a href="./README_zh-CN.md"><b>中文</b></a>
</p>

</p>
## 环境准备
### 1.Go语言环境以及gomobile环境搭建

**1. 安装Go语言环境**

- 安装[GO语言（1.18以上版本）](https://go.dev/dl/)。

**2. 配置GOPATH环境变量**

- 通过 go env 查看 Go 的环境变量设置。以下是典型的输出示例：

  ```bash
  > go env
  set GOPATH=D:\Go\gopath
  set GOMODCACHE=D:\Go\gopath\pkg\mod
  ```
- 根据 go env 的输出，将 GOPATH 添加到系统的环境变量中：
    - **Windows**: 右键点击“此电脑”，选择“属性 -> 高级系统设置 -> 环境变量”，将 GOPATH 的路径（如 D:\Go\gopath）添加到 Path 变量中。
    - **Mac/Linux**: 在终端中编辑 ~/.bashrc 或 ~/.zshrc，添加以下内容：
      ```bash
         export GOPATH=/Users/youruser/go  # 根据你的实际路径替换
         export PATH=$PATH:$GOPATH/bin
      ```

**3. 安装 gomobile 和 gobind**

在 Go 1.18 或更高版本中，执行以下命令安装最新的 gomobile 和 gobind：
```bash
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
```
**4. 初始化 gomobile**

执行以下命令完成 gomobile 的初始化：
```bash
gomobile init
```
### 在 Windows 平台编译 Android AAR 包

#### 环境要求

1. **确保 Go 和 gomobile 配置正确**：执行 `gomobile version` 验证工具是否安装成功。

2. **安装 Android 开发环境**：确保已安装最新版本的 **Android Studio**。

3. **配置 Android NDK**：下载适用于 Windows 的 NDK（推荐版本：`r20b`），将其解压到 Android SDK 的 `ndk-bundle` 目录下。例如：

   ```
   C:\Users\Admin\AppData\Local\Android\Sdk\ndk-bundle
   ```

4. **配置 Make 命令支持（可选）**：

    - Windows 不自带 `make`，可以安装 MinGW64：

        1. 下载 MinGW64 并安装。
        2. 安装后，将 MinGW 的 `bin` 目录（如 `C:\mingw64\bin`）添加到系统的环境变量中。
- 如果无法安装 `make`，你可以直接使用 `gomobile bind` 命令完成编译。

#### 编译AAR包

进入你的项目根目录，例如 `openim-sdk-core`，运行以下命令以编译 Android AAR 包：
```bash
gomobile bind -v -trimpath -ldflags="-s -w" -o ./open_im_sdk.aar -target=android ./open_im_sdk/ ./open_im_sdk_callback/
```

##### **注意事项**：

1. 确保网络畅通，编译过程中需要从 GitHub 下载依赖包。首次运行可能需要较长时间。
2. 如果使用 MinGW64，可以执行：
   ```bash
    mingw32-make android
   ```
3. 编译完成后，你会看到类似以下的输出：
    ```bash
    aar: jni/armeabi-v7a/libgojni.so
    aar: jni/arm64-v8a/libgojni.so
    aar: jni/x86/libgojni.so
    aar: jni/x86_64/libgojni.so
    aar: R.txt
    aar: res/
    ...
    ```

​       `open_im_sdk.aar` 文件会生成在当前目录。

4. 将生成的 AAR 包通过 Android Studio 的本地导入方式引入到项目中，即可使用导出的函数和回调接口。



### 在 macOS 平台编译 Android AAR 包和 iOS xcframework

#### 环境要求

1. 安装 Xcode：确保已安装 Xcode（建议版本：15.4 或更高）。
2. 安装 Android Studio：确保已安装并配置好 Android SDK 和 NDK（Mac 推荐 NDK 版本为：20.0.5594570）。

#### 编译Android AAR包

进入项目根目录（如 `openim-sdk-core`），运行以下命令：


```bash
make android
```
编译完成后，将生成的 AAR 包导入到 Android Studio 项目中。

#### 编译iOS xcframework库

1. 进入项目根目录。

2. 执行以下命令以编译 iOS 的 xcframework：
   ```bash
    make ios
   ```

3. 编译完成后，将生成的 `.xcframework` 文件导入到 Xcode 项目中。



### 常见问题及解决方案

1. **卡在写入 `go.mod` 文件**

- 可能是网络原因导致依赖下载失败。尝试设置 Go 的代理：
  ```bash
  go env -w GOPROXY=https://proxy.golang.org,direct
  ```
  如果在国内，可以使用：
    ```bash
  go env -w GOPROXY=https://goproxy.cn,direct
    ```

2. **找不到 `make` 命令**

- 在 Windows 上，确保安装 MinGW64 并使用 `mingw32-make` 命令代替 `make`。
- 或直接运行 `gomobile bind` 命令。

3. **NDK 版本兼容问题**

- 如果遇到 NDK 版本兼容问题，尝试切换到推荐的版本（`r20b` 或 `20.0.5594570`），并确保路径正确。



通过 `gomobile`，可以轻松将 Go 语言编写的openim-sdk-core代码打包为 Android 的 AAR 包或 iOS 的 xcframework 库，方便在移动平台中集成和使用。
