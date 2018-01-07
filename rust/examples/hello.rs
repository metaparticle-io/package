extern crate metaparticle;

fn run() {
    println!("Hello World!");
}

fn main() {
    let runtime = metaparticle::Runtime{
        ports: Some(80),
        ..Default::default()
    };
    let package = metaparticle::Package{
        name: "hello".to_string(),
        ..Default::default()
    };
    metaparticle::containerize(run, runtime, package)
}