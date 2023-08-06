sudo docker build --target release -t ghcr.io/snarechops/stream-bot:$GITHUB_RUN_ID .
echo $GITHUB_TOKEN | sudo docker login ghcr.io -u $GITHUB_ACTOR --password-stdin
sudo docker push ghcr.io/snarechops/stream-bot:$GITHUB_RUN_ID
