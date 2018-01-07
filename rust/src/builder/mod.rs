pub mod docker;
use super::run_docker_process;


pub trait Builder {
    fn build(&self, dir: String, image: String);
    fn push(&self, image: String);
}