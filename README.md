# NAME
:sushi:(salmon) - salmon is .* bot

# SYNOPSIS
```
# run migration subcommand
salmon flake

# run plugin maker for bot
cd .../salmon/core/command
salmon flake -m <commandname>

# run register plugin for bot
cd .../salmon/core/command
salmon flake -r

# run plugin on cil
salmon flake -e <command> <args>

# run slack mode
salmon slack

```
# Description
***THIS IS EXPERIMENTAL!***  
salmon is multiplatform bot. LINE, Slack, etc...  
but, salmon have subcommand for migration. So, you can try plugin that you have been created.  

# FUTURE
I want to create server for machine learning with python(django?). In order to  use machine learning for salmon. 

# TASK
- [ ] add auto deploy plugin(update.go)
- [ ] support Slack
- [ ] support LINE
- [ ] add more plugins
- [ ] create machine learning api
- [ ] use database
- [ ] use docker component
