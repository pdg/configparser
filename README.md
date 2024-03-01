# Configparser

This implements a simple config format parser. The format is similar to the one used by the `caddy` webservers [Caddyfile](https://caddyserver.com/docs/caddyfile). 

## Sample Usage

~~~go

package main

import (
  "encoding/json"
  "log"
  "strings"
)

func main() {

  input := strings.NewReader(`
  
    distribution "debian" {

      suite "stable"
      architecture "amd64"

      repository {
        security
        backports
        updates
      }

    }

    packages {
      
      openssh-server {
        permit-root-login false
        password-auth false
        pubkey-auth true
      }
    
      keyboard-configuration {
        layout "de"
      }

    }

  `)

  config, err := Parse(input)
  if err != nil {
    log.Fatalln(err)
  }

  j, err := json.MarshalIndent(config, "", "  ")
  if err != nil {
    log.Fatalln(err)
  }

  println(string(j))

}

~~~

## License 

MIT License
