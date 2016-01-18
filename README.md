#go-nsot-api

An *experimental* golang client for NSOT (http://github.com/dropbox/nsot)

## Example
```
	import (
        "fmt"
        "github.com/sarguru/nsot"
	)

	func main() {
		c,_ := nsot.NewClient("EMAIL","SECRET","http://NSOTSERVER:PORT/api")
        	newSite := &nsot.SiteOpts{
                	Name: "sarguru.me",
                	Desc: "bla",
        	}
		site,err := c.CreateSite(newSite)
                if err != nil {
			panic(err)
		}
		fmt.Println(site)
	}
```
## Installing

	Use ``` go get ``` after importing in your program.

## Reference
	See [the godoc](http://www.godoc.org/github.com/sarguru/go-nsot-api)

# Supported Features

## Sites
	* Create
	* Update
	* Retrieve
	* Delete

## Networks
	* Create
	* Retrieve
	* Delete
  

# Support and Contributions

	This is an *experimental* library and will continue to evolve until it stabilizes. That said, I'll try
my best to respond to reported issues as soon as I could and respond to pull requests. Patches/Suggestions are
always welcome !

	Please feel free to contact me via GH issues or Twitter or irc (sarguru_ in freenode) if you've any issues
with the code or have any suggestions or patches!

# Copyright

	Copyright (c) Sargurunathan Mohan 2016

# LICENSE

	Apache2 - see the included LICENSE file for more information
