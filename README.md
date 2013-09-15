# Overview

Have a UserVoice account and need to send support requests to
developers? This is a one-way ticket creator for Sprint.ly.

[![Build Status](https://drone.io/github.com/sprintly/UserVoice-Import/status.png)](https://drone.io/github.com/sprintly/UserVoice-Import/latest)

# How to install

1. Create a configuration file in the format denoted in the "Configuration File" format below.
2. Download the [latest binary](https://drone.io/github.com/sprintly/UserVoice-Import/files/uservoice_import)
3. Move it to a place on your `PATH`

# Usage

To import a Uservoice ticket as a new Sprint.ly item, execute the
command below. This will create a new Sprint.ly ticket (who's URL is
printed after it's created) with the contents of the initial user
complaint. It will make a note on the uservoice ticket with the URL to
the Sprint.ly item.

`uservoice_import --uservoice_ticket_number 1234`

To link a support request to an existing Sprint.ly item, try this
command. It will link the uservoice ticket (as a note) to the
Sprint.ly item and will create a similar link in Sprint.ly.

`uservoice_import --uservoice_ticket_number 1234 --link_to=5465`

# Configuration File

Create a configuation file (default location in
`~/.config/sprintly/uservoice_import.ini`) with the following format:

```ini
[Uservoice]
subdomain=USERVOICE_SUBDOMAIN_HERE
apikey=USERVOICE_API_KEY_HERE
apisecret=USERVOICE_SECRET_KEY_HERE
[Sprintly]
baseurl=https://sprint.ly/api/
email=SPRINTLY_EMAIL_ADDRESS
apikey=SPRINTLY_API_KEY
productid=YOUR_PRODUCT_ID
```

It is important that the API key you generate in Uservoice is a
trusted API client. This allows us the access to fetch items and make
notes.

# Development

This project is managed through a standard github pull-request
workflow. Each change should be accompanied by tests which follow a
similar style (if applicable) to those already in the repository.

If you have any questions, please get in touch!