package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	// Start ChromeDriver service manually before running the script
	const seleniumURL = "http://localhost:4444"

	// Set Chrome capabilities
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Args: []string{"--disable-gpu", "--no-sandbox"},
	}
	caps.AddChrome(chromeCaps)

	// Connect to WebDriver
	driver, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Fatal("Error connecting to WebDriver:", err)
	}
	defer driver.Quit()

	// Open the login page
	loginURL := "http://localhost:3000/login"
	fmt.Println("Opening login page:", loginURL)
	err = driver.Get(loginURL)
	if err != nil {
		log.Fatal("Error opening login page:", err)
	}

	// Wait for the page to load
	time.Sleep(3 * time.Second)

	// Find and fill the username field
	usernameField, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Username, phone or email']")
	if err != nil {
		log.Fatal("Error finding username field:", err)
	}
	usernameField.SendKeys("user") // Replace with actual username

	// Find and fill the password field
	passwordField, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Password']")
	if err != nil {
		log.Fatal("Error finding password field:", err)
	}
	passwordField.SendKeys("userpass") // Replace with actual password

	// Find and click the login button
	loginButton, err := driver.FindElement(selenium.ByCSSSelector, ".auth-button")
	if err != nil {
		log.Fatal("Error finding login button:", err)
	}
	loginButton.Click()

	// Wait for login to complete
	time.Sleep(5 * time.Second)

	fmt.Println("Login successful!")
}
