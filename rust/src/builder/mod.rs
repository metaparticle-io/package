pub mod docker;
use super::run_docker_process;


pub trait Builder {
    fn build(&self, dir: &str, image: &str);
    fn push(&self, image: &str);
}