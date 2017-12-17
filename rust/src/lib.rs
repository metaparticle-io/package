use std::env;
use std::error::Error;
use std::fs::File;
use std::ffi::OsStr;
use std::io;
use std::io::prelude::*;
use std::path::Path;


#[derive(Debug)]
pub struct Runtime {
    pub replicas: Option<u64>,
    pub shards: Option<u64>,
    pub url_shard_pattern: Option<String>,
    pub executor: Option<String>,
    pub ports: Option<u64>,
    pub public_address: Option<bool>,
}


impl Default for Runtime {
    fn default() -> Runtime {
        Runtime {
            replicas: Some(1),
            shards: Some(2),
            url_shard_pattern: Some("something".to_string()),
            executor: Some("exec".to_string()),
            ports: Some(3),
            public_address: Some(false)
        }
    }
}


#[derive(Debug)]
pub struct Package {
	pub repository: Option<String>,
	pub verbose:    Option<bool>,
	pub quiet:      Option<bool>,
	pub builder:    Option<String>,
	pub publish:    Option<bool>,
}


impl Default for Package {
    fn default() -> Package {
        Package {
            repository: Some("repository".to_string()),
            verbose: Some(false), 
            quiet:  Some(false),
            builder: Some("builder".to_string()),
            publish: Some(false)
        }
    }
}

pub trait Executor {
    fn placeholder(&self);
}



struct DockerExecutor{}

struct FakeExecutor{}

impl Executor for DockerExecutor {
    fn placeholder(&self) {
        println!("Docker executor");
    }
}

impl Executor for FakeExecutor {
    fn placeholder(&self) {
        println!("Fake executor");
    }
}

pub trait Builder {
    fn placeholder(&self);
}

struct DockerBuilder{}

impl Builder for DockerBuilder {
    fn placeholder(&self) {
        println!("Docker builder!");
    }
}

struct FakeBuilder{}

impl Builder for FakeBuilder {
    fn placeholder(&self) {
        println!("Fake builder!");
    }
}

fn in_docker_container() -> bool {
    let env_var_key = OsStr::new("METAPARTICLE_IN_CONTAINER");
    let env_var = env::var(env_var_key);
    if let Ok(value) = env_var {
        return value == "true" || value == "1";
    }

    let mut buffer  = String::with_capacity(256); // kind of a guess on initial capacity
    let mut file_result =  File::open("/proc/1/cgroup");

    if let Ok(mut file) = File::open("/proc/1/cgroup") {
        if let Ok(_) = file.read_to_string(&mut buffer) {
            return buffer.contains("docker");
        }
    }
    false
}

fn executor_from_runtime(executor_name: Option<String>) -> Box<Executor> {
    if let Some(name) = executor_name {
        let executor : Box<Executor> = match name.as_ref() {
            "docker" => Box::new(DockerExecutor{}),
            "fake" => Box::new(FakeExecutor{}),
            _ => panic!("Unsupported executor type {}", name),
        };
        return executor;
    }
    Box::new(DockerExecutor{})
}

fn build_from_runtime(builder_name: Option<String>) -> Box<Builder> {
    if let Some(name) = builder_name {
        let builder : Box<Builder> = match name.as_ref() { 
            "docker" => Box::new(DockerBuilder{}),
            "fake" => Box::new(FakeBuilder{}),
            _ => panic!("Unsupported builder type {}", name),
        };
        builder
    } else {
        Box::new(DockerBuilder{})
    }
}

fn write_dockerfile(name: &str) {
    let dockerfile = &b"FROM rust:latest
    FOO
    HELLO WORLD
    "[..];
    let path = Path::new("Dockerfile");
    let display = path.display();

    let mut file = match File::create(&path) {
        Err(why) => panic!("couldn't create {}: {}",
                           display,
                           why.description()),
        Ok(file) => file,
    };

    if let Err(why) = file.write_all(dockerfile) {
        panic!("Could not write dockerfile at {} because {}", display, why.description());
    }
}

pub fn containerize<F>(f: F, runtime: Runtime, package: Package) where F: Fn() {
    if(in_docker_container()) {
        f();
    } else {
        let executor = executor_from_runtime(runtime.executor.clone());
        executor.placeholder();

        let builder = build_from_runtime(package.builder.clone());
        builder.placeholder();
        write_dockerfile("foo");
    }
}