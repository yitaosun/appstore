package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
)

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

var allowedPlatform = map[string]struct{}{
	"ios":     struct{}{},
	"android": struct{}{},
	"":        struct{}{},
}
var iosBundleRegex = regexp.MustCompile(`\d{9,}`)

func guessPlatform(p string) string {
	if iosBundleRegex.MatchString(p) {
		return "ios"
	}
	return "android"
}

func main() {
	var platform string
	flag.StringVar(&platform, "p", "", "Platform: ios or android")
	flag.Parse()
	if _, ok := allowedPlatform[platform]; !ok {
		log.Fatal("Invalid platform", platform)
	}
	for _, bundle := range flag.Args() {
		bundlePlatform := platform
		if bundlePlatform == "" {
			bundlePlatform = guessPlatform(bundle)
		}
		switch bundlePlatform {
		case "ios":
			openBrowser("https://itunes.apple.com/us/app/id" + bundle)
		case "android":
			openBrowser("https://play.google.com/store/apps/details?id=" + bundle)
		}
	}
}
