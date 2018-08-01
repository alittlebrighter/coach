## Coach [![Build Status](https://travis-ci.org/alittlebrighter/coach.svg?branch=master)](http://travis-ci.org/alittlebrighter/coach)

Coach helps you document Ops processes by recording all of your commands and prompting you to save an alias and document frequently run commands.  This 
documentation can then be queried later on the command line.

### Install
`go get -u github.com/alittlebrighter/coach/cmd/coach`

To use `coach` to monitor your command line usage, add the following to your `.bashrc` file: 
```
function prompt {
    coach history --record "$(history 1)" # you do actually need this
}

PROMPT_COMMAND=prompt
```
Run `coach --help` to see available options. 

### Notes
Try it out, submit an issue and/or PR if you see room for improvement.

- Only works with `bash` at the moment, more support coming.
- `coach` history ignores commands starting with 'coach'
- Saving a script with the same alias as a previous script will destroy the previous script.
- Multi-line scripts are now supported.  Run `coach doc -l [n] [alias] [tags] [documentation...]` to pull in the previous n lines of history into the documented script.  Or you can also run `coach doc -e [alias]` to edit the script inside of the text editor set in `$EDITOR`.

### Contributing

1. Fork it
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create new Pull Request