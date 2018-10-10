## Coach [![Build Status](https://travis-ci.org/alittlebrighter/coach.svg?branch=master)](http://travis-ci.org/alittlebrighter/coach)

Coach helps you maximize you and your team's performance on the command line.  

### Install
`go get -u github.com/alittlebrighter/coach/cmd/coach`

To use `coach` to monitor your command line usage, add the following to your `.bashrc` file: 
```
function prompt {
    coach history --record "$(history 1)"
}

PROMPT_COMMAND=prompt
```

There is currently no way to track command history in any other shell (I'm open to ideas on how to resolve that).

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

### Docs

Until I get proper documentation up for now you can use the following:
- `coach` - run `coach --help` 
- `coach-grpc-server` - the `.proto` files that the gRPC service implements can be found in the `protobuf` directory.  There are some comments there but further documentation will be provided once I have an interface to plug in security components.
- `coach-grpc-web` - documentation waiting on security components

### Roadmap
- interfaces to add authentication via Go plugins
- authorization
- interface to allow launching deeper analysis of command history

### Contributing

1. Fork it
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create new Pull Request