package io.metaparticle;

import java.io.OutputStream;

public interface Builder {
    boolean build(String dir, String image, OutputStream stdout, OutputStream stderr);

    boolean push(String image, OutputStream stdout, OutputStream stderr);
}