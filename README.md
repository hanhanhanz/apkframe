# Apkframe


Tools for finding framework used inside an APK for recon process in apk exploitation. Pinpointing the right framework of an APK will affect the exploit methodology. The signature used is specifically created based on multiple APK comparation with various frameworks

### Installation

```sh
$ go build apkframe.go
$ ./apkframe
```
### Requirement

Apktool

### Usage

```sh
$ ./apkframe.go -h
Usage of apkframe:
  -a    decompile apk with 
  -d string
        specify output directory of apktool
```



### Example
```sh
$ go run apkframe.go -a -d Spaces_by_Wix_base.apk
1 reactNative signature found
13 corona signature found

```

```sh
$ ls
AndroidManifest.xml  apktool.yml  assets  kotlin  original  res  smali  smali_classes2  unknown
$
$ go run apkframe.go -d Gyroscope_base
3 reactNative signature found
```
