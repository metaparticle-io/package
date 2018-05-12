pub mod docker;
use super::Runtime;
use super::run_docker_process;

pub trait Executor {
    fn cancel(&self, name: &str);
    fn logs(&self, name: &str);
    fn run(&self, image: &str, name: &str, config: Runtime);
}