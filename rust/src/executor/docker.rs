use super::run_docker_process;
use super::Runtime;
use Executor;

pub struct DockerExecutor{}

impl Executor for DockerExecutor {
    fn cancel(&self, name: &str) {
        run_docker_process(vec!["docker", "stop", name]);
        run_docker_process(vec!["docker", "rm", "-f", name]);
    }

    fn logs(&self, name: &str) {
        run_docker_process(vec!["docker", "logs", "-f", name]);
    }

    fn run(&self, image: &str, name: &str, config: Runtime) {
        let mut ports = String::new();
        let mut args = vec!["docker", "run", "-d", "--rm", "--name", name];
        
        if let Some(port) = config.ports {
            ports.push_str(&format!("-p {port}", port=port));
            args.push(&ports);
        }

        args.push(image);

        run_docker_process(args);
    }
}

