//#![feature(custom_attribute)]
extern crate metaparticle;

// static runtime: metaparticle::Runtime = metaparticle::Runtime{
//     replicas: Some(1444444),
//     ..Default::default()
// };

// static package: metaparticle::Package = metaparticle::Package{
//     verbose: Some(true), 
//     ..Default::default()
// };

// #[containerize(runtime, package)]
// fn main() {
//     println!("Hello World!");
//     metaparticle::nothing()
// }

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