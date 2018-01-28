
mod builder;
mod executor;

use builder::Builder;
use executor::Executor;

use std::env;
use std::error::Error;
use std::fs::File;
use std::ffi::OsStr;
use std::io::prelude::*;
use std::path::Path;
use std::process;

#[derive(Debug)]
pub struct Runtime {
    pub replicas: Option<u64>,
    pub shards: Option<u64>,
    pub url_shard_pattern: Option<String>,
    pub executor: String,
    pub ports: Option<u64>,
    pub public_address: Option<bool>,
}


impl Default for Runtime {
    fn default() -> Runtime {
        Runtime {
            replicas: None,
            shards: None,
            url_shard_pattern: Some("".to_string()),
            executor: "docker".to_string(),
            ports: None,
            public_address: Some(false)
        }
    }
}

#[derive(Debug)]
pub struct Package {
    pub name:       String,
	pub repository: String,
	pub verbose:    Option<bool>,
	pub quiet:      Option<bool>,
	pub builder:    String,
	pub publish:    Option<bool>,
}


impl Default for Package {
    fn default() -> Package {
        Package {
            name: "".to_string(),
            repository: "".to_string(),
            verbose: Some(false), 
            quiet:  Some(false),
            builder: "docker".to_string(),
            publish: Some(false)
        }
    }
}

pub fn run_docker_process(args: Vec<&str>) {
    let name = args[0].clone();
    let mut child = process::Command::new("docker")
        .args(&args)
        .spawn()
        .expect(&format!("failed to execute 'docker {name}'", name=name));

    let status = child.wait()
        .ok().expect(&format!("couldn't wait for 'docker {name}'", name=name));

    if !status.success() {
        match status.code() {
            Some(code) => panic!("'docker {}' failed with code {:?}", name, code),
            None       => panic!("'docker {}' failed", name)
        }
    }
}

fn in_docker_container() -> bool {
    let env_var_key = OsStr::new("METAPARTICLE_IN_CONTAINER");
    let env_var = env::var(env_var_key);
    if let Ok(value) = env_var {
        return value == "true" || value == "1";
    }

    let mut buffer  = String::with_capacity(256); // kind of a guess on initial capacity

    if let Ok(mut file) = File::open("/proc/1/cgroup") {
        if let Ok(_) = file.read_to_string(&mut buffer) {
            return buffer.contains("docker");
        }
    }
    false
}

fn executor_from_runtime(executor_name: String) -> Box<Executor> {
    let executor : Box<Executor> = match executor_name.as_ref() {
        "docker" => Box::new(executor::docker::DockerExecutor{}),
        _ => panic!("Unsupported executor type {}", executor_name),
    };
    return executor;
}

fn build_from_runtime(builder_name: String) -> Box<Builder> {
    let builder : Box<Builder> = match builder_name.as_ref() { 
        "docker" => Box::new(builder::docker::DockerBuilder{}),
        _ => panic!("Unsupported builder type {}", builder_name),
    };
    builder

}

fn write_dockerfile(name: &str) {
    let dockerfile = &format!("FROM ubuntu:16.04
    COPY ./{name} /tmp/{name}
    CMD /tmp/{name}
    ", name=name);
    let path = Path::new("Dockerfile");
    let display = path.display();

    let mut file = match File::create(&path) {
        Err(why) => panic!("couldn't create {}: {}",
                           display,
                           why.description()),
        Ok(file) => file,
    };

    if let Err(why) = file.write_all(dockerfile.as_bytes()) {
        panic!("Could not write dockerfile at {} because {}", display, why.description());
    }
}

pub fn containerize<F>(f: F, runtime: Runtime, package: Package) where F: Fn() {
    if in_docker_container() {
        f();
    } else {

        if package.repository.len() == 0 {
            panic!("A package must be given a 'repository' value");
        }
        if package.name.len() == 0 {
            panic!("A package must be given a 'name' value");
        }
        let image = &format!("{repo}/{name}:latest", repo=package.repository, name=package.name);
        
        write_dockerfile(&package.name);
        let builder = build_from_runtime(package.builder.clone());

        let arg_0 = env::args().nth(0).unwrap();
        let path = Path::new(&arg_0);
        let docker_context = Path::new(path.parent().unwrap());

        builder.build(docker_context.to_str().unwrap(), image);

        let executor = executor_from_runtime(runtime.executor.clone());
        executor.run(image, &package.name, runtime);
        executor.logs(&package.name);
    }
}