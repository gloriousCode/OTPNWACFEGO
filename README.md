# OTPNWACFEGO
## One Time Password, Now With A Cool Front End in GO
A OTP generator for Windows, MacOS and Linux so you don't have to get out your phone every two seconds in this 2FA world

Unlike [OTPCMDGO](https://github.com/gloriousCode/OTPCMDGO), this will have the ability to add new entries from the application.

![OTPNWACFEGO](https://i.imgur.com/KUYITRR.gif)

## Requirements
- Chrome installed
- Golang installed
- Your 2FA secret codes

## Installation
- Run build-{{your-operating-system}}
- Run the output executable

## Configuration
- Codes are stored in `data.json`
- Enter the name and the secret for your code. eg
`{"Name":"Hello", "Secret": "QBGN657ZHIQ34QIA"}`

## Current Features
- GUI for viewing your 2FA OTPs
- A countdown timer
- A button to click copy

## Upcoming Features
- Encryption of your secrets
- Add new keys in-app
- Delete existing keys in-app

## Final notes
This isn't looking to be a rock-solid secure application. If you are at all concerned with your 2FA secrets being stolen - which is a legitimate concern - then do not use this. Thing will become slightly better when `data.json` encryption is implemented, but I am not responsible for anything happening to your secrets or OTPs being stolen

### Built with:
- https://github.com/pquerna/otp
- https://github.com/zserge/lorca
