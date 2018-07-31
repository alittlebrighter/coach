## Coach (alpha)

Coach helps you document Ops processes by recording all of your commands and prompting you to save an alias and document frequently run commands.  This 
documentation can then be queried later on the command line.

### Install
`go get -u github.com/alittlebrighter/coach/cmd/coach`

To use `coach` to monitor your command line usage, add the following to your `.bashrc` file: 
```
function prompt {
    coach history --record "$(history 1)"
}

PROMPT_COMMAND=prompt
```
Run `coach --help` to see available options. 

### Notes
This is an alpha quality WIP.  Try it out,  submit an issue and/or PR if you see room for improvement.

Only works with `bash` at the moment, more support coming.

### Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Added some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request