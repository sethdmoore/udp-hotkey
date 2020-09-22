# Introduction
Send hard coded hotkeys from Windows to Linux over UDP.
Very fast, but very rigid without recompiling.

# Security?
This application is not secure and probably never will be.  
**Do not expose the application's UDP port 1111 to the internet.**  
The hotkey-receiver (UDP listener) will simply run any and all hotkeys as 
if they were actual keyboard key presses.
