# How to use
First fill your server detail first in .env.example file, and then rename it into .env

Let say you want forward your port 80 to your server so it can be accessed by public then, you can use
tobira.exe -local=80 -remote=6550

or you if you want random remote port you can just use
tobira.exe -local=80 

# to build this 
go get
go build