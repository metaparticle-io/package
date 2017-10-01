package io.metaparticle.examples.web;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;
import static io.metaparticle.Metaparticle.Containerize;

import java.util.function.Function;

public class Main {
    @Runtime(ports={5000},
             executor="aci")
    @Package(repository="brendanburns",
             jarFile="target/metaparticle-package-0.1-SNAPSHOT-jar-with-dependencies.jar")
    public static void main(String[] args) {
        Containerize(() -> {
            while (true) {
                System.out.println(args + "Hello World!");
                try {
                    Thread.sleep(10 * 1000);
                } catch (InterruptedException ex) {}
            }
        });
    }
}