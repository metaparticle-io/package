package io.metaparticle.annotations;

public @interface Runtime {
    public int replicas() default 1;

    public String executor() default "docker";

    public int[] ports() default {};

    public boolean publicAddress() default false;
}
