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


pub fn nothing() {
    let package = Package {
        repository: Some("repository".to_string()),
        verbose: Some(false), 
        quiet:  Some(false),
        builder: Some("builder".to_string()),
        publish: Some(false)
    };
    println!("Nothing... {:?}", package);
}



// pub fn containerize(runtime: Runtime, package: Package) {
//     println!("Hello Metaparticle! package: {:?}", package);
//     println!("Hello Metaparticle! runtime: {:?}", runtime);
// }

pub fn containerize<F>(f: F, runtime: Runtime, package: Package) where F: Fn() {
    println!("Hello Metaparticle! package: {:?}", package);
    println!("Hello Metaparticle! runtime: {:?}", runtime);
    f()
}