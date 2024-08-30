# Meeting

This feature is described in [../problem.feature](../problem.feature) file in the *Scenario: Problem discussion*.
Below are some implementation details of how we want this work. 

## Client

People log into our system using a web browser.

They start a meeting by calling each other. 
Right now the tool for setting up the call is Jitsi Meet, 
but we might change it later or provide several options.

When a person starts the call our client app starts 
streaming their audio to our server.

## Server

The server authenticates users.

It receives and processes audio streams from user calls.

The server groups received audio streams into meetings. 
It creates a transcript for each meeting.
It also creates a summary of each meeting.
