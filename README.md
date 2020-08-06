# runescape-pwfill-go
An implementation of my Python RuneScape password automation in Go

Automates filling your password into the RS3 client

If you're like me and use a password manager (which you really should be doing these days), having to manually type those generated passwords can be a real annoyance. Most sites/applications allow you to copy/paste EXCEPT the Runescape NXT client. 

This tool will allow you to read your password from a file or password manager and autofill this into the RS3 client.

### Requirements
rspw requires Go 1.6 or above due to the [ini package](https://github.com/go-ini/ini) dependency

this tool also requires a KeePass database to be configured with a structure similar to the following:

```  
  Root
    General
      runescape password entry
```
Note that root and general are folders in KeePass

#### Known limitations
* The Go version of this tool is working mainly in Windows at the moment
* Linux support is being worked on
