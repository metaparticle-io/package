extern crate metaparticle;

fn run() {
    println!("Hello World!");
}

fn main() {
    let runtime = metaparticle::Runtime{
        executor: Some("docker".to_string()),
        ..Default::default()
    };
    let package = metaparticle::Package{
        builder: Some("also docker".to_string()),
        ..Default::default()
    };
    metaparticle::containerize(run, runtime, package)
}