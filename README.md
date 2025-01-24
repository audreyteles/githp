# GITHP

This is a TUI (Terminal User Interface) app designed to simplify your daily commit process.

---

### Building

To build locally, you'll need to clone this repository by running:

```shell
git clone https://github.com/audreyteles/githp.git
```

Next, install Go on your machine. To see how to install it, [click here](https://go.dev/doc/install).

In the terminal, access the **githp** folder, and then you should be able to build:

```shell
cd githp && go build -o githp cmd/main.go
```

A file named `githp` will be created in your directory.

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

Create an alias for the application:

```shell
nano ~/.bashrc
```

or

```shell
nano ~/.zshrc
```

Add the following line:

```
alias githp='/usr/local/bin/githp'
```

After that, you can use it in your git directories:

```shell
githp .
```
