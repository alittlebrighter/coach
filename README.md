## Coach

Coach is a tool to help you build a queryable library of documented scripts for everything you do on the command line.  Coach records your terminal sessions and prompts you to save and document any commands/groups of commands that are run multiple times.  This allows you to tag and write documentation for your scripts while you are in the right context.  Once you have a library built, you (or your coworkers) can search for scripts and save time reading man pages or a trip to StackOverflow.  

You have a vocabulary that you prefer when doing your work but the tools you have to use don't necessarily operate with the same vocabulary.  When you have to look up how to do whatever you're trying to do you know you should save a script so you don't have to go through the research hassle again but usually once the job is accomplished it's on to the next task and you end up looking up the task again and again.  Coach prompts you to save and document these tasks in the moment so you can document and save the script while you are thinking about it.

To use `coach` add the following to your .bashrc file: 
```
function prompt {
    coach history --record "$(history 1)"
}

PROMPT_COMMAND=prompt
```

TODO:
- [ ] intelligently parse commands (split commands joined with &&/;, group arguments surrounded by quotes/parenteses, etc.)
- [ ] add ability to interactively add arguments to commands when run
- [ ] record the directory that a command was run from