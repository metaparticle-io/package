use std::process;
use super::Runtime;
use Executor;

pub struct DockerExecutor{}

impl Executor for DockerExecutor {
    fn placeholder(&self) {
        println!("Docker Executor");
    }


    fn run(&self, image: String, name: String, config: Runtime) {

    }
}

