
- Layouts seem to be an issue: https://github.com/golang-standards/project-layout/issues/117
- uuuh

lets just use
- cmd (for cobra commands)
- internal (for your multipurpose "package code")
- main.go (for the entrypoint)

And get going. You can check/add more stuff later, I don't think you need more folders right now

cobra add scan
- scan <single-repo>
- scan --config=as #uh, how do you do the scan from config thing?

### Copy & paste notes
- pkg -> internal
- cmd -> ???
- main.go -> ???

### Usage
- mdscanner single <repo>
- mdscanner --from-config=
- mdscanner scan-all --config=