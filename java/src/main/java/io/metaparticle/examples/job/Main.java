package io.metaparticle.examples.job;

import static io.metaparticle.Metaparticle.Containerize;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;

public class Main {
    @Runtime(iterations=4,
             executor="metaparticle")
    @Package(repository="brendanburns",
             verbose=true,
             jarFile="target/metaparticle-package-0.1-SNAPSHOT-jar-with-dependencies.jar")
    public static void main(String[] args) {
        Containerize(() -> {
            System.out.println("Hello batch job!");
        });
    }
}