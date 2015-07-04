#JSONIF - Simple GO JSON Config Manager
## Overview
Simple configuration manager for GOlang. Supports loading multiple JSON files. Supports overrides, i.e. default, dev, staging, prod.

The last loaded JSON files overrides previous added file.
### Syntax


#### See it in action:

	package main
	import (
		"github.com/samphomsopha/jsonfig"
		"fmt"
		)
	
	func main() {
		config := jsonfig.New()
		config.loadFile("config-files/default.json")
		config.loadFile("config-files/prod.json")
		
		//read config
		config.Get("facebook.api_key")
	}

