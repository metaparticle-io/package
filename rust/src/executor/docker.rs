use super::run_docker_process;
use super::Runtime;
use Executor;

use std::iter::Iterator;

pub struct DockerExecutor{}

impl Executor for DockerExecutor {
    fn cancel(&self, name: &str) {
        run_docker_process(vec!["stop", name]);
        run_docker_process(vec!["rm", "-f", name]);
    }

    fn logs(&self, name: &str) {
        run_docker_process(vec!["logs", "-f", name]);
    }

    fn run(&self, image: &str, name: &str, config: Runtime) {
        let mut args = vec![ "run".to_string(),
            "-d".to_string(),
            "--rm".to_string(), 
            "--name".to_string(), 
            name.to_string()];

        if let Some(port) = config.ports {
            args.push(format!("-p {port}", port=port));
        }

        args.push(image.to_string());
        let args_refs = args.iter()
                .map(|a| a.as_ref())
                .collect();

        run_docker_process(args_refs);
    }
}

