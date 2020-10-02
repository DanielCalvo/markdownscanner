### Actionable items
- Explicitly rename things like "projectRepos" for "projectRepoUrls" when they're different types
- Set up AWS user with proper permissions (terraform this part too)
    - Cron set up
- Set up a working Dockerfile
- Is it possible to write a Prometheus Exporter?
- Check for stray print messages
- Is variable shadowing OK?
- Check why some links don't get templated properly like this one (see the 404 in here): https://mdscanner.dcalvo.dev/kubernetes/sample-apiserver.html 
- Stick to a coding standard pls (is it c or config? r or repo?)
- Order the functions in the packages (put close together what happens close together)
- Refine HTTP ERR (check the error and take it from there, it's more valuable to know that something is a timeout rather than a 404)
- Can you have a status page? You can have an independent routine running that updates a status page
- You are still missing tests!

### Code TODOs
- Redo the Dockerfile for the executable thing
- Careful with what you export (https://blog.golang.org/organizing-go-code)
- https://peter.bourgon.org/go-best-practices-2016/#repository-structure
- https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1

- About Cobra: https://github.com/spf13/cobra
    - I think you don't need cobra for now as you'll only run with one parameter, but maybe investigate more functionality later
    - Like running the program with a parameter (--project) or to save results locally and don't upload them to S3 (--dont-upload)
    - Viper remains to be explored later as well: https://github.com/spf13/viper

### Things to do
- Scan projects that have not been scanned the longest first
- Check how having using a more sophisticated logging method could improve the program
- Set a timeout / handle timeouts on HTTP checks properly. Further info here: http://networkbit.ch/golang-http-client/#minimal
- Create usage instructions
- Don't forget to document which links get ignored somewhere
- Find a way to implement header checks for markdown files. That might be tough...
- Recheck go.mod and general install (create instructions)
- There are no unit tests. You should create some!
- Remove the excessive newlines on the html if you can, more info here: https://github.com/golang/go/issues/9969

#### Terraform
- Create the following:
    - S3 bucket (not public)
    - User to access the bucket
    - A policy for the user granting access only to that bucket

### md scanner ideas
- Create a cloudwatch event that launches an instance, does something and then exits?

- How do I stop and start Amazon EC2 instances at regular intervals using Lambda?
    - https://aws.amazon.com/premiumsupport/knowledge-center/start-stop-lambda-cloudwatch/

- How to connect to S3?
    - https://docs.aws.amazon.com/AmazonS3/latest/dev/security-best-practices.html

- Can you use Lambda pre signed URLs to upload files securely?
    - https://medium.com/@lakshmanLD/upload-file-to-s3-using-lambda-the-pre-signed-url-way-158f074cda6c

- How is Lambda local development done again? Can you launch a function that uploads a mock file to s3?

#### Maybe
- Just run things on an ec2 instance for now, worry about migrating to t2 micro later

#### Possible feature ideas
- Can you take the feature toggle approach as defined in the cloud ops book?
    - Application could save something either to S3 or local filesystem
    - "How to feature toggle golang"
    - https://featureflags.io/go-feature-flags/