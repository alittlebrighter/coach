## Coach [![Build Status](https://travis-ci.org/alittlebrighter/coach.svg?branch=master)](http://travis-ci.org/alittlebrighter/coach)

Coach helps you document Ops processes by recording all of your commands and prompting you to save an alias and document frequently run commands.  This 
documentation can then be queried later on the command line.

Go beyond autosuggestion tools, quickly attach queryable tags and documentation to all of your scripts without leaving the command line.

### Install
`go get -u github.com/alittlebrighter/coach/cmd/coach`

To use `coach` to monitor your command line usage, add the following to your `.bashrc` file: 
```
function prompt {
    coach history --record "$(history 1)" # you do actually need this
}

PROMPT_COMMAND=prompt
```

### Usage
Once you've started a new session just continue using your terminal as you normally would and `coach` will prompt you to save frequently run commands.
![Terminal Usage](https://i.imgur.com/ear5FUW.jpg)

To add more lines you can run:
```
$ coach doc -e get-weather
```
and `coach` will pull up your script along with its metadata inside of the editor specified by `$EDITOR`
![Edit Script](https://i.imgur.com/QOUR1UY.png)
Just save and quit to keep your changes in `coach`.

Run `coach --help` to see other available options. 

### Notes
Try it out, submit an issue and/or PR if you see room for improvement.

- Return all documented scripts with the tag "?" or just run `coach doc`
- Basic Windows support has been implemented.  You can document, save, and run batch and Powershell scripts.  You can manually add lines to your history but a good way of automatically tracking history in the Windows command prompt or Powershell hasn't been found yet.
- Rudimentary support for other interpreters.  So long as the interpreter takes a file name as a first argument, and you enter the `SHELL` value as it appears on your `$PATH` it should work.  `bash` is the default.
- `coach` history ignores commands starting with 'coach'

### Contributing

1. Fork it
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create new Pull Request