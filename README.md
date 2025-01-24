# GITHP

This is a TUI (Terminal User Interface) app designed to simplify your daily commit process.

---

### Building

To build locally, you'll need to clone this repository by running:
```shell
 git clone https://github.com/audreyteles/githp.git
```

Next, install golang at your machine, to see how to install, [click here](https://go.dev/doc/install).

In terminal, access the **githp** folder, and then you should be able to build:

 ```shell
  cd githp && go build -o githp cmd/main.go 
 ```

A file named githp, has been created on your directory.

---

### Configuring

Once the app is on your machine, you need to give it executable permissions:

  ```shell
  sudo chmod +x githp  
  ```

Next, move the app to your local binary folder:
  ```shell
  mv githp /usr/local/bin/
  ```
Create an alias to the application:

`nano ~/.bashrc` or `nano ~/.zshrc` 

Add the following line:
```
alias githp='/usr/local/bin/githp'
```

And then. you can use in your git directories:
```
githp .
```