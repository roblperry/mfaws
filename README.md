## mfaws 
A little tool to help with updating session profiles generated using aws sts get-session-token

### Why

Because `aws sts get-session-token` is a pain and why let needing
to update your session token interrupt your groove for a moment longer than necessary?

### How to install

1. Install go
2. run `go get github.com/roblperry/mfaws`

### How to run

`go run github.com/roblperry/mfaws help`

or if you put ~/go/bin in your path, just `mfaws help`

### Theory of operation

mfaws is gonna call get-session-token for you and then update/creates a named profile.
You can then use that named profile to get your work done.

### Opinionated How To

Create a named profile with your static `aws_access_key_id` and 
`aws_secret_access_key`. For that sake of this conversation, let's call it `my_profile`.

Set a couple of environment variable so that all your aws commands are just
going to work for you. I use direnv, but I'm only willing to go so far in telling 
you how to live your life.

    export AWS_PROFILE=my_profile_session
    export AWS_REGION=cc-whatever-you-are-using-01

I hope it obvious that `cc-whatever-you-are-using-01` is not real and should instead
be set to whatever region you are working in, but where did `my_profile_session`
come from?  That is just the original profile named `my_profile` with `_session`
appended, which is a convention mfaws follows by default.

Now all you have to do is run mfaws.

    mfaws --profile my_profile

mfaws should prompt you for your mfa and then update your ~/.aws/credentials.


Yay, now when aws libraries see your AWS_PROFILE they will be able to pull your 
aws_access_key_id, aws_secret_access_key, and aws_session_token from your
~/.aws/credentials.

### How else can I do it?

Grrr....lots of ways.  I'll never think of them all, but I'll come back 
and open a few more doors for you. 



