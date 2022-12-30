### Markdownscanner (TM)
Are your links broken? Let's find out: https://mdscanner.dcalvo.dev/ 
Please note that this project is under (sporadic) development and it's not finished. I just need to get around polishing a few things...

### Okay but now for real
While signing up to contribute to k8s, I found a broken link on the sign up process. [This was then my my first contribution.](https://github.com/kubernetes/community/pull/4304)  
I then wondered: How many other markdown links are broken in open source projects? As it turns out, a lot of them.  
This tool will hopefully help me find and fix these links.

## TODOs:
- Fix this first: Why do some links show as 404 in the report even though they are not 404s?
    - The etcd repo has a few occurences of this
    - You should be able to replica this fairly easily! (Maybe you can even create a cobra command named "check link" to see which result your application will return for that internally)
- Hey test accessing the S3 bucket before launching the program, if you scan everything but can't upload the results, that's terrible!
    - Perhaps put that in the root command together with reading the config too!
- Do a ctrl+f for "deprecated", a few of your functions became deprecated!
- Document your functions!
- Document the settings that config.yaml accepts!
- Put the S3 and templating stuff in other file away from repository.go? HmmMmMmm
- Check for the git command when starting mdscanner, you need it! (for all use cases?)
- Amahgad remove your hardcoding of "tmp" filesystem path on the config file!
- Write tests man, `shame_bell_gameofthrones.gif`
- On the init() function for all commands, you're checking for the config.yaml file by copying and pasting it around. Any way to do this for all commands? Maybe on root.go?

## Random thougths
- It seems that having functions be part of the markdown link type make them a bit inflexible