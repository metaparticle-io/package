use super::run_docker_process;
use super::Runtime;
use Executor;

pub struct DockerExecutor{}

impl Executor for DockerExecutor {
    fn cancel(&self, name: String) {
        run_docker_process(vec!["stop".to_string(), name.clone()]);
        run_docker_process(vec!["rm".to_string(), "-f".to_string(), name.clone()]);
    }

    fn logs(&self, name: String) {
        run_docker_process(vec!["logs".to_string(), "-f".to_string(), name]);
    }

    fn run(&self, image: String, name: String, config: Runtime) {
        let mut args = vec!["run".to_string(), "-d".to_string(), "--name".to_string(), name];
        if let Some(port) = config.ports {
            args.extend(vec!["-p".to_string(), port.to_string()]);
        }
        args.extend(vec![image]);
        run_docker_process(args);
    }
}

