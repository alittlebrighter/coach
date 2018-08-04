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

### Usage
Once you've started a new session just continue using your terminal as you normally would and `coach` will prompt you to save frequently run commands.
![Terminal Usage](https://lh3.googleusercontent.com/7g05U_6JQ1P0baLvZENCKbwmA4RaMam7gHkNYBKuKGKjR68WN7y5gl5yWeCuJo9DYGwQ1xcMsQjEH1m4bMLS7sR_e3ZUf3VIeXaGDw4UzbOEjekDELXlo4Su7HbXY3_R5uzVCw3A36KzCyNKEkM7ynji3-tS9ASCb_HpvvPl0tu-HVPO0nxikuDhTvZTU_oDDh8GWPjwZPRmqCLDyBY61wllsRq8_cSX9wQoktYjuURY3IN_dCCXWKLE0Amw82ffrlSuPwHVYyWB0461OyfCEJPi5VWANJye9QL4V9B-4oGtL7x7zDB4Fzc5fmZvq6AmPHUDQFW6evgpWZ-Ggt5kv7L7iIoStL-mZQde7UjXiy3HBywkvatbysNWgfJVoDPvDXmA32hb4pSNspi6Bc-hBcY0ZSd43fLq_qZ1eGbg3iYJOivONy_TbjXx29vqIF8yfnQXRa6ScwPN3xFWZXgppmtMoW9JIyldKVgfpaK9DF0kL4C8RDgWfk_Z49S0wDXsPAtm7raiyIW0a1UBuaNVoGDxKpLV7JkyGxilDxIRApWYPQEl7dXuWEqKnNts8x5U3yn0_RAlgaouVS6sG52OlR7mzz1Zl0zixAkpFBg=w1310-h549-no)

To add more lines you can run:
```
$ coach doc -e get-weather
```
and `coach` will pull up your script along with its metadata inside of the editor specified by `$EDITOR`
![Edit Script](https://lh3.googleusercontent.com/nGKyXKpqv6iFCmovdoDF-Gk5WDrng4Hygn_U-C3wrvOBsrvZt7cV2kSg1pTXQMrTXJmsn4z8A9kOA-_sdabz3x1if2qfWYOnXrBcm6QnWc-fJdTgwZpD63uiebYru6uNUGt4sOkf6jeUz_Ux6rqea2yZtq5EbFhnhllBJuHIGefrRUIud6EIZjRNZtaCCECziGKiJyF6bJy6GSWSUxRYa-1xhxsgRja2MT7GQcZkVOFwsGEblDOpARzykz5Ke44E7gZ-iyjjB53vqtUrTvQFlQd13MtS09bg4M1kwyvgbmvvTkR55u9zyUqedpuW0rtVUmqN2MDHzwejafU8sYEVrUmtGSfUWP78lnbDK1T6vJ0oZqEoULRvDROMs1YuRClYGCkJnTgJz9LhMnfw8B-oui6H1Y0Px6f4AG6C0qxT7bvIiPDQP45E1j3mFdfj2rRwK3XNvYUL7_UyyqrYBENZqctDcJyGM8oR_-0dNhP66g2lgR0IiYsDehdTalNOq8LePZdk036fT6fsOZAsesPlpnPmkknmXfze0RvMTI_KudEhMlYB8NPieG4V6ov68GVQiUlJG_kTpZkuJykcZFEGdM3iUgLzv-FcnyCjVqE=w1531-h201-no)
Just save and quit to keep your changes in `coach`.

Run `coach --help` to see other available options. 

### Notes
Try it out, submit an issue and/or PR if you see room for improvement.

- Rudimentary support for other interpreters.  So long as the interpreter takes a file name as a first argument, and you enter the `SHELL` value as it appears on your `$PATH` it should work.  `bash` is the default.
- `coach` history ignores commands starting with 'coach'
- Multi-line scripts are now supported.  Run `coach doc -l [n] [alias] [tags] [documentation...]` to pull in the previous n lines of history into the documented script.  Or you can also run `coach doc -e [alias]` to edit the script inside of the text editor set in `$EDITOR`.

### Contributing

1. Fork it
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create new Pull Request