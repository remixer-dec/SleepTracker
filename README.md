## Public habit / mood /sleep tracker service  
- 5MB Golang executable to run in a distroless container (8mb total)
- Multiple widgets to track different stats (mood, sleep quality, hours, yes-no questions)
- Streaks and anti-streaks, leveling system, flame animation and XP adjusted to current streak
- Minimum dependencies: jwt and bluntdb

### Setup:
```sh
# 1) clone the repo with git clone and cd into it.
# 2) make sure that the container can access and persist local data direcrory 
# (it is a distroless container with nonroot user and cannot access mounts by default)
sudo install -d -o 65532 -g 65532 -m 700 ./data
# 3) run the compose command to build and deploy the service
docker compose up -d
# 4) to log in, run this command to create a link
docker exec sleeptracker /app/sleeptracker -create-join-link
# visit the link to log in
# to deploy with existing reverse proxy modify the network and ports sections in docker-compose.yml
```

#### Disclaimer
Most of the code was implemented with assistance of Claude Code, I guided it to implement requested features in a way that aligns with my vision of the project and fixed the bugs that it made.
