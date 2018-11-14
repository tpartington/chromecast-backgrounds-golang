# chromecast-backgrounds-golang

A simple golang port of [dconnolly/chromecast-backgrounds](https://github.com/dconnolly/chromecast-backgrounds)

(simple because it doesn't have all the functionality and never will).

## Usage

```
Usage of ./chromecast-backgrounds-golang:
  -dir string
        the directory to download the images to
  -url string
        the chromecast homepage (default "https://clients3.google.com/cast/chromecast/home")
```

example

```
./chromecast-backgrounds-golang -dir=downloads
```

An exit code of 0 means images were downloaded.  
An exit code of 1 means no images were downloaded.  
