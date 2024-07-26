# Batch-worker for manuspect written in go

#### Operating System

- __Linux:__ (Prefered) Ubuntu 22.04, Debian Latest
- __MacOS:__ 10.14 or later
- __Windows:__ (10, 11 + WSL2) Use WSL for development

#### Tools

- __[Docker Engine](https://docs.docker.com/engine/install/)__ latest
- __[Make](https://www.gnu.org/software/make/#download)__ latest
- __[Golang](https://go.dev/dl/)__ >= 1.22


### Fast start (ubuntu x86_64 and WSL2 (ubuntu) only)

For fast install deps (except Docker) use this script

```sh
sudo apt update
sudo apt install curl unzip zip tar git make

curl -LO https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
#sudo rm -rf ~/go
tar -C ~/ -xzf go1.22.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
echo 'export GOPATH=~/.go' >> ~/.bashrc
rm go1.22.4.linux-amd64.tar.gz
source ~/.bashrc

go version

make dev
```

### Before start

Create `.env` in working dir using `.env.sample`:

```sh
cp ./.env.sample ./.env
```

### Building

To build a binary file, run the command in the console in the project directory:

```shell
make build
```

## Data flow

1. Get data-message from queue (NATS)

    - userID

    - recordId

    - timestamp (milis)

    - filePath

2. Download data (batch) from S3 using filePath

3. Unpack batch

4. Convert video to images

5. Upload images to S3

6. Send requests to processing service (each requestr consint of image and meta data)

    - get mata data from `.csv` file

7. Send acknowledge to queue if OK

8. Send to another queue processed files names