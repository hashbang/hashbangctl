User config tool for #!

Lets users change ldap details, ssh keys, etc.


## Documentation 

Typing ``hashbangctl`` when you're connected to one of the hashbang servers will start the programm and give you some details about your user.
Follow the instructions on the screen to change your 

- username 
- shell(`cat /etc/shells` will list all available shells) 
- add or delete SSH Keys
- import SSH Keys from github

Remember to save your changes, quitting the program will discard your changes.


## Installation

You don't need to this on hashbang hosts as it's already installed, otherwise `pip install git+https://github.com/hashbang/hashbangctl`

Upgrading is managed via Ansible (see [admin-tools](https://github.com/hashbang/admin-tools)).
