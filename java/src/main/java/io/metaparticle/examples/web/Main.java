package io.metaparticle.examples.web;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;
import static io.metaparticle.Metaparticle.Containerize;

import java.util.function.Function;

public class Main {
    @Runtime(ports={5000})
    @Package(repository="brendanburns")
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