package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/getlantern/systray"
	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"

	"log"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Define the structure for the YAML file
type Link struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Sub  []Link `yaml:"sub"`
}

type Links []Link

func main() {
	onExit := func() {
		// Cleaning logic here
	}
	systray.Run(onReady, onExit)
}

func onReady() {
	// Load the PNG image
	img, err := loadPNGImage("icon.png")
	if err != nil {
		systray.Quit()
		panic(errors.Wrap(err, "loading PNG image failed"))
	}

	// Set the image as the icon
	systray.SetIcon(img)
	// systray.SetTitle("Əlçatan")
	systray.SetTooltip("Əlçatan")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
		fmt.Println("Quit now...")
	}()

	// Add other menu items here

	systray.AddSeparator()

	// Read the contents of the YAML file
	data, err := ioutil.ReadFile("links.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Unmarshal the YAML contents into the Links struct
	var links Links
	err = yaml.Unmarshal(data, &links)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML data: %v", err)
	}

	// Print the contents of the Links struct
	for _, link := range links {
		fmt.Printf("Name: %s\nURL: %s, %v\n", link.Name, link.URL, link.Sub)

		mShowMsg := systray.AddMenuItem(link.Name, link.URL)

		if len(link.Sub) > 0 {
			for _, linkSub := range link.Sub {
				mSub := mShowMsg.AddSubMenuItem(linkSub.Name, linkSub.URL)
				go func() {
					<-mSub.ClickedCh
					open.Start(linkSub.URL)
				}()
			}
		} else {
			go func() {
				<-mShowMsg.ClickedCh
				open.Start(link.URL)
			}()
		}
	}

}

func loadPNGImage(path string) ([]byte, error) {
	// Open the PNG file
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "opening PNG file failed")
	}
	defer f.Close()

	// Decode the PNG image
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "decoding PNG image failed")
	}

	// Encode the image as PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, errors.Wrap(err, "encoding image as PNG failed")
	}

	return buf.Bytes(), nil
}
