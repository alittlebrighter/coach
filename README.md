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

```
adam@devops-1:~/$ svcName=$(journalctl --since today | grep error \
    | jq .service.name); systemctl restart $svcName
'users-api' restarted
...
adam@devops-1:~/$ svcName=$(journalctl --since today | grep error \
    | jq .service.name); systemctl restart $svcName
'widget-svc' restarted
...
adam@devops-1:~/$ svcName=$(journalctl --since today | grep error \
    | jq .service.name); systemctl restart $svcName
'webhooks' restarted

---
This command has been used 3+ times.
`coach lib [alias] [tags] [comment...]` to save and document this command.
`coach ignore` to silence this output for this command.
adam@devops-1:~/$ coach lib prod.findAndFix prod,services,restart Finds \
    any services logging errors and restarts them.
```

To add more lines you can run:
```
adam@devops-1:~/$ coach lib -e prod.findAndFix
```
and `coach` will pull up your script along with its metadata inside of the editor specified by `$EDITOR`
Just save and quit to keep your changes in `coach`.

Run `coach --help` to see other available options. 

### Docs

Until I get proper documentation up for now you can use the following:
- `coach` - run `coach --help` 
- `coach-grpc-server` - the `.proto` files that the gRPC service implements can be found in the `protobuf` directory.  There are some comments there but further documentation will be provided once I have an interface to plug in security components.
- `coach-grpc-web` - documentation waiting on security components

### Contributing

1. Fork it
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create new Pull Request