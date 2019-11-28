# retype

re-types your standard input

```
echo -e "passwd\noldpassword\nnewpassword\nnewpassword\n" | ./retype -interval 100
```

# Installation

If you have go
```
go get github.com/flashvoid/retype
```

If you have docker
```
mkdir /tmp/output
docker run -it -v /tmp/output:/go/bin golang go get github.com/flashvoid/retype
```

# Keymaps
This program takes characters from command line and types them on virtual keyboard, however,
not all keyboards are born equal, retype needs a way to map characters from standard input on
keyboard buttons this is accomplished with keymaps file.

Retype has standard en-us (pc-105) keymap file embedded into it which can be viewed with
`retype -dump-keymap`.  
If it doesn't work for you then make/find your own and point retype to it using either `retype -keymap myfile` or `KEYMAPS=myfile retype`.

# Permissions
Retype uses `/dev/uinput` to create virtual keyboards, it will need read/write access to do so.  

* Example udev rules to make it work
```
echo KERNEL==\"uinput\", GROUP=\"$USER\", MODE:=\"0660\" | sudo tee /etc/udev/rules.d/99-$USER.rules
sudo udevadm trigger
```

or just chmod

# Credits
This is a simple tool that should have been written many years ago by people much smarter then me.  
Credits to 
* linux community for ./dev uinput  
* github.com/bendahl for all the real work.

# TODO
* Support more non-printable  characters
* Function keys?

# GPL2

# Contibutions
Welcome
