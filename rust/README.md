# Metaparticle for Python
Metaparticle/Package is a collection of libraries intended to
make building and deploying containers a seamless and idiomatic
experience for developers.

This is the implementation for Rust.

## Introduction
Metaparticle/Package simplifies and centralizes the task of
building and deploying a container image.

Here is a quick example.

Consider this simple Python application:
```rust
fn main() {
    println!("Hello World!");
}
```

To containerize this application, you need to use the `metaparticle` crate and
the `containerize` wrapper function like this:

```rust
fn run() {
    println!("Hello World!");
}

fn main() {
    let runtime = metaparticle::Runtime{
        ..Default::default()
    };
    let package = metaparticle::Package{
        name: "hello".to_string(),
        ..Default::default()
    };
    metaparticle::containerize(run, runtime, package)
}
```

When you run this application, instead of printing "Hello world", it first packages itself as a container, and
then (optionally) deploys itself inside that container.

## Tutorial

```bash
git clone https://github.com/metaparticle-io/package/
cd package/rust

cargo build --example hello
./target/debug/examples/hello
```