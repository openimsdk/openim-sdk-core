# Using gomobile to compile openim-sdk-core for Android/iOS

<p align="center">
    <a href="./gomobile-android-ios-setup.md"><b>English</b></a> •
    <a href="./gomobile-android-ios-setup-cn.md"><b>中文</b></a>
</p>

## Environment Setup

### 1. Go Language Environment and gomobile Setup

**1. Install Go Language Environment**

- Install [Go language (version 1.18 or higher)](https://go.dev/dl/).

**2. Configure GOPATH Environment Variable**

- Use `go env` to check Go's environment variable settings. Here's a typical output example:

  ```bash
  > go env
  set GOPATH=D:\Go\gopath
  set GOMODCACHE=D:\Go\gopath\pkg\mod
  ```

- Based on the `go env` output, add GOPATH to your system's environment variables:
  - **Windows**: Right-click "This PC", select "Properties -> Advanced system settings -> Environment Variables", add the GOPATH path (e.g., D:\Go\gopath) to the Path variable.
  - **Mac/Linux**: Edit ~/.bashrc or ~/.zshrc in terminal, add the following content:
    ```bash
       export GOPATH=/Users/youruser/go  # Replace with your actual path
       export PATH=$PATH:$GOPATH/bin
    ```

**3. Install gomobile and gobind**

In Go 1.18 or higher, execute the following commands to install the latest gomobile and gobind:

```bash
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
```

**4. Initialize gomobile**

Execute the following command to complete gomobile initialization:

```bash
gomobile init
```

### Compiling Android AAR Package on Windows Platform

#### Environment Requirements

1. **Ensure Go and gomobile are configured correctly**: Execute `gomobile version` to verify if the tools are installed successfully.

2. **Install Android Development Environment**: Ensure the latest version of **Android Studio** is installed.

3. **Configure Android NDK**: Download NDK for Windows (recommended version: `20.1.5948944`(`r20b`) ), extract it to the `ndk-bundle` directory of Android SDK. For example:

   ```
   C:\Users\Admin\AppData\Local\Android\Sdk\ndk-bundle
   ```

4. **Configure Make Command Support (Optional)**:

   - Windows doesn't come with `make`, you can install MinGW64:

     1. Download and install MinGW64.
     2. After installation, add MinGW's `bin` directory (e.g., `C:\mingw64\bin`) to system environment variables.

- If you can't install `make`, you can directly use the `gomobile bind` command for compilation.

#### Compiling AAR Package

Navigate to your project root directory, such as `openim-sdk-core`, and run the following command to compile the Android AAR package:

```bash
gomobile bind -v -trimpath -ldflags="-s -w" -o ./open_im_sdk.aar -target=android ./open_im_sdk/ ./open_im_sdk_callback/
```

##### **Notes**:

1. Ensure network connectivity, as dependency packages need to be downloaded from GitHub during compilation. First run may take a long time.
2. If using MinGW64, you can execute:
   ```bash
    mingw32-make android
   ```
3. After compilation is complete, you'll see output similar to:
   ```bash
   aar: jni/armeabi-v7a/libgojni.so
   aar: jni/arm64-v8a/libgojni.so
   aar: jni/x86/libgojni.so
   aar: jni/x86_64/libgojni.so
   aar: R.txt
   aar: res/
   ...
   ```

​ The `open_im_sdk.aar` file will be generated in the current directory.

4. Import the generated AAR package into your project through Android Studio's local import method to use the exported functions and callback interfaces.

### Compiling Android AAR Package and iOS xcframework on macOS Platform

#### Environment Requirements

1. Install Xcode: Ensure Xcode is installed (recommended version: 15.4 or higher).
2. Install Android Studio: Ensure Android SDK and NDK are installed and configured (Mac recommended NDK version: `20.1.5948944`(`r20b`) or `20.0.5594570`).

#### Compiling Android AAR Package

Navigate to the project root directory (e.g., `openim-sdk-core`) and run:

```bash
make android
```

After compilation is complete, import the generated AAR package into your Android Studio project.

#### Compiling iOS xcframework Library

1. Navigate to the project root directory.

2. Execute the following command to compile iOS xcframework:

   ```bash
    make ios
   ```

3. After compilation is complete, import the generated `.xcframework` file into your Xcode project.

### Common Issues and Solutions

1. **Stuck on writing `go.mod` file**

- This might be due to network issues causing dependency download failures. Try setting Go proxy:
  ```bash
  go env -w GOPROXY=https://proxy.golang.org,direct
  ```
  If you're in China, you can use:
  ```bash
  go env -w GOPROXY=https://goproxy.cn,direct
  ```

2. **Cannot find `make` command**

- On Windows, ensure MinGW64 is installed and use `mingw32-make` command instead of `make`.
- Or directly run the `gomobile bind` command.

3. **NDK Version Compatibility Issues**

- If you encounter NDK version compatibility issues, try switching to the recommended version (`20.1.5948944`(`r20b`) or `20.0.5594570`) and ensure the path is correct.

Through `gomobile`, you can easily package Go language written openim-sdk-core code into Android AAR packages or iOS xcframework libraries, making it convenient for integration and use on mobile
