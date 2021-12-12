# 'Appened
'Appened (rhymes with Happened) is a minimalist append-only note service, meant to be paired a variety of clients.

## Running 'Appened

You can either 'Appened as local processes or inside of docker containers using the included scripts.

### Setup

All that needs to be done is to create a place to keep your data and to set your authentication token.

1. Create a directory called `data/` in the root of the project, this is where your folios will be stored as CSVs.
2. Set an environment variable `APPENED_AUTH_TOKEN`, if you're using the docker scripts in `containers.sh` place it inside a file named `.env`

### Docker Scripts

Optionally, you may use the inlcluded docker scripts. These are just wrappers to simplify the boilerplate when running them. 

Note that `run` will stop and remove any running container of that name.

```sh
./containers.sh build NAME_HERE
./containers.sh run NAME_HERE
```
## Go SDK

This library includes a simple library that wraps the REST API. 

## Twilio Client

There is a simple twilio client to allow interacting with 'Appened over SMS.

Note that it is set up such that only one phone number is whitelisted. It will only work with one phone number.

### Running The Client

To run the client

1. Create a file called `config.json` in `appened/clients/twilio`.
2. Fill it as follows:

```json
{
        "accountSid": "TWILIO_ACCOUNT_SID_HERE",
        "authToken": "TWILIO_ACCOUNT_AUTH_TOKEN_HERE",
        "clientNumber": "YOUR_PHONE_NUMBER",
        "twilioNumber": "YOUR_TWILIO_PHONE_NUMBER",
        "appenedToken": "APPENED_AUTH_TOKEN",
        "appenedURL": "APPENED_HOST"
}
```

### Usage

```
h: help message
lf: list folios
cf <folioName>: create folio
df <folioName>: delete folio
ln <folioName>: list notes in folio
lna <folioName>: list all notes in folio, including done
lnd <folioName>: list all done notes in folio
dn <folioName> <number>: Toggle done on note at number
a <folioName> <msg>: append note to folio
```

