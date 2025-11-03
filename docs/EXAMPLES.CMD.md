
**Get the logs for a job**:
```bash
docker compose -f 07.compose.ci.yml logs stress-test
```


What you can do:
- the pipeline is a set of jobs and dependencies between them
- each job is a service in the docker-compose file
- you can run a specific job by using the `docker compose up <service-name>` command
- you can view the logs of a specific job by using the `docker compose logs <service-name>` command
- you can view the status of all jobs by using the `docker compose ps` command


You can get a lot of information with Docker Desktop as well.

You can use Docker Model Runner to run AI models as part of your CI/CD pipeline.

You can use Compose Watch to automatically restart jobs when files change.