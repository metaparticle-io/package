extern crate metaparticle;

fn run() {
    println!("Hello World!");
}

fn main() {
    let runtime = metaparticle::Runtime{
        ..Default::default()
    };
    let package = metaparticle::Package{
        ..Default::default()
    };
    metaparticle::containerize(run, runtime, package)
}