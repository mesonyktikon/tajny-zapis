# Tajny Zapis

This is the backend for [Tajny Zapis](https://tajnyzapis.dev), a toy
project I made for the simple pleasure of designing and building a
software system.

Four years ago I began writing aws lambdas in javascript and then moved onto
typescript. With this project I learned to love writing lambdas in go.
It's very enjoyable. The code compiles into a single binary. With the js/ts
ecosystem I resorted to webpacking. The go approach is far superior.
No bundling system to tweak, compilation is a single easy command,
compilation is super fast, single binary, no runtime to bootstrap on warmup,
no jit delay.

## Design

I like to do this by hand, with paper and pen, and I don't want to take
pictures and upload them since I don't like the look of it. So you're going
to have to take my word for it that the design is drawn out. That being said,
when time and desire permit, I like to draw on the computer.

Here is the design of creating a tajny zapis (secret note).
![](encrypt.drawio.svg)

Honestly, it is unlikely that I'll draw out the decrypt flow because making
pretty images for fun takes time, and my time is limited. But who knows,
perhaps one day you'll return and you'll find it here 😉

## Deploying

I have this in my `.zshrc` and it works well enough.
```
function deploy-tajny-zapis () {
  GOOS=linux GOARCH=arm64 go build -o dist/bootstrap main.go && \
  zip dist/lambda-handler.zip dist/bootstrap && \
  aws lambda update-function-code --function-name=tajny-zapis --zip-file fileb://dist/lambda-handler.zip
}
```

```
aws s3 sync frontend s3://tajnyzapis.dev
```