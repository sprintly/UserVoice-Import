# Overview

Have a UserVoice account and need to send support requests to
developers? This is a one-way ticket creator for Sprint.ly.

# How to use

1. Copy settings_example.json to settings.json
2. Change the values to be correct for you. 
    - Your Sprint.ly API key is available on your user profile page.
    - Your product id is the number in your url, for example 1 in the url:   https://sprint.ly/product/1/items/
    - To create a UserVoice API setup, its under Settings > Channels > API
    - It is important your UserVoice API key is marked as Trusted.
3. Install the requirements in a virtualenv with pip. (See: http://guide.python-distribute.org/pip.html )
4. `python import.py 4938` where 4938 is the ticket number you want to import.


