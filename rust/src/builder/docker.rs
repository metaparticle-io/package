use super::run_docker_process;
use Builder;

pub struct DockerBuilder{}

impl Builder for DockerBuilder {
    fn build(&self, dir: &str, image: &str) {
        run_docker_process(vec!["build", &*format!("-t{}", image), dir]);
    }
    fn push(&self, image: &str) {
        run_docker_process(vec!["push", image]);
    }
}