[![Build Status](https://travis-ci.org/ssalvatori/rest2command.svg?branch=master)](https://travis-ci.org/ssalvatori/rest2command)

[![Coverage Status](https://coveralls.io/repos/github/ssalvatori/rest2command/badge.svg?branch=master)](https://coveralls.io/github/ssalvatori/rest2command?branch=master)

# rest2command
Http server written in golang to execute some command

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

# Files
* dist/rest2command.sh (init.d file)
