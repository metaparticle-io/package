use super::run_docker_process;
use Builder;

pub struct DockerBuilder{}

impl Builder for DockerBuilder {
    fn build(&self, dir: String, image: String) {
        run_docker_process(vec!["build".to_string(), "-t".to_string()+&image, dir]);
    }
    fn push(&self, image: String) {
        run_docker_process(vec!["push".to_string(), image.to_string()]);
    }
}