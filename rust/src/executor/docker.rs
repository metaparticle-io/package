use super::run_docker_process;
use super::Runtime;
use Executor;

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
        let mut args = vec!["run", "-d", "--name", name];
        if let Some(port) = config.ports {
            args.extend(vec!["-p",&*format!("-p{}", port)]);
        }
        args.extend(vec![image]);
        run_docker_process(args);
    }
}

