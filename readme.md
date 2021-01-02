# How to use

First, you should build this app/download the binary file. (Executable binary files may be provided in future)

Then copy the url of the novel you want to download

Run the application with an argument like `./kakuspider https://kakuyomu.jp/works/117735405488213152482131524`

If you do not give out the url, the app will run in interactive mode

Wait for a while and enjoy the novel!

provided flags:

| name | description                             | default |
|------|-----------------------------------------|---------|
| -n   | use number as filename instead of title | false   |
| -j   | numbers of goroutines to download       | 5       |
