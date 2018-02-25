package io.metaparticle;

import java.io.OutputStream;

import io.metaparticle.annotations.Runtime;

public interface Executor {
    public boolean run(String image, String name, Runtime config, OutputStream stdout, OutputStream stderr);
    public boolean logs(String name, OutputStream stdout, OutputStream stderr);
    public boolean cancel(String name);
}