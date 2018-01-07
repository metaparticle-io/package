pub mod docker;
use super::Runtime;
use super::run_docker_process;

pub trait Executor {
    fn cancel(&self, name: String);
    fn logs(&self, name: String);
    fn run(&self, image: String, name: String, config: Runtime);
}