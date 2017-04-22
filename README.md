[![Build Status](https://travis-ci.org/ssalvatori/http_server_go.svg?branch=master)](https://travis-ci.org/ssalvatori/http_server_go)
[![Coverage Status](https://coveralls.io/repos/github/ssalvatori/http_server_go/badge.svg?branch=master)](https://coveralls.io/github/ssalvatori/http_server_go?branch=master)

# http_server_go
Http server written in golang to execute some commands 

# Environment variables

```bash
LOG_LEVEL (debgu|info|warn|error|fatal)
PORT (default 8891)
FILE_CONFIGURATION (default ./configuration.json)
```

# Configuration file
```json
[
  {
  "url": "/change/opt1",
  "command": "/opt/opt1.sh"
  },
  {
    "url": "/change/opt2",
    "command": "/opt/opt2.sh"
  }
]
```
